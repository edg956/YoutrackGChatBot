# YoutrackGChatBot
An exploration project into golang that implements a Google Chat bot to interact with Youtrack.

## Motivation
Google Chat and Youtrack are both fine tools for development teams and companies allowing communication and scrum/kanban board management for projects, but there is currently no integration between them. I.E: In a Chat someone might suggest creating an issue as a result of a discussion, which requires someone manually creating a ticket in their project's current sprint or backlog and hides the information logged into Youtrack to all the participants, which may need to go onto the website and check what was created, update if needed and notify through chat.

This aims to make these actions transparent by allowing issues to be created in the same chat they're proposed through bot commands to interact with Youtrack.

Also, I just wanted to program something in Go, even though a small webserver :).

## TODO
- [ ] Implement handlers for MESSAGE, ADDED, and REMOVED events
- [ ] Implement Card Message to configure settings for a project
- [ ] Implement web form that takes in the settings from the card message
- [ ] Add in database models for storing spaces' configuration
- [ ] Implement commands to list sprints of a project's board
- [ ] Implement commands to manage issues (create, update, delete, assign)
- [ ] Implement functionality to update when a new sprint is created (including updating configuration)
- [ ] Implement commands to assign a sprint as default to create issues, etc
- [ ] Implement CI/CD
- [ ] Ideas?
