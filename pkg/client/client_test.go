package client

import (
	"testing"
)

func TestRedisClient_AddCredentials(t *testing.T) {
	mockUsername := "api-key"
	mockPassword := "api-token"

	client := NewClient(mockUsername, mockPassword, "", "")

	if client.Username != mockUsername {
		t.Errorf("Set username failed. Expected %s, got %s", mockUsername, client.Username)
	}
	if client.Password != mockPassword {
		t.Errorf("Set password failed. Expected %s, got %s", mockPassword, client.Password)
	}
}
