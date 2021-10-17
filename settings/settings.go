package settings

import (
	logging "YoutrackGChatBot/logging"
	"errors"
	"os"
	"sync"
)

type Settings struct {
	YOUTRACK_TOKEN         string
	GCHAT_ISSUER           string
	PUBLIC_CERT_URL_PREFIX string
	GCHAT_AUDIENCE         string
}

var (
	singleton Settings
	once      sync.Once
)

func GetSettings() (*Settings, error) {
	// We run this once to make sure there are no race conditions
	once.Do(func() {
		logger := logging.GetLogger()
		err := settingsFromEnvironment()

		if err != nil {
			logger.Println(err)
		}
	})

	if singleton == (Settings{}) {
		return (&Settings{}), errors.New("Could not load Settings from environment.")
	}

	return &singleton, nil
}

func settingsFromEnvironment() error {
	logger := logging.GetLogger()

	logger.Println("Discovering settings from environment")

	youtrackToken := os.Getenv("YOUTRACK_TOKEN")
	if youtrackToken == "" {
		return errors.New("Missing YOUTRACK_TOKEN.")
	}

	gchatAudience := os.Getenv("GCHAT_AUDIENCE")
	if gchatAudience == "" {
		return errors.New("Missing GCHAT_AUDIENCE.")
	}

	gchatIssuer := os.Getenv("GCHAT_ISSUER")
	if gchatIssuer == "" {
		gchatIssuer = "chat@system.gserviceaccount.com"
	}

	gchatCertUrl := os.Getenv("PUBLIC_CERT_URL_PREFIX")
	if gchatCertUrl == "" {
		gchatCertUrl = "https://www.googleapis.com/service_accounts/v1/metadata/x509/"
	}

	singleton = Settings{
		YOUTRACK_TOKEN:         youtrackToken,
		GCHAT_AUDIENCE:         gchatAudience,
		GCHAT_ISSUER:           gchatIssuer,
		PUBLIC_CERT_URL_PREFIX: gchatCertUrl,
	}
	return nil
}
