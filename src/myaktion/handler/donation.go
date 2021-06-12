package handler

import (
	"encoding/json"
	"github.com/MJ7898/myaktion-go/src/myaktion/model"
	"github.com/MJ7898/myaktion-go/src/myaktion/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func AddDonation(w http.ResponseWriter, r *http.Request)  {
	campaignId, err := getId(r)
	if err != nil {
		log.Errorf("Error getting ID: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	donation, err := getDonation(r)
	if err != nil {
		log. Errorf("Can't serialize body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//_, err = ...
	//campaign, err := service.AddDonation(campaignId, donation)
	// TODO: if the campaign doesn't exist, return 404 - don't show FK error
	err = service.AddDonation(campaignId, donation)
	if err != nil {
		log.Errorf("Failure adding donation to campaign with ID %v: %v", campaignId, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, donation) // donation instead
}

func getDonation(r *http.Request) (*model.Donation, error) {
	var donation model.Donation
	err := json.NewDecoder(r.Body).Decode(&donation)
	if err != nil {
		log.Errorf("Can't serialize request body to donation struct: %v", err)
		return nil, err
	}
	return &donation, nil
}