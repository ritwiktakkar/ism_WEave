package main

import "time"

// Asset: CottonBale
type CottonBale struct {
	Approval     bool      `json:"Approval"`
	AssemblyDate time.Time `json:"AssemblyDate"`
	CreatorID    string    `json:"CreatorID"` // programmatically updated
	FlagReason   string    `json:"FlagReason"`
	ID           string    `json:"ID"`
	IsFlagged    bool      `json:"IsFlagged"`
	Notes        string    `json:"Notes"`
	Origin       string    `json:"Origin"`
	QualityGrade string    `json:"QualityGrade"`
	TotalWeight  float32   `json:"TotalWeight"` // inputted by the user
	UpdatedAt    time.Time `json:"UpdatedAt"`   // programmatically updated
}

// Asset: Lot
type Lot struct {
	AllAssetsApproved bool      `json:"AllAssetsApproved"` // programmatically updated
	AssemblyDate      time.Time `json:"AssemblyDate"`
	AssetIDPrefix     string    `json:"AssetIDPrefix"`
	Content           []string  `json:"Content"`       // IDs of asset with AssetIDPrefix
	ContentWeight     float32   `json:"ContentWeight"` // sum of the content weights programmatically updated
	CreatorID         string    `json:"CreatorID"`     // programmatically updated
	Destination       string    `json:"Destination"`
	FlagReason        string    `json:"FlagReason"`
	ID                string    `json:"ID"`
	IsFlagged         bool      `json:"IsFlagged"`
	Notes             string    `json:"Notes"`
	Origin            string    `json:"Origin"`
	Owner             string    `json:"Owner"`
	PreviousOwner     string    `json:"PreviousOwner"`
	Quantity          int       `json:"Quantity"`         // programmatically updated
	TotalWeight       float32   `json:"TotalWeight"`      // inputted by the user
	UpdatedAt         time.Time `json:"UpdatedAt"`        // programmatically updated
	WeightDifference  float32   `json:"WeightDifference"` // percentage difference between the total and content weight programmatically updated
}

// Asset: CottonYarn
type CottonYarn struct {
	Approval         bool      `json:"Approval"`
	AssemblyDate     time.Time `json:"AssemblyDate"`
	Content          []string  `json:"Content"`       // IDs of CottonBale Lots
	ContentWeight    float32   `json:"ContentWeight"` // sum of the content weights programmatically updated
	CreatorID        string    `json:"CreatorID"`     // programmatically updated
	FlagReason       string    `json:"FlagReason"`
	ID               string    `json:"ID"`
	IsFlagged        bool      `json:"IsFlagged"`
	Notes            string    `json:"Notes"`
	Origin           string    `json:"Origin"`
	TotalWeight      float32   `json:"TotalWeight"`      // inputted by the user
	UpdatedAt        time.Time `json:"UpdatedAt"`        // programmatically updated
	WeightDifference float32   `json:"WeightDifference"` // percentage difference between the total and content weight programmatically updated
	YarnCount        int       `json:"YarnCount"`
}

// Asset: UnfinishedFabric
type UnfinishedFabric struct {
	Approval         bool      `json:"Approval"`
	AssemblyDate     time.Time `json:"AssemblyDate"`
	Content          []string  `json:"Content"`       // IDs of CottonYarn Lots
	ContentWeight    float32   `json:"ContentWeight"` // sum of the content weights programmatically updated
	CreatorID        string    `json:"CreatorID"`     // programmatically updated
	FlagReason       string    `json:"FlagReason"`
	ID               string    `json:"ID"`
	IsFlagged        bool      `json:"IsFlagged"`
	Length           float32   `json:"Length"`
	Notes            string    `json:"Notes"`
	Origin           string    `json:"Origin"`
	TotalWeight      float32   `json:"TotalWeight"`      // inputted by the user
	UpdatedAt        time.Time `json:"UpdatedAt"`        // programmatically updated
	WeightDifference float32   `json:"WeightDifference"` // percentage difference between the total and content weight programmatically updated
	Width            float32   `json:"Width"`
}

// Asset: FinishedFabric
type FinishedFabric struct {
	Approval         bool      `json:"Approval"`
	AssemblyDate     time.Time `json:"AssemblyDate"`
	Content          []string  `json:"Content"`       // IDs of UnfinishedFabric Lots
	ContentWeight    float32   `json:"ContentWeight"` // sum of the content weights programmatically updated
	CreatorID        string    `json:"CreatorID"`     // programmatically updated
	FlagReason       string    `json:"FlagReason"`
	ID               string    `json:"ID"`
	IsFlagged        bool      `json:"IsFlagged"`
	Length           float32   `json:"Length"`
	Notes            string    `json:"Notes"`
	Origin           string    `json:"Origin"`
	TotalWeight      float32   `json:"TotalWeight"`      // inputted by the user
	UpdatedAt        time.Time `json:"UpdatedAt"`        // programmatically updated
	WeightDifference float32   `json:"WeightDifference"` // percentage difference between the total and content weight programmatically updated
	Width            float32   `json:"Width"`
}

// Asset: CutPart
type CutPart struct {
	Approval         bool      `json:"Approval"`
	AssemblyDate     time.Time `json:"AssemblyDate"`
	Content          []string  `json:"Content"`       // IDs of FinishedFabric Lots
	ContentWeight    float32   `json:"ContentWeight"` // sum of the content weights programmatically updated
	CreatorID        string    `json:"CreatorID"`     // programmatically updated
	FlagReason       string    `json:"FlagReason"`
	ID               string    `json:"ID"`
	IsFlagged        bool      `json:"IsFlagged"`
	Notes            string    `json:"Notes"`
	Origin           string    `json:"Origin"`
	PatternPiece     string    `json:"PatternPiece"`
	TotalWeight      float32   `json:"TotalWeight"`      // inputted by the user
	UpdatedAt        time.Time `json:"UpdatedAt"`        // programmatically updated
	WeightDifference float32   `json:"WeightDifference"` // percentage difference between the total and content weight programmatically updated
}

// Asset: Button
type Button struct {
	Approval     bool      `json:"Approval"`
	AssemblyDate time.Time `json:"AssemblyDate"`
	CreatorID    string    `json:"CreatorID"` // programmatically updated
	FlagReason   string    `json:"FlagReason"`
	ID           string    `json:"ID"`
	IsFlagged    bool      `json:"IsFlagged"`
	Notes        string    `json:"Notes"`
	Origin       string    `json:"Origin"`
	TotalWeight  float32   `json:"TotalWeight"`
	UpdatedAt    time.Time `json:"UpdatedAt"` // programmatically updated
}

// Asset: AssembledGarment
type AssembledGarment struct {
	Approval         bool      `json:"Approval"`
	AssemblyDate     time.Time `json:"AssemblyDate"`
	Buttons          []string  `json:"Buttons"`       // IDs of Buttons
	ContentWeight    float32   `json:"ContentWeight"` // sum of the content weights programmatically updated
	CreatorID        string    `json:"CreatorID"`     // programmatically updated
	CutParts         []string  `json:"CutParts"`      // IDs of CutParts
	FlagReason       string    `json:"FlagReason"`
	ID               string    `json:"ID"`
	IsFlagged        bool      `json:"IsFlagged"`
	Notes            string    `json:"Notes"`
	Origin           string    `json:"Origin"`
	TotalWeight      float32   `json:"TotalWeight"`      // inputted by the user
	UpdatedAt        time.Time `json:"UpdatedAt"`        // programmatically updated
	WeightDifference float32   `json:"WeightDifference"` // percentage difference between the total and content weight programmatically updated
}

// Asset: Carton
type Carton struct {
	AllAssetsApproved bool      `json:"AllAssetsApproved"`
	AssemblyDate      time.Time `json:"AssemblyDate"`
	Content           []string  `json:"Content"`       // IDs of AssembledGarments
	ContentWeight     float32   `json:"ContentWeight"` // sum of the content weights programmatically updated
	CreatorID         string    `json:"CreatorID"`     // programmatically updated
	CustomerID        string    `json:"CustomerID"`
	FlagReason        string    `json:"FlagReason"`
	ID                string    `json:"ID"`
	IsFlagged         bool      `json:"IsFlagged"`
	Notes             string    `json:"Notes"`
	Origin            string    `json:"Origin"`
	Owner             string    `json:"Owner"`
	PreviousOwner     string    `json:"PreviousOwner"`
	Quantity          int       `json:"Quantity"`
	TotalWeight       float32   `json:"TotalWeight"`      // inputted by the user
	UpdatedAt         time.Time `json:"UpdatedAt"`        // programmatically updated
	WeightDifference  float32   `json:"WeightDifference"` // percentage difference between the total and content weight programmatically updated
}

// Asset: Container
type Container struct {
	AssetType        string    `json:"AssetType"`
	Content          []string  `json:"Content"`       // IDs of Cartons
	ContentWeight    float32   `json:"ContentWeight"` // sum of the content weights programmatically updated
	CreatorID        string    `json:"CreatorID"`     // programmatically updated
	DestinationPort  string    `json:"DestinationPort"`
	FlagReason       string    `json:"FlagReason"`
	ID               string    `json:"ID"`
	IsFlagged        bool      `json:"IsFlagged"`
	LoadedAt         time.Time `json:"LoadedAt"`
	OriginPort       string    `json:"OriginPort"`
	TotalWeight      float32   `json:"TotalWeight"` // inputted by the user
	UpdatedAt        time.Time `json:"UpdatedAt"`   // programmatically updated
	Vessel           string    `json:"Vessel"`
	WeightDifference float32   `json:"WeightDifference"` // percentage difference between the total and content weight programmatically updated
}

// Asset: BillOfLading
type BillOfLading struct {
	Consignee     string    `json:"Consignee"`
	DeliveryPlace string    `json:"DeliveryPlace"`
	DischargePort string    `json:"DischargePort"`
	DocumentID    string    `json:"DocumentID"`
	FlagReason    string    `json:"FlagReason"`
	FreightTerms  string    `json:"FreightTerms"`
	GrossWeight   float32   `json:"GrossWeight"`
	ID            string    `json:"ID"`
	IsFlagged     bool      `json:"IsFlagged"`
	IssueDate     time.Time `json:"IssueDate"`
	LoadingPort   string    `json:"LoadingPort"`
	ReceiptPlace  string    `json:"ReceiptPlace"`
	SealNumber    uint8     `json:"SealNumber"`
	Shipper       string    `json:"Shipper"`
	URL           string    `json:"URL"`
	Vessel        string    `json:"Vessel"`
}
