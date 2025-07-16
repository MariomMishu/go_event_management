package domain

import (
	"ems/models"
)

// Campaign represents an campaign entity.
type Campaign struct {
	ID          int
	Title       string
	Description string
	StartTime   string // Consider using time.Time in real applications
	EndTime     string
	Remarks     string
	Status      string
}

// CampaignRepository defines the interface for campaign data operations.
type CampaignRepository interface {
	CreateCampaign(campaign *models.Campaign) (*models.Campaign, error)
	ReadCampaignById(id int) (*models.Campaign, error)
	ReadCampaignByTitle(title string) (*models.Campaign, error)
	UpdateCampaign(campaign *models.Campaign) error
	DeleteCampaign(id int) error
	ApproveCampaign(id int) error
	ListCampaigns() ([]*models.Campaign, error)
}

// CampaignService defines the business logic for campaign.
type CampaignService interface {
	CreateCampaign(campaign *models.Campaign) (*models.Campaign, error)
	GetCampaignByID(id int) (*models.Campaign, error)
	UpdateCampaign(campaign *models.Campaign) error
	DeleteCampaign(id int) error
	ListCampaigns() ([]*models.Campaign, error)
}
