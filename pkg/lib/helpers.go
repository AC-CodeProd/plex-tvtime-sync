package lib

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

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

func (h *Helpers) FuncNameAndFile() (string, error) {
	const names = "__helpers.go__ : FuncNameAndFile"
	fpcs := make([]uintptr, 1)

	// Skip 2 levels to get the caller
	n := runtime.Callers(2, fpcs)
	if n == 0 {
		h.logger.Error(fmt.Sprintf("%s | %s", names, "MSG: NO CALLER"))
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		h.logger.Error(fmt.Sprintf("%s | %s", names, "MSG CALLER WAS NIL"))
	}
	// Get the full name of the function
	fullName := caller.Name()

	// Split the full name into parts
	parts := strings.Split(fullName, ".")

	// Get the name of the function
	funcName := parts[len(parts)-1]
	// Get the file information
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		h.logger.Error(fmt.Sprintf("%s | %s", names, "MSG: Could not get file info"))
	}

	return fmt.Sprintf("__%s__ : %s", filepath.Base(file), funcName), nil
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
