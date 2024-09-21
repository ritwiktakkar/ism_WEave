package main

import (
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

// SPEC_IsReadyforApproval checks if the asset status is ready for approval based on the provided conditions
func SPEC_IsReadyforApproval(conditions ...bool) error {
	for i, condition := range conditions {
		if !condition {
			return fmt.Errorf("the asset's status is not ready for approval: condition %d failed", i+1)
		}
	}
	return nil
}
