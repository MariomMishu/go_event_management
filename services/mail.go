package services

import (
	"ems/domain"
	"fmt"
	"github.com/labstack/gommon/log"
)

type MailService struct {
	userRepo domain.UserRepository
}

func NewMailService(userRepo domain.UserRepository) *MailService {
	return &MailService{
		userRepo: userRepo,
	}
}
func (s *MailService) SendMail(userIds []int) error {
	users, err := s.userRepo.ReadUsers(userIds)
	if err != nil {
		return err
	}
	for _, user := range users {
		log.Info(fmt.Sprintf("Send Invitation Mail to %s", user.Email))
	}
	return nil
}
