package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var ErrNotificationExists = errors.New("notification already exists")
var ErrNotificationNotFound = errors.New("notification not found")

type scheduledNotification struct {
	userID    string
	cropGroup CropGroup
	cropName  string
	cropValue string
	notifyAt  time.Time
	timer     *time.Timer
}

type Scheduler struct {
	mu            sync.RWMutex
	discord       *discordgo.Session
	thumbnails    *wikiThumbnailService
	notifications map[string]*scheduledNotification
}

func NewScheduler(discord *discordgo.Session) *Scheduler {
	return &Scheduler{
		discord:       discord,
		thumbnails:    newWikiThumbnailService(),
		notifications: make(map[string]*scheduledNotification),
	}
}

func (s *Scheduler) Schedule(req NotificationRequest) (NotificationResponse, error) {
	key := notificationKey(req.UserID, req.CropGroup)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.notifications[key]; exists {
		return NotificationResponse{}, ErrNotificationExists
	}

	return s.scheduleLocked(req, "scheduled"), nil
}

func (s *Scheduler) Reschedule(req NotificationRequest) (NotificationResponse, error) {
	key := notificationKey(req.UserID, req.CropGroup)

	s.mu.Lock()
	defer s.mu.Unlock()

	status := "scheduled"
	if existing, exists := s.notifications[key]; exists {
		existing.timer.Stop()
		delete(s.notifications, key)
		status = "rescheduled"
	}

	return s.scheduleLocked(req, status), nil
}

func (s *Scheduler) Get(userID string, cropGroup CropGroup) (NotificationResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	notification, exists := s.notifications[notificationKey(userID, cropGroup)]
	if !exists {
		return NotificationResponse{}, ErrNotificationNotFound
	}

	return NotificationResponse{
		UserID:       notification.userID,
		CropGroup:    notification.cropGroup,
		ScheduledFor: notification.notifyAt,
		Status:       "scheduled",
	}, nil
}

func (s *Scheduler) Cancel(userID string, cropGroup CropGroup) error {
	key := notificationKey(userID, cropGroup)

	s.mu.Lock()
	defer s.mu.Unlock()

	notification, exists := s.notifications[key]
	if !exists {
		return ErrNotificationNotFound
	}

	notification.timer.Stop()
	delete(s.notifications, key)

	return nil
}

func (s *Scheduler) scheduleLocked(req NotificationRequest, status string) NotificationResponse {
	notifyAt := time.Now().UTC().Add(time.Duration(req.NotifyInMinutes) * time.Minute)
	key := notificationKey(req.UserID, req.CropGroup)

	notification := &scheduledNotification{
		userID:    req.UserID,
		cropGroup: req.CropGroup,
		cropName:  req.CropName,
		cropValue: req.CropValue,
		notifyAt:  notifyAt,
	}

	notification.timer = time.AfterFunc(time.Until(notifyAt), func() {
		s.fire(key, req.UserID, req.CropGroup, notifyAt)
	})

	s.notifications[key] = notification

	return NotificationResponse{
		UserID:       req.UserID,
		CropGroup:    req.CropGroup,
		ScheduledFor: notifyAt,
		Status:       status,
	}
}

func (s *Scheduler) fire(key string, userID string, cropGroup CropGroup, notifyAt time.Time) {
	s.mu.Lock()
	notification, exists := s.notifications[key]
	if !exists || !notification.notifyAt.Equal(notifyAt) {
		s.mu.Unlock()
		return
	}
	delete(s.notifications, key)
	s.mu.Unlock()

	if err := s.sendHarvestReadyDM(notification); err != nil {
		log.Printf("failed sending notification to user %s for %s: %v", userID, cropGroup, err)
	}
}

func (s *Scheduler) sendHarvestReadyDM(notification *scheduledNotification) error {
	channel, err := s.discord.UserChannelCreate(notification.userID)
	if err != nil {
		return fmt.Errorf("create DM channel: %w", err)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Harvest Ready",
		Description: fmt.Sprintf("Your %s is ready to harvest", notification.displayCropName()),
		Color:       0x4CAF50,
	}

	if crop, ok := cropForGroup(notification.cropGroup, notification.cropValue); ok {
		thumbnailURL, err := s.thumbnails.ThumbnailURL(crop)
		if err != nil {
			log.Printf("failed to fetch wiki thumbnail for %s: %v", crop.Name, err)
		} else if thumbnailURL != "" {
			embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: thumbnailURL}
		}
	}

	if _, err := s.discord.ChannelMessageSendEmbed(channel.ID, embed); err != nil {
		return fmt.Errorf("send embed DM: %w", err)
	}

	return nil
}

func (n *scheduledNotification) displayCropName() string {
	if n.cropName != "" {
		return n.cropName
	}

	return n.cropGroup.DisplayName()
}

func notificationKey(userID string, cropGroup CropGroup) string {
	return userID + ":" + string(cropGroup)
}
