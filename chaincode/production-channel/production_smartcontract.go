package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// CreateCottonBale issues a new asset (CottonBale) to the state with select attributes. Contains the following specifications: 1) SPEC_IDPrefix, 2) SPEC_IsNewAsset, 3) SPEC_IsInvokedByAllowedOrg, 4) SPEC_IsValidFlag, 5) SPEC_Chronology
func (s *SmartContract) CreateCottonBale(ctx contractapi.TransactionContextInterface, approval bool, assemblyDate time.Time, flagReason string, cottonBaleID string, isFlagged bool, notes string, origin string, qualityGrade string, totalWeight float32) error {
	// Ensure the id begins with "cottonbale_"
	if err := SPEC_IDPrefix(cottonBaleID, "cottonbale_"); err != nil {
		return err
	}
	// Ensure the cottonbale does not already exist
	if err := SPEC_IsNewAsset(ctx, cottonBaleID); err != nil {
		return err
	}
	// Ensure that the function is invoked by an allowed organization
	if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org4MSP"); err != nil {
		return err
	}
	// Ensure that the flagReason is provided if isFlagged is true
	if err := SPEC_IsValidFlag(isFlagged, flagReason); err != nil {
		return err
	}

	// Retrieve the invoking organization's MSP ID
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSPID: %v", err)
	}

	cottonBale := CottonBale{
		Approval:     approval,
		AssemblyDate: assemblyDate,
		CreatorID:    clientMSPID,
		FlagReason:   flagReason,
		ID:           cottonBaleID,
		IsFlagged:    isFlagged,
		Notes:        notes,
		Origin:       origin,
		QualityGrade: qualityGrade,
		TotalWeight:  totalWeight,
		UpdatedAt:    time.Now(),
	}

	// Ensure that the dates are in chronological order
	if err := SPEC_Chronology(cottonBale.AssemblyDate, cottonBale.UpdatedAt); err != nil {
		return err
	}

	// Convert cottonBale to JSON
	cottonBaleJSON, err := json.Marshal(cottonBale)
	if err != nil {
		return err
	}

	// Save the cottonBale to the world state
	return ctx.GetStub().PutState(cottonBaleID, cottonBaleJSON)
}

// CreateLot issues a new asset (Lot) to the state with select attributes. Contains the following specifications: 1) SPEC_IDPrefix, 2) SPEC_IsNewAsset, 3) SPEC_IsInvokedByAllowedOrg, 4) SPEC_IsValidFlag, 5) SPEC_IsNotFlagged, 6) SPEC_LotConsistency, 7) SPEC_NoDuplicateAssetInThisLot, 8) SPEC_NoDuplicateAssetInState, 9) SPEC_Chronology
func (s *SmartContract) CreateLot(ctx contractapi.TransactionContextInterface, assemblyDate time.Time, assetIDPrefix string, content []string, destination string, flagReason string, lotID string, isFlagged bool, notes string, origin string, owner string, totalWeight float32) error {
	// Ensure the id begins with "lot_"
	if err := SPEC_IDPrefix(lotID, "lot_"); err != nil {
		return err
	}
	// Ensure the lot does not already exist
	if err := SPEC_IsNewAsset(ctx, lotID); err != nil {
		return err
	}

	// Switch on the assetIDPrefix to determine which organizations are allowed to invoke the function
	switch assetIDPrefix {
	case "cottonbale_":
		if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP", "Org4MSP"); err != nil {
			return err
		}

	case "cottonyarn_":
		if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP", "Org4MSP"); err != nil {
			return err
		}

	case "unfinishedfabric_":
		if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP", "Org5MSP"); err != nil {
			return err
		}

	case "finishedfabric_":
		if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP", "Org5MSP"); err != nil {
			return err
		}

	case "cutpart_":
		if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP", "Org6MSP"); err != nil {
			return err
		}

	case "button_":
		if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP", "Org6MSP"); err != nil {
			return err
		}

	case "assembledgarment_":
		if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP", "Org6MSP"); err != nil {
			return err
		}
	}

	// Ensure that the flagReason is provided if isFlagged is true
	if err := SPEC_IsValidFlag(isFlagged, flagReason); err != nil {
		return err
	}
	// Ensure that not a single asset in the content list is flagged
	for _, assetID := range content {
		if err := SPEC_IsNotFlagged(ctx, assetID); err != nil {
			return err
		}
	}
	// Ensure that each asset in the content list is unique
	if err := SPEC_NoDuplicateAssetInThisLot(content); err != nil {
		return err
	}
	// Ensure that each asset in the content list exists, and that the content list is consistent with the asset type
	if err := SPEC_LotConsistency(ctx, content, assetIDPrefix); err != nil {
		return err
	}
	// Ensure that each asset in this lot is not already part of another lot
	if err := SPEC_NoDuplicateAssetInState(ctx, lotID, content); err != nil {
		return err
	}
	// Ensure that all assets in this lot are approved
	if err := SPEC_CheckAssetsApproval(ctx, content); err != nil {
		return err
	}

	// Calculate the content weight by summing asset weights in content
	contentWeight, err := GetContentWeight(ctx, content)
	if err != nil {
		return err
	}
	// Calculate percentage difference between total weight and content weight
	var weightDifference float32 = GetPercentageDifference(contentWeight, totalWeight)

	// Retrieve the invoking organization's MSP ID
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSPID: %v", err)
	}

	lot := Lot{
		AllAssetsApproved: true,
		AssemblyDate:      assemblyDate,
		AssetIDPrefix:     assetIDPrefix,
		Content:           content,
		ContentWeight:     contentWeight,
		CreatorID:         clientMSPID,
		Destination:       destination,
		FlagReason:        flagReason,
		ID:                lotID,
		IsFlagged:         isFlagged,
		Notes:             notes,
		Origin:            origin,
		Owner:             clientMSPID,
		PreviousOwner:     "Updated when ownership changes",
		Quantity:          len(content),
		TotalWeight:       totalWeight,
		UpdatedAt:         time.Now(),
		WeightDifference:  weightDifference,
	}

	// Ensure that the dates are in chronological order
	if err := SPEC_Chronology(lot.AssemblyDate, lot.UpdatedAt); err != nil {
		return err
	}

	// Convert lot to JSON
	lotJSON, err := json.Marshal(lot)
	if err != nil {
		return err
	}

	// Save the lot to the world state
	return ctx.GetStub().PutState(lotID, lotJSON)
}

// CreateCottonYarn issues a new asset (CottonYarn) to the state with select attributes. Contains the following specifications: 1) SPEC_IDPrefix, 2) SPEC_IsNewAsset, 3) SPEC_IsInvokedByAllowedOrg, 4) SPEC_IsValidFlag, 5) SPEC_NoDuplicateAssetInThisLot, 6) SPEC_CheckLotAssetType, 7) SPEC_Chronology
func (s *SmartContract) CreateCottonYarn(ctx contractapi.TransactionContextInterface, approval bool, assemblyDate time.Time, content []string, flagReason string, cottonYarnID string, isFlagged bool, notes string, origin string, totalWeight float32, yarnCount int) error {
	// Ensure the id begins with "cottonyarn_"
	if err := SPEC_IDPrefix(cottonYarnID, "cottonyarn_"); err != nil {
		return err
	}
	// Ensure the cottonyarn does not already exist
	if err := SPEC_IsNewAsset(ctx, cottonYarnID); err != nil {
		return err
	}
	// Ensure that the function is invoked by an allowed organization
	if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org4MSP"); err != nil {
		return err
	}
	// Ensure that the flagReason is provided if isFlagged is true
	if err := SPEC_IsValidFlag(isFlagged, flagReason); err != nil {
		return err
	}
	// Ensure that each asset in the content list is unique
	if err := SPEC_NoDuplicateAssetInThisLot(content); err != nil {
		return err
	}
	// Ensure that each raw material in content lot(s) is of type cottonbale
	for _, lotID := range content {
		if err := SPEC_CheckLotAssetType(ctx, lotID, "cottonbale_"); err != nil {
			return err
		}
	}

	// Calculate the content weight by summing asset weights in content
	contentWeight, err := GetContentWeight(ctx, content)
	if err != nil {
		return err
	}
	// Calculate percentage difference between total weight and content weight
	var weightDifference float32 = GetPercentageDifference(contentWeight, totalWeight)

	// Retrieve the invoking organization's MSP ID
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSPID: %v", err)
	}

	cottonYarn := CottonYarn{
		Approval:         approval,
		AssemblyDate:     assemblyDate,
		Content:          content,
		ContentWeight:    contentWeight,
		CreatorID:        clientMSPID,
		FlagReason:       flagReason,
		ID:               cottonYarnID,
		IsFlagged:        isFlagged,
		Notes:            notes,
		Origin:           origin,
		TotalWeight:      totalWeight,
		UpdatedAt:        time.Now(),
		WeightDifference: weightDifference,
		YarnCount:        yarnCount,
	}

	// Ensure that the dates are in chronological order
	if err := SPEC_Chronology(cottonYarn.AssemblyDate, cottonYarn.UpdatedAt); err != nil {
		return err
	}

	// Convert cottonYarn to JSON
	cottonYarnJSON, err := json.Marshal(cottonYarn)
	if err != nil {
		return err
	}

	// Save the cottonYarn to the world state
	return ctx.GetStub().PutState(cottonYarnID, cottonYarnJSON)
}

// CreateUnfinishedFabric issues a new asset (UnfinishedFabric) to the state with select attributes. Contains the following specifications: 1) SPEC_IDPrefix, 2) SPEC_IsNewAsset, 3) SPEC_IsInvokedByAllowedOrg, 4) SPEC_IsValidFlag, 5) SPEC_NoDuplicateAssetInThisLot, 6) SPEC_CheckLotAssetType, 7) SPEC_Chronology
func (s *SmartContract) CreateUnfinishedFabric(ctx contractapi.TransactionContextInterface, approval bool, assemblyDate time.Time, content []string, flagReason string, unfinishedFabricID string, isFlagged bool, notes string, origin string, length float32, totalWeight float32, width float32) error {
	// Ensure the id begins with "unfinishedfabric_"
	if err := SPEC_IDPrefix(unfinishedFabricID, "unfinishedfabric_"); err != nil {
		return err
	}
	// Ensure the unfinishedFabric does not already exist
	if err := SPEC_IsNewAsset(ctx, unfinishedFabricID); err != nil {
		return err
	}
	// Ensure that the function is invoked by an allowed organization
	if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org5MSP"); err != nil {
		return err
	}
	// Ensure that the flagReason is provided if isFlagged is true
	if err := SPEC_IsValidFlag(isFlagged, flagReason); err != nil {
		return err
	}
	// Ensure that each asset in the content list is unique
	if err := SPEC_NoDuplicateAssetInThisLot(content); err != nil {
		return err
	}
	// Ensure that each raw material in content lot(s) is of type cottonyarn
	for _, lotID := range content {
		if err := SPEC_CheckLotAssetType(ctx, lotID, "cottonyarn_"); err != nil {
			return err
		}
	}

	// Calculate the content weight by summing asset weights in content
	contentWeight, err := GetContentWeight(ctx, content)
	if err != nil {
		return err
	}
	// Calculate percentage difference between total weight and content weight
	var weightDifference float32 = GetPercentageDifference(contentWeight, totalWeight)

	// Retrieve the invoking organization's MSP ID
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSPID: %v", err)
	}

	unfinishedFabric := UnfinishedFabric{
		Approval:         approval,
		AssemblyDate:     assemblyDate,
		Content:          content,
		ContentWeight:    contentWeight,
		CreatorID:        clientMSPID,
		FlagReason:       flagReason,
		ID:               unfinishedFabricID,
		IsFlagged:        isFlagged,
		Length:           length,
		Notes:            notes,
		Origin:           origin,
		TotalWeight:      totalWeight,
		UpdatedAt:        time.Now(),
		WeightDifference: weightDifference,
		Width:            width,
	}

	// Ensure that the dates are in chronological order
	if err := SPEC_Chronology(unfinishedFabric.AssemblyDate, unfinishedFabric.UpdatedAt); err != nil {
		return err
	}

	// Convert unfinishedFabric to JSON
	unfinishedFabricJSON, err := json.Marshal(unfinishedFabric)
	if err != nil {
		return err
	}

	// Save the unfinishedFabric to the world state
	return ctx.GetStub().PutState(unfinishedFabricID, unfinishedFabricJSON)
}

// CreateFinishedFabric issues a new asset (FinishedFabric) to the state with select attributes. Contains the following specifications: 1) SPEC_IDPrefix, 2) SPEC_IsNewAsset, 3) SPEC_IsInvokedByAllowedOrg, 4) SPEC_IsValidFlag, 5) SPEC_NoDuplicateAssetInThisLot, 6) SPEC_CheckLotAssetType, 7) SPEC_Chronology
func (s *SmartContract) CreateFinishedFabric(ctx contractapi.TransactionContextInterface, approval bool, assemblyDate time.Time, content []string, finishedFabricID string, flagReason string, isFlagged bool, length float32, notes string, origin string, totalWeight float32, width float32) error {
	// Ensure the id begins with "finishedfabric_"
	if err := SPEC_IDPrefix(finishedFabricID, "finishedfabric_"); err != nil {
		return err
	}
	// Ensure the finishedFabric does not already exist
	if err := SPEC_IsNewAsset(ctx, finishedFabricID); err != nil {
		return err
	}
	// Ensure that the function is invoked by an allowed organization
	if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org5MSP"); err != nil {
		return err
	}
	// Ensure that the flagReason is provided if isFlagged is true
	if err := SPEC_IsValidFlag(isFlagged, flagReason); err != nil {
		return err
	}
	// Ensure that each asset in the content list is unique
	if err := SPEC_NoDuplicateAssetInThisLot(content); err != nil {
		return err
	}
	// Ensure that each raw material in content lot(s) is of type unfinishedfabric
	for _, lotID := range content {
		if err := SPEC_CheckLotAssetType(ctx, lotID, "unfinishedfabric_"); err != nil {
			return err
		}
	}

	// Calculate the content weight by summing asset weights in content
	contentWeight, err := GetContentWeight(ctx, content)
	if err != nil {
		return err
	}
	// Calculate percentage difference between total weight and content weight
	var weightDifference float32 = GetPercentageDifference(contentWeight, totalWeight)

	// Retrieve the invoking organization's MSP ID
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSPID: %v", err)
	}

	finishedFabric := FinishedFabric{
		Approval:         approval,
		AssemblyDate:     assemblyDate,
		Content:          content,
		ContentWeight:    contentWeight,
		CreatorID:        clientMSPID,
		FlagReason:       flagReason,
		ID:               finishedFabricID,
		IsFlagged:        isFlagged,
		Length:           length,
		Notes:            notes,
		Origin:           origin,
		TotalWeight:      totalWeight,
		UpdatedAt:        time.Now(),
		WeightDifference: weightDifference,
		Width:            width,
	}

	// Ensure that the dates are in chronological order
	if err := SPEC_Chronology(finishedFabric.AssemblyDate, finishedFabric.UpdatedAt); err != nil {
		return err
	}

	// Convert finishedFabric to JSON
	finishedFabricJSON, err := json.Marshal(finishedFabric)
	if err != nil {
		return err
	}

	// Save the finishedFabric to the world state
	return ctx.GetStub().PutState(finishedFabricID, finishedFabricJSON)
}

// CreateCutPart issues a new asset (CutPart) to the state with select attributes. Contains the following specifications: 1) SPEC_IDPrefix, 2) SPEC_IsNewAsset, 3) SPEC_IsInvokedByAllowedOrg, 4) SPEC_IsValidFlag, 5) SPEC_NoDuplicateAssetInThisLot, 6) SPEC_CheckLotAssetType, 7) SPEC_Chronology
func (s *SmartContract) CreateCutPart(ctx contractapi.TransactionContextInterface, approval bool, assemblyDate time.Time, content []string, flagReason string, cutPartID string, isFlagged bool, notes string, origin string, patternPiece string, totalWeight float32) error {
	// Ensure the id begins with "cutpart_"
	if err := SPEC_IDPrefix(cutPartID, "cutpart_"); err != nil {
		return err
	}
	// Ensure the cutPart does not already exist
	if err := SPEC_IsNewAsset(ctx, cutPartID); err != nil {
		return err
	}
	// Ensure that the function is invoked by an allowed organization
	if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org6MSP"); err != nil {
		return err
	}
	// Ensure that the flagReason is provided if isFlagged is true
	if err := SPEC_IsValidFlag(isFlagged, flagReason); err != nil {
		return err
	}
	// Ensure that each asset in the content list is unique
	if err := SPEC_NoDuplicateAssetInThisLot(content); err != nil {
		return err
	}
	// Ensure that each raw material in content lot(s) is of type finishedfabric
	for _, lotID := range content {
		if err := SPEC_CheckLotAssetType(ctx, lotID, "finishedfabric_"); err != nil {
			return err
		}
	}
	// Calculate the content weight by summing asset weights in content
	contentWeight, err := GetContentWeight(ctx, content)
	if err != nil {
		return err
	}
	// Calculate percentage difference between total weight and content weight
	var weightDifference float32 = GetPercentageDifference(contentWeight, totalWeight)

	// Retrieve the invoking organization's MSP ID
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSPID: %v", err)
	}

	cutPart := CutPart{
		Approval:         approval,
		AssemblyDate:     assemblyDate,
		Content:          content,
		ContentWeight:    contentWeight,
		CreatorID:        clientMSPID,
		FlagReason:       flagReason,
		ID:               cutPartID,
		IsFlagged:        isFlagged,
		Notes:            notes,
		Origin:           origin,
		PatternPiece:     patternPiece,
		TotalWeight:      totalWeight,
		UpdatedAt:        time.Now(),
		WeightDifference: weightDifference,
	}

	// Ensure that the dates are in chronological order
	if err := SPEC_Chronology(cutPart.AssemblyDate, cutPart.UpdatedAt); err != nil {
		return err
	}

	// Convert cutPart to JSON
	cutPartsJSON, err := json.Marshal(cutPart)
	if err != nil {
		return err
	}

	// Save the cutPart to the world state
	return ctx.GetStub().PutState(cutPartID, cutPartsJSON)
}

// CreateButton issues a new asset (Button) to the state with select attributes. Contains the following specifications: 1) SPEC_IDPrefix, 2) SPEC_IsNewAsset, 3) SPEC_IsInvokedByAllowedOrg, 4) SPEC_IsValidFlag, 5) SPEC_Chronology
func (s *SmartContract) CreateButton(ctx contractapi.TransactionContextInterface, approval bool, assemblyDate time.Time, flagReason string, buttonID string, isFlagged bool, notes string, origin string, totalWeight float32) error {
	// Ensure the id begins with "button_"
	if err := SPEC_IDPrefix(buttonID, "button_"); err != nil {
		return err
	}
	// Ensure the button does not already exist
	if err := SPEC_IsNewAsset(ctx, buttonID); err != nil {
		return err
	}
	// Ensure that the function is invoked by an allowed organization
	if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org6MSP"); err != nil {
		return err
	}
	// Ensure that the flagReason is provided if isFlagged is true
	if err := SPEC_IsValidFlag(isFlagged, flagReason); err != nil {
		return err
	}

	// Retrieve the invoking organization's MSP ID
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSPID: %v", err)
	}

	button := Button{
		Approval:     approval,
		AssemblyDate: assemblyDate,
		CreatorID:    clientMSPID,
		FlagReason:   flagReason,
		ID:           buttonID,
		IsFlagged:    isFlagged,
		Notes:        notes,
		Origin:       origin,
		TotalWeight:  totalWeight,
		UpdatedAt:    time.Now(),
	}

	// Ensure that the dates are in chronological order
	if err := SPEC_Chronology(button.AssemblyDate, button.UpdatedAt); err != nil {
		return err
	}

	// Convert button to JSON
	buttonJSON, err := json.Marshal(button)
	if err != nil {
		return err
	}

	// Save the button to the world state
	return ctx.GetStub().PutState(buttonID, buttonJSON)
}

// CreateAssembledGarment issues a new asset (AssembledGarment) to the state with select attributes. Contains the following specifications: 1) SPEC_IDPrefix, 2) SPEC_IsNewAsset, 3) SPEC_IsInvokedByAllowedOrg, 4) SPEC_IsValidFlag, 5) SPEC_NoDuplicateAssetInThisLot, 6) SPEC_LotConsistency, 7) SPEC_LotConsistency, 8) SPEC_Chronology
func (s *SmartContract) CreateAssembledGarment(ctx contractapi.TransactionContextInterface, approval bool, assemblyDate time.Time, buttons []string, cutParts []string, flagReason string, assembledGarmentID string, isFlagged bool, notes string, origin string, totalWeight float32) error {
	// Ensure the id begins with "assembledgarment_"
	if err := SPEC_IDPrefix(assembledGarmentID, "assembledgarment_"); err != nil {
		return err
	}
	// Ensure the assembledGarment does not already exist
	if err := SPEC_IsNewAsset(ctx, assembledGarmentID); err != nil {
		return err
	}
	// Ensure that the function is invoked by an allowed organization
	if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org6MSP"); err != nil {
		return err
	}
	// Ensure that the flagReason is provided if isFlagged is true
	if err := SPEC_IsValidFlag(isFlagged, flagReason); err != nil {
		return err
	}
	// Ensure that each asset in the content list is unique
	if err := SPEC_NoDuplicateAssetInThisLot(buttons); err != nil {
		return err
	}
	// Ensure that each asset in the content list is unique
	if err := SPEC_NoDuplicateAssetInThisLot(cutParts); err != nil {
		return err
	}
	// Ensure that each asset in the content list exists, and that the content list is consistent with the asset type
	if err := SPEC_LotConsistency(ctx, buttons, "button_"); err != nil {
		return err
	}
	if err := SPEC_LotConsistency(ctx, cutParts, "cutpart_"); err != nil {
		return err
	}

	// Calculate the content weight by summing asset weights in content
	buttonsWeight, err := GetContentWeight(ctx, buttons)
	if err != nil {
		return err
	}
	cutPartsWeight, err := GetContentWeight(ctx, cutParts)
	if err != nil {
		return err
	}
	contentWeight := buttonsWeight + cutPartsWeight
	// Calculate percentage difference between total weight and content weight
	var weightDifference float32 = GetPercentageDifference(contentWeight, totalWeight)

	// Retrieve the invoking organization's MSP ID
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSPID: %v", err)
	}

	assembledGarment := AssembledGarment{
		Approval:         approval,
		AssemblyDate:     assemblyDate,
		Buttons:          buttons,
		ContentWeight:    contentWeight,
		CreatorID:        clientMSPID,
		CutParts:         cutParts,
		FlagReason:       flagReason,
		ID:               assembledGarmentID,
		IsFlagged:        isFlagged,
		Notes:            notes,
		Origin:           origin,
		TotalWeight:      totalWeight,
		UpdatedAt:        time.Now(),
		WeightDifference: weightDifference,
	}

	// Ensure that the dates are in chronological order
	if err := SPEC_Chronology(assembledGarment.AssemblyDate, assembledGarment.UpdatedAt); err != nil {
		return err
	}

	// Convert assembledGarment to JSON
	assembledGarmentJSON, err := json.Marshal(assembledGarment)
	if err != nil {
		return err
	}

	// Save the assembledGarment to the world state
	return ctx.GetStub().PutState(assembledGarmentID, assembledGarmentJSON)
}

// CreateCarton issues a new asset (Carton) to the state with select attributes. Contains the following specifications: 1) SPEC_IDPrefix, 2) SPEC_IsNewAsset, 3) SPEC_IsInvokedByAllowedOrg, 4) SPEC_IsValidFlag, 5) SPEC_NoDuplicateAssetInThisLot, 6) SPEC_LotConsistency, 7) SPEC_Chronology
func (s *SmartContract) CreateCarton(ctx contractapi.TransactionContextInterface, allAssetsApproved bool, assemblyDate time.Time, content []string, customerID string, flagReason string, cartonID string, isFlagged bool, notes string, origin string, owner string, totalWeight float32) error {
	// Ensure the id begins with "carton_"
	if err := SPEC_IDPrefix(cartonID, "carton_"); err != nil {
		return err
	}
	// Ensure the carton does not already exist
	if err := SPEC_IsNewAsset(ctx, cartonID); err != nil {
		return err
	}
	// Ensure that the function is invoked by an allowed organization
	if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org6MSP"); err != nil {
		return err
	}
	// Ensure that the flagReason is provided if isFlagged is true
	if err := SPEC_IsValidFlag(isFlagged, flagReason); err != nil {
		return err
	}
	// Ensure that each asset in the content list is unique
	if err := SPEC_NoDuplicateAssetInThisLot(content); err != nil {
		return err
	}
	if err := SPEC_LotConsistency(ctx, content, "assembledgarment_"); err != nil {
		return err
	}

	// Calculate the content weight by summing asset weights in content
	contentWeight, err := GetContentWeight(ctx, content)
	if err != nil {
		return err
	}
	// Calculate percentage difference between total weight and content weight
	var weightDifference float32 = GetPercentageDifference(contentWeight, totalWeight)

	//

	// Retrieve the invoking organization's MSP ID
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSPID: %v", err)
	}

	carton := Carton{
		AllAssetsApproved: allAssetsApproved,
		AssemblyDate:      assemblyDate,
		Content:           content,
		ContentWeight:     contentWeight,
		CreatorID:         clientMSPID,
		CustomerID:        customerID,
		FlagReason:        flagReason,
		ID:                cartonID,
		IsFlagged:         isFlagged,
		Notes:             notes,
		Origin:            origin,
		Owner:             owner,
		PreviousOwner:     "Updated when ownership changes",
		Quantity:          len(content),
		TotalWeight:       totalWeight,
		UpdatedAt:         time.Now(),
		WeightDifference:  weightDifference,
	}

	// Ensure that the dates are in chronological order
	if err := SPEC_Chronology(carton.AssemblyDate, carton.UpdatedAt); err != nil {
		return err
	}

	// Convert carton to JSON
	cartonJSON, err := json.Marshal(carton)
	if err != nil {
		return err
	}

	// Save the carton to the world state
	return ctx.GetStub().PutState(cartonID, cartonJSON)
}

// CreateContainer issues a new asset (Container) to the state with select attributes. Contains the following specifications: 1) SPEC_IDPrefix, 2) SPEC_IsNewAsset, 3) SPEC_IsInvokedByAllowedOrg, 4) SPEC_IsValidFlag, 5) SPEC_oDuplicateAssetInThisLot, 6) SPEC_LotConsistency, 7) SPEC_Chronology
func (s *SmartContract) CreateContainer(ctx contractapi.TransactionContextInterface, content []string, destinationPort string, flagReason string, containerID string, isFlagged bool, loadedAt time.Time, originPort string, totalWeight float32, vessel string) error {
	// Ensure the id begins with "container_"
	if err := SPEC_IDPrefix(containerID, "container_"); err != nil {
		return err
	}
	// Ensure the container does not already exist
	if err := SPEC_IsNewAsset(ctx, containerID); err != nil {
		return err
	}
	// Ensure that the function is invoked by an allowed organization
	if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org6MSP"); err != nil {
		return err
	}
	// Ensure that the flagReason is provided if isFlagged is true
	if err := SPEC_IsValidFlag(isFlagged, flagReason); err != nil {
		return err
	}
	// Ensure that each asset in the content list is unique
	if err := SPEC_NoDuplicateAssetInThisLot(content); err != nil {
		return err
	}
	// Ensure that each asset in the content list exists, and that the content list is consistent with the asset type
	if err := SPEC_LotConsistency(ctx, content, "carton_"); err != nil {
		return err
	}

	// Calculate the content weight by summing asset weights in content
	contentWeight, err := GetContentWeight(ctx, content)
	if err != nil {
		return err
	}
	// Calculate percentage difference between total weight and content weight
	var weightDifference float32 = GetPercentageDifference(contentWeight, totalWeight)

	// Retrieve the invoking organization's MSP ID
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSPID: %v", err)
	}

	container := Container{
		AssetType:        "carton",
		Content:          content,
		ContentWeight:    contentWeight,
		CreatorID:        clientMSPID,
		DestinationPort:  destinationPort,
		FlagReason:       flagReason,
		ID:               containerID,
		IsFlagged:        isFlagged,
		LoadedAt:         loadedAt,
		OriginPort:       originPort,
		TotalWeight:      totalWeight,
		UpdatedAt:        time.Now(),
		Vessel:           vessel,
		WeightDifference: weightDifference,
	}

	// Ensure that the dates are in chronological order
	if err := SPEC_Chronology(container.LoadedAt, container.UpdatedAt); err != nil {
		return err
	}

	// Convert container to JSON
	containerJSON, err := json.Marshal(container)
	if err != nil {
		return err
	}

	// Save the container to the world state
	return ctx.GetStub().PutState(containerID, containerJSON)
}

// UpdateLotOwner updates the owner field of an asset in the world state. Contains the following specifications: 1) SPEC_IDPrefix, 2) SPEC_IsAssetOwner, 3) SPEC_IsInvokedByAllowedOrg, 4) SPEC_IsAllowedToOwn, 5) SPEC_Chronology
func (s *SmartContract) UpdateLotOwner(ctx contractapi.TransactionContextInterface, lotID string, newOwner string) error {

	// Retrieve the asset from the world state using the provided ID
	assetJSON, err := ctx.GetStub().GetState(lotID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return fmt.Errorf("the asset %s does not exist", lotID)
	}
	var asset map[string]interface{}
	// Unmarshal the JSON into a map
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Switch on the assetIDPrefix to determine whether: 1) the function is invoked by an allowed organization, and 2) the new owner is allowed to own the lot based on the underlying asset
	switch asset["AssetIDPrefix"] {
	case "cottonbale_":
		if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP", "Org4MSP"); err != nil {
			return err
		}
		if err := SPEC_IsAllowedToOwn(ctx, newOwner, "Org1MSP", "Org2MSP", "Org4MSP"); err != nil {
			return err
		}

	case "cottonyarn_":
		if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP", "Org4MSP", "Org5MSP"); err != nil {
			return err
		}
		if err := SPEC_IsAllowedToOwn(ctx, newOwner, "Org1MSP", "Org2MSP", "Org4MSP", "Org5MSP"); err != nil {
			return err
		}

	case "unfinishedfabric_":
		if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP", "Org5MSP"); err != nil {
			return err
		}
		if err := SPEC_IsAllowedToOwn(ctx, newOwner, "Org1MSP", "Org2MSP", "Org5MSP"); err != nil {
			return err
		}

	case "finishedfabric_":
		if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP", "Org5MSP", "Org6MSP"); err != nil {
			return err
		}
		if err := SPEC_IsAllowedToOwn(ctx, newOwner, "Org1MSP", "Org2MSP", "Org5MSP", "Org6MSP"); err != nil {
			return err
		}

	case "cutpart_":
		if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP", "Org6MSP"); err != nil {
			return err
		}
		if err := SPEC_IsAllowedToOwn(ctx, newOwner, "Org1MSP", "Org2MSP", "Org6MSP"); err != nil {
			return err
		}

	case "button_":
		if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP", "Org6MSP"); err != nil {
			return err
		}
		if err := SPEC_IsAllowedToOwn(ctx, newOwner, "Org1MSP", "Org2MSP", "Org6MSP"); err != nil {
			return err
		}

	case "assembledgarment_":
		if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP", "Org6MSP"); err != nil {
			return err
		}
		if err := SPEC_IsAllowedToOwn(ctx, newOwner, "Org1MSP", "Org2MSP", "Org6MSP"); err != nil {
			return err
		}
	}

	asset["PreviousOwner"] = asset["Owner"]
	asset["Owner"] = newOwner
	asset["UpdatedAt"] = time.Now()

	// Marshal the updated asset back to JSON
	updatedAssetJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	// Save the updated asset to the world state
	return ctx.GetStub().PutState(lotID, updatedAssetJSON)
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating production-channel chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting production-channel chaincode: %v", err)
	}
}
