package services

import (
	"bytes"
	"ems/config"
	"ems/domain"
	"ems/models"
	"ems/types"
	"ems/worker"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
)

type MailService struct {
	userRepo    domain.UserRepository
	emailClient *http.Client
	workerPool  *worker.Pool
}

func NewMailService(userRepo domain.UserRepository, emailClient *http.Client, workerPool *worker.Pool) *MailService {
	return &MailService{
		userRepo:    userRepo,
		emailClient: emailClient,
		workerPool:  workerPool,
	}
}

func (m *MailService) SendEmail(requestData types.EmailPayload) error {
	reqURL := config.Email().Url

	reqByte, err := json.Marshal(requestData)
	if err != nil {
		log.Printf("Failed to marshal email payload: %v", err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(reqByte))
	if err != nil {
		log.Printf("Failed to create email request: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.emailClient.Do(req)
	if err != nil {
		log.Printf("Failed to send email request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Email service returned non-success status: %v", resp.Status)
		return fmt.Errorf("email service returned non-success status: %v", resp.Status)
	}

	log.Infof("Email service responded with status [%v] for recipient: %v", resp.StatusCode, requestData.MailTo)
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
