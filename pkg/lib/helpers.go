package lib

import (
	"go.uber.org/fx"
)

type Helpers struct {
	config Config
	logger Logger
}

// Params defines the base objects for a storage.
type HelpersParams struct {
	fx.In
	Config Config
	Logger Logger
}

// Result defines the objects that the storage module provides.
type HelpersResult struct {
	fx.Out

	Helpers Helpers
}

func NewHelpers(hP HelpersParams) HelpersResult {
	helpers := Helpers{
		config: hP.Config,
		logger: hP.Logger,
	}
	return HelpersResult{
		Helpers: helpers,
	}
}

// func (h *Helpers) CloseFile(f io.ReadCloser) {
// 	if f != nil {
// 		f.Close()
// 	}
// }
