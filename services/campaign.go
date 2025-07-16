package services

import (
	"ems/domain"
	"ems/types"
	"ems/utils/errutil"
)

type CampaignServiceImpl struct {
	repo domain.CampaignRepository
}

func NewCampaignServiceImpl(campaignRepo domain.CampaignRepository) *CampaignServiceImpl {
	return &CampaignServiceImpl{
		repo: campaignRepo,
	}
}
func (svc *CampaignServiceImpl) CreateCampaign(req *types.CampaignCreateRequest) (*types.CampaignCreateResponse, error) {
	// Check if a campaign with the same title already exists
	isExist, err := svc.IsCampaignExist(req.Title)
	if err != nil {
		return nil, err
	}
	if isExist {
		return nil, errutil.ErrAlreadyExists
	}

	// Convert request to campaign model
	campaign := req.ToCampaignModel()

	// Create campaign using repository
	createdCampaign, err := svc.repo.CreateCampaign(campaign)
	if err != nil {
		return nil, err
	}

	// Return successful response
	return &types.CampaignCreateResponse{
		Message:  "Campaign created successfully",
		Campaign: createdCampaign,
	}, nil
}

// IsCampaignExist checks whether a campaign with the given title exists
func (svc *CampaignServiceImpl) IsCampaignExist(title string) (bool, error) {
	existing, err := svc.repo.ReadCampaignByTitle(title)
	if err != nil {
		return false, err
	}
	return existing != nil, nil
}

// Define methods for user service
