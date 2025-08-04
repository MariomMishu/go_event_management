package services

import (
	"ems/config"
	"ems/domain"
	"ems/models"
	"ems/types"
	"ems/worker"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/smtp"
	"strings"
	"time"
)

type MailService struct {
	userRepo   domain.UserRepository
	mailRepo   domain.MailRepository
	workerPool *worker.Pool
}

func NewMailService(userRepo domain.UserRepository, mailRepo domain.MailRepository, workerPool *worker.Pool) *MailService {
	return &MailService{
		userRepo:   userRepo,
		mailRepo:   mailRepo,
		workerPool: workerPool,
	}
}

func (m *MailService) SendEmail(requestData types.EmailPayload) error {
	err := m.mailRepo.SendEmail(&requestData)
	if err != nil {
		log.Error("Failed to send email: ", err)
		return err
	}
	return nil
}

func (m *MailService) SendCampaignEmail(roleIds []int, campaign *models.Campaign) error {
	users, err := m.userRepo.ReadUsers(roleIds)
	if err != nil {
		log.Printf("Failed to read users: %v", err)
		return err
	}

	for _, user := range users {
		emailPayload := types.EmailPayload{
			MailTo:  user.Email,
			Subject: "Test Email For Campaign",
			Body: map[string]interface{}{
				"Title":       campaign.Title,
				"Description": campaign.Description,
				"Remarks":     campaign.Remarks,
			},
		}

		//if err := m.SendEmail(emailPayload); err != nil {
		//	log.Printf("Failed to send email to %v: %v", user.Email, err)
		//	// Optionally, you can continue sending emails to others instead of returning immediately
		//	return err
		//}
		//log.Infof("Sent Invitation Mail to %s", user.Email)

		//Add the email sending task to the worker pool
		task := worker.NewTask(func() error {
			return m.SendEmail(emailPayload)
		}, func(err error) {
			log.Error("Failed to send email: ", err, "to user: ", user.Email)
		}, 3)
		m.workerPool.AddTask(task)
	}
	return nil
}
func SendEmailSMTP(to []string, subject, body string) error {
	emailCfg := config.Email()

	auth := smtp.PlainAuth("", emailCfg.Username, emailCfg.Password, emailCfg.Host)

	msg := "From: " + emailCfg.Username + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	addr := emailCfg.Host + ":" + emailCfg.Port

	// Optional timeout handling (basic)
	done := make(chan error, 1)
	go func() {
		err := smtp.SendMail(addr, auth, emailCfg.Username, to, []byte(msg))
		done <- err
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(emailCfg.Timeout):
		return fmt.Errorf("SMTP send timeout")
	}
}
