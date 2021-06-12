package handler

import (
	"encoding/json"
	"github.com/MJ7898/myaktion-go/src/myaktion/model"
	"github.com/MJ7898/myaktion-go/src/myaktion/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CreateCampaign(w http.ResponseWriter, r *http.Request) {
	campaign, err := getCampaign(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Infof("campaign: %v", campaign)
	if err := service.CreateCampaign(campaign); err != nil {
		log.Errorf("Error calling service CreateCampaign: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, campaign)
}

func getCampaign(r *http.Request) (*model.Campaign, error) {
	var campaign model.Campaign
	// TODO: log http body as middleware
	err := json.NewDecoder(r.Body).Decode(&campaign)
	if err != nil {
		log.Errorf("Can't serialize request body to campaign struct: %v", err)
		return nil, err
	}
	return &campaign, nil
}

func GetCampaigns(w http.ResponseWriter, _ *http.Request) {
	campaigns, err := service.GetCampaigns()
	if err != nil {
		log.Errorf("Error calling service GetCampaigns: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(campaigns); err != nil {
		log.Errorf("Failure encoding value to JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	sendJson(w, campaigns)
}

// GetCampaign Getcampaign-Handler function to get an single campaign with id/**
func GetCampaign(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		log.Errorf("Error calling service GetSingleCampaign: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	campaign, err := service.GetCampaign(id)
	if campaign == nil {
		http.Error(w, "404 Campaign not found", http.StatusNotFound)
		return
	}
	sendJson(w, campaign)
}

func UpdateCampaign(w http.ResponseWriter, r *http.Request)  {
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	campaign, err := getCampaign(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	campaign, err = service.UpdateCampaign(id, campaign)
	if err != nil {
		log.Errorf("Failure updating campaign with ID %v: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if campaign == nil {
		http.Error(w, "404 campaign not found", http.StatusNotFound)
		return
	}
	sendJson(w, campaign)
}

func DeleteCampaign(w http.ResponseWriter, r *http.Request)  {
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	campaign, err := service.DeleteCampaign(id)

	if err != nil {
		log.Errorf("Failure updating campaign with ID %v: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if campaign == nil {
		http.Error(w, "404 campaign not found", http.StatusNotFound)
		return
	}
	sendJson(w, result{Success: "Success (Ok)"})
}
