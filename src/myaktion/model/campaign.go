package model

type Campaign struct {
	ID              uint `json:"id"`
	Name            string `json:"name"`
	OrganizerName   string `json:"organizerName"`
	DonationMinimum float64 `json:"donationMinimum"`
	TargetAmount    float64 `json:"targetAmount"`
	Account         Account `json:"account"`
	Donations       []Donation `json:"donations"`
}
