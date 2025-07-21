package services

import (
	"bytes"
	"ems/config"
	"ems/domain"
	"ems/models"
	"ems/types"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
)

type MailService struct {
	userRepo    domain.UserRepository
	emailClient *http.Client
}

func NewMailService(userRepo domain.UserRepository, emailClient *http.Client) *MailService {
	return &MailService{
		userRepo:    userRepo,
		emailClient: emailClient,
	}
}
func (m *MailService) SendEmail(requestData types.EmailPayload) error {
	reqURL := config.Email().Url

	reqByte, _ := json.Marshal(requestData)

	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(reqByte))
	if err != nil {
		return err
	}
	resp, err := m.emailClient.Do(req)
	log.Printf("reqURL : %v, resp : %v", reqURL, resp)
	if err != nil {
		log.Printf("err : %v", err)
		return err
	}
	fmt.Printf("email service status code [%v] after email send to: %v \n", resp.StatusCode, requestData.MailTo)
	return nil
}
func (m *MailService) SendCampaignEmail(roleIds []int, campaign *models.Campaign) error {
	users, err := m.userRepo.ReadUsers(roleIds)
	if err != nil {
		return err
	}
	for _, user := range users {
		emailPayload := types.EmailPayload{
			MailTo:  "testEms@yopmail.com",
			Subject: "Test Email FoR Campaign",
			Body: map[string]interface{}{
				"Title":       campaign.Title,
				"Description": campaign.Description,
				"Remarks":     campaign.Remarks,
			},
		}
		err := m.SendEmail(emailPayload)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("Send Invitation Mail to %s", user.Email))
	}
	return nil
}
