package session

import "testing"

func TestStateMachineLifecycle(t *testing.T) {
	sm := NewStateMachine()
	chatID := int64(123)

	if sm.IsActive(chatID) {
		t.Fatalf("expected inactive session")
	}

	if !sm.StartSession(chatID) {
		t.Fatalf("expected StartSession to succeed")
	}

	if sm.StartSession(chatID) {
		t.Fatalf("expected duplicate StartSession to fail")
	}

	sm.AppendBlock(chatID, TextBlock{})
	if !sm.ClearSession(chatID) {
		t.Fatalf("expected ClearSession to succeed")
	}

	session := sm.GetSession(chatID)
	if session == nil || len(session.Blocks) != 0 {
		t.Fatalf("expected cleared session")
	}

	ended, ok := sm.EndSession(chatID)
	if !ok || ended == nil {
		t.Fatalf("expected EndSession to return session")
	}

	if sm.IsActive(chatID) {
		t.Fatalf("expected inactive after EndSession")
	}

	if sm.DiscardSession(chatID) {
		t.Fatalf("expected DiscardSession to fail on inactive")
	}
}
