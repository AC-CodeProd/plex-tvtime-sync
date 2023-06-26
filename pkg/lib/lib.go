package lib

import "go.uber.org/fx"

// exports libraries dependency
var Module = fx.Options(
	fx.Provide(GetLogger),
	fx.Provide(GetConfig),
	fx.Provide(GetJsonStorage),
	fx.Provide(NewHelpers),
)
