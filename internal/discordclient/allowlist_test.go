package discordclient

import "testing"

func TestIsAllowedUser(t *testing.T) {
	allowed := []string{"user-1", "user-2"}

	if !IsAllowedUser(allowed, "user-1") {
		t.Error("expected user-1 to be allowed")
	}
	if IsAllowedUser(allowed, "user-3") {
		t.Error("expected user-3 to be denied")
	}
	if IsAllowedUser(allowed, "") {
		t.Error("expected empty user id to be denied")
	}
}
