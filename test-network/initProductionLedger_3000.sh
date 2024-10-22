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

infoln "Deploying chaincode 'production' on production-channel..."
./network.sh deployCC -c production-channel -ccn production -ccp ../chaincode/production-channel/ -ccl go -ccv 1.0 > /dev/null 2>&1
check_status "Deploying production-channel chaincode"

# TRACE:   
#   1. Add cotton bale to ledger as raw materials supplier (org4)
#   1.1 Assemble cotton bale into lot (org4)
#   2. Add cotton yarn to ledger as raw materials supplier (org4)
#   2.1 Assemble cotton yarn into lots (org4)
#   3. Update cotton yarn lots owner to textiles (org5)
#   4. Add unfinished fabric to ledger (org5)
#   4.1 Assemble unfinished fabric into lots (org5)
#   5. Add finished fabric to ledger as textiles (org5)
#   5.1 Assemble finished fabric into lots (org6)
#   6. Update finished fabric lots owner to fps (org6)
#   7. Add cut parts to ledger as fps (org6)
#   8. Add buttons to ledger as fps (org6)
#   9. Assemble cut parts and buttons into shirts (org6)
#   10. Pack assembled garments into cartons (org6)
#   11. Place cartons in container (org6)

#MANUAL
ORDER_QUANTITY=3000
CONES_OF_YARN=4650

infoln "1/11. Adding cotton bale to ledger as raw materials supplier (org4)..."
setGlobals 4
# Add cotton bale 
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 4 1 2 3 5 6 -c '{"Args":["CreateCottonBale","true","2024-07-01T10:00:00Z","","cottonbale_1","false","One bale of cotton; Staple length 34; Strength 29 g/tex; Micronaire 4.0; Color grade 41; Leaf grade 3-4; Uniformity 81%","Vadodara, Gujarat, India","Medium","480.00"]}'

sleep 2s

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 4 1 2 3 5 6 -c '{"Args":["CreateCottonBale","true","2024-07-01T10:05:10Z","","cottonbale_2","false","One bale of cotton; Staple length 34; Strength 29 g/tex; Micronaire 4.0; Color grade 41; Leaf grade 3-4; Uniformity 81%","Vadodara, Gujarat, India","Medium","480.00"]}'

sleep 2s

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 4 1 2 3 5 6 -c '{"Args":["CreateCottonBale","true","2024-07-01T10:10:10Z","","cottonbale_3","false","One bale of cotton; Staple length 34; Strength 29 g/tex; Micronaire 4.0; Color grade 41; Leaf grade 3-4; Uniformity 81%","Vadodara, Gujarat, India","Medium","480.00"]}'

sleep 2s

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 4 1 2 3 5 6 -c '{"Args":["CreateCottonBale","true","2024-07-01T10:15:10Z","","cottonbale_4","false","One bale of cotton; Staple length 34; Strength 29 g/tex; Micronaire 4.0; Color grade 41; Leaf grade 3-4; Uniformity 81%","Vadodara, Gujarat, India","Medium","480.00"]}'

sleep 2s

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 4 1 2 3 5 6 -c '{"Args":["CreateCottonBale","true","2024-07-01T10:20:10Z","","cottonbale_5","false","One bale of cotton; Staple length 34; Strength 29 g/tex; Micronaire 4.0; Color grade 41; Leaf grade 3-4; Uniformity 81%","Vadodara, Gujarat, India","Medium","480.00"]}'

sleep 2s

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 4 1 2 3 5 6 -c '{"Args":["CreateCottonBale","true","2024-07-01T10:25:10Z","","cottonbale_6","false","One bale of cotton; Staple length 34; Strength 29 g/tex; Micronaire 4.0; Color grade 41; Leaf grade 3-4; Uniformity 81%","Vadodara, Gujarat, India","Medium","480.00"]}'

check_status "Creating one full bale of cotton suitable for order quantity of $ORDER_QUANTITY shirts"

sleep 8s

infoln "1.1/11. Assembling cotton bale into lot (org4)..."
setGlobals 4
# Assemble cotton bale into lot
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 4 1 2 3 5 6 -c '{"Args":["CreateLot","2024-07-01T11:00:00Z","cottonbale_","[\"cottonbale_1\",\"cottonbale_2\",\"cottonbale_3\",\"cottonbale_4\",\"cottonbale_5\",\"cottonbale_6\"]","Vadodara, Gujarat, India","","lot_1","false", "", "Vadodara, Gujarat, India", "Org4MSP", "2880.00"]}'

check_status "Assembling cotton bales into 1 lot, i.e., lot_1"

sleep 5s

infoln "2/11. Adding cotton yarn to ledger as raw materials supplier (org4)..."
setGlobals 4
# Add cotton yarn
for i in $(seq 1 $CONES_OF_YARN) 
do
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 4 1 2 3 5 6 -c "{\"Args\":[\"CreateCottonYarn\",\"true\",\"2024-07-03T10:00:00Z\",\"[\\\"lot_1\\\"]\",\"\",\"cottonyarn_$i\",\"false\",\"\",\"Vadodara, Gujarat, India\",\"0.397\", \"30\"]}"
done
check_status "Creating $CONES_OF_YARN cones of cotton yarn"

sleep 5s

infoln "2.1/11. Assembling cotton yarn into lots (org4)..."
setGlobals 4
# Assemble cotton yarn into lots
generate_yarn_list() {
    local start=$1
    local end=$2
    local list=""
    for i in $(seq $start $end); do
        list+="\\\"cottonyarn_$i\\\","
    done
    echo "${list%,}"  # Remove the trailing comma
}
# First invocation for lot with first half of cones
yarn_list_1=$(generate_yarn_list 1 $((CONES_OF_YARN / 2)))
yarn_lot_1_weight=$(echo "scale=2; $CONES_OF_YARN / 2 * 0.397" | bc)
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 4 1 2 3 5 6 -c "{\"Args\":[\"CreateLot\",\"2024-07-04T11:00:05Z\",\"cottonyarn_\",\"[$yarn_list_1]\",\"Vadodara, Gujarat, India\",\"\",\"lot_2\",\"false\", \"\", \"Vadodara, Gujarat, India\", \"Org4MSP\", \"$yarn_lot_1_weight\"]}"

sleep 5s

# Second invocation for lot with remaining half of cones
yarn_list_2=$(generate_yarn_list $((CONES_OF_YARN / 2 + 1)) $CONES_OF_YARN)
yarn_lot_2_weight=$(echo "scale=2; $CONES_OF_YARN * 0.397 - $yarn_lot_1_weight" | bc)
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 4 1 2 3 5 6 -c "{\"Args\":[\"CreateLot\",\"2024-07-04T11:10:00Z\",\"cottonyarn_\",\"[$yarn_list_2]\",\"Vadodara, Gujarat, India\",\"\",\"lot_3\",\"false\", \"\", \"Vadodara, Gujarat, India\", \"Org4MSP\", \"$yarn_lot_2_weight\"]}"
check_status "Assembling cotton yarn into 2 lots, i.e., lot_2 and lot_3"

sleep 10s

infoln "3/11. Update cotton yarn lots ownership to textiles manufacturer (org5)..."
setGlobals 5
# Update cotton yarn lots owner to textiles (org5)
for i in {2..3}
do
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 5 1 2 3 4 6 -c "{\"Args\":[\"UpdateLotOwner\",\"lot_$i\",\"Org5MSP\"]}"
done
check_status "Updating cotton yarn lots ownership to textiles manufacturer (org5), i.e., lot_2 and lot_3"

sleep 5s

infoln "4/11. Adding unfinished fabric to ledger (org5)..."
setGlobals 5
# Add unfinished fabric
# MANUALLY INPUT RANGE SUCH THAT SPLIT EVENLY (I.E., 16.74 * UPPER BOUND < WEIGHT[LOT2,3])
for i in {1..54}
do
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 5 1 2 3 4 6 -c "{\"Args\":[\"CreateUnfinishedFabric\",\"true\",\"2024-07-06T10:00:00Z\",\"[\\\"lot_2\\\"]\",\"\",\"unfinishedfabric_$i\",\"false\",\"Length in linear yards; Weight in lbs; Width in inches.\",\"Vadodara, Gujarat, India\",\"50\", \"16.74\", \"60\"]}"
done

sleep 5s

# MANUALLY INPUT RANGE SUCH THAT SPLIT EVENLY (I.E., 16.74 * UPPER BOUND < WEIGHT[LOT2,3])
for i in {55..108}
do
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 5 1 2 3 4 6 -c "{\"Args\":[\"CreateUnfinishedFabric\",\"true\",\"2024-07-06T12:00:00Z\",\"[\\\"lot_3\\\"]\",\"\",\"unfinishedfabric_$i\",\"false\",\"Length in linear yards; Width in inches; Weight in lbs.\",\"Vadodara, Gujarat, India\",\"50\", \"16.74\", \"60\"]}"
done
check_status "Creating 108 rolls of unfinished fabric"

sleep 5s

infoln "4.1/11. Assembling unfinished fabric into lots (org5)..."
setGlobals 5
# Assemble unfinished fabric into lots
generate_unfinished_fabric_list() {
    local start=$1
    local end=$2
    local list=""
    for i in $(seq $start $end); do
        list+="\\\"unfinishedfabric_$i\\\","
    done
    echo "${list%,}"  # Remove the trailing comma
}
# First invocation for lot with first half unfinished fabric rolls
# MANUALLY INPUT RANGE FROM STEP 4
unfinished_fabric_list_1=$(generate_unfinished_fabric_list 1 54)
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 5 1 2 3 4 6 -c "{\"Args\":[\"CreateLot\",\"2024-07-07T11:00:05Z\",\"unfinishedfabric_\",\"[$unfinished_fabric_list_1]\",\"Vadodara, Gujarat, India\",\"\",\"lot_4\",\"false\", \"\", \"Vadodara, Gujarat, India\", \"Org5MSP\", \"903.96\"]}" # MANUAL: last variable = 16.74 * UPPER BOUND FROM #4

sleep 5s

# Second invocation for lot with remaining unfinished fabric rolls
# MANUALLY INPUT RANGE FROM STEP 4
unfinished_fabric_list_2=$(generate_unfinished_fabric_list 55 108)
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 5 1 2 3 4 6 -c "{\"Args\":[\"CreateLot\",\"2024-07-07T11:10:00Z\",\"unfinishedfabric_\",\"[$unfinished_fabric_list_2]\",\"Vadodara, Gujarat, India\",\"\",\"lot_5\",\"false\", \"\", \"Vadodara, Gujarat, India\", \"Org5MSP\", \"903.96\"]}" # MANUAL: last variable = 16.74 * UPPER BOUND FROM #4
check_status "Assembling unfinished fabric into 2 lots, i.e., lot_4 and lot_5"

sleep 5s

infoln "5/11. Adding finished fabric to ledger as textiles manufacturer (org5)..."
setGlobals 5
# Add finished fabric
# MANUALLY INPUT RANGE FROM STEP 4
for i in {1..54}
do
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 5 1 2 3 4 6 -c "{\"Args\":[\"CreateFinishedFabric\",\"true\",\"2024-07-09T10:00:00Z\",\"[\\\"lot_4\\\"]\",\"finishedfabric_$i\",\"\",\"false\",\"47.5\",\"Length in linear yards; Weight in lbs; Width in inches.\",\"Vadodara, Gujarat, India\",\"15.9025\",\"58.8\"]}"
done

sleep 5s

# MANUALLY INPUT RANGE FROM STEP 4
for i in {55..108}
do
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 5 1 2 3 4 6 -c "{\"Args\":[\"CreateFinishedFabric\",\"true\",\"2024-07-09T12:00:00Z\",\"[\\\"lot_5\\\"]\",\"finishedfabric_$i\",\"\",\"false\",\"47.5\",\"Length in linear yards; Weight in lbs; Width in inches.\",\"Vadodara, Gujarat, India\",\"15.9025\",\"58.8\"]}"
done
check_status "Creating 108 rolls of finished fabric"

sleep 5s

infoln "5.1/11. Assembling finished fabric into lots (org5)..."
setGlobals 5
# Assemble finished fabric into lots
generate_finished_fabric_list() {
    local start=$1
    local end=$2
    local list=""
    for i in $(seq $start $end); do
        list+="\\\"finishedfabric_$i\\\","
    done
    echo "${list%,}"  # Remove the trailing comma
}

# First invocation for lot with first finished fabric rolls
# MANUALLY INPUT RANGE FROM STEP 4
finished_fabric_list_1=$(generate_finished_fabric_list 1 54)
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 5 1 2 3 4 6 -c "{\"Args\":[\"CreateLot\",\"2024-07-10T11:00:05Z\",\"finishedfabric_\",\"[$finished_fabric_list_1]\",\"Ashulia, Bangladesh\",\"\",\"lot_6\",\"false\", \"\", \"Vadodara, Gujarat, India\", \"Org5MSP\", \"858.735\"]}" # MANUAL: last variable = 15.9025 * UPPER BOUND FROM #4

sleep 5s

# Second invocation for lot with remaining finished fabric rolls
# MANUALLY INPUT RANGE FROM STEP 4
finished_fabric_list_2=$(generate_finished_fabric_list 55 108)
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 5 1 2 3 4 6 -c "{\"Args\":[\"CreateLot\",\"2024-07-10T11:10:00Z\",\"finishedfabric_\",\"[$finished_fabric_list_2]\",\"Ashulia, Bangladesh\",\"\",\"lot_7\",\"false\", \"\", \"Vadodara, Gujarat, India\", \"Org5MSP\", \"858.735\"]}" # MANUAL: last variable = 15.9025 * UPPER BOUND FROM #4
check_status "Assembling finished fabric into 2 lots, i.e., lot_6 and lot_7"

sleep 5s

infoln "6/11. Update finished fabric lots ownership to fps (org6)..."
setGlobals 6
# Update finished fabric lots owner to fps (org6)
for i in {6..7}
do
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 6 1 2 3 4 5 -c "{\"Args\":[\"UpdateLotOwner\",\"lot_$i\",\"Org6MSP\"]}"
done
check_status "Updating finished fabric lots ownership to fps (org6), i.e., lot_6 and lot_7"

sleep 5s

infoln "7/11. Adding cut parts to ledger as fps (org6)..."
setGlobals 6
# Array of parts
parts=(
  "front_panel"
  "back_panel"
  "left_sleeve"
  "right_sleeve"
  "collar"
  "front_pocket"
)
# Add cut parts
total_parts=$((ORDER_QUANTITY * 6))
# for i in {1..3000}
for i in $(seq 1 $total_parts)
do
  # Calculate the index for the dates array
  part_index=$(( (i-1) / $ORDER_QUANTITY ))
  # Get the appropriate part
  current_part=${parts[$part_index]}
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 6 1 2 3 4 5 -c "{\"Args\":[\"CreateCutPart\",\"true\",\"2024-07-13T11:10:00Z\",\"[\\\"lot_6\\\", \\\"lot_7\\\"]\",\"\",\"cutpart_$i\",\"false\",\"Weight in lbs.\",\"Ashulia, Bangladesh\",\"$current_part\", \"0.091\"]}"
done
check_status "Creating $total_parts cut parts for shirts"

sleep 5s

infoln "8/11. Adding buttons to ledger as fps (org6)..."
setGlobals 6
# Add buttons
total_buttons=$((ORDER_QUANTITY * 10))
for i in $(seq 1 $total_buttons)
do
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 6 1 2 3 4 5 -c "{\"Args\":[\"CreateButton\",\"true\",\"2024-07-14T11:00:00Z\",\"\",\"button_$i\",\"false\",\"Weight in lbs.\",\"Ashulia, Bangladesh\",\"0.00165\"]}"
done

check_status "Creating $total_buttons buttons for shirts"

sleep 5s

infoln "9/11. Assembling cut parts and buttons into shirts (org6)..."
setGlobals 6
# Array of dates for the assembled garments
dates=(
  "2024-07-16T10:00:00Z"
  "2024-07-17T10:00:00Z"
)

# Function to generate a list of cut parts for a shirt
generate_cut_parts() {
    local shirt_num=$1
    local cut_parts=""
    local offsets=(0 $ORDER_QUANTITY $((ORDER_QUANTITY * 2)) $((ORDER_QUANTITY * 3)) $((ORDER_QUANTITY * 4)) $((ORDER_QUANTITY * 5)))

    for offset in "${offsets[@]}"
    do
        cut_parts+="\\\"cutpart_$((offset + shirt_num))\\\","
    done
    echo "${cut_parts%,}"  # Remove trailing comma
}

# Function to generate a list of buttons for a shirt
generate_buttons() {
    local shirt_num=$1
    local buttons=""
    for i in {1..10}
    do
        buttons+="\\\"button_$(( (shirt_num-1)*10 + i ))\\\","
    done
    echo "${buttons%,}"  # Remove trailing comma
}

# Add assembled garments
# for shirt_num in {1..500}
for shirt_num in $(seq 1 $ORDER_QUANTITY)
do
    # Calculate the index for the dates array
    date_index=$(( (shirt_num-1) / ($ORDER_QUANTITY/2) ))
    current_date=${dates[$date_index]}

    # Generate cut parts and buttons lists
    cut_parts=$(generate_cut_parts $shirt_num)
    buttons=$(generate_buttons $shirt_num)

    # Invoke the chaincode
    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 6 1 2 3 4 5 -c "{\"Args\":[\"CreateAssembledGarment\",\"true\",\"$current_date\",\"[${buttons}]\",\"[${cut_parts}]\",\"\",\"assembledgarment_$shirt_num\",\"false\",\"Weight in lbs.\",\"Ashulia, Bangladesh\",\"0.554\"]}"
done
check_status "Creating $ORDER_QUANTITY shirts"

sleep 5s

infoln "10/11. Pack assembled garments into cartons (org6)..."
setGlobals 6

carton_weight=32.0
shirts_per_carton=50
total_cartons=$((ORDER_QUANTITY / shirts_per_carton))

# Pack assembled garments into cartons
# for carton_num in {1..10}
for carton_num in $(seq 1 $total_cartons)
do
    # Calculate the range of shirts in this carton
    start_shirt=$(( (carton_num - 1) * shirts_per_carton + 1 ))
    end_shirt=$(( carton_num * shirts_per_carton ))

    # Generate the list of assembled garments in this carton
    garments=""
    for i in $(seq $start_shirt $end_shirt)
    do
        garments+="\\\"assembledgarment_$i\\\","
    done
    garments=${garments%,}  # Remove trailing comma

    # Invoke the chaincode to create a carton
    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 6 1 2 3 4 5 -c "{\"Args\":[\"CreateCarton\",\"true\",\"2024-07-18T10:00:00Z\",\"[$garments]\",\"Org1MSP\",\"\",\"carton_$carton_num\",\"false\",\"Weight in lbs.\",\"Ashulia, Bangladesh\",\"Org6MSP\",\"$carton_weight\"]}"
done
check_status "Packing $ORDER_QUANTITY shirts into $total_cartons cartons"

sleep 5s

infoln "11/11. Place cartons in container (org6)..."
setGlobals 6

# Set the required variables
destination_port="Los Angeles, California, USA"
origin_port="Chittagong, Bangladesh"
loaded_at="2024-07-20T08:00:00Z"
vessel="EXAMPLE Hong Kong"
container_id="container_1"

# Generate the content list (total cartons)
content=""
# for i in {1..10}
for i in $(seq 1 $total_cartons)
do
    content+="\\\"carton_$i\\\","
done
content=${content%,}  # Remove trailing comma

# Invoke the chaincode
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C production-channel -n production 6 1 2 3 4 5 -c "{\"Args\":[\"CreateContainer\",\"[$content]\",\"$destination_port\",\"\",\"$container_id\",\"\",\"$loaded_at\",\"$origin_port\",\"38448\",\"$vessel\"]}"
check_status "Creating a container with $total_cartons cartons of shirts"

sleep 10s

infoln ""

# Get blockchain info
peer channel getinfo -c production-channel > ${ORDER_QUANTITY}_ledger_info.txt
# Get ledger size
docker exec peer0.org6.example.com du -sb /var/hyperledger/production/ledgersData/chains/chains/production-channel > ${ORDER_QUANTITY}_ledger_size.txt
# Get state size
docker exec peer0.org6.example.com du -sb /var/hyperledger/production/ledgersData/stateLeveldb > ${ORDER_QUANTITY}_state_size.txt

infoln "Blockchain info saved to 5 files starting with ${ORDER_QUANTITY}_*.txt"

# View container asset in the ledger
infoln "Here is the container AFTER execution of all init commands:"
infoln "----------"
peer chaincode query -C production-channel -n production -c '{"Args":["GetAsset","container_1"]}' | jq .
infoln "----------"
infoln "Feel free to invoke other chaincode functions to interact with the ledger, such as those defined in chaincode/production-channel/production_audit_functions.go..." 
