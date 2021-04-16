package service

import (
	"fmt"
	"github.com/MJ7898/myaktion-go/src/myaktion/model"
	log "github.com/sirupsen/logrus"
)


var (
	campaignStore map[uint]*model.Campaign
	actCampaignId uint = 1 )

func init() {
	campaignStore = make(map[uint]*model.Campaign)
}

func CreateCampaign(campaign *model.Campaign) error {
	campaign.ID = actCampaignId
	campaignStore[actCampaignId] = campaign
	actCampaignId += 1
	log.Infoln("Successfully stored new campaign with ID %v in database.", campaign.ID)
	log.Tracef("Stored: %v", campaign)
	return nil
}

func GetCampaigns() ([]model.Campaign, error) {
	var campaigns []model.Campaign
	for _, campaign := range campaignStore {
		campaigns = append(campaigns, *campaign) }
	log.Tracef("Retrieved: %v", campaigns)
	return campaigns, nil
}

// GetCampaignById function to get an single campaign with a given id/***
func GetCampaignById(id uint) (*model.Campaign, error) {
	if campaign, ok := campaignStore[id]; ok {
		return campaign, nil
	}
	log.Errorf("Campaign with ID %v not found.", id)
	return nil, nil
}

func UpdateCampaign(id uint, campaign *model.Campaign) (*model.Campaign, error)  {
	existingCampaign, err := GetCampaignById(id)
	if err != nil {
		return existingCampaign, err
	}
	existingCampaign.Name = campaign.Name
	existingCampaign.OrganizerName = campaign.OrganizerName
	existingCampaign.TargetAmount = campaign.TargetAmount
	existingCampaign.DonationMinimum = campaign.DonationMinimum

	entry := log.WithField("ID", id)
	entry.Info("Successfully updated campaign.")
	entry.Tracef("Updated: %v", existingCampaign)
	return existingCampaign, nil
}

func DeleteCampaign(id uint) (*model.Campaign, error)  {
	campaign:= campaignStore[id]
	if campaign == nil {
		log.Tracef("404 Campaign not found")
		return nil, fmt.Errorf("no campaign with ID %d", id)
	}
	delete(campaignStore,id)
	entry := log.WithField("ID", id)
	entry.Info("Successfully deleted campaign.")
	log.Printf("Successfully deleted campaign with ID %d", id)
	return  campaign, nil
}