package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type apiHandler struct {
	scheduler *Scheduler
}

func NewAPIHandler(scheduler *Scheduler) http.Handler {
	api := &apiHandler{scheduler: scheduler}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", api.handleHealth)
	mux.HandleFunc("/api/v1/notifications", api.handleNotifications)
	mux.HandleFunc("/api/v1/notifications/", api.handleNotificationByGroup)

	return withLogging(mux)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("API %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func (a *apiHandler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *apiHandler) handleNotifications(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v1/notifications" {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	req, err := decodeNotificationRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := a.scheduler.Schedule(req)
	if err != nil {
		if errors.Is(err, ErrNotificationExists) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to schedule notification")
		return
	}

	writeJSON(w, http.StatusCreated, response)
}

func (a *apiHandler) handleNotificationByGroup(w http.ResponseWriter, r *http.Request) {
	cropGroup, err := cropGroupFromPath(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	if err := cropGroup.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	switch r.Method {
	case http.MethodPut:
		a.handleReschedule(w, r, cropGroup)
	case http.MethodGet:
		a.handleGet(w, r, cropGroup)
	case http.MethodDelete:
		a.handleDelete(w, r, cropGroup)
	default:
		writeMethodNotAllowed(w, http.MethodPut, http.MethodGet, http.MethodDelete)
	}
}

func (a *apiHandler) handleReschedule(w http.ResponseWriter, r *http.Request, cropGroup CropGroup) {
	req, err := decodeNotificationRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.CropGroup = cropGroup

	response, err := a.scheduler.Reschedule(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to reschedule notification")
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (a *apiHandler) handleGet(w http.ResponseWriter, r *http.Request, cropGroup CropGroup) {
	userID := strings.TrimSpace(r.URL.Query().Get("userId"))
	if userID == "" {
		writeError(w, http.StatusBadRequest, "userId is required")
		return
	}

	response, err := a.scheduler.Get(userID, cropGroup)
	if err != nil {
		if errors.Is(err, ErrNotificationNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to fetch notification")
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (a *apiHandler) handleDelete(w http.ResponseWriter, r *http.Request, cropGroup CropGroup) {
	userID := strings.TrimSpace(r.URL.Query().Get("userId"))
	if userID == "" {
		writeError(w, http.StatusBadRequest, "userId is required")
		return
	}

	if err := a.scheduler.Cancel(userID, cropGroup); err != nil {
		if errors.Is(err, ErrNotificationNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to cancel notification")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func decodeNotificationRequest(r *http.Request) (NotificationRequest, error) {
	defer r.Body.Close()

	var req NotificationRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		return NotificationRequest{}, fmt.Errorf("invalid request body: %w", err)
	}

	if strings.TrimSpace(req.UserID) == "" {
		return NotificationRequest{}, errors.New("userId is required")
	}

	if req.NotifyInMinutes <= 0 {
		return NotificationRequest{}, errors.New("notifyInMinutes must be greater than 0")
	}

	if req.CropGroup != "" {
		if err := req.CropGroup.Validate(); err != nil {
			return NotificationRequest{}, err
		}
	}

	if req.GameMode != "" {
		if err := req.GameMode.Validate(); err != nil {
			return NotificationRequest{}, err
		}
	}

	if req.NotifyMode != "" {
		if err := req.NotifyMode.Validate(); err != nil {
			return NotificationRequest{}, err
		}
	}

	for _, patch := range req.Patches {
		if !patch.Location.Validate() {
			return NotificationRequest{}, fmt.Errorf("unknown patch location %q", patch.Location)
		}
	}

	return req, nil
}

func cropGroupFromPath(path string) (CropGroup, error) {
	const prefix = "/api/v1/notifications/"
	if !strings.HasPrefix(path, prefix) {
		return "", errors.New("not found")
	}

	group := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	if group == "" || strings.Contains(group, "/") {
		return "", errors.New("not found")
	}

	return CropGroup(group), nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"error":             message,
		"allowedCropGroups": allowedCropGroups(),
	})
}

func writeMethodNotAllowed(w http.ResponseWriter, methods ...string) {
	w.Header().Set("Allow", strings.Join(methods, ", "))
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}
