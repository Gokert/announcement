package session

import (
	"context"
	"filmoteka/pkg/models"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

func GetTestClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
}

func TestCheckActiveSession(t *testing.T) {
	ctx := context.Background()
	activeSession := models.Session{
		SID:   "session_id",
		Login: "user_login",
	}

	client := GetTestClient()
	defer client.Close()

	repo := &SessionRepo{
		DB: client,
	}

	err := repo.DB.Set(ctx, activeSession.SID, activeSession.Login, 24*time.Hour).Err()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	added, err := repo.CheckActiveSession(ctx, activeSession.SID)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if !added {
		t.Errorf("session not added")
	}

	err = repo.DB.Set(ctx, activeSession.SID, activeSession.Login, 24*time.Hour).Err()
	if err != nil {
		t.Errorf("expected error, got nil")
	}

	added, err = repo.CheckActiveSession(ctx, "new_sid")
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if added {
		t.Errorf("session not added")
	}
}

func TestAddSession(t *testing.T) {
	ctx := context.Background()
	activeSession := models.Session{
		SID:   "session_id",
		Login: "user_login",
	}

	client := GetTestClient()
	repo := &SessionRepo{
		DB: client,
	}

	added, err := repo.AddSession(ctx, activeSession)
	if err != nil {
		t.Errorf("expected error, got nil")
	}

	if !added {
		t.Errorf("session not added")
	}

	client.Close()
	_, err = repo.AddSession(ctx, activeSession)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

}

func TestGetUserLogin(t *testing.T) {
	ctx := context.Background()
	activeSession := models.Session{
		SID:   "session_id",
		Login: "user_login",
	}

	client := GetTestClient()

	repo := &SessionRepo{
		DB: client,
	}

	result, err := repo.GetUserLogin(ctx, activeSession.SID)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if result == "" {
		t.Errorf("session not found")
	}

	nonExistentSessionID := "non_existent_session_id"
	result, err = repo.GetUserLogin(ctx, nonExistentSessionID)
	if err == nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if result != "" {
		t.Errorf("expected empty result for non-existent session")
	}

	client.Close()
	_, err = repo.GetUserLogin(ctx, activeSession.SID)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestDeleteSession(t *testing.T) {
	ctx := context.Background()
	activeSession := models.Session{
		SID:   "session_id",
		Login: "user_login",
	}

	client := GetTestClient()

	repo := &SessionRepo{
		DB: client,
	}

	removed, err := repo.DeleteSession(ctx, activeSession.SID)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if !removed {
		t.Errorf("session not deleres")
	}

	client.Close()
	_, err = repo.DeleteSession(ctx, activeSession.SID)
	if err == nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

}
