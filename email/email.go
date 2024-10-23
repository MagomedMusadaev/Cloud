package email

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendConfirmationEmail отправляет электронное письмо с кодом подтверждения на указанный адрес.
// to - адрес электронной почты получателя.
// code - код подтверждения, который будет отправлен в письме.
func SendConfirmationEmail(to, code string) error {
	// Получаем переменные окружения для настройки почты
	from := os.Getenv("MAIL_FROM")         // Адрес электронной почты отправителя
	password := os.Getenv("MAIL_PASSWORD") // Пароль для SMTP-сервера
	smtpHost := os.Getenv("MAIL_HOST")     // Хост SMTP-сервера
	smtpPort := os.Getenv("MAIL_PORT")     // Порт SMTP-сервера

	// Формируем сообщение, включая заголовок и тело письма
	msg := []byte(fmt.Sprintf("Subject: Подтверждение регистрации\n\nВаш код подтверждения: %s", code))

	// Настраиваем аутентификацию для отправки почты
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Отправляем почту с использованием smtp.SendMail
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return err // Возвращаем ошибку, если отправка не удалась
	}

	return nil // Возвращаем nil, если отправка прошла успешно
}
