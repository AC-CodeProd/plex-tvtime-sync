package bootstrap

import (
	"plex-tvtime-sync/domain/repositories"
	"plex-tvtime-sync/domain/usecases"
	"plex-tvtime-sync/infrastructure/api"
	"plex-tvtime-sync/infrastructure/mailer"
	"plex-tvtime-sync/infrastructure/storage"
	"plex-tvtime-sync/interfaces"
	"plex-tvtime-sync/interfaces/process"
	"plex-tvtime-sync/pkg/lib"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	api.Module,
	mailer.Module,
	storage.Module,
	process.Module,
	interfaces.Module,
	lib.Module,
	usecases.Module,
	repositories.Module,
)
