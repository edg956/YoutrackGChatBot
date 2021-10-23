package settings

import (
	logging "YoutrackGChatBot/logging"
	"errors"
	"os"
)

type Settings struct {
	YOUTRACK_TOKEN         string
	GCHAT_ISSUER           string
	PUBLIC_CERT_URL_PREFIX string
	GCHAT_AUDIENCE         string
	DATABASE_PATH          string
}

func GetSettings() (Settings, error) {
	logger := logging.GetLogger()

	logger.Println("Discovering settings from environment")

	youtrackToken := os.Getenv("YOUTRACK_TOKEN")
	if youtrackToken == "" {
		return Settings{}, errors.New("Missing YOUTRACK_TOKEN.")
	}

	gchatAudience := os.Getenv("GCHAT_AUDIENCE")
	if gchatAudience == "" {
		return Settings{}, errors.New("Missing GCHAT_AUDIENCE.")
	}

	gchatIssuer := os.Getenv("GCHAT_ISSUER")
	if gchatIssuer == "" {
		gchatIssuer = "chat@system.gserviceaccount.com"
	}

	gchatCertUrl := os.Getenv("PUBLIC_CERT_URL_PREFIX")
	if gchatCertUrl == "" {
		gchatCertUrl = "https://www.googleapis.com/service_accounts/v1/metadata/x509/"
	}

	databasePath := os.Getenv("DATABASE_PATH")
	if databasePath == "" {
		databasePath = "/var/lib/youtrack-gchat-bot/data"
	}

	return Settings{
		YOUTRACK_TOKEN:         youtrackToken,
		GCHAT_AUDIENCE:         gchatAudience,
		GCHAT_ISSUER:           gchatIssuer,
		PUBLIC_CERT_URL_PREFIX: gchatCertUrl,
		DATABASE_PATH:          databasePath,
	}, nil

}
