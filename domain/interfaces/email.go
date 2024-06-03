package interfaces

import "plex-tvtime-sync/domain/entities"

type IEmail interface {
	Send(email *entities.Email) error
	RenderTemplate(templateName string, data interface{}) (string, error)
}

type IEmailUsecase interface {
	SendEmail(email *entities.Email) error
	SendEmailWithTemplate(email *entities.Email, data interface{}) error
}
