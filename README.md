# Plex TVTime Sync
This project is a Go-based tool that syncs your Plex viewing history with TVTime. It was designed to allow users to easily keep track of the TV shows and movies they've watched on Plex and have that reflected on TVTime.

- [Plex TVTime Sync](#plex-tvtime-sync)
  - [Prerequisites](#prerequisites)
  - [Project Structure](#project-structure)
  - [Configuration](#configuration)
  - [Run](#run)
    - [Development](#development)
    - [Production](#production)
  - [Credits](#credits)
---
## Prerequisites
- Golang 1.20 or higher is required to build and run the project. You can find the installer on
  the official Golang [download](https://go.dev/doc/install) page.
- Plex Account
- TVTime Account
## Project Structure

Here's a high-level overview of the project structure:
```plaintext
.
├── bootstrap: Contains the application's main initialization logic.
├── build: Contains build scripts, if any.
├── commands: Contains the application's command handling code.
├── config: Contains configuration files for different environments (development, production, staging).
├── domain: Contains the application's domain logic, including entities, repositories, and use cases.
│   ├── entities: Defines the application's entities.
│   ├── repositories: Contains the repository layer of the application.
│   └── usecases: Contains the use case layer of the application.
├── files: Contains JSON files for storing data.
├── infrastructure: Contains code related to infrastructure, such as APIs and processes.
│   ├── api: Contains the API layer of the application.
│   └── process: Contains the process layer of the application.
├── logs: Contains log files.
└── pkg: Contains library-wide functions, variables, and configuration logic.
    └── lib: Contains helper functions for handling commands, configuration, logging, and storage.
```

```plaintext
.
├── bootstrap
│   ├── app.go: Contains the application's main initialization logic.
│   └── modules.go: Contains code to initialize the different modules used in the application.
├── build: Contains build scripts, if any.
├── commands
│   ├── commands.go: Defines the application's command interface.
│   └── run.go: Contains the implementation of the run command.
├── config
│   ├── development.json: Configuration file for the development environment.
│   ├── production.json: Configuration file for the production environment.
│   └── staging.json: Configuration file for the staging environment.
├── docker-compose.yml: Docker Compose file to set up your development environment.
├── Dockerfile: Dockerfile for building your application's Docker image.
├── domain
│   ├── entities
│   │   ├── episode.go: Defines the Episode entity.
│   │   ├── history.go: Defines the History entity.
│   │   ├── search.go: Defines the Search entity.
│   │   ├── season.go: Defines the Season entity.
│   │   └── show.go: Defines the Show entity.
│   ├── repositories
│   │   ├── plex_repository.go: Contains the repository layer for Plex.
│   │   ├── repositories.go: Defines the repository interface.
│   │   └── tvtime_repository.go: Contains the repository layer for TVTime.
│   └── usecases
│       ├── plex_usecase.go: Contains the use case layer for Plex.
│       ├── tvtime_usecase.go: Contains the use case layer for TVTime.
│       └── usecases.go: Defines the use case interface.
├── files
│   └── link.json: JSON file to store links.
├── go.mod and go.sum: Go module and sum files.
├── infrastructure
│   ├── api
│   │   ├── api.go: Defines the API interface.
│   │   ├── plex_api.go: Contains the API layer for Plex.
│   │   └── tvtime_api.go: Contains the API layer for TVTime.
│   └── process
│       ├── process.go: Defines the process interface.
│       └── sync_process.go: Contains the implementation of the sync process.
├── logs
│   └── sync.log: Log file for sync operations.
├── main.go: The application's entry point.
├── Makefile: Contains build commands.
├── pkg
│   └── lib
│       ├── cmd.go: Contains helper functions for handling commands.
│       ├── config.go: Contains the application's configuration logic.
│       ├── lib.go: Contains library-wide functions or variables.
│       ├── logger.go: Contains the application's logging logic.
│       └── storage.go: Contains the application's storage logic.
├── README.md: This file.
├── RELEASE_NOTES.md: Contains release notes.
└── VERSION: Contains the current version of the application.
```
## Configuration

This application uses a JSON configuration file to set up various parameters. Here's a brief overview of the configuration fields:

| Key | Description | Default Value |
| --- | ----------- | ------------- |
| `environment` | Sets the environment mode for the application. | `development` |
| `log_output` | Specifies the file where the application's logs will be written. | `logs/sync.log` |
| `log_level` | Sets the level of logs that will be written. | `debug` |
| `tz` | Sets the timezone for the application. | `Europe/Paris` |
| `tv_time.token.symfony` | The `symfony` token needed to interact with the TVTime API. | `""` |
| `tv_time.token.tvst_remember` | The `tvst_remember` token needed to interact with the TVTime API. | `""` |
| `tv_time.accept_language` | Sets the `Accept-Language` header for TVTime API requests. | `""` |
| `plex.base_url` | The base URL of your Plex server. | `Your Plex url` |
| `plex.token` | Your Plex access token. | `Your Plex token` |
| `plex.account_id` | Your Plex account ID. | `1` |
| `plex.init_viewed_at` | Sets the initial date for viewing history. | `""` |
| `timer` | Sets the interval (in minute) at which the application syncs with Plex and TVTime. | `30` |
| `file_storage.filename` | The file where the application stores persistent data. | `files/link.json` |

You should adjust these settings according to your needs and environment. Please ensure to replace all the placeholders with your actual data.
## Run

### Development
Run the following
```sh
$ git clone git@github.com:AC-CodeProd/plex-tvtime-sync.git
$ cd plex-tvtime-syn
$ mv config/{{ENVIRONMENT}}.json.sample config/{{ENVIRONMENT}}.json
```
With docker-compose
```sh
$ docker-compose -f docker-compose.yml up -d
```
OR Run the Go program in live-reload mode using the 'air'
```sh
$ make run-live
```
### Production
Download binaries for Linux, macOS, and Windows are available as [Github Releases](https://github.com/AC-CodeProd/plex-tvtime-sync/releases/latest).
Using binarie:
```sh
$ curl -o {{ENVIRONMENT}}.json https://raw.githubusercontent.com/AC-CodeProd/plex-tvtime-sync/main/config/{{ENVIRONMENT}}.json.sample 
$ plex-tvtime-sync-v{{VERSION}}-{{ARCHITECTURES}} run -c {{ENVIRONMENT}}.json
```
OR Build binarie
```sh
$ git clone git@github.com:AC-CodeProd/plex-tvtime-sync.git
$ cd plex-tvtime-syn
$ make build
$ cd build
$ plex-tvtime-sync-v{{VERSION}}-{{ARCHITECTURES}} run -c ../config/{{ENVIRONMENT}}.json
```
## Credits
This work was inspired by <a href="https://github.com/Paypito/plextvtimesync" target="_blank">Paypito</a>