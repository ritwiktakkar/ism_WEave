package main

import "time"

// Asset: Order
type Order struct {
	CreatedAt       time.Time `json:"CreatedAt"`
	CreatorID       string    `json:"CreatorID"`
	DeliveryDate    time.Time `json:"DeliveryDate"`
	FlagReason      string    `json:"FlagReason"`
	ID              string    `json:"ID"`
	IsAccepted      bool      `json:"IsAccepted"`
	IsFlagged       bool      `json:"IsFlagged"`
	Notes           string    `json:"Notes"`
	PaymentTerms    string    `json:"PaymentTerms"`
	PlanID          string    `json:"PlanID"`
	ProductDetails  string    `json:"ProductDetails"`
	ReceiverID      string    `json:"ReceiverID"`
	Status          string    `json:"Status"`
	TotalOrderValue float32   `json:"TotalOrderValue"`
	UpdatedAt       time.Time `json:"UpdatedAt"`
}

// Asset: Plan
type Plan struct {
	AllFactoriesApproved bool      `json:"AllFactoriesApproved"`
	CreatedAt            time.Time `json:"CreatedAt"`
	CreatorID            string    `json:"CreatorID"`
	Factories            []string  `json:"Factories"`
	FlagReason           string    `json:"FlagReason"`
	ID                   string    `json:"ID"`
	IsAuditorApproved    bool      `json:"IsAuditorApproved"`
	IsFlagged            bool      `json:"IsFlagged"`
	IsRetailerApproved   bool      `json:"IsRetailerApproved"`
	Notes                string    `json:"Notes"`
	OrderID              string    `json:"OrderID"`
	ProductionPlan       string    `json:"ProductionPlan"`
	Status               string    `json:"Status"`
	UpdatedAt            time.Time `json:"UpdatedAt"`
}

// Asset: Factory
type Factory struct {
	FactoryOwner       string    `json:"FactoryOwner"`
	FlagReason         string    `json:"FlagReason"`
	ID                 string    `json:"ID"`
	IsAuditorApproved  bool      `json:"IsAuditorApproved"`
	IsFlagged          bool      `json:"IsFlagged"`
	IsRetailerApproved bool      `json:"IsRetailerApproved"`
	Location           string    `json:"Location"`
	Name               string    `json:"Name"`
	Notes              string    `json:"Notes"`
	PastFulfillment    bool      `json:"PastFulfillment"`
	StartDate          time.Time `json:"StartDate"`
	Status             string    `json:"Status"`
	UpdatedAt          time.Time `json:"UpdatedAt"`
}
