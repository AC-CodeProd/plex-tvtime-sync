package mailer

import (
	"go.uber.org/fx"
)

// exports api dependency
var Module = fx.Options(
	fx.Provide(NewSMTPMailer),
)
