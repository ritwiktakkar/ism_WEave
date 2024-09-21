package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// SetFlag sets the isFlagged field of an asset with a specific ID given the flagReason. Contains the following checks: 1) SPEC_IsValidFlag
func (s *SmartContract) SetFlag(ctx contractapi.TransactionContextInterface, id string, isFlagged bool, flagReason string) error {
	// Retrieve the asset from the world state using the provided ID
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return fmt.Errorf("the asset %s does not exist", id)
	}
	var asset map[string]interface{}
	// Unmarshal the JSON into a map
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	// Update the isFlagged and flagReason fields
	asset["IsFlagged"] = isFlagged
	asset["FlagReason"] = flagReason
	// Ensure that the flagReason is provided if isFlagged is true
	if err := SPEC_IsValidFlag(isFlagged, flagReason); err != nil {
		return err
	}
	// Update the updatedAt field
	asset["UpdatedAt"] = time.Now()
	// Marshal the updated asset back to JSON
	updatedAssetJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	// Save the updated asset to the world state
	return ctx.GetStub().PutState(id, updatedAssetJSON)
}

// SetNotes sets the notes field of an asset with a specific ID
func (s *SmartContract) SetNotes(ctx contractapi.TransactionContextInterface, id string, notes string) error {
	// Retrieve the asset from the world state using the provided ID
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return fmt.Errorf("the asset %s does not exist", id)
	}
	var asset map[string]interface{}
	// Unmarshal the JSON into a map
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	// Update the notes field
	asset["Notes"] = notes
	// Update the updatedAt field
	asset["UpdatedAt"] = time.Now()
	// Marshal the updated asset back to JSON
	updatedAssetJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	// Save the updated asset to the world state
	return ctx.GetStub().PutState(id, updatedAssetJSON)
}

// GetAsset retrieves an asset with a specific ID from the world state
func (s *SmartContract) GetAsset(ctx contractapi.TransactionContextInterface, id string) (map[string]interface{}, error) {
	// Retrieve the asset from the world state using the provided ID
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}
	var asset map[string]interface{}
	// Unmarshal the JSON into a map
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return asset, nil
}

// GetAllAssets retrieves all records from the world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]map[string]interface{}, error) {
	// Define a composite key prefix that includes the document type
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var records []map[string]interface{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var record map[string]interface{}
		err = json.Unmarshal(queryResponse.Value, &record)
		if err != nil {
			return nil, err
		}
		// Add the record to the result list
		records = append(records, record)
	}
	return records, nil
}

// GetPercentageDifference calculates the percentage difference between two float32 values
func GetPercentageDifference(original, new float32) float32 {
	if original == 0 {
		return 0
	}
	percentageDifference := ((new - original) / original) * 100
	return percentageDifference
}

// GetAllAssetsCount retrieves the amount of assets in the world state
func (s *SmartContract) GetAllAssetsCount(ctx contractapi.TransactionContextInterface) (int, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return 0, err
	}
	defer resultsIterator.Close()

	count := 0

	for resultsIterator.HasNext() {
		_, err := resultsIterator.Next()
		if err != nil {
			return 0, err
		}
		count++
	}

	return count, nil
}

// GetAllAssetsOfType is a helper function not meant for direct client invocation as it retrieves all records from the world state that contain the specified substring in their keys but returns the records as an interface
func (s *SmartContract) GetAllAssetsOfType(ctx contractapi.TransactionContextInterface, assetIDPrefix string) (interface{}, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var cottonBales []*CottonBale
	var lots []*Lot
	var cottonYarns []*CottonYarn
	var unfinishedFabrics []*UnfinishedFabric
	var finishedFabrics []*FinishedFabric
	var cutParts []*CutPart
	var buttons []*Button
	var assembledGarments []*AssembledGarment
	var cartons []*Carton
	var containers []*Container

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// Check if the key contains the specified substring
		if !strings.Contains(queryResponse.Key, assetIDPrefix) {
			continue
		}

		switch assetIDPrefix {

		case "cottonBale_":
			var cottonBale CottonBale
			err = json.Unmarshal(queryResponse.Value, &cottonBale)
			if err != nil {
				return nil, err
			}
			cottonBales = append(cottonBales, &cottonBale)

		case "lot_":
			var lot Lot
			err = json.Unmarshal(queryResponse.Value, &lot)
			if err != nil {
				return nil, err
			}
			lots = append(lots, &lot)

		case "cottonYarn_":
			var cottonYarn CottonYarn
			err = json.Unmarshal(queryResponse.Value, &cottonYarn)
			if err != nil {
				return nil, err
			}
			cottonYarns = append(cottonYarns, &cottonYarn)

		case "unfinishedFabric_":
			var unfinishedFabric UnfinishedFabric
			err = json.Unmarshal(queryResponse.Value, &unfinishedFabric)
			if err != nil {
				return nil, err
			}
			unfinishedFabrics = append(unfinishedFabrics, &unfinishedFabric)

		case "finishedFabric_":
			var finishedFabric FinishedFabric
			err = json.Unmarshal(queryResponse.Value, &finishedFabric)
			if err != nil {
				return nil, err
			}
			finishedFabrics = append(finishedFabrics, &finishedFabric)

		case "cutPart_":
			var cutPart CutPart
			err = json.Unmarshal(queryResponse.Value, &cutPart)
			if err != nil {
				return nil, err
			}
			cutParts = append(cutParts, &cutPart)

		case "button_":
			var button Button
			err = json.Unmarshal(queryResponse.Value, &button)
			if err != nil {
				return nil, err
			}
			buttons = append(buttons, &button)

		case "assembledGarment_":
			var assembledGarment AssembledGarment
			err = json.Unmarshal(queryResponse.Value, &assembledGarment)
			if err != nil {
				return nil, err
			}
			assembledGarments = append(assembledGarments, &assembledGarment)

		case "carton_":
			var carton Carton
			err = json.Unmarshal(queryResponse.Value, &carton)
			if err != nil {
				return nil, err
			}
			cartons = append(cartons, &carton)

		case "container_":
			var container Container
			err = json.Unmarshal(queryResponse.Value, &container)
			if err != nil {
				return nil, err
			}
			containers = append(containers, &container)

		// Add additional cases for other record types
		default:
			return nil, fmt.Errorf("invalid record type: %s", assetIDPrefix)
		}
	}

	switch assetIDPrefix {
	case "cottonbale_":
		return cottonBales, nil
	case "lot_":
		return lots, nil
	case "cottonyarn_":
		return cottonYarns, nil
	case "unfinishedfabric_":
		return unfinishedFabrics, nil
	case "finishedfabric_":
		return finishedFabrics, nil
	case "cutpart_":
		return cutParts, nil
	case "button_":
		return buttons, nil
	case "assembledgarment_":
		return assembledGarments, nil
	case "carton_":
		return cartons, nil
	case "container_":
		return containers, nil
	default:
		return nil, fmt.Errorf("invalid record type: %s", assetIDPrefix)
	}
}

// GetAllCottonBales retrieves all cotton bales from the world state
func (s *SmartContract) GetAllCottonBales(ctx contractapi.TransactionContextInterface) ([]*CottonBale, error) {
	records, err := s.GetAllAssetsOfType(ctx, "cottonbale_")
	if err != nil {
		return nil, err
	}
	return records.([]*CottonBale), nil
}

// GetAllLots retrieves all lots from the world state
func (s *SmartContract) GetAllLots(ctx contractapi.TransactionContextInterface) ([]*Lot, error) {
	records, err := s.GetAllAssetsOfType(ctx, "lot_")
	if err != nil {
		return nil, err
	}
	return records.([]*Lot), nil
}

// GetAllCottonYarns retrieves all cotton yarns from the world state
func (s *SmartContract) GetAllCottonYarns(ctx contractapi.TransactionContextInterface) ([]*CottonYarn, error) {
	records, err := s.GetAllAssetsOfType(ctx, "cottonyarn_")
	if err != nil {
		return nil, err
	}
	return records.([]*CottonYarn), nil
}

// GetAllUnfinishedFabrics retrieves all unfinished fabrics from the world state
func (s *SmartContract) GetAllUnfinishedFabrics(ctx contractapi.TransactionContextInterface) ([]*UnfinishedFabric, error) {
	records, err := s.GetAllAssetsOfType(ctx, "unfinishedfabric_")
	if err != nil {
		return nil, err
	}
	return records.([]*UnfinishedFabric), nil
}

// GetAllFinishedFabrics retrieves all finished fabrics from the world state
func (s *SmartContract) GetAllFinishedFabrics(ctx contractapi.TransactionContextInterface) ([]*FinishedFabric, error) {
	records, err := s.GetAllAssetsOfType(ctx, "finishedfabric_")
	if err != nil {
		return nil, err
	}
	return records.([]*FinishedFabric), nil
}

// GetAllCutParts retrieves all cut parts from the world state
func (s *SmartContract) GetAllCutParts(ctx contractapi.TransactionContextInterface) ([]*CutPart, error) {
	records, err := s.GetAllAssetsOfType(ctx, "cutpart_")
	if err != nil {
		return nil, err
	}
	return records.([]*CutPart), nil
}

// GetAllButtons retrieves all buttons from the world state
func (s *SmartContract) GetAllButtons(ctx contractapi.TransactionContextInterface) ([]*Button, error) {
	records, err := s.GetAllAssetsOfType(ctx, "button_")
	if err != nil {
		return nil, err
	}
	return records.([]*Button), nil
}

// GetAllAssembledGarments retrieves all assembled garments from the world state
func (s *SmartContract) GetAllAssembledGarments(ctx contractapi.TransactionContextInterface) ([]*AssembledGarment, error) {
	records, err := s.GetAllAssetsOfType(ctx, "assembledgarment_")
	if err != nil {
		return nil, err
	}
	return records.([]*AssembledGarment), nil
}

// GetAllCartons retrieves all cartons from the world state
func (s *SmartContract) GetAllCartons(ctx contractapi.TransactionContextInterface) ([]*Carton, error) {
	records, err := s.GetAllAssetsOfType(ctx, "carton_")
	if err != nil {
		return nil, err
	}
	return records.([]*Carton), nil
}

// GetAllContainers retrieves all containers from the world state
func (s *SmartContract) GetAllContainers(ctx contractapi.TransactionContextInterface) ([]*Container, error) {
	records, err := s.GetAllAssetsOfType(ctx, "container_")
	if err != nil {
		return nil, err
	}
	return records.([]*Container), nil
}

// GetAllAssetsOfTypeCount retrieves the count of records from the world state that contain the specified substring in their keys
func (s *SmartContract) GetAllAssetsOfTypeCount(ctx contractapi.TransactionContextInterface, assetIDPrefix string) (int, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return 0, err
	}
	defer resultsIterator.Close()

	count := 0

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return 0, err
		}

		// Check if the key contains the specified substring
		if strings.Contains(queryResponse.Key, assetIDPrefix) {
			count++
		}
	}

	return count, nil
}

// GetContentWeight retrieves the TotalWeight for each asset ID in content and sums them up
func GetContentWeight(ctx contractapi.TransactionContextInterface, content []string) (float32, error) {
	var contentWeight float32

	for _, assetID := range content {
		// Retrieve the asset from the ledger
		assetJSON, err := ctx.GetStub().GetState(assetID)
		if err != nil {
			return 0, fmt.Errorf("failed to read from world state for asset %s: %v", assetID, err)
		}
		if assetJSON == nil {
			return 0, fmt.Errorf("asset %s does not exist", assetID)
		}

		// Unmarshal the asset JSON into a generic map
		var assetMap map[string]interface{}
		err = json.Unmarshal(assetJSON, &assetMap)
		if err != nil {
			return 0, fmt.Errorf("failed to unmarshal asset JSON for asset %s: %v", assetID, err)
		}

		// Extract the TotalWeight attribute
		totalWeight, ok := assetMap["TotalWeight"].(float64)
		if !ok {
			return 0, fmt.Errorf("asset %s does not have a valid TotalWeight attribute", assetID)
		}

		// Add the asset's TotalWeight to the total contentWeight
		contentWeight += float32(totalWeight)
	}

	return contentWeight, nil
}

// GetLotsWithPrefixCount returns the total number of assets with IDs prefixed by "lot_" further filtered by AssetIDPrefix
func GetLotsWithPrefixCount(ctx contractapi.TransactionContextInterface, assetIDPrefix string) (int, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return 0, fmt.Errorf("failed to get state by range: %v", err)
	}
	defer resultsIterator.Close()

	count := 0

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return 0, fmt.Errorf("failed to iterate through results: %v", err)
		}

		if strings.HasPrefix(queryResponse.Key, "lot_") {
			var lot Lot
			err = json.Unmarshal(queryResponse.Value, &lot)
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal lot: %v", err)
			}

			if lot.AssetIDPrefix == assetIDPrefix {
				count++
			}
		}
	}

	return count, nil
}
