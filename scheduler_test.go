package main

import (
	"testing"
	"time"
)

func notificationRequest(overrides ...func(*NotificationRequest)) NotificationRequest {
	req := NotificationRequest{
		UserID:          "testuser",
		CropGroup:       CropGroupHerb,
		NotifyInMinutes: 100000,
		CropName:        "Ranarr",
		CropValue:       "ranarr",
	}
	for _, fn := range overrides {
		fn(&req)
	}
	return req
}

func TestNotificationKey(t *testing.T) {
	tests := []struct {
		userID    string
		cropGroup CropGroup
		want      string
	}{
		{"user123", CropGroupHerb, "user123:herb"},
		{"alice", CropGroupFruitTree, "alice:fruit_tree"},
		{"", CropGroupHerb, ":herb"},
	}

	for _, tt := range tests {
		got := notificationKey(tt.userID, tt.cropGroup)
		if got != tt.want {
			t.Errorf("notificationKey(%q, %q) = %q, want %q", tt.userID, tt.cropGroup, got, tt.want)
		}
	}
}

func TestSchedulerSchedule_CreatesNotification(t *testing.T) {
	s := NewScheduler(nil)
	req := notificationRequest()

	resp, err := s.Schedule(req)
	if err != nil {
		t.Fatalf("Schedule() error = %v", err)
	}

	if resp.Status != "scheduled" {
		t.Errorf("Status = %q, want %q", resp.Status, "scheduled")
	}
	if resp.UserID != "testuser" {
		t.Errorf("UserID = %q, want %q", resp.UserID, "testuser")
	}
	if resp.CropGroup != CropGroupHerb {
		t.Errorf("CropGroup = %q, want %q", resp.CropGroup, CropGroupHerb)
	}
	if resp.ScheduledFor.IsZero() {
		t.Error("ScheduledFor is zero")
	}

	expectedTime := time.Now().UTC().Add(100000 * time.Minute)
	diff := resp.ScheduledFor.Sub(expectedTime)
	if diff < -time.Second || diff > time.Second {
		t.Errorf("ScheduledFor = %v, want ~%v (diff %v)", resp.ScheduledFor, expectedTime, diff)
	}
}

func TestSchedulerSchedule_DuplicateReturnsError(t *testing.T) {
	s := NewScheduler(nil)
	req := notificationRequest()

	if _, err := s.Schedule(req); err != nil {
		t.Fatalf("first Schedule() error = %v", err)
	}

	_, err := s.Schedule(req)
	if err != ErrNotificationExists {
		t.Errorf("second Schedule() error = %v, want %v", err, ErrNotificationExists)
	}
}

func TestSchedulerSchedule_DifferentUsersNoConflict(t *testing.T) {
	s := NewScheduler(nil)
	req1 := notificationRequest(func(r *NotificationRequest) { r.UserID = "user1" })
	req2 := notificationRequest(func(r *NotificationRequest) { r.UserID = "user2" })

	if _, err := s.Schedule(req1); err != nil {
		t.Fatalf("schedule user1 error = %v", err)
	}
	if _, err := s.Schedule(req2); err != nil {
		t.Fatalf("schedule user2 error = %v", err)
	}
}

func TestSchedulerSchedule_DifferentGroupsNoConflict(t *testing.T) {
	s := NewScheduler(nil)
	req1 := notificationRequest(func(r *NotificationRequest) { r.CropGroup = CropGroupHerb })
	req2 := notificationRequest(func(r *NotificationRequest) { r.CropGroup = CropGroupTree })

	if _, err := s.Schedule(req1); err != nil {
		t.Fatalf("schedule herb error = %v", err)
	}
	if _, err := s.Schedule(req2); err != nil {
		t.Fatalf("schedule tree error = %v", err)
	}
}

func TestSchedulerReschedule_CreatesWhenNotExists(t *testing.T) {
	s := NewScheduler(nil)
	req := notificationRequest()

	resp, err := s.Reschedule(req)
	if err != nil {
		t.Fatalf("Reschedule() error = %v", err)
	}
	if resp.Status != "scheduled" {
		t.Errorf("Status = %q, want %q", resp.Status, "scheduled")
	}
}

func TestSchedulerReschedule_ReschedulesWhenExists(t *testing.T) {
	s := NewScheduler(nil)
	req1 := notificationRequest(func(r *NotificationRequest) { r.NotifyInMinutes = 100000 })
	req2 := notificationRequest(func(r *NotificationRequest) { r.NotifyInMinutes = 200000 })

	resp1, _ := s.Reschedule(req1)
	resp2, err := s.Reschedule(req2)
	if err != nil {
		t.Fatalf("second Reschedule() error = %v", err)
	}

	if resp2.Status != "rescheduled" {
		t.Errorf("Status = %q, want %q", resp2.Status, "rescheduled")
	}
	if resp2.ScheduledFor.Equal(resp1.ScheduledFor) {
		t.Error("ScheduledFor should have changed after reschedule")
	}
}

func TestSchedulerReschedule_DifferentUsersNotConfused(t *testing.T) {
	s := NewScheduler(nil)
	req1 := notificationRequest(func(r *NotificationRequest) { r.UserID = "user1"; r.NotifyInMinutes = 100 })
	req2 := notificationRequest(func(r *NotificationRequest) { r.UserID = "user2"; r.NotifyInMinutes = 200 })

	s.Reschedule(req1)
	resp, _ := s.Reschedule(req2)
	if resp.Status != "scheduled" {
		t.Errorf("second user should get 'scheduled', got %q", resp.Status)
	}
}

func TestSchedulerGet_Found(t *testing.T) {
	s := NewScheduler(nil)
	req := notificationRequest()

	s.Schedule(req)
	resp, err := s.Get("testuser", CropGroupHerb)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if resp.UserID != "testuser" {
		t.Errorf("UserID = %q, want %q", resp.UserID, "testuser")
	}
}

func TestSchedulerGet_NotFound(t *testing.T) {
	s := NewScheduler(nil)

	_, err := s.Get("testuser", CropGroupHerb)
	if err != ErrNotificationNotFound {
		t.Errorf("Get() error = %v, want %v", err, ErrNotificationNotFound)
	}
}

func TestSchedulerGet_AfterCancelReturnsNotFound(t *testing.T) {
	s := NewScheduler(nil)
	req := notificationRequest()

	s.Schedule(req)
	s.Cancel("testuser", CropGroupHerb)

	_, err := s.Get("testuser", CropGroupHerb)
	if err != ErrNotificationNotFound {
		t.Errorf("Get() after cancel error = %v, want %v", err, ErrNotificationNotFound)
	}
}

func TestSchedulerCancel_Found(t *testing.T) {
	s := NewScheduler(nil)
	req := notificationRequest()

	s.Schedule(req)
	err := s.Cancel("testuser", CropGroupHerb)
	if err != nil {
		t.Fatalf("Cancel() error = %v", err)
	}
}

func TestSchedulerCancel_NotFound(t *testing.T) {
	s := NewScheduler(nil)

	err := s.Cancel("testuser", CropGroupHerb)
	if err != ErrNotificationNotFound {
		t.Errorf("Cancel() error = %v, want %v", err, ErrNotificationNotFound)
	}
}

func TestSchedulerSchedule_EmptyCropValueFillsDefault(t *testing.T) {
	s := NewScheduler(nil)
	req := NotificationRequest{
		UserID:          "testuser",
		CropGroup:       CropGroupHerb,
		NotifyInMinutes: 100000,
	}

	if _, err := s.Schedule(req); err != nil {
		t.Fatalf("Schedule() error = %v", err)
	}

	n := s.notifications["testuser:herb"]
	if n == nil {
		t.Fatal("notification not stored")
	}
	if n.cropValue != "guam" {
		t.Errorf("cropValue = %q, want %q", n.cropValue, "guam")
	}
	if n.cropName != "Guam" {
		t.Errorf("cropName = %q, want %q", n.cropName, "Guam")
	}
}

func TestSchedulerReschedule_EmptyCropValueFillsDefault(t *testing.T) {
	s := NewScheduler(nil)
	req := NotificationRequest{
		UserID:          "testuser",
		CropGroup:       CropGroupTree,
		NotifyInMinutes: 100000,
	}

	if _, err := s.Reschedule(req); err != nil {
		t.Fatalf("Reschedule() error = %v", err)
	}

	n := s.notifications["testuser:tree"]
	if n == nil {
		t.Fatal("notification not stored")
	}
	if n.cropValue != "oak" {
		t.Errorf("cropValue = %q, want %q", n.cropValue, "oak")
	}
	if n.cropName != "Oak" {
		t.Errorf("cropName = %q, want %q", n.cropName, "Oak")
	}
}

func TestSchedulerSchedule_StopsPreviousTimerOnReschedule(t *testing.T) {
	s := NewScheduler(nil)
	req1 := notificationRequest(func(r *NotificationRequest) { r.NotifyInMinutes = 100000 })

	s.Schedule(req1)

	n := s.notifications["testuser:herb"]
	if n == nil || n.timer == nil {
		t.Fatal("expected timer to be set")
	}
	initialTimer := n.timer

	req2 := notificationRequest(func(r *NotificationRequest) { r.NotifyInMinutes = 200000 })
	s.Reschedule(req2)

	if initialTimer.Stop() {
		t.Error("original timer should have been stopped and not return true from Stop()")
	}
}

func TestSchedulerDisplayCropName(t *testing.T) {
	tests := []struct {
		name  string
		notif *scheduledNotification
		want  string
	}{
		{"uses crop name when set", &scheduledNotification{cropName: "Ranarr", cropGroup: CropGroupHerb}, "Ranarr"},
		{"falls back to group display name", &scheduledNotification{cropGroup: CropGroupHerb}, "Herb"},
		{"fruit tree fallback", &scheduledNotification{cropGroup: CropGroupFruitTree}, "Fruit Tree"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.notif.displayCropName()
			if got != tt.want {
				t.Errorf("displayCropName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSchedulerBuildHarvestEmbed(t *testing.T) {
	s := NewScheduler(nil)

	notif := &scheduledNotification{
		userID:    "testuser",
		cropGroup: CropGroupHerb,
		cropName:  "Ranarr",
		cropValue: "ranarr",
	}

	embed := s.buildHarvestEmbed(notif)
	if embed == nil {
		t.Fatal("buildHarvestEmbed() returned nil")
	}
	if embed.Title != "Harvest Ready" {
		t.Errorf("Title = %q, want %q", embed.Title, "Harvest Ready")
	}
	if embed.Description != "Your Ranarr is ready to harvest" {
		t.Errorf("Description = %q, want %q", embed.Description, "Your Ranarr is ready to harvest")
	}
	if embed.Color != 0x4CAF50 {
		t.Errorf("Color = %#x, want %#x", embed.Color, 0x4CAF50)
	}
}

func TestSchedulerBuildHarvestEmbed_WithPatches(t *testing.T) {
	s := NewScheduler(nil)

	notif := &scheduledNotification{
		userID:    "testuser",
		cropGroup: CropGroupHerb,
		patches: []PatchInfo{
			{Crop: "ranarr", Location: "Farming Guild"},
			{Crop: "irit", Location: "Falador"},
		},
	}

	embed := s.buildHarvestEmbed(notif)
	if embed == nil {
		t.Fatal("buildHarvestEmbed() returned nil")
	}
	if embed.Title != "Harvest Ready" {
		t.Errorf("Title = %q, want %q", embed.Title, "Harvest Ready")
	}
	want := "Your herbs are ready:\n- Ranarr at Farming Guild\n- Irit at Falador"
	if embed.Description != want {
		t.Errorf("Description = %q, want %q", embed.Description, want)
	}
	if embed.Color != 0x4CAF50 {
		t.Errorf("Color = %#x, want %#x", embed.Color, 0x4CAF50)
	}
}

func TestSchedulerBuildHarvestEmbed_EmptyCropValueNotCached(t *testing.T) {
	s := NewScheduler(nil)
	s.thumbnails.cache["Guam"] = "https://example.com/guam.png"

	notif := &scheduledNotification{
		userID:    "testuser",
		cropGroup: CropGroupHerb,
		cropName:  "Guam",
		cropValue: "guam",
	}

	embed := s.buildHarvestEmbed(notif)
	if embed == nil {
		t.Fatal("buildHarvestEmbed() returned nil")
	}
	if embed.Thumbnail == nil {
		t.Fatal("expected thumbnail to be set from cache")
	}
	if embed.Thumbnail.URL != "https://example.com/guam.png" {
		t.Errorf("Thumbnail.URL = %q, want %q", embed.Thumbnail.URL, "https://example.com/guam.png")
	}
}
