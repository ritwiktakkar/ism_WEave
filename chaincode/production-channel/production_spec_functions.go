package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SPEC_IsNewAsset ensures that the asset does not already exist
func SPEC_IsNewAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("the asset %s already exists", id)
	}
	return nil
}

// SPEC_AssetExists ensures that the asset exists
func SPEC_AssetExists(ctx contractapi.TransactionContextInterface, assetID string) error {
	assetJSON, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return fmt.Errorf("the asset %s does not exist", assetID)
	}
	return nil
}

// SPEC_IDPrefix ensures that the id starts with the correct prefix
func SPEC_IDPrefix(id string, prefix string) error {
	if !strings.HasPrefix(id, prefix) {
		return fmt.Errorf("the id '%s' must start with '%s'", id, prefix)
	}

	return nil
}

// SPEC_IsValidFlag ensures that (the flagReason is provided if isFlagged is true) and (the flagReason is empty or N/A if isFlagged is false)
func SPEC_IsValidFlag(isFlagged bool, flagReason string) error {
	if isFlagged && (len(flagReason) == 0 || flagReason == "N/A") {
		return fmt.Errorf("flagReason must be provided if isFlagged is true")
	} else if !isFlagged && (len(flagReason) != 0 && flagReason != "N/A") {
		return fmt.Errorf("flagReason must be empty or N/A if isFlagged is false")
	} else {
		return nil
	}
}

// SPEC_IsNotFlagged ensures that the asset is not flagged
func SPEC_IsNotFlagged(ctx contractapi.TransactionContextInterface, assetID string) error {
	assetJSON, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return fmt.Errorf("the asset %s does not exist", assetID)
	}

	var asset map[string]interface{}
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	isFlagged, ok := asset["IsFlagged"].(bool)
	if !ok {
		return fmt.Errorf("IsFlagged field missing or not a boolean for asset %s", assetID)
	}
	if isFlagged {
		return fmt.Errorf("the asset %s is flagged", assetID)
	}

	return nil
}

// SPEC_Chronoloy ensures that the provided dates are in chronological order
func SPEC_Chronology(dates ...time.Time) error {
	for i := 1; i < len(dates); i++ {
		if !dates[i].After(dates[i-1]) {
			return fmt.Errorf("date %v is not after date %v", dates[i], dates[i-1])
		}
	}
	return nil
}

// SPEC_IsInvokedByAllowedOrg checks if the function is invoked by one of the allowed OrgMSPIDs
func SPEC_IsInvokedByAllowedOrg(ctx contractapi.TransactionContextInterface, allowedOrgMSPIDs ...string) error {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSPID: %v", err)
	}

	for _, orgMSPID := range allowedOrgMSPIDs {
		if clientMSPID == orgMSPID {
			return nil
		}
	}
	return fmt.Errorf("the function is not invoked by an allowed organization. Invoked by: %s. Allowed organizations: %v", clientMSPID, allowedOrgMSPIDs)
}

// SPEC_LotConsistency ensures that the content list is not empty, that each asset in the content list exists in the ledger, and that each asset has the correct prefix
func SPEC_LotConsistency(ctx contractapi.TransactionContextInterface, content []string, assetIDPrefix string) error {
	// Check if content list is empty
	if len(content) == 0 {
		return fmt.Errorf("content list cannot be empty")
	}

	// Check if each asset has the correct prefix and exists in the ledger
	for _, assetID := range content {
		// Check if assetID has the correct prefix
		if !strings.HasPrefix(assetID, assetIDPrefix) {
			return fmt.Errorf("asset %s does not have the correct prefix %s", assetID, assetIDPrefix)
		}

		// Check if the asset exists in the ledger
		if err := SPEC_AssetExists(ctx, assetID); err != nil {
			return fmt.Errorf("asset %s does not exist or could not be accessed: %v", assetID, err)
		}
	}
	return nil
}

// SPEC_NoDuplicateAssetInLots ensures that each asset in the given lot is not stored in any other lot's contents, i.e., the state
func SPEC_NoDuplicateAssetInState(ctx contractapi.TransactionContextInterface, currentLotID string, content []string) error {
	// Get all lots from the ledger
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return fmt.Errorf("failed to get all lots from world state: %v", err)
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return fmt.Errorf("failed to iterate through lots: %v", err)
		}

		var checkLot Lot
		err = json.Unmarshal(queryResponse.Value, &checkLot)
		if err != nil {
			return fmt.Errorf("failed to unmarshal lot: %v", err)
		}

		// Skip the lot being checked
		if checkLot.ID == currentLotID {
			continue
		}

		// Check for duplicate assets in state's existing other lots
		for _, assetID := range content {
			for _, existingAssetID := range checkLot.Content {
				if assetID == existingAssetID {
					return fmt.Errorf("asset %s cannot be placed in lot %s because it is already stored in lot %s", assetID, currentLotID, checkLot.ID)
				}
			}
		}
	}

	return nil
}

func SPEC_CheckLotAssetType(ctx contractapi.TransactionContextInterface, lotID string, prefix string) error {
	// Retrieve the lot from the ledger
	lotJSON, err := ctx.GetStub().GetState(lotID)
	if err != nil {
		return fmt.Errorf("failed to read from world state for lot %s: %v", lotID, err)
	}
	if lotJSON == nil {
		return fmt.Errorf("lot %s does not exist", lotID)
	}

	// Unmarshal the lot JSON into a generic map
	var lotMap map[string]interface{}
	err = json.Unmarshal(lotJSON, &lotMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal lot JSON for lot %s: %v", lotID, err)
	}

	// Check the content field
	content, ok := lotMap["Content"].([]interface{})
	if !ok {
		return fmt.Errorf("content field missing or not an array for lot %s", lotID)
	}

	// Check the asset type
	for _, assetID := range content {
		if !strings.HasPrefix(assetID.(string), prefix) {
			return fmt.Errorf("asset %s in lot %s does not have the correct prefix %s", assetID, lotID, prefix)
		}
	}

	return nil
}

// SPEC_NoDuplicateAssetInThisLot ensures that there are no duplicate asset IDs in the content list
func SPEC_NoDuplicateAssetInThisLot(content []string) error {
	assetMap := make(map[string]bool)

	for _, assetID := range content {
		if _, exists := assetMap[assetID]; exists {
			return fmt.Errorf("duplicate asset ID found: %s", assetID)
		}
		assetMap[assetID] = true
	}

	return nil
}

// SPEC_CheckAssetsApproval checks if all assets in the content list have their approval field set to true
func SPEC_CheckAssetsApproval(ctx contractapi.TransactionContextInterface, content []string) error {
	for _, assetID := range content {
		// Retrieve the asset from the ledger
		assetJSON, err := ctx.GetStub().GetState(assetID)
		if err != nil {
			return fmt.Errorf("failed to read from world state for asset %s: %v", assetID, err)
		}
		if assetJSON == nil {
			return fmt.Errorf("asset %s does not exist", assetID)
		}

		// Unmarshal the asset JSON into a generic map
		var assetMap map[string]interface{}
		err = json.Unmarshal(assetJSON, &assetMap)
		if err != nil {
			return fmt.Errorf("failed to unmarshal asset JSON for asset %s: %v", assetID, err)
		}

		// Check the approval field
		approval, ok := assetMap["Approval"].(bool)
		if !ok {
			return fmt.Errorf("approval field missing or not a boolean for asset %s", assetID)
		}
		if !approval {
			return fmt.Errorf("asset %s has approval set to false", assetID)
		}
	}

	return nil
}

// SPEC_IsAllowedToOwn checks if the new owner is allowed to own the asset based on the allowedOrgMSPIDs
func SPEC_IsAllowedToOwn(ctx contractapi.TransactionContextInterface, newOwner string, allowedOrgMSPIDs ...string) error {

	for _, orgMSPID := range allowedOrgMSPIDs {
		if newOwner == orgMSPID {
			return nil
		}
	}
	return fmt.Errorf("the proposed owner, %s, is not allowed to own this lot.  Allowed organizations: %v", newOwner, allowedOrgMSPIDs)

}
