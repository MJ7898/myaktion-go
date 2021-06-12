package service

import (
	"errors"
	"github.com/MJ7898/myaktion-go/src/myaktion/db"
	"github.com/MJ7898/myaktion-go/src/myaktion/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateCampaign(campaign *model.Campaign) error {
	// campaign.ID = actCampaignId -> Art und Weise für in Memory
	// campaignStore[actCampaignId] = campaign
	// actCampaignId += 1 -> abgelöst durch DB-Anbindung über gorm
	result := db.DB.Create(campaign)
	if result.Error != nil {
		return result.Error
	}
	log.Infoln("Successfully stored new campaign with ID %v in database.", campaign.ID)
	log.Tracef("Stored: %v", campaign)
	return nil
}

func GetCampaigns() ([]model.Campaign, error) {
	var campaigns []model.Campaign
	/**for _, campaign := range campaignStore { In-Memory Lösung und abfrage der vorhandenen campaigns
		campaigns = append(campaigns, *campaign) }**/ // -> abgelöst durch die DB-Lösung
	result := db.DB.Preload("Donations").Find(&campaigns)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Tracef("Retrieved: %v", campaigns)
	return campaigns, nil
}

func GetCampaign(id uint) (*model.Campaign, error) {
	campaign := new(model.Campaign)
	result := db.DB.Preload("Donations").First(campaign, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	log.Tracef("Retrieved: %v", campaign)
	return campaign, nil
}

func UpdateCampaign(id uint, campaign *model.Campaign) (*model.Campaign, error)  {
	existingCampaign, err := GetCampaign(id)
	if  existingCampaign != nil|| err != nil {
		return existingCampaign, err
	}
	existingCampaign.Name = campaign.Name
	existingCampaign.OrganizerName = campaign.OrganizerName
	existingCampaign.TargetAmount = campaign.TargetAmount
	existingCampaign.DonationMinimum = campaign.DonationMinimum
	existingCampaign.Account = campaign.Account
	result := db.DB.Save(existingCampaign)
	if result.Error != nil {
		return nil, result.Error
	}
	entry := log.WithField("ID", id)
	entry.Info("Successfully updated campaign.")
	entry.Tracef("Updated: %v", existingCampaign)
	return existingCampaign, nil
}

// UpdateCampaignNew newer function from the UpdateCampaign
func UpdateCampaignNew(id uint, campaign *model.Campaign) (*model.Campaign, error) {
	//var campaign model.Campaign
	existingCampaign, _ := GetCampaign(id)
	result := db.DB.Model(&existingCampaign).Updates(campaign)
	if result.Error != nil {
		return nil, result.Error
	}
	entry := log.WithField("ID", id)
	entry.Info("Successfully updated campaign.")
	entry.Tracef("Updated: %v", &campaign)
	return existingCampaign, nil

}

func DeleteCampaign(id uint) (*model.Campaign, error)  {
	campaign, err := GetCampaign(id)
	// campaignResult := db.DB.Preload("Campaign").Find(&campaignById)
	// campaign:= campaignStore[id] -> IN-Memory-Lösung
	if err == nil {
		return campaign, nil
		// log.Tracef("404 Campaign not found")
		// return nil, fmt.Errorf("no campaign with ID %d", id)
	}
	deleteResult := db.DB.Delete(&campaign)
	if deleteResult != nil {
		return nil, deleteResult.Error
	}
	// delete(campaignStore,id) -> In-Memory-Lösung
	entry := log.WithField("ID", id)
	entry.Info("Successfully deleted campaign.")
	log.Printf("Successfully deleted campaign with ID %d", id)
	return  campaign, nil
}