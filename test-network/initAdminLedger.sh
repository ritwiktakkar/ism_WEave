#!/bin/bash

# Get the current working directory
CURRENT_DIR=$(pwd)

# Check if the current directory ends with "test-network"
if [[ "$CURRENT_DIR" != */test-network ]]; then
    echo "Error: This script must be run from a directory ending with 'test-network'"
    exit 1
fi

. scripts/utils.sh

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}/../config

. scripts/envVar.sh 

infoln "Deploying chaincode 'admin' on admin-channel..."
./network.sh deployCC -c admin-channel -ccn admin -ccp ../chaincode/admin-channel/ -ccl go -ccv 1.0 > /dev/null 2>&1
check_status "Deploying admin-channel chaincode"

# TRACE:   
# 1. Generating Order as retailer (org1)
# 2. Entering upstream factories to the world state as buying agent (org2)
# 3. Approving factories as retailer (org1)
# 4. Approving factories as auditor (org3)
# 5. Setting factory status now that they are approved by retailer and auditor
# 6. Issuing Plan as buying agent (org2)
# 7. Approving Plan as retailer (org1)
# 8. Approving Plan as auditor (org3)
# 9. Updating plan status field now that it is approved by retailer and auditor
# 10. Accepting Order as buying agent now that plan is approved (org2)
# 11. Updating Order status field now that it is accepted by buying agent

infoln "1/11. Generating Order as retailer (org1)..."
setGlobals 1
# Generate order
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 1 2 3 -c '{"Args":["CreateOrder","2024-05-10T10:00:00Z","2025-02-01T10:00:00Z","N/A","order_1","false","For Spring 2025","Net 60","200 shirts","Org2MSP","450000.00"]}'
check_status "Creating order"

sleep 2.5s

infoln "2/11. Entering upstream factories to the world state as buying agent (org2)..."
setGlobals 2
# Add raw materials supplier (org4) factory
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 2 1 3 -c '{"Args":["CreateFactory", "Org4MSP", "", "factory_1", "false", "Vadodara, Gujarat, India", "Example Cotton Mills", "", "true", "2008-11-22T10:00:00Z"]}'
check_status "Created factory 1"
# Add textiles (org5) factory
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 2 1 3 -c '{"Args":["CreateFactory", "Org5MSP", "", "factory_2", "false", "Vadodara, Gujarat, India", "Imaginary Textiles", "", "true", "2008-11-23T11:00:00Z","2021-01-28T11:30:00Z"]}'
check_status "Created factory 2"
# Add fps (org6) factory
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 2 1 3 -c '{"Args":["CreateFactory", "Org6MSP", "", "factory_3", "false", "Ashulia, Bangladesh", "Notareal Group", "", "true", "2010-11-24T12:00:00Z","2024-05-28T11:45:00Z"]}' 
check_status "Created factory 3"

sleep 2.5s

infoln "3/11. Approving factories as retailer (org1)..."
setGlobals 1
# Approve raw materials supplier (org4) factory
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 1 2 3 -c '{"Args":["SetFactoryApproval","factory_1","true"]}'
check_status "Approved factory 1 as retailer"
# Approve textiles (org5) factory
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 1 2 3 -c '{"Args":["SetFactoryApproval","factory_2","true"]}'
check_status "Approved factory 2 as retailer"
# Approve fps (org6) factory
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 1 2 3 -c '{"Args":["SetFactoryApproval","factory_3","true"]}'
check_status "Approved factory 3 as retailer"

sleep 2.5s

infoln "4/11. Approving factories as auditor (org3)..."
setGlobals 3
# Approve raw materials supplier org4) factory
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 3 2 1 -c '{"Args":["SetFactoryApproval","factory_1","true"]}'
check_status "Approved factory 1 as auditor"
# Approve textiles (org5) factory
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 3 2 1 -c '{"Args":["SetFactoryApproval","factory_2","true"]}'
check_status "Approved factory 2 as auditor"
# Approve fps (org6) factory
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 3 2 1 -c '{"Args":["SetFactoryApproval","factory_3","true"]}'
check_status "Approved factory 3 as auditor"

sleep 2.5s

infoln "5/11. Setting factory status now that they are approved by retailer and auditor..."
setGlobals 2
# Update status field in factory
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 2 1 3 -c '{"Args":["SetFactoryStatus","factory_1"]}'
check_status "Updating factory_1 status"
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 2 1 3 -c '{"Args":["SetFactoryStatus","factory_2"]}'
check_status "Updating factory_2 status"
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 2 1 3 -c '{"Args":["SetFactoryStatus","factory_3"]}'
check_status "Updating factory_3 status"

sleep 2.5s

infoln "6/11. Issuing Plan as buying agent (org2)..."
setGlobals 2
# Create plan
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 1 2 3 -c '{"Args":["CreatePlan","2024-05-21T10:00:00Z","[\"factory_1\",\"factory_2\",\"factory_3\"]","N/A","plan_1","false","","order_1","Add raw materials supplier, textiles manufacturer, and full-package supplier"]}'
check_status "Creating plan as buying agent"

sleep 2.5s

infoln "7/11. Approving Plan as retailer (org1)..."
setGlobals 1
# Update IsRetailerApproved field in plan
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 1 2 3 -c '{"Args":["SetPlanApproval","plan_1","true"]}'
check_status "Approving plan as retailer"

sleep 2.5s

infoln "8/11. Approving Plan as auditor (org3)..."
setGlobals 3
# Update IsAuditorApproved field in plan
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 3 2 1 -c '{"Args":["SetPlanApproval","plan_1","true"]}'
check_status "Approving plan as auditor"

sleep 2.5s

infoln "9/11. Updating plan status field now that it is approved by retailer and auditor..."
setGlobals 2
# Update status field in plan
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 2 1 3 -c '{"Args":["SetPlanStatus","plan_1"]}'
check_status "Updating plan status"

sleep 2.5s

infoln "10/11. Accepting Order as buying agent (org2)..."
setGlobals 2
# Update IsAccepted field in order
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 2 1 3 -c '{"Args":["SetOrderAcceptance","order_1","plan_1","true"]}'
check_status "Accepting order"

sleep 2.5s

infoln "11/11. Updating Order status field now that it is accepted by buying agent..."
setGlobals 1
# Update status field in order
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C admin-channel -n admin 1 2 3 -c '{"Args":["SetOrderStatus","order_1","accepted"]}'
check_status "Updating order status"

sleep 2.5s

# View all assets in the ledger
infoln ""
infoln "Here are all assets in the ledger AFTER execution of all init commands:"
infoln "----------"
peer chaincode query -C admin-channel -n admin -c '{"Args":["GetAllAssets"]}' | jq .
infoln "----------"
infoln "Feel free to invoke other chaincode functions to interact with the ledger, such as those defined in chaincode/admin-channel/admin_audit_functions.go..." 