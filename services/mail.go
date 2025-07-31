package services

import (
	"bytes"
	"ems/config"
	"ems/consts"
	"ems/domain"
	"ems/models"
	"ems/types"
	"ems/utils/errutil"
	"ems/worker"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
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
		return fmt.Errorf("error sending email to %s: %w", requestData.MailTo, err)
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

func (m *MailService) EnqueueReminderEmailNotification(roleIds []int, campaign *models.Campaign) error {
	startTime := campaign.StartTime
	reminderTime := startTime.Add(-consts.ReminderInterval)
	now := time.Now()

	//Reminder time: 9:00 AM
	//Current time: 9:30 AM
	//Reminder time is already past, so skip sending the email.

	if reminderTime.Before(now) {
		log.Info("Reminder time is in the past, skipping reminder email :", campaign.Title)
		return errutil.ErrReminderEmailNotEnqueued
	}
	timeLeftToSendReminderEmail := reminderTime.Sub(now)
	schedular := worker.NewScheduler(timeLeftToSendReminderEmail)
	schedular.Start(func() {
		m.SendReminderEmail(roleIds, campaign)
	})
	time.Sleep(time.Duration(timeLeftToSendReminderEmail))
	schedular.Stop()
	return nil
}

func (m *MailService) SendReminderEmail(roleIds []int, campaign *models.Campaign) {
	users, err := m.userRepo.ReadUsers(roleIds)
	if err != nil {
		if err == errutil.ErrRecordNotFound {
			log.Printf("Failed to read users: %v", err)
			return
		}
		log.Error(fmt.Sprintf("Failed to read users: %v", err))
		return
	}
	for _, user := range users {
		emailPayload := types.EmailPayload{
			MailTo:  user.Email,
			Subject: "Reminder: " + campaign.Title,
			Body: map[string]interface{}{
				"campaign_title":       campaign.Title,
				"campaign_description": campaign.Description,
				"campaign_remarks":     campaign.Remarks,
				"campaign_startTime":   campaign.StartTime,
				"campaign_endTime":     campaign.EndTime,
			},
		}
		task := worker.NewTask(func() error {
			return m.SendEmail(emailPayload)
		}, func(err error) {
			log.Error("Failed to send reminder email: ", err, "to user: ", user.Email)
		}, 0)
		m.workerPool.AddTask(task)
	}
}
