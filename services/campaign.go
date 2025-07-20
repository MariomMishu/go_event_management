package services

import (
	"ems/domain"
	"ems/types"
	"ems/utils/errutil"
	"ems/utils/msgutil"
	"errors"
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
		Message:  msgutil.CampaignCreatedSuccessfully(),
		Campaign: createdCampaign,
	}, nil
}

// IsCampaignExist checks whether a campaign with the given title exists
func (svc *CampaignServiceImpl) IsCampaignExist(title string) (bool, error) {
	return svc.repo.ReadCampaignByTitle(title)
}

func (svc *CampaignServiceImpl) GetCampaignByID(campaignID int) (*types.CampaignCommonResponse, error) {
	campaign, err := svc.repo.ReadCampaignById(campaignID)
	if err != nil {
		return nil, err
	}
	return &types.CampaignCommonResponse{
		Message:  "Campaign Fetched successfully",
		Campaign: campaign,
	}, nil
}

func (svc *CampaignServiceImpl) UpdateCampaign(campaign *types.CampaignUpdateRequest) (*types.CampaignUpdateResponse, error) {
	existingCampaign, err := svc.repo.ReadCampaignById(campaign.ID)
	if err != nil {
		return nil, err
	}
	if existingCampaign == nil {
		return nil, errutil.ErrRecordNotFound
	}
	req := campaign.ToCampaignModel()
	updatedCampaign, err := svc.repo.UpdateCampaign(req)
	if err != nil {
		return nil, err
	}
	return &types.CampaignUpdateResponse{
		Message:  "Campaign Updated Successfully",
		Campaign: updatedCampaign,
	}, nil
}

func (svc *CampaignServiceImpl) DeleteCampaign(campaignID int) (*types.CampaignDeleteResponse, error) {
	err := svc.repo.DeleteCampaign(campaignID)
	if err != nil {
		return nil, err
	}
	return &types.CampaignDeleteResponse{
		Message: "Campaign Deleted Successfully",
	}, nil
}

func (svc *CampaignServiceImpl) ApproveRejectCampaign(campaignID int, updatedBy int) (*types.CampaignApproveRejectResponse, error) {
	err := svc.repo.ApproveRejectCampaign(campaignID, updatedBy)
	if err != nil {
		return nil, err
	}
	return &types.CampaignApproveRejectResponse{
		Message: "Event Successful",
	}, nil
}

func (svc *CampaignServiceImpl) ListCampaigns() (*types.CampaignCommonResponseList, error) {
	campaigns, err := svc.repo.ListCampaigns()
	if errors.Is(err, errutil.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &types.CampaignCommonResponseList{
		Message:  "Campaign List Fetched Successfully",
		Campaign: campaigns,
	}, nil
}
