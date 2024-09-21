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

// CreateOrder issues a new asset (order) to the state with select attributes. Contains the following 4 specifications: 1) SPEC_IDPrefix, 2) SPEC_IsNewAsset, 3) SPEC_IsValidFlag, 4) SPEC_Chronology
func (s *SmartContract) CreateOrder(ctx contractapi.TransactionContextInterface, createdAt time.Time, deliveryDate time.Time, flagReason string, orderID string, isFlagged bool, notes string, paymentTerms string, productDetails string, receiverID string, totalOrderValue float32) error {
	// Ensure the id begins with "order_"
	if err := SPEC_IDPrefix(orderID, "order_"); err != nil {
		return err
	}
	// Ensure the order does not already exist
	if err := SPEC_IsNewAsset(ctx, orderID); err != nil {
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
	order := Order{
		CreatedAt:       createdAt,
		CreatorID:       clientMSPID,
		DeliveryDate:    deliveryDate,
		FlagReason:      flagReason,
		ID:              orderID,
		IsAccepted:      false,
		IsFlagged:       isFlagged,
		Notes:           notes,
		PaymentTerms:    paymentTerms,
		PlanID:          "",
		ProductDetails:  productDetails,
		ReceiverID:      receiverID,
		Status:          "issued",
		TotalOrderValue: totalOrderValue,
		UpdatedAt:       time.Now(),
	}
	// Ensure that the dates are in chronological order
	if err := SPEC_Chronology(order.CreatedAt, order.UpdatedAt, order.DeliveryDate); err != nil {
		return err
	}
	// Convert order to JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}
	// Save the order to the world state
	return ctx.GetStub().PutState(orderID, orderJSON)
}

// CreatePlan issues a new asset (plan) to the state with select attributes. Contains the following 5 specifications: 1) SPEC_IDPrefix, 2) SPEC_IsNewAsset, 3) SPEC_AssetExists, 4) SPEC_IsValidFlag, 5) SPEC_Chronology
func (s *SmartContract) CreatePlan(ctx contractapi.TransactionContextInterface, createdAt time.Time, factoryIDs []string, flagReason string, planID string, isFlagged bool, notes string, orderID string, productionPlan string) error {
	// Ensure the id begins with "plan_"
	if err := SPEC_IDPrefix(planID, "plan_"); err != nil {
		return err
	}
	// Ensure the plan does not already exist
	if err := SPEC_IsNewAsset(ctx, planID); err != nil {
		return err
	}
	// Ensure the order exists
	if err := SPEC_AssetExists(ctx, orderID); err != nil {
		return err
	}
	// Ensure that the flagReason is provided if isFlagged is true
	if err := SPEC_IsValidFlag(isFlagged, flagReason); err != nil {
		return err
	}
	// Ensure that all factories exist
	for _, factoryID := range factoryIDs {
		if err := SPEC_AssetExists(ctx, factoryID); err != nil {
			return err
		}
	}
	// Retrieve the invoking organization's MSP ID
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSPID: %v", err)
	}
	plan := Plan{
		AllFactoriesApproved: false,
		CreatedAt:            createdAt,
		CreatorID:            clientMSPID,
		Factories:            factoryIDs,
		FlagReason:           flagReason,
		ID:                   planID,
		IsAuditorApproved:    false,
		IsFlagged:            isFlagged,
		IsRetailerApproved:   false,
		Notes:                notes,
		OrderID:              orderID,
		ProductionPlan:       productionPlan,
		Status:               "issued",
		UpdatedAt:            time.Now(),
	}
	// Ensure that the dates are in chronological order
	if err := SPEC_Chronology(plan.CreatedAt, plan.UpdatedAt); err != nil {
		return err
	}
	// Convert plan to JSON
	planJSON, err := json.Marshal(plan)
	if err != nil {
		return err
	}
	// Save the plan to the world state
	return ctx.GetStub().PutState(planID, planJSON)
}

// CreateFactory issues a new asset (factory) to the world state with select attributes. Contains the following 4 specifications: 1) SPEC_IDPrefix, 2) SPEC_IsNewAsset, 3) SPEC_IsValidFlag, 4) SPEC_Chronology
func (s *SmartContract) CreateFactory(ctx contractapi.TransactionContextInterface, factoryOwner string, flagReason string, factoryID string, isFlagged bool, location string, name string, notes string, pastFulfillment bool, startDate time.Time) error {
	// Ensure the id begins with "factory_"
	if err := SPEC_IDPrefix(factoryID, "factory_"); err != nil {
		return err
	}
	// Ensure the factory does not already exist
	if err := SPEC_IsNewAsset(ctx, factoryID); err != nil {
		return err
	}
	// Ensure that the flagReason is provided if isFlagged is true
	if err := SPEC_IsValidFlag(isFlagged, flagReason); err != nil {
		return err
	}
	factory := Factory{
		FactoryOwner:       factoryOwner,
		FlagReason:         flagReason,
		ID:                 factoryID,
		IsAuditorApproved:  false,
		IsRetailerApproved: false,
		IsFlagged:          isFlagged,
		Location:           location,
		Name:               name,
		Notes:              notes,
		PastFulfillment:    pastFulfillment,
		StartDate:          startDate,
		Status:             "pending",
		UpdatedAt:          time.Now(),
	}
	// Ensure that the dates are in chronological order
	if err := SPEC_Chronology(factory.StartDate, factory.UpdatedAt); err != nil {
		return err
	}
	// Convert factory to JSON
	factoryJSON, err := json.Marshal(factory)
	if err != nil {
		return err
	}
	// Save the factory to the world state
	return ctx.GetStub().PutState(factoryID, factoryJSON)
}

// SetOrderAcceptance updates the IsAccepted field of an Order. Contains the following 3 specifications: 1) SPEC_IsInvokedByAllowedOrg, 2) SPEC_AssetExists, 3) SPEC_Chronology
func (s *SmartContract) SetOrderAcceptance(ctx contractapi.TransactionContextInterface, orderID string, planID string, acceptance bool) error {
	// Retrieve the order from the world state
	orderJSON, err := ctx.GetStub().GetState(orderID)
	if err != nil {
		return fmt.Errorf("failed to read order: %v", err)
	}
	if orderJSON == nil {
		return fmt.Errorf("order %s does not exist", orderID)
	}
	var order Order
	err = json.Unmarshal(orderJSON, &order)
	if err != nil {
		return fmt.Errorf("failed to unmarshal order: %v", err)
	}
	// Ensure that the function is invoked by an allowed organization
	if err := SPEC_IsInvokedByAllowedOrg(ctx, order.ReceiverID); err != nil {
		return err
	}
	// Ensure the plan exists
	if err := SPEC_AssetExists(ctx, planID); err != nil {
		return err
	}
	// Update the Plan field
	order.PlanID = planID
	// Update the IsAccepted field
	order.IsAccepted = acceptance
	// Update the updatedAt field
	order.UpdatedAt = time.Now()
	// Ensure that the dates are in chronological order
	if err := SPEC_Chronology(order.CreatedAt, order.UpdatedAt, order.DeliveryDate); err != nil {
		return err
	}
	// Marshal the updated order and put it back in the world state
	updatedOrderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal updated order: %v", err)
	}
	err = ctx.GetStub().PutState(orderID, updatedOrderJSON)
	if err != nil {
		return fmt.Errorf("failed to update order: %v", err)
	}

	return nil
}

// SetOrderStatus updates the Status field of an Order based on specific conditions. Contains the following 2 specifications: 1) SPEC_IsInvokedByAllowedOrg, 2) SPEC_IsReadyforApproval
func (s *SmartContract) SetOrderStatus(ctx contractapi.TransactionContextInterface, orderID string, newStatus string) error {
	// Ensure that the function is invoked by an allowed organization
	if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org2MSP"); err != nil {
		return err
	}
	// Retrieve the order from the world state
	orderJSON, err := ctx.GetStub().GetState(orderID)
	if err != nil {
		return fmt.Errorf("failed to read order: %v", err)
	}
	if orderJSON == nil {
		return fmt.Errorf("order %s does not exist", orderID)
	}
	var order Order
	err = json.Unmarshal(orderJSON, &order)
	if err != nil {
		return fmt.Errorf("failed to unmarshal order: %v", err)
	}
	// Retrieve the plan from the world state
	planJSON, err := ctx.GetStub().GetState(order.PlanID)
	if err != nil {
		return fmt.Errorf("failed to read plan: %v", err)
	}
	if planJSON == nil {
		return fmt.Errorf("plan %s does not exist", order.PlanID)
	}
	var plan Plan
	err = json.Unmarshal(planJSON, &plan)
	if err != nil {
		return fmt.Errorf("failed to unmarshal plan: %v", err)
	}
	// Ensure all conditions are met for approval of plan
	if err := SPEC_IsReadyforApproval(order.IsAccepted, !order.IsFlagged, (plan.Status == "approved")); err != nil {
		return err
	}
	// Ensure that the new status is either "accepted" or "cancelled"
	if newStatus != "accepted" && newStatus != "cancelled" && newStatus != "rejected" {
		return fmt.Errorf("new status must be 'accepted', 'cancelled', or 'rejected', got: '%s' instead", newStatus)
	}
	// Update the order status and updatedAt
	order.Status = newStatus
	order.UpdatedAt = time.Now()
	// Marshal the updated order and put it back in the world state
	updatedOrderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal updated order: %v", err)
	}
	err = ctx.GetStub().PutState(orderID, updatedOrderJSON)
	if err != nil {
		return fmt.Errorf("failed to update order: %v", err)
	}

	return nil
}

// SetPlanApproval updates the relevant fields of a plan based on whether the retailer or auditor invokes it. Contains the following 1 specification: 1) SPEC_IsInvokedByAllowedOrg
func (s *SmartContract) SetPlanApproval(ctx contractapi.TransactionContextInterface, planID string, approval bool) error {
	// Ensure that the function is invoked by an allowed organization
	if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org3MSP"); err != nil {
		return err
	}
	// Retrieve the plan from the world state
	planJSON, err := ctx.GetStub().GetState(planID)
	if err != nil {
		return fmt.Errorf("failed to read plan: %v", err)
	}
	if planJSON == nil {
		return fmt.Errorf("plan %s does not exist", planID)
	}
	var plan Plan
	err = json.Unmarshal(planJSON, &plan)
	if err != nil {
		return fmt.Errorf("failed to unmarshal plan: %v", err)
	}
	// Retrieve the invoking organization's MSP ID
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSPID: %v", err)
	}
	// Check who is invoking and whether to update the approval status
	if clientMSPID == "Org1MSP" && approval {
		plan.IsRetailerApproved = true
	}
	if clientMSPID == "Org3MSP" && approval {
		plan.IsAuditorApproved = true
	}
	// Update the updatedAt field
	plan.UpdatedAt = time.Now()
	// Marshal the updated order and put it back in the world state
	updatedPlanJSON, err := json.Marshal(plan)
	if err != nil {
		return fmt.Errorf("failed to marshal updated plan: %v", err)
	}
	err = ctx.GetStub().PutState(planID, updatedPlanJSON)
	if err != nil {
		return fmt.Errorf("failed to update plan: %v", err)
	}
	return nil
}

// SetPlanStatus updates the Status field of a Plan based on whether it is approved by the retailer and auditor, and also if all factories are approved. Contains the following 1 specification: 1) SPEC_IsReadyforApproval
func (s *SmartContract) SetPlanStatus(ctx contractapi.TransactionContextInterface, planID string) error {
	// Retrieve the plan from the world state
	planJSON, err := ctx.GetStub().GetState(planID)
	if err != nil {
		return fmt.Errorf("failed to read plan: %v", err)
	}
	if planJSON == nil {
		return fmt.Errorf("plan %s does not exist", planID)
	}
	var plan Plan
	err = json.Unmarshal(planJSON, &plan)
	if err != nil {
		return fmt.Errorf("failed to unmarshal plan: %v", err)
	}
	// Check whether all factories are approved
	allFactoriesApproved := true
	for _, factoryID := range plan.Factories {
		factoryJSON, err := ctx.GetStub().GetState(factoryID)
		if err != nil {
			return fmt.Errorf("failed to read factory %s: %v", factoryID, err)
		}
		if factoryJSON == nil {
			return fmt.Errorf("factory %s does not exist", factoryID)
		}
		var factory Factory
		err = json.Unmarshal(factoryJSON, &factory)
		if err != nil {
			return fmt.Errorf("failed to unmarshal factory %s: %v", factoryID, err)
		}
		if factory.Status != "approved" {
			allFactoriesApproved = false
			break
		}
	}
	plan.AllFactoriesApproved = allFactoriesApproved
	// Ensure all conditions are met for approval of plan
	if err := SPEC_IsReadyforApproval(plan.AllFactoriesApproved, plan.IsAuditorApproved, plan.IsRetailerApproved, !plan.IsFlagged); err != nil {
		return err
	}
	// Set the status to "approved" if all conditions are met
	plan.Status = "approved"
	// Update the updatedAt field
	plan.UpdatedAt = time.Now()
	// Marshal the updated plan and put it back in the world state
	updatedPlanJSON, err := json.Marshal(plan)
	if err != nil {
		return fmt.Errorf("failed to marshal updated plan: %v", err)
	}
	err = ctx.GetStub().PutState(planID, updatedPlanJSON)
	if err != nil {
		return fmt.Errorf("failed to update plan: %v", err)
	}
	return nil
}

// SetFactoryApproval updates the IsRetailerApproved and IsAuditorApproved fields of a Factory based on specific conditions. Contains the following 1 specification: 1) SPEC_IsInvokedByAllowedOrg. Only works for Org1MSP and Org3MSP.
func (s *SmartContract) SetFactoryApproval(ctx contractapi.TransactionContextInterface, factoryID string, approval bool) error {
	// Ensure that the function is invoked by an allowed organization
	if err := SPEC_IsInvokedByAllowedOrg(ctx, "Org1MSP", "Org3MSP"); err != nil {
		return err
	}
	// Retrieve the factory from the world state
	factoryJSON, err := ctx.GetStub().GetState(factoryID)
	if err != nil {
		return fmt.Errorf("failed to read factory: %v", err)
	}
	if factoryJSON == nil {
		return fmt.Errorf("factory %s does not exist", factoryID)
	}
	var factory Factory
	err = json.Unmarshal(factoryJSON, &factory)
	if err != nil {
		return fmt.Errorf("failed to unmarshal factory: %v", err)
	}
	// Retrieve the invoking organization's MSP ID
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSPID: %v", err)
	}
	// Check who is invoking and whether to update the approval status
	if clientMSPID == "Org1MSP" && approval {
		factory.IsRetailerApproved = true
	}
	if clientMSPID == "Org3MSP" && approval {
		factory.IsAuditorApproved = true
	}
	// Update the updatedAt field
	factory.UpdatedAt = time.Now()
	// Marshal the updated factory and put it back in the world state
	updatedFactoryJSON, err := json.Marshal(factory)
	if err != nil {
		return fmt.Errorf("failed to marshal updated factory: %v", err)
	}
	err = ctx.GetStub().PutState(factoryID, updatedFactoryJSON)
	if err != nil {
		return fmt.Errorf("failed to update factory: %v", err)
	}
	return nil
}

// SetFactoryStatus updates the Status field of a Factory based on whether it was approved by both: auditor and retailer. Contains the following 1 specification: 1) SPEC_IsReadyforApproval
func (s *SmartContract) SetFactoryStatus(ctx contractapi.TransactionContextInterface, factoryID string) error {
	// Retrieve the factory from the world state
	factoryJSON, err := ctx.GetStub().GetState(factoryID)
	if err != nil {
		return fmt.Errorf("failed to read factory: %v", err)
	}
	if factoryJSON == nil {
		return fmt.Errorf("factory %s does not exist", factoryID)
	}
	var factory Factory
	err = json.Unmarshal(factoryJSON, &factory)
	if err != nil {
		return fmt.Errorf("failed to unmarshal factory: %v", err)
	}
	// Ensure all conditions are met for approval of factory
	if err := SPEC_IsReadyforApproval(factory.IsAuditorApproved, factory.IsRetailerApproved, !factory.IsFlagged); err != nil {
		return err
	}
	// Set the status to "approved" if all conditions are met
	factory.Status = "approved"
	// Update the updatedAt field
	factory.UpdatedAt = time.Now()
	// Marshal the updated factory and put it back in the world state
	updatedFactoryJSON, err := json.Marshal(factory)
	if err != nil {
		return fmt.Errorf("failed to marshal updated factory: %v", err)
	}
	err = ctx.GetStub().PutState(factoryID, updatedFactoryJSON)
	if err != nil {
		return fmt.Errorf("failed to update factory: %v", err)
	}
	return nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating admin-channel chaincode: %v", err)
	}
	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting admin-channel chaincode: %v", err)
	}
}
