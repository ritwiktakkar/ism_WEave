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

// SetFlag sets the isFlagged field of an asset with a specific ID given the flagReason
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
func (s *SmartContract) GetAllAssetsOfType(ctx contractapi.TransactionContextInterface, recordType string) (interface{}, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var orders []*Order
	var plans []*Plan
	var factories []*Factory

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// Check if the key contains the specified substring
		if !strings.Contains(queryResponse.Key, recordType) {
			continue
		}

		switch recordType {
		case "order_":
			var order Order
			err = json.Unmarshal(queryResponse.Value, &order)
			if err != nil {
				return nil, err
			}
			orders = append(orders, &order)

		case "plan_":
			var plan Plan
			err = json.Unmarshal(queryResponse.Value, &plan)
			if err != nil {
				return nil, err
			}
			plans = append(plans, &plan)

		case "factory_":
			var factory Factory
			err = json.Unmarshal(queryResponse.Value, &factory)
			if err != nil {
				return nil, err
			}
			factories = append(factories, &factory)

		// Add additional cases for other record types
		default:
			return nil, fmt.Errorf("invalid record type: %s", recordType)
		}
	}

	switch recordType {
	case "order_":
		return orders, nil
	case "plan_":
		return plans, nil
	case "factory_":
		return factories, nil
	default:
		return nil, fmt.Errorf("invalid record type: %s", recordType)
	}
}

// GetAllOrders retrieves all orders from the world state
func (s *SmartContract) GetAllOrders(ctx contractapi.TransactionContextInterface) ([]*Order, error) {
	records, err := s.GetAllAssetsOfType(ctx, "order_")
	if err != nil {
		return nil, err
	}
	return records.([]*Order), nil
}

// GetAllPlans retrieves all plans from the world state
func (s *SmartContract) GetAllPlans(ctx contractapi.TransactionContextInterface) ([]*Plan, error) {
	records, err := s.GetAllAssetsOfType(ctx, "plan_")
	if err != nil {
		return nil, err
	}
	return records.([]*Plan), nil
}

// GetAllFactories retrieves all factories from the world state
func (s *SmartContract) GetAllFactories(ctx contractapi.TransactionContextInterface) ([]*Factory, error) {
	records, err := s.GetAllAssetsOfType(ctx, "factory_")
	if err != nil {
		return nil, err
	}
	return records.([]*Factory), nil
}

// GetAllAssetsOfTypeCount retrieves the count of records from the world state that contain the specified substring in their keys
func (s *SmartContract) GetAllAssetsOfTypeCount(ctx contractapi.TransactionContextInterface, recordType string) (int, error) {
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
		if strings.Contains(queryResponse.Key, recordType) {
			count++
		}
	}

	return count, nil
}
