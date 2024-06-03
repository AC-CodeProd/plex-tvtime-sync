package usecases

import (
	"plex-tvtime-sync/domain/entities"
	"plex-tvtime-sync/domain/interfaces"
	"plex-tvtime-sync/pkg/lib"

	"go.uber.org/fx"
)

type EmailUseCaseParams struct {
	fx.In

	Logger lib.Logger
	Mailer interfaces.IEmail
}

type emailUseCase struct {
	logger lib.Logger
	mailer interfaces.IEmail
}

func NewEmailUseCase(puP EmailUseCaseParams) interfaces.IEmailUsecase {
	return &emailUseCase{
		logger: puP.Logger,
		mailer: puP.Mailer,
	}
}

func (uc *emailUseCase) SendEmail(email *entities.Email) error {
	return uc.mailer.Send(email)
}

func (uc *emailUseCase) SendEmailWithTemplate(email *entities.Email, data interface{}) error {
	body, err := uc.mailer.RenderTemplate(email.TemplateName, data)
	if err != nil {
		return err
	}
	email.Body = body
	return uc.mailer.Send(email)
}
