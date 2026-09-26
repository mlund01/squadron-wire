package protocol

import (
	"encoding/json"
	"testing"
)

func TestMissionRuntimeResultRoundTrip(t *testing.T) {
	result := GetMissionResultPayload{Runtime: &MissionRuntimeInfo{
		Revision: 3, State: "running",
		Commander: MissionRuntimeCommanderInfo{ID: "run:commander", State: "active", Messages: []MissionRuntimeMessageInfo{{
			ID: "message-1", Role: "assistant", Parts: []MissionRuntimeContentPart{{Type: "tool_use", ToolName: "dispatch_task", ToolUseID: "call-1", ToolInputJSON: `{"task":"research"}`}},
		}}},
		Tasks: []MissionRuntimeTaskInfo{{Name: "research", State: "running", Attempt: 1, Attempts: []MissionRuntimeAttemptInfo{{ID: "attempt-1", Model: "models.fast"}}}},
	}}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	var decoded GetMissionResultPayload
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Runtime == nil || decoded.Runtime.Revision != 3 || decoded.Runtime.Commander.Messages[0].Parts[0].ToolName != "dispatch_task" || decoded.Runtime.Tasks[0].Attempts[0].ID != "attempt-1" {
		t.Fatalf("runtime response changed in transit: %#v", decoded.Runtime)
	}
}

func TestOlderMissionResultHasNoRuntimeProjection(t *testing.T) {
	var decoded GetMissionResultPayload
	if err := json.Unmarshal([]byte(`{"mission":{"id":"run-1","name":"legacy","status":"completed","startedAt":"2026-01-01T00:00:00Z"},"tasks":[]}`), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Runtime != nil {
		t.Fatalf("unexpected runtime projection for older payload: %#v", decoded.Runtime)
	}
}
