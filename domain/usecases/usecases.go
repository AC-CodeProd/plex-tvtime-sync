package usecases

import "go.uber.org/fx"

// exports usecases dependency
var Module = fx.Options(
	fx.Provide(NewTvTimeUsecase),
	fx.Provide(NewPlexUsecase),
	fx.Provide(NewEmailUseCase),
	fx.Provide(NewStorageUseCase),
)
