package services

import (
	"ems/domain"
	"ems/models"
)

type MailService struct {
	userRepo domain.UserRepository
}

func NewMailService(userRepo domain.UserRepository) *MailService {
	return &MailService{
		userRepo: userRepo,
	}
}
func (m *MailService) SendMail(userIds []int, campaign *models.Event) error {
	//users, err := m.userRepo.ReadUsers(userIds)
	//if err != nil {
	//	return err
	//}
	//for _, user := range users {
	//	emailPayload := types.EmailPayload{
	//		MailTo:  user.Email,
	//		Subject: "Email Subject" + campaign.Title,
	//		Body: map[string]interface{}{
	//			"campaign":  campaign,
	//			"rsvp_link": fmt.Sprintf("http://127.0.0.1:8080/v1/campaign/%d/rsvp", campaign.ID),
	//		},
	//	}
	//	m.SendMail(emailPayload)
	//	log.Info(fmt.Sprintf("Send Invitation Mail to %s", user.Email))
	//}
	return nil
}
