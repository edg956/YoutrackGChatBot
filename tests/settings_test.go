package test

import (
	"YoutrackGChatBot/settings"
	"os"
	"testing"
)

func TestGetSettingsWithoutEnvironmentSetFails(t *testing.T) {
	if _, err := settings.GetSettings(); err == nil {
		t.Errorf("GetSettings should fail if envvars are not set")
	}
}

func TestGetSettingsWithEnvironmentSetWorks(t *testing.T) {
	testYoutrackToken := "SOMEVALUE"
	testAudience := "SOMEOTHERVALUE"
	os.Setenv("YOUTRACK_TOKEN", testYoutrackToken)
	os.Setenv("GCHAT_AUDIENCE", testAudience)

	if settings, err := settings.GetSettings(); err == nil {
		if settings.YOUTRACK_TOKEN != testYoutrackToken {
			t.Errorf("YOUTRACK_TOKEN expected to be %s, got %s instead.\n", testYoutrackToken, settings.YOUTRACK_TOKEN)
		}

		if settings.GCHAT_AUDIENCE != testAudience {
			t.Errorf("GCHAT_AUDIENCE expected to be %s, got %s instead.\n", testAudience, settings.GCHAT_AUDIENCE)
		}
	} else {
		t.Errorf("Unexpected error %v\n", err)
	}
}
