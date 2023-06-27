package api

import (
	"io"

	"go.uber.org/fx"
)

// exports api dependency
var Module = fx.Options(
	fx.Provide(GetTVTimeApi),
	fx.Provide(GetPlexApi),
)

func closeFile(f io.ReadCloser) {
	if f != nil {
		f.Close()
	}
}
