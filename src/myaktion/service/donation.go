package service

import (
	"github.com/MJ7898/myaktion-go/src/myaktion/db"
	"github.com/MJ7898/myaktion-go/src/myaktion/model"
	log "github.com/sirupsen/logrus"
)

/**func AddDonation(id uint, donation *model.Donation) (*model.Campaign, error) {
	var existingCampaign *model.Campaign
	existingCampaignById := db.DB.Preload("Campaigns").Find(&existingCampaign)
	if existingCampaignById != nil {
		existingCampaignById.Update(existingCampaign.Name, &donation)
		entry := log.WithField("ID", id)
		entry.Info("Successfully updated campaign.")
		entry.Tracef("Updated: %v", existingCampaign)
		return existingCampaign, nil
	}
	if existingCampaign, ok := campaignStore[id]; ok { -> In-Memory-Lösung
		existingCampaign.Donations = append(existingCampaign.Donations, *donation)
		entry := log.WithField("ID", id)
		entry.Info("Successfully updated campaign.")
		entry.Tracef("Updated: %v", existingCampaign)
		return existingCampaign, nil
	}
	return nil, fmt.Errorf("Campaign for id not found: %d", id)
}**/

func AddDonation(campaignId uint, donation *model.Donation) (*model.Campaign, error) {
	donation.CampaignID = campaignId
	result := db.DB.Create(donation)
	if result.Error != nil {
		return nil, result.Error
	}
	entry := log.WithField("ID", campaignId)
	entry.Info("Successfully added new donation to campaign in database.")
	entry.Tracef("Stored: %v", donation)
	return GetCampaign(campaignId)
}