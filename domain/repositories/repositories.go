package repositories

import "go.uber.org/fx"

// exports repositories dependency
var Module = fx.Options(
	fx.Provide(NewTvTimeRepository),
	fx.Provide(NewPlexRepository),
)
