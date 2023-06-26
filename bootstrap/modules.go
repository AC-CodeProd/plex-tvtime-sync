package bootstrap

import (
	"plex-tvtime-sync/domain/process"
	"plex-tvtime-sync/domain/repositories"
	"plex-tvtime-sync/domain/usecases"
	"plex-tvtime-sync/infrastructure/api"
	"plex-tvtime-sync/pkg/lib"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	api.Module,
	process.Module,
	lib.Module,
	usecases.Module,
	repositories.Module,
)
