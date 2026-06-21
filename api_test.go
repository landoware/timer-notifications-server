package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupAPI(t *testing.T) (http.Handler, *Scheduler) {
	t.Helper()
	s := NewScheduler(nil)
	return NewAPIHandler(s), s
}

func apiRequest(method, path string, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestHealthEndpoint(t *testing.T) {
	handler, _ := setupAPI(t)
	req := apiRequest("GET", "/api/v1/health", "")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	if body["status"] != "ok" {
		t.Errorf(`body["status"] = %q, want "ok"`, body["status"])
	}
}

func TestCreateNotification_Success(t *testing.T) {
	handler, _ := setupAPI(t)
	body := `{"userId":"user1","cropGroup":"herb","notifyInMinutes":80,"crop":"ranarr"}`
	req := apiRequest("POST", "/api/v1/notifications", body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}

	var resp NotificationResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Status != "scheduled" {
		t.Errorf("Status = %q, want %q", resp.Status, "scheduled")
	}
	if resp.UserID != "user1" {
		t.Errorf("UserID = %q, want %q", resp.UserID, "user1")
	}
	if resp.CropGroup != CropGroupHerb {
		t.Errorf("CropGroup = %q, want %q", resp.CropGroup, CropGroupHerb)
	}
}

func TestCreateNotification_DuplicateReturnsConflict(t *testing.T) {
	handler, _ := setupAPI(t)
	body := `{"userId":"user1","cropGroup":"herb","notifyInMinutes":80}`
	req1 := apiRequest("POST", "/api/v1/notifications", body)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)
	if w1.Code != http.StatusCreated {
		t.Fatalf("first request status = %d", w1.Code)
	}

	req2 := apiRequest("POST", "/api/v1/notifications", body)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)
	if w2.Code != http.StatusConflict {
		t.Errorf("duplicate status = %d, want %d", w2.Code, http.StatusConflict)
	}

	var errResp map[string]any
	json.NewDecoder(w2.Body).Decode(&errResp)
	if errResp["error"] != ErrNotificationExists.Error() {
		t.Errorf(`error = %q, want %q`, errResp["error"], ErrNotificationExists.Error())
	}
}

func TestCreateNotification_MissingUserId(t *testing.T) {
	handler, _ := setupAPI(t)
	body := `{"cropGroup":"herb","notifyInMinutes":80}`
	req := apiRequest("POST", "/api/v1/notifications", body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateNotification_InvalidNotifyInMinutes(t *testing.T) {
	handler, _ := setupAPI(t)
	tests := []string{
		`{"userId":"user1","cropGroup":"herb","notifyInMinutes":0}`,
		`{"userId":"user1","cropGroup":"herb","notifyInMinutes":-1}`,
	}
	for _, body := range tests {
		req := apiRequest("POST", "/api/v1/notifications", body)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("for body %q: status = %d, want %d", body, w.Code, http.StatusBadRequest)
		}
	}
}

func TestCreateNotification_UnknownPatchLocation(t *testing.T) {
	handler, _ := setupAPI(t)
	body := `{"userId":"user1","cropGroup":"herb","notifyInMinutes":80,"patches":[{"crop":"ranarr","location":"Nowhere"}]}`
	req := apiRequest("POST", "/api/v1/notifications", body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateNotification_UnknownCropGroup(t *testing.T) {
	handler, _ := setupAPI(t)
	body := `{"userId":"user1","cropGroup":"invalid","notifyInMinutes":80}`
	req := apiRequest("POST", "/api/v1/notifications", body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateNotification_UnknownFieldsRejected(t *testing.T) {
	handler, _ := setupAPI(t)
	body := `{"userId":"user1","cropGroup":"herb","notifyInMinutes":80,"extraField":"value"}`
	req := apiRequest("POST", "/api/v1/notifications", body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateNotification_InvalidJSON(t *testing.T) {
	handler, _ := setupAPI(t)
	body := `not json`
	req := apiRequest("POST", "/api/v1/notifications", body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateNotification_MethodNotAllowed(t *testing.T) {
	handler, _ := setupAPI(t)
	methods := []string{"GET", "PUT", "DELETE", "PATCH"}
	for _, method := range methods {
		req := apiRequest(method, "/api/v1/notifications", "")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s /api/v1/notifications: status = %d, want %d", method, w.Code, http.StatusMethodNotAllowed)
		}
	}
}

func TestCreateNotification_ErrorResponseIncludesAllowedGroups(t *testing.T) {
	handler, _ := setupAPI(t)
	body := `{"userId":"user1","cropGroup":"invalid","notifyInMinutes":80}`
	req := apiRequest("POST", "/api/v1/notifications", body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	var errResp map[string]any
	json.NewDecoder(w.Body).Decode(&errResp)
	groups, ok := errResp["allowedCropGroups"]
	if !ok {
		t.Fatal("response missing allowedCropGroups")
	}
	groupList, ok := groups.([]any)
	if !ok {
		t.Fatal("allowedCropGroups is not a list")
	}
	if len(groupList) != len(validCropGroups) {
		t.Errorf("allowedCropGroups has %d entries, want %d", len(groupList), len(validCropGroups))
	}
}

func TestCreateNotification_WithPatches(t *testing.T) {
	handler, _ := setupAPI(t)
	body := `{"userId":"user1","cropGroup":"herb","notifyInMinutes":80,"patches":[{"crop":"ranarr","location":"Farming Guild"},{"crop":"irit","location":"Falador"}]}`
	req := apiRequest("POST", "/api/v1/notifications", body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}

	var resp NotificationResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(resp.Patches) != 2 {
		t.Fatalf("got %d patches, want 2", len(resp.Patches))
	}
	if resp.Patches[0].Crop != "ranarr" {
		t.Errorf("Patches[0].Crop = %q, want %q", resp.Patches[0].Crop, "ranarr")
	}
	if resp.Patches[0].Location != PatchLocation("Farming Guild") {
		t.Errorf("Patches[0].Location = %q, want %q", resp.Patches[0].Location, "Farming Guild")
	}
	if resp.Patches[1].Crop != "irit" {
		t.Errorf("Patches[1].Crop = %q, want %q", resp.Patches[1].Crop, "irit")
	}
	if resp.Patches[1].Location != PatchLocation("Falador") {
		t.Errorf("Patches[1].Location = %q, want %q", resp.Patches[1].Location, "Falador")
	}
}

func TestRescheduleNotification_Success(t *testing.T) {
	handler, _ := setupAPI(t)
	body := `{"userId":"user1","notifyInMinutes":80,"crop":"ranarr"}`
	req := apiRequest("PUT", "/api/v1/notifications/herb", body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp NotificationResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Status != "scheduled" {
		t.Errorf("Status = %q, want %q", resp.Status, "scheduled")
	}
	if resp.CropGroup != CropGroupHerb {
		t.Errorf("CropGroup = %q, want %q", resp.CropGroup, CropGroupHerb)
	}
}

func TestRescheduleNotification_ReschedulesExisting(t *testing.T) {
	handler, s := setupAPI(t)
	body := `{"userId":"user1","notifyInMinutes":80,"crop":"ranarr"}`

	req1 := apiRequest("PUT", "/api/v1/notifications/herb", body)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)

	req2 := apiRequest("PUT", "/api/v1/notifications/herb", body)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	var resp NotificationResponse
	json.NewDecoder(w2.Body).Decode(&resp)
	if resp.Status != "rescheduled" {
		t.Errorf("Status = %q, want %q", resp.Status, "rescheduled")
	}

	_ = s
}

func TestCreateAndGetNotification_WithPatches(t *testing.T) {
	handler, _ := setupAPI(t)

	body := `{"userId":"user1","cropGroup":"herb","notifyInMinutes":80,"patches":[{"crop":"ranarr","location":"Farming Guild"},{"crop":"irit","location":"Falador"}]}`
	req := apiRequest("POST", "/api/v1/notifications", body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("create status = %d", w.Code)
	}

	getReq := apiRequest("GET", "/api/v1/notifications/herb?userId=user1", "")
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, getReq)

	if w2.Code != http.StatusOK {
		t.Fatalf("get status = %d", w2.Code)
	}

	var resp NotificationResponse
	json.NewDecoder(w2.Body).Decode(&resp)
	if len(resp.Patches) != 2 {
		t.Fatalf("got %d patches on get, want 2", len(resp.Patches))
	}
	if resp.Patches[0].Crop != "ranarr" || resp.Patches[0].Location != PatchLocation("Farming Guild") {
		t.Errorf("patch 0 = %+v", resp.Patches[0])
	}
	if resp.Patches[1].Crop != "irit" || resp.Patches[1].Location != PatchLocation("Falador") {
		t.Errorf("patch 1 = %+v", resp.Patches[1])
	}
}

func TestGetNotification_Success(t *testing.T) {
	handler, s := setupAPI(t)

	s.Schedule(NotificationRequest{
		UserID:          "user1",
		CropGroup:       CropGroupHerb,
		NotifyInMinutes: 100000,
		CropValue:       "ranarr",
		CropName:        "Ranarr",
	})

	req := apiRequest("GET", "/api/v1/notifications/herb?userId=user1", "")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp NotificationResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.UserID != "user1" {
		t.Errorf("UserID = %q, want %q", resp.UserID, "user1")
	}
	if resp.CropGroup != CropGroupHerb {
		t.Errorf("CropGroup = %q, want %q", resp.CropGroup, CropGroupHerb)
	}
}

func TestGetNotification_NotFound(t *testing.T) {
	handler, _ := setupAPI(t)
	req := apiRequest("GET", "/api/v1/notifications/herb?userId=nonexistent", "")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestGetNotification_MissingUserId(t *testing.T) {
	handler, _ := setupAPI(t)
	req := apiRequest("GET", "/api/v1/notifications/herb", "")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetNotification_InvalidCropGroup(t *testing.T) {
	handler, _ := setupAPI(t)
	req := apiRequest("GET", "/api/v1/notifications/invalid_group?userId=user1", "")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestDeleteNotification_Success(t *testing.T) {
	handler, s := setupAPI(t)

	s.Schedule(NotificationRequest{
		UserID:          "user1",
		CropGroup:       CropGroupHerb,
		NotifyInMinutes: 100000,
		CropValue:       "ranarr",
		CropName:        "Ranarr",
	})

	req := apiRequest("DELETE", "/api/v1/notifications/herb?userId=user1", "")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestDeleteNotification_NotFound(t *testing.T) {
	handler, _ := setupAPI(t)
	req := apiRequest("DELETE", "/api/v1/notifications/herb?userId=user1", "")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestDeleteNotification_MissingUserId(t *testing.T) {
	handler, _ := setupAPI(t)
	req := apiRequest("DELETE", "/api/v1/notifications/herb", "")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestNotificationByGroup_MethodNotAllowed(t *testing.T) {
	handler, _ := setupAPI(t)
	methods := []string{"POST", "PATCH"}
	for _, method := range methods {
		req := apiRequest(method, "/api/v1/notifications/herb?userId=user1", "")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s /api/v1/notifications/herb: status = %d, want %d", method, w.Code, http.StatusMethodNotAllowed)
		}
	}
}

func TestCropGroupFromPath(t *testing.T) {
	tests := []struct {
		path      string
		want      CropGroup
		wantError bool
	}{
		{"/api/v1/notifications/herb", CropGroupHerb, false},
		{"/api/v1/notifications/fruit_tree", CropGroupFruitTree, false},
		{"/api/v1/notifications/", "", true},
		{"/api/v1/notifications", "", true},
		{"/api/v1/notifications/herb/extra", "", true},
		{"/api/v1/other", "", true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("path=%q", tt.path), func(t *testing.T) {
			got, err := cropGroupFromPath(tt.path)
			if (err != nil) != tt.wantError {
				t.Fatalf("cropGroupFromPath(%q) error = %v, wantError = %v", tt.path, err, tt.wantError)
			}
			if got != tt.want {
				t.Errorf("cropGroupFromPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
