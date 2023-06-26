package process

import "go.uber.org/fx"

// exports process dependency
var Module = fx.Options(
	fx.Provide(NewSyncProcess),
)
