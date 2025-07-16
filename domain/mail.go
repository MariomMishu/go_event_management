package domain

type (
	MailService interface {
		SendMail(userIds []int) error
	}
)
