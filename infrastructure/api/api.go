package api

import (
	"plex-tvtime-sync/infrastructure/api/plex"
	"plex-tvtime-sync/infrastructure/api/tvtime"

	"go.uber.org/fx"
)

// exports api dependency
var Module = fx.Options(
	fx.Provide(tvtime.NewTVTimeApi),
	fx.Provide(plex.NewPlexApi),
)
