package mailer

import (
	"bytes"
	"fmt"
	"net/smtp"
	"path/filepath"
	"plex-tvtime-sync/domain/entities"
	"plex-tvtime-sync/domain/interfaces"
	"text/template"
)

func NewSMTPMailer(smp SMTPMailerParams) interfaces.IEmail {
	// return &SMTPMailer{Host: host, Port: port, Username: username, Password: password}
	return &smtpMailer{
		logger: smp.Logger,
		config: smp.Config,
	}
}

func (sm *smtpMailer) loadTemplate() error {
	tmpl, err := template.ParseGlob(filepath.Join(sm.config.Mailer.TemplateDir, "*.html"))
	if err != nil {
		return err
	}
	sm.templates = tmpl
	return nil
}

func (sm *smtpMailer) Send(email *entities.Email) error {
	from := sm.config.Mailer.SMTP.Username
	password := sm.config.Mailer.SMTP.Password
	to := email.Recipient
	smtpHost := sm.config.Mailer.SMTP.Host
	smtpPort := sm.config.Mailer.SMTP.Port
	auth := smtp.PlainAuth("", from, password, smtpHost)
	// msg := []byte("To: " + email.Recipient + "\r\n" +
	// 	"Subject: " + email.Subject + "\r\n" +
	// 	"MIME-version: 1.0;\r\n" +
	// 	"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
	// 	email.Body)

	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	boundary := "unique-boundary-1"
	header := make(map[string]string)
	header["From"] = from
	header["To"] = to
	header["Subject"] = email.Subject
	header["MIME-Version"] = "1.0"
	// header["Content-Type"] = "text/plain; charset=\"UTF-8\";\n\n"
	header["Content-Type"] = "multipart/mixed; boundary=\"" + boundary + "\""
	// header := m.header(strings.Join(to, ", "), subject, "multipart/mixed; boundary=\""+boundary+"\"")
	var msg bytes.Buffer
	for k, v := range header {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n--" + boundary + "\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
	msg.WriteString(email.Body)
	if email.SectionSuccessEmails != nil {
		for _, row := range email.SectionSuccessEmails {
			for _, image := range row {
				msg.WriteString("\r\n--" + boundary + "\r\n")
				msg.WriteString(fmt.Sprintf("Content-ID: <%s>\r\n", image.CID))
				msg.WriteString(fmt.Sprintf("Content-Disposition: inline; filename=\"%s\"\r\n", image.Title))
				msg.WriteString("Content-Transfer-Encoding: base64\r\n")
				msg.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n\r\n", "image/jpeg", image.Title))
				msg.WriteString(image.Data)
			}
		}
	}
	msg.WriteString("\r\n--" + boundary + "--")
	return smtp.SendMail(addr, auth, from, []string{to}, msg.Bytes())
}

func (sm *smtpMailer) RenderTemplate(templateName string, data interface{}) (string, error) {
	var tpl bytes.Buffer
	if sm.templates == nil {
		sm.loadTemplate()
	}
	if templateName == "" {
		templateName = "default.html"
	}
	if err := sm.templates.ExecuteTemplate(&tpl, templateName, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
