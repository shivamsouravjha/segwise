package helpers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"segwise/clients/postgres"
	redis_client "segwise/clients/redis"
	"segwise/models"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
	cron "github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Scheduler struct {
	Cron  *cron.Cron
	Jobs  map[string]cron.EntryID // Store triggerID -> Job ID
	Mutex sync.Mutex              // Prevent race conditions
	DB    *gorm.DB
	Redis *redis.Client
}

// StartScheduler initializes the scheduler and continuously checks for updates
func StartScheduler() *Scheduler {
	db := postgres.GetDB()
	redisClient := redis_client.RedisSession()
	s := &Scheduler{
		Cron:  cron.New(),
		Jobs:  make(map[string]cron.EntryID),
		DB:    db,
		Redis: redisClient,
	}

	// Load initial triggers
	s.syncTriggers()

	// Start a goroutine to refresh triggers every 1 minute
	go func() {
		for {
			time.Sleep(60 * time.Second)
			s.syncTriggers()
		}
	}()

	// Start the cron scheduler
	s.Cron.Start()
	return s
}

// Sync database triggers with running jobs
func (s *Scheduler) syncTriggers() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	var triggers []models.Trigger
	s.DB.Where("type = ?", "scheduled").Find(&triggers)

	activeTriggerIDs := make(map[string]bool)

	// Add or update triggers
	for _, trigger := range triggers {
		activeTriggerIDs[trigger.ID] = true

		// If trigger is new, add it
		if _, exists := s.Jobs[trigger.ID]; !exists {
			if trigger.OneTime {
				duration, err := ParseOneTimeSchedule(trigger.Schedule)
				if err == nil {
					go func(t models.Trigger) {
						time.Sleep(duration)
						s.executeScheduledTrigger(t)
					}(trigger)
					zap.L().Info("Scheduled one-time trigger: %s in %s", zap.Any("triggerID", trigger.ID), zap.Any("scheduler", trigger.Schedule))
				} else {
					zap.L().Error("Invalid one-time schedule",
						zap.String("schedule", trigger.Schedule),
						zap.Error(err))
				}
			} else {
				jobID, err := s.Cron.AddFunc(trigger.Schedule, func() {
					s.executeScheduledTrigger(trigger)
				})
				if err != nil {
					zap.L().Error("Failed to schedule trigger:", zap.Error(err))
					continue
				}
				s.Jobs[trigger.ID] = jobID
				zap.L().Info("Scheduled new trigger:", zap.Any("triggerID", trigger.ID), zap.Any("scheduler", trigger.Schedule))
			}
		}
	}

	// Remove deleted triggers
	for triggerID, jobID := range s.Jobs {
		if !activeTriggerIDs[triggerID] {
			s.Cron.Remove(jobID)
			delete(s.Jobs, triggerID)
			zap.L().Info("Removed trigger: %s", zap.Any("triggerID", triggerID))
		}
	}
}

// Execute the scheduled trigger
func (s *Scheduler) executeScheduledTrigger(trigger models.Trigger) {
	event := models.EventLog{
		TriggerID:   trigger.ID,
		TriggeredAt: time.Now(),
		Payload:     trigger.Payload,
		Type:        "scheduled",
		Status:      "active",
	}

	s.DB.Create(&event)

	// Cache in Redis for fast lookup
	eventJSON, _ := json.Marshal(event)
	s.Redis.Set("event:"+trigger.ID, eventJSON, 2*time.Hour)

	zap.L().Info("Executed scheduled trigger: %s", zap.Any("triggerID", trigger.ID))
	if trigger.OneTime {
		s.Mutex.Lock()
		defer s.Mutex.Unlock()
		if jobID, exists := s.Jobs[trigger.ID]; exists {
			s.Cron.Remove(jobID)
			delete(s.Jobs, trigger.ID)
			zap.L().Info("One-time trigger removed: %s", zap.Any("triggerID", trigger.ID))
		}
	}
}

// parseOneTimeSchedule converts different formats like "in 10 seconds" or "10s" into time.Duration
func ParseOneTimeSchedule(schedule string) (time.Duration, error) {
	duration, err := time.ParseDuration(schedule)
	if err == nil {
		return duration, nil
	}

	re := regexp.MustCompile(`in (\d+) (seconds|minutes|hours)`)
	matches := re.FindStringSubmatch(schedule)
	if len(matches) == 3 {
		amount, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, fmt.Errorf("invalid number in schedule")
		}

		unit := matches[2]
		switch unit {
		case "seconds":
			return time.Duration(amount) * time.Second, nil
		case "minutes":
			return time.Duration(amount) * time.Minute, nil
		case "hours":
			return time.Duration(amount) * time.Hour, nil
		default:
			return 0, fmt.Errorf("unsupported time unit")
		}
	}

	return 0, fmt.Errorf("invalid one-time schedule format: %s", schedule)
}
