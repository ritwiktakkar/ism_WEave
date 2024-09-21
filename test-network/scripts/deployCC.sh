#!/bin/bash

source scripts/utils.sh

CHANNEL_NAME=${1:-"mychannel"}
CC_NAME=${2}
CC_SRC_PATH=${3}
CC_SRC_LANGUAGE=${4}
CC_VERSION=${5:-"1.0"}
CC_SEQUENCE=${6:-"1"}
CC_INIT_FCN=${7:-"NA"}
CC_END_POLICY=${8:-"NA"}
CC_COLL_CONFIG=${9:-"NA"}
DELAY=${10:-"3"}
MAX_RETRY=${11:-"5"}
VERBOSE=${12:-"false"}

println "executing with the following"
println "- CHANNEL_NAME: ${C_GREEN}${CHANNEL_NAME}${C_RESET}"
println "- CC_NAME: ${C_GREEN}${CC_NAME}${C_RESET}"
println "- CC_SRC_PATH: ${C_GREEN}${CC_SRC_PATH}${C_RESET}"
println "- CC_SRC_LANGUAGE: ${C_GREEN}${CC_SRC_LANGUAGE}${C_RESET}"
println "- CC_VERSION: ${C_GREEN}${CC_VERSION}${C_RESET}"
println "- CC_SEQUENCE: ${C_GREEN}${CC_SEQUENCE}${C_RESET}"
println "- CC_END_POLICY: ${C_GREEN}${CC_END_POLICY}${C_RESET}"
println "- CC_COLL_CONFIG: ${C_GREEN}${CC_COLL_CONFIG}${C_RESET}"
println "- CC_INIT_FCN: ${C_GREEN}${CC_INIT_FCN}${C_RESET}"
println "- DELAY: ${C_GREEN}${DELAY}${C_RESET}"
println "- MAX_RETRY: ${C_GREEN}${MAX_RETRY}${C_RESET}"
println "- VERBOSE: ${C_GREEN}${VERBOSE}${C_RESET}"

INIT_REQUIRED="--init-required"
# check if the init fcn should be called
if [ "$CC_INIT_FCN" = "NA" ]; then
  INIT_REQUIRED=""
fi

if [ "$CC_END_POLICY" = "NA" ]; then
  CC_END_POLICY=""
else
  CC_END_POLICY="--signature-policy $CC_END_POLICY"
fi

if [ "$CC_COLL_CONFIG" = "NA" ]; then
  CC_COLL_CONFIG=""
else
  CC_COLL_CONFIG="--collections-config $CC_COLL_CONFIG"
fi

FABRIC_CFG_PATH=$PWD/../config/

# import utils
. scripts/envVar.sh
. scripts/ccutils.sh

function checkPrereqs() {
  jq --version > /dev/null 2>&1

  if [[ $? -ne 0 ]]; then
    errorln "jq command not found..."
    errorln
    errorln "Follow the instructions in the Fabric docs to install the prereqs"
    errorln "https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html"
    exit 1
  fi
}

#check for prerequisites
checkPrereqs

## package the chaincode
./scripts/packageCC.sh $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION

PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid ${CC_NAME}.tar.gz)

## Install chaincode
if [ $CHANNEL_NAME = "admin-channel" ]; then 
  infoln "Installing chaincode on peer0.org1..."
  installChaincode 1
  infoln "Install chaincode on peer0.org2..."
  installChaincode 2
  infoln "Installing chaincode on peer0.org3..."
  installChaincode 3

  resolveSequence
elif [ $CHANNEL_NAME = "production-channel" ]; then
  infoln "Installing chaincode on peer0.org1..."
  installChaincode 1
  infoln "Install chaincode on peer0.org2..."
  installChaincode 2
  infoln "Installing chaincode on peer0.org3..."
  installChaincode 3
  infoln "Install chaincode on peer0.org4..."
  installChaincode 4
  infoln "Install chaincode on peer0.org5..."
  installChaincode 5
  infoln "Install chaincode on peer0.org6..."
  installChaincode 6

  resolveSequence
else 
  infoln "Installing chaincode on peer0.org1..."
  installChaincode 1
  infoln "Install chaincode on peer0.org2..."
  installChaincode 2

  resolveSequence
fi

if [ $CHANNEL_NAME = "admin-channel" ]; then
  ## query whether the chaincode is installed
  queryInstalled 1
  ## approve the definition for org1
  approveForMyOrg 1
  ## check whether the chaincode definition is ready to be committed
  # checkCommitReadiness 1 "\"Org1MSP\": true" "\"Org2MSP\": false\"\"Org3MSP\": false"
  ## query whether the chaincode is installed
  queryInstalled 2
  ## approve the definition for org2
  approveForMyOrg 2
  ## check whether the chaincode definition is ready to be committed
  # checkCommitReadiness 2 "\"Org1MSP\": true" "\"Org2MSP\": true \"\"Org3MSP\": false"
  ## query whether the chaincode is installed
  queryInstalled 3
  ## approve the definition for org3
  approveForMyOrg 3
  ## check whether the chaincode definition is ready to be committed
  checkCommitReadiness 3 "\"Org1MSP\": true" "\"Org2MSP\": true" "\"Org3MSP\": true"
elif [ $CHANNEL_NAME = "production-channel" ]; then
  ## query whether the chaincode is installed
  queryInstalled 1
  ## approve the definition for org1
  approveForMyOrg 1
  ## check whether the chaincode definition is ready to be committed
  checkCommitReadiness 1 "\"Org1MSP\": true" "\"Org2MSP\": false" "\"Org3MSP\": false" "\"Org4MSP\": false" "\"Org5MSP\": false" "\"Org6MSP\": false"
  ## query whether the chaincode is installed
  queryInstalled 2
  ## approve the definition for org2
  approveForMyOrg 2
  ## check whether the chaincode definition is ready to be committed
  checkCommitReadiness 2 "\"Org1MSP\": true" "\"Org2MSP\": true" "\"Org3MSP\": false" "\"Org4MSP\": false" "\"Org5MSP\": false" "\"Org6MSP\": false"
  ## query whether the chaincode is installed
  queryInstalled 3
  ## approve the definition for org3
  approveForMyOrg 3
  ## check whether the chaincode definition is ready to be committed
  checkCommitReadiness 3 "\"Org1MSP\": true" "\"Org2MSP\": true" "\"Org3MSP\": true" "\"Org4MSP\": false" "\"Org5MSP\": false" "\"Org6MSP\": false"
  ## query whether the chaincode is installed
  queryInstalled 4
  ## approve the definition for org4
  approveForMyOrg 4
  ## check whether the chaincode definition is ready to be committed
  checkCommitReadiness 4 "\"Org1MSP\": true" "\"Org2MSP\": true" "\"Org3MSP\": true" "\"Org4MSP\": true" "\"Org5MSP\": false" "\"Org6MSP\": false"
  ## query whether the chaincode is installed
  queryInstalled 5
  ## approve the definition for org5
  approveForMyOrg 5
  ## check whether the chaincode definition is ready to be committed
  checkCommitReadiness 5 "\"Org1MSP\": true" "\"Org2MSP\": true" "\"Org3MSP\": true" "\"Org4MSP\": true" "\"Org5MSP\": true" "\"Org6MSP\": false"
  ## query whether the chaincode is installed
  queryInstalled 6
  ## approve the definition for org6
  approveForMyOrg 6
  ## check whether the chaincode definition is ready to be committed
  checkCommitReadiness 6 "\"Org1MSP\": true" "\"Org2MSP\": true" "\"Org3MSP\": true" "\"Org4MSP\": true" "\"Org5MSP\": true" "\"Org6MSP\": true"
else
  ## query whether the chaincode is installed
  queryInstalled 1
  ## approve the definition for org1
  approveForMyOrg 1
  ## check whether the chaincode definition is ready to be committed
  ## expect org1 to have approved and org2 not to
  checkCommitReadiness 1 "\"Org1MSP\": true" "\"Org2MSP\": false"
  checkCommitReadiness 2 "\"Org1MSP\": true" "\"Org2MSP\": false"
  ## now approve also for org2
  approveForMyOrg 2
  ## check whether the chaincode definition is ready to be committed
  ## expect them both to have approved
  checkCommitReadiness 1 "\"Org1MSP\": true" "\"Org2MSP\": true"
  checkCommitReadiness 2 "\"Org1MSP\": true" "\"Org2MSP\": true"
fi

## now that we know for sure both orgs have approved, commit the definition
if [ $CHANNEL_NAME = "admin-channel" ]; then 
  commitChaincodeDefinition 1 2 3
elif [ $CHANNEL_NAME = "production-channel" ]; then
  commitChaincodeDefinition 1 2 3 4 5 6
else
  commitChaincodeDefinition 1 2
fi

if [ $CHANNEL_NAME = "admin-channel" ]; then
  ## query on 3 orgs to see that the definition committed successfully
  queryCommitted 1
  queryCommitted 2
  queryCommitted 3
elif [ $CHANNEL_NAME = "production-channel" ]; then
  ## query on 6 orgs to see that the definition committed successfully
  queryCommitted 1
  queryCommitted 2
  queryCommitted 3
  queryCommitted 4
  queryCommitted 5
  queryCommitted 6
else
  ## query on both orgs to see that the definition committed successfully
  queryCommitted 1
  queryCommitted 2
fi

## Invoke the chaincode - this does require that the chaincode have the 'initLedger'
## method defined
if [ "$CC_INIT_FCN" = "NA" ]; then
  infoln "Chaincode initialization is not required"
else
  chaincodeInvokeInit 1 2
fi

exit 0
