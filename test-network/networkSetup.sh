#!/bin/bash

# Get the current working directory
CURRENT_DIR=$(pwd)

# Check if the current directory ends with "test-network"
if [[ "$CURRENT_DIR" != */test-network ]]; then
    echo "Error: This script must be run from a directory ending with 'test-network'"
    exit 1
fi

echo "Script is running from the correct directory (i.e., */test-network)."

. scripts/utils.sh

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}/../config

# Set default values
createProductionChannel=${1:-0}

# Define the containers' name
ORG1_PEER_CONTAINERS="peer0.org1.example.com"
ORG2_PEER_CONTAINERS="peer0.org2.example.com"
ORG3_PEER_CONTAINERS="peer0.org3.example.com"
ORG4_PEER_CONTAINERS="peer0.org4.example.com"
ORG5_PEER_CONTAINERS="peer0.org5.example.com"
ORG6_PEER_CONTAINERS="peer0.org6.example.com"
ORDERER_PEER_CONTAINER="orderer.example.com"

create_containers_and_admin_channel() {
    ./network.sh up > /dev/null 2>&1
    check_status "Starting network"
    infoln "Creating auditor (org3)..."
    cd addOrg3/
    ./addOrg3.sh up > /dev/null 2>&1
    successln "Created containers for retailer (org1), buying agent (org2), auditor (org3), and the orderer."
    infoln "Creating channel 'admin-channel' with retailer, buying agent, and auditor..."
    cd ../
    ./network.sh createAdminChannel > /dev/null 2>&1
    successln "Created admin-channel"
}

create_containers_and_admin_production_channel() {
    ADMIN_ALREADY_CREATED=$1
    if [ ${ADMIN_ALREADY_CREATED} != 1 ]; then
        create_containers_and_admin_channel
    fi   
    infoln "Creating raw materials supplier (org4)..."
    cd addOrg4/
    ./addOrg4.sh up > /dev/null 2>&1
    infoln "Creating textiles manufacturer (org5)..."
    cd ../addOrg5/
    ./addOrg5.sh up > /dev/null 2>&1
    infoln "Creating full-package supplier (org6)..."
    cd ../addOrg6/
    ./addOrg6.sh up > /dev/null 2>&1
    successln "Created containers for raw materials supplier (org4), textiles manufacturer (org5), and full-package supplier (org6)."
    infoln "Creating channel 'production-channel' with all organizations..."
    cd ../
    ./network.sh createProductionChannel > /dev/null 2>&1
    check_status "Creating production-channel"
}

if [ ${createProductionChannel} = 0 ]; then
    # Check if Org1, Org2, and Org3 containers are already running
    if [ "$(docker ps -aq -f status=running -f name=$ORG1_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=running -f name=$ORG2_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=running -f name=$ORG3_PEER_CONTAINERS)" ]; then
        echo "Containers $ORG1_PEER_CONTAINERS, $ORG2_PEER_CONTAINERS, and $ORG3_PEER_CONTAINERS are already running."
    else
        # Check if the containers are paused (but not running)
        if [ "$(docker ps -aq -f status=paused -f name=$ORG1_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=paused -f name=$ORG2_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=paused -f name=$ORG3_PEER_CONTAINERS)" ]; then
            echo "Containers $ORG1_PEER_CONTAINERS, $ORG2_PEER_CONTAINERS, and $ORG3_PEER_CONTAINERS paused. Unpausing the containers."
            # Unpause the containers
            # Find all paused containers matching the pattern and unpause them
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG1_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG2_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG3_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            docker unpause $ORDERER_PEER_CONTAINER
            if [ "$(docker ps -aq -f status=exited)" ]; then
                for container_id in $(docker ps -aq -f status=exited); do
                    echo "FAIL: container with ID $container_id exited thus removed."
                    docker rm $container_id
                done
            else
                successln "All containers unpaused"
            fi
        else
            echo "Containers $ORG1_PEER_CONTAINERS, $ORG2_PEER_CONTAINERS, and $ORG3_PEER_CONTAINERS do not exist."
            infoln "Creating 3 containers (Orgs1-3) and the orderer... Then creating admin-channel."
            create_containers_and_admin_channel
        fi
    fi
else
    # Check if Org1, Org2, Org3, Org4, Org5, and Org6 containers are already running
    if [ "$(docker ps -aq -f status=running -f name=$ORG1_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=running -f name=$ORG2_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=running -f name=$ORG3_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=running -f name=$ORG4_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=running -f name=$ORG5_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=running -f name=$ORG6_PEER_CONTAINERS)" ]; then
        echo "Containers $ORG1_PEER_CONTAINERS, $ORG2_PEER_CONTAINERS, $ORG3_PEER_CONTAINERS, $ORG4_PEER_CONTAINERS, $ORG5_PEER_CONTAINERS, and $ORG6_PEER_CONTAINERS are already running."
    else
        # Check if all 6 containers are paused (but not running)
        if [ "$(docker ps -aq -f status=paused -f name=$ORG1_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=paused -f name=$ORG2_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=paused -f name=$ORG3_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=paused -f name=$ORG4_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=paused -f name=$ORG5_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=paused -f name=$ORG6_PEER_CONTAINERS)" ]; then
            echo "Containers $ORG1_PEER_CONTAINERS, $ORG2_PEER_CONTAINERS, $ORG3_PEER_CONTAINERS, $ORG4_PEER_CONTAINERS, $ORG5_PEER_CONTAINERS and $ORG6_PEER_CONTAINERS paused. Unpausing the containers."
            # Unpause the containers
            # Find all paused containers matching the pattern and unpause them
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG1_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG2_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG3_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG4_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG5_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG6_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            docker unpause $ORDERER_PEER_CONTAINER
            if [ "$(docker ps -aq -f status=exited)" ]; then
                for container_id in $(docker ps -aq -f status=exited); do
                    echo "FAIL: container with ID $container_id exited thus removed."
                    docker rm $container_id
                done
            else
                successln "All containers unpaused"
            fi
        # check if only org1-3 containers are paused and org4-6 containers don't exist
        elif [ "$(docker ps -aq -f status=paused -f name=$ORG1_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=paused -f name=$ORG2_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=paused -f name=$ORG3_PEER_CONTAINERS)" ] && [ -z "$(docker ps -aq -f name=$ORG4_PEER_CONTAINERS)" ] && [ -z "$(docker ps -aq -f name=$ORG5_PEER_CONTAINERS)" ] && [ -z "$(docker ps -aq -f name=$ORG6_PEER_CONTAINERS)" ]; then
            echo "Containers $ORG1_PEER_CONTAINERS, $ORG2_PEER_CONTAINERS, and $ORG3_PEER_CONTAINERS paused - unpausing them and starting additional containers for orgs4-6."
            # Unpause the containers
            # Find all paused containers matching the pattern and unpause them
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG1_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG2_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG3_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            docker unpause $ORDERER_PEER_CONTAINER
            if [ "$(docker ps -aq -f status=exited)" ]; then
                for container_id in $(docker ps -aq -f status=exited); do
                    echo "FAIL: container with ID $container_id exited thus removed."
                    docker rm $container_id
                done
            else
                check_status "All 3 containers unpaused (orgs1-3) and orderer"
            fi
            create_containers_and_admin_production_channel 1
        # check if only org1-3 containers are running and org4-6 containers don't exist
        elif [ "$(docker ps -aq -f status=running -f name=$ORG1_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=running -f name=$ORG2_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=running -f name=$ORG3_PEER_CONTAINERS)" ] && [ -z "$(docker ps -aq -f name=$ORG4_PEER_CONTAINERS)" ] && [ -z "$(docker ps -aq -f name=$ORG5_PEER_CONTAINERS)" ] && [ -z "$(docker ps -aq -f name=$ORG6_PEER_CONTAINERS)" ]; then
            echo "Containers $ORG1_PEER_CONTAINERS, $ORG2_PEER_CONTAINERS, and $ORG3_PEER_CONTAINERS already running. Starting additional containers for orgs4-6."
            create_containers_and_admin_production_channel 1
        # check if only org1-3 containers are running and org4-6 containers are paused
        elif [ "$(docker ps -aq -f status=running -f name=$ORG1_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=running -f name=$ORG2_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=running -f name=$ORG3_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=paused -f name=$ORG4_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=paused -f name=$ORG5_PEER_CONTAINERS)" ] && [ "$(docker ps -aq -f status=paused -f name=$ORG6_PEER_CONTAINERS)" ]; then
            echo "Containers $ORG1_PEER_CONTAINERS, $ORG2_PEER_CONTAINERS, and $ORG3_PEER_CONTAINERS already running. Unpausing containers for orgs4-6."
            # Unpause the containers
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG4_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG5_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            for container_id in $(docker ps -a -f status=paused -q --filter "name=$ORG6_PEER_CONTAINERS"); do
                echo "Unpausing container with ID $container_id"
                docker unpause $container_id
            done
            if [ "$(docker ps -aq -f status=exited)" ]; then
                for container_id in $(docker ps -aq -f status=exited); do
                    echo "FAIL: container with ID $container_id exited thus removed."
                    docker rm $container_id
                done
            else
                successln "All containers unpaused"
            fi
        else
            echo "Containers $ORG1_PEER_CONTAINERS, $ORG2_PEER_CONTAINERS, $ORG3_PEER_CONTAINERS, $ORG4_PEER_CONTAINERS, $ORG5_PEER_CONTAINERS and $ORG6_PEER_CONTAINERS do not exist."
            infoln "Creating 6 containers (Orgs1-6) and the orderer... Then creating admin-channel, and finally, production-channel."
            create_containers_and_admin_production_channel 0
        fi
    fi
fi
