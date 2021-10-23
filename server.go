package main

import (
	handlers "YoutrackGChatBot/handlers"
	logging "YoutrackGChatBot/logging"
	middlewares "YoutrackGChatBot/middleware"
	settings "YoutrackGChatBot/settings"
	"net/http"
)

func main() {
	logger := logging.GetLogger()
	// Init settings
	_, err := settings.GetSettings()

	if err != nil {
		logger.Fatal(err)
		return
	}

	server := http.NewServeMux()

	server.HandleFunc("/", middlewares.ApplyMiddleware(handlers.Hello))

	listeningAddress := ":8080"
	logger.Printf("Server listening on address: %s", listeningAddress)
	http.ListenAndServe(listeningAddress, server)
}
