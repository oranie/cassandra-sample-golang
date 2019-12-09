package chat

import "testing"

func TestInsertData(t *testing.T) {
	env, session, chatData := initApp()

	if result != expext {
		t.Error("\n実際： ", result, "\n理想： ", expext)
	}

	t.Log("Test done")
}
