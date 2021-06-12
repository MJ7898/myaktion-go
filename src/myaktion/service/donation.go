package service

import (
	"github.com/MJ7898/myaktion-go/src/myaktion/client"
	"github.com/MJ7898/myaktion-go/src/myaktion/client/banktransfer"
	"github.com/MJ7898/myaktion-go/src/myaktion/db"
	"github.com/MJ7898/myaktion-go/src/myaktion/model"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
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
}

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
}**/


func AddDonation(campaignId uint, donation *model.Donation) error {
	campaign, err := GetCampaign(campaignId)
	if err != nil {
		return err
	}

	donation.CampaignID = campaignId
	result := db.DB.Create(donation)
	if result.Error != nil {
		return result.Error
	}

	conn, err := client.GetBankTransferConnection()
	if err != nil {
		log.Errorf("error connecting to the banktransfer service: %v", err)
		deleteDonation(donation)
		return err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	banktransferClient := banktransfer.NewBankTransferClient(conn)
	_, err = banktransferClient.TransferMoney(ctx, &banktransfer.Transaction{
		DonationId: int32(donation.ID),
		Amount: float32(donation.Amount), Reference: "Donation",
		FromAccount: convertAccount(&donation.Account), ToAccount: convertAccount(&campaign.Account),
	})
	if err != nil {
		log.Errorf("error calling the banktransfer service: %v", err)
		deleteDonation(donation)
		return err
	}

entry := log.WithField("ID", campaignId)
entry.Info("Successfully added new donation to campaign in database.")
entry.Tracef("Stored: %v", donation)
return nil
}

func convertAccount(account *model.Account) *banktransfer.Account {
	return &banktransfer.Account{
		Name: account.Name, BankName: account.BankName, Number: account.Number,
	}
}
func deleteDonation(donation *model.Donation) error {
	entry := log.WithField("donationID", donation.ID)
	entry.Info("Trying to delete donation to make state consistent.")
	result := db.DB.Delete(donation)
	if result.Error != nil {
		// Note: configure logger to raise an alarm to compensate inconsistent state
		entry.WithField("alarm", true).Error("")
		return result.Error
	}
	entry.Info("Successfully deleted campaign.")
	return nil
}

func MarkDonation(id uint) error {
	entry := log.WithField("id", id)
	donation := new(model.Donation)
	result := db.DB.First(donation, id)
	if result.Error != nil {
		entry.WithError(result.Error).Error("Can't retrieve donation")
		return result.Error
	}
	entry = entry.WithField("donation", donation)
	entry.Trace("Retrieved donation")
	donation.Status = model.TRANSFERRED
	result = db.DB.Save(donation)
	if result.Error != nil {
		entry.WithError(result.Error).Error("Can't update donation")
		return result.Error
	}
	entry.Info("Successfully update status of donation.")
	return nil
}