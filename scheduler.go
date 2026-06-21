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
	userID     string
	cropGroup  CropGroup
	cropName   string
	cropValue  string
	gameMode   GameMode
	notifyMode NotifyMode
	patches    []PatchInfo
	notifyAt   time.Time
	timer      *time.Timer
}

type Scheduler struct {
	mu            sync.RWMutex
	discord       *discordgo.Session
	thumbnails    *wikiThumbnailService
	notifications map[string]*scheduledNotification
	store         *store
}

func NewScheduler(discord *discordgo.Session, store *store) *Scheduler {
	s := &Scheduler{
		discord:       discord,
		thumbnails:    newWikiThumbnailService(),
		notifications: make(map[string]*scheduledNotification),
		store:         store,
	}

	stored, err := store.GetAll()
	if err != nil {
		log.Printf("failed to load stored notifications: %v", err)
		return s
	}

	now := time.Now().UTC()
	for _, sn := range stored {
		if sn.notifyAt.Before(now) {
			log.Printf("skipping past-due notification for %s:%s (was due at %s)",
				sn.userID, sn.cropGroup, sn.notifyAt.Format(time.RFC3339))
			if err := store.Delete(sn.userID, sn.cropGroup); err != nil {
				log.Printf("failed to delete past-due notification: %v", err)
			}
			continue
		}

		key := notificationKey(sn.userID, sn.cropGroup)
		notification := &scheduledNotification{
			userID:     sn.userID,
			cropGroup:  sn.cropGroup,
			cropName:   sn.cropName,
			cropValue:  sn.cropValue,
			gameMode:   sn.gameMode,
			notifyMode: sn.notifyMode,
			patches:    sn.patches,
			notifyAt:   sn.notifyAt,
		}

		notification.timer = time.AfterFunc(time.Until(sn.notifyAt), func() {
			s.fire(key, sn.userID, sn.cropGroup, sn.notifyAt)
		})

		s.notifications[key] = notification
		log.Printf("restored notification for %s:%s, fires at %s", sn.userID, sn.cropGroup, sn.notifyAt.Format(time.RFC3339))
	}

	return s
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

	if existing, exists := s.notifications[key]; exists && len(req.Patches) > 0 {
		return s.mergeLocked(req, existing), nil
	}

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
		GameMode:     notification.gameMode,
		NotifyMode:   notification.notifyMode,
		Patches:      notification.patches,
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

	if err := s.store.Delete(userID, cropGroup); err != nil {
		log.Printf("failed to delete notification from store: %v", err)
	}

	return nil
}

func (s *Scheduler) scheduleLocked(req NotificationRequest, status string) NotificationResponse {
	if req.CropValue == "" {
		if crop, ok := defaultCropForGroup(req.CropGroup); ok {
			req.CropValue = crop.Value
			req.CropName = crop.Name
		}
	}

	if req.GameMode == "" {
		req.GameMode = GameModeStandard
	}

	if req.NotifyMode == "" {
		req.NotifyMode = NotifyModeFirstReady
	}

	duration := gameModeDuration(time.Duration(req.NotifyInMinutes)*time.Minute, req.GameMode)
	notifyAt := time.Now().UTC().Add(duration)
	key := notificationKey(req.UserID, req.CropGroup)

	notification := &scheduledNotification{
		userID:     req.UserID,
		cropGroup:  req.CropGroup,
		cropName:   req.CropName,
		cropValue:  req.CropValue,
		gameMode:   req.GameMode,
		notifyMode: req.NotifyMode,
		patches:    req.Patches,
		notifyAt:   notifyAt,
	}

	if err := s.store.Insert(*notification); err != nil {
		log.Printf("failed to persist notification: %v", err)
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
		GameMode:     req.GameMode,
		NotifyMode:   req.NotifyMode,
		Patches:      req.Patches,
	}
}

func (s *Scheduler) mergeLocked(req NotificationRequest, existing *scheduledNotification) NotificationResponse {
	key := notificationKey(req.UserID, req.CropGroup)

	notifyMode := req.NotifyMode
	if notifyMode == "" {
		notifyMode = NotifyModeFirstReady
	}
	existing.notifyMode = notifyMode

	merged := mergePatches(existing.patches, req.Patches)
	existing.patches = merged

	if len(merged) > 0 {
		existing.cropValue = merged[0].Crop
	}

	if notifyMode == NotifyModeAllReady {
		duration := gameModeDuration(time.Duration(req.NotifyInMinutes)*time.Minute, req.GameMode)
		potentialNotifyAt := time.Now().UTC().Add(duration)
		if potentialNotifyAt.After(existing.notifyAt) {
			existing.timer.Stop()
			existing.notifyAt = potentialNotifyAt
			existing.timer = time.AfterFunc(time.Until(potentialNotifyAt), func() {
				s.fire(key, req.UserID, req.CropGroup, potentialNotifyAt)
			})
		}
	}

	if err := s.store.Insert(*existing); err != nil {
		log.Printf("failed to persist merged notification: %v", err)
	}

	s.notifications[key] = existing

	return NotificationResponse{
		UserID:       existing.userID,
		CropGroup:    existing.cropGroup,
		ScheduledFor: existing.notifyAt,
		Status:       "merged",
		GameMode:     existing.gameMode,
		NotifyMode:   existing.notifyMode,
		Patches:      existing.patches,
	}
}

func mergePatches(existing, incoming []PatchInfo) []PatchInfo {
	byLocation := make(map[PatchLocation]PatchInfo, len(existing))
	for _, p := range existing {
		byLocation[p.Location] = p
	}
	for _, p := range incoming {
		byLocation[p.Location] = p
	}
	result := make([]PatchInfo, 0, len(byLocation))
	for _, p := range byLocation {
		result = append(result, p)
	}
	return result
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

	if err := s.store.Delete(userID, cropGroup); err != nil {
		log.Printf("failed to delete fired notification from store: %v", err)
	}

	if err := s.sendHarvestReadyDM(notification); err != nil {
		log.Printf("failed sending notification to user %s for %s: %v", userID, cropGroup, err)
	}
}

func (s *Scheduler) buildHarvestEmbed(notification *scheduledNotification) *discordgo.MessageEmbed {
	title := notification.cropGroup.DisplayNamePluralTitle() + " ready"
	if notification.gameMode != "" && notification.gameMode != GameModeStandard {
		title = fmt.Sprintf("%s ready [%s]", notification.cropGroup.DisplayNamePluralTitle(), notification.gameMode)
	}

	embed := &discordgo.MessageEmbed{
		Title: title,
		Color: 0x4CAF50,
	}

	if notification.cropGroup == CropGroupFarmingContract {
		embed.Title = "Farming Contract"
		if notification.gameMode != "" && notification.gameMode != GameModeStandard {
			embed.Title = fmt.Sprintf("Farming Contract [%s]", notification.gameMode)
		}
		embed.Color = 0xFFA500

		cropName := notification.displayCropName()
		if len(notification.patches) > 0 {
			embed.Description = fmt.Sprintf("Your contracted %s at %s is ready to harvest!", cropName, notification.patches[0].Location)
		} else {
			embed.Description = fmt.Sprintf("Your contracted %s is ready to harvest!", cropName)
		}

		if crop, ok := cropForGroup(notification.cropGroup, notification.cropValue); ok {
			thumbnailURL, err := s.thumbnails.ThumbnailURL(crop)
			if err != nil {
				log.Printf("failed to fetch wiki thumbnail for %s: %v", crop.Name, err)
			} else if thumbnailURL != "" {
				embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: thumbnailURL}
			}
		}
	} else if len(notification.patches) > 0 {
		desc := fmt.Sprintf("Your %s are ready:", notification.cropGroup.DisplayNamePlural())
		for _, p := range notification.patches {
			cropDisplay := p.Crop
			if crop, ok := cropForGroup(notification.cropGroup, p.Crop); ok {
				cropDisplay = crop.Name
			}
			desc += fmt.Sprintf("\n- %s at %s", cropDisplay, p.Location)
		}
		embed.Description = desc

		if crop, ok := cropForGroup(notification.cropGroup, notification.patches[0].Crop); ok {
			thumbnailURL, err := s.thumbnails.ThumbnailURL(crop)
			if err != nil {
				log.Printf("failed to fetch wiki thumbnail for %s: %v", crop.Name, err)
			} else if thumbnailURL != "" {
				embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: thumbnailURL}
			}
		}
	} else {
		embed.Description = fmt.Sprintf("Your %s is ready to harvest", notification.displayCropName())

		if crop, ok := cropForGroup(notification.cropGroup, notification.cropValue); ok {
			thumbnailURL, err := s.thumbnails.ThumbnailURL(crop)
			if err != nil {
				log.Printf("failed to fetch wiki thumbnail for %s: %v", crop.Name, err)
			} else if thumbnailURL != "" {
				embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: thumbnailURL}
			}
		}
	}

	return embed
}

func (s *Scheduler) sendHarvestReadyDM(notification *scheduledNotification) error {
	channel, err := s.discord.UserChannelCreate(notification.userID)
	if err != nil {
		return fmt.Errorf("create DM channel: %w", err)
	}

	cropValue := notification.cropValue
	if cropValue == "" && len(notification.patches) > 0 {
		cropValue = notification.patches[0].Crop
	}

	customID := fmt.Sprintf("reschedule:%s:%s", notification.cropGroup, cropValue)
	_, err = s.discord.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{s.buildHarvestEmbed(notification)},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "I replanted",
						Style:    discordgo.SuccessButton,
						CustomID: customID,
					},
				},
			},
		},
	})
	if err != nil {
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
