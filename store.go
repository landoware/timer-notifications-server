package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type store struct {
	db *sql.DB
}

func newStore(dbPath string) (*store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		return nil, fmt.Errorf("set WAL mode: %w", err)
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS notifications (
			user_id    TEXT NOT NULL,
			crop_group TEXT NOT NULL,
			crop_name  TEXT NOT NULL DEFAULT '',
			crop_value TEXT NOT NULL DEFAULT '',
			game_mode  TEXT NOT NULL DEFAULT 'standard',
			notify_mode TEXT NOT NULL DEFAULT 'first_ready',
			patches    TEXT NOT NULL DEFAULT '[]',
			notify_at  INTEGER NOT NULL,
			PRIMARY KEY (user_id, crop_group)
		)
	`); err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	// Migration: add notify_mode column if missing (existing databases)
	if _, err := db.Exec(`ALTER TABLE notifications ADD COLUMN notify_mode TEXT NOT NULL DEFAULT 'first_ready'`); err != nil {
		// column already exists — ignore
	}

	return &store{db: db}, nil
}

func (s *store) Close() error {
	return s.db.Close()
}

func (s *store) Insert(n scheduledNotification) error {
	patchesJSON, err := json.Marshal(n.patches)
	if err != nil {
		return fmt.Errorf("marshal patches: %w", err)
	}

	_, err = s.db.Exec(
		`INSERT OR REPLACE INTO notifications (user_id, crop_group, crop_name, crop_value, game_mode, notify_mode, patches, notify_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		n.userID, string(n.cropGroup), n.cropName, n.cropValue, string(n.gameMode),
		string(n.notifyMode), string(patchesJSON), n.notifyAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("insert notification: %w", err)
	}

	return nil
}

func (s *store) Delete(userID string, cropGroup CropGroup) error {
	_, err := s.db.Exec(
		`DELETE FROM notifications WHERE user_id = ? AND crop_group = ?`,
		userID, string(cropGroup),
	)
	if err != nil {
		return fmt.Errorf("delete notification: %w", err)
	}

	return nil
}

type storedNotification struct {
	userID     string
	cropGroup  CropGroup
	cropName   string
	cropValue  string
	gameMode   GameMode
	notifyMode NotifyMode
	patches    []PatchInfo
	notifyAt   time.Time
}

func (s *store) GetAll() ([]storedNotification, error) {
	rows, err := s.db.Query(
		`SELECT user_id, crop_group, crop_name, crop_value, game_mode, notify_mode, patches, notify_at FROM notifications`,
	)
	if err != nil {
		return nil, fmt.Errorf("query notifications: %w", err)
	}
	defer rows.Close()

	var result []storedNotification
	for rows.Next() {
		var (
			userID, cropGroup, cropName, cropValue, gameMode, notifyMode, patchesStr string
			notifyAtUnix                                                             int64
		)
		if err := rows.Scan(&userID, &cropGroup, &cropName, &cropValue, &gameMode, &notifyMode, &patchesStr, &notifyAtUnix); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		var patches []PatchInfo
		if patchesStr != "" && patchesStr != "[]" {
			if err := json.Unmarshal([]byte(patchesStr), &patches); err != nil {
				return nil, fmt.Errorf("unmarshal patches: %w", err)
			}
		}

		if notifyMode == "" {
			notifyMode = "first_ready"
		}

		result = append(result, storedNotification{
			userID:     userID,
			cropGroup:  CropGroup(cropGroup),
			cropName:   cropName,
			cropValue:  cropValue,
			gameMode:   GameMode(gameMode),
			notifyMode: NotifyMode(notifyMode),
			patches:    patches,
			notifyAt:   time.Unix(notifyAtUnix, 0).UTC(),
		})
	}

	return result, rows.Err()
}

func (s *store) DeleteAll() error {
	_, err := s.db.Exec(`DELETE FROM notifications`)
	if err != nil {
		return fmt.Errorf("delete all notifications: %w", err)
	}

	return nil
}
