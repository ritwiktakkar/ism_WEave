#!/bin/bash

# Fixed order size
ORDER_SIZE=10000

# Number of simulations to run (variable foo)
foo=3

# Dynamically create the CSV filename based on foo and order size
output_file="repeated_${foo}_results_${ORDER_SIZE}.csv"

# Function to run the simulation and collect metrics
run_simulation() {
    # bring network down to known state and restart necessary containers and channel
    ./network.sh down
    ./networkSetup.sh 1

    local order_size=$1
    local run_number=$2
    local remaining_runs=$((foo - run_number))

    # Print the current run and the remaining runs
    echo "Running simulation $run_number of $foo for $order_size shirts"
    echo "$remaining_runs run(s) remaining"
    
    # Record start time
    start_time=$(date +%s.%N)
    
    # Running simulation command here, using the order size in the filename
    ./initProductionLedger_${order_size}.sh &
    sim_pid=$!
    
    # Start CPU and memory monitoring in the background
    (
        peak_memory=0
        while kill -0 $sim_pid 2>/dev/null; do
            current_memory=$(grep VmRSS /proc/$sim_pid/status 2>/dev/null | awk '{print $2}')
            if [[ $current_memory -gt $peak_memory ]]; then
                peak_memory=$current_memory
            fi
            sleep 0.1
        done
        echo $peak_memory > ${order_size}_peak_memory.txt
    ) &
    monitor_pid=$!
    
    pidstat -p $sim_pid 1 > ${order_size}_cpu_usage.txt &
    pidstat_pid=$!
    
    # Wait for the simulation to complete
    wait $sim_pid
    
    # Stop monitoring
    kill $pidstat_pid
    wait $monitor_pid
    
    # Record end time
    end_time=$(date +%s.%N)
    
    # Calculate total time
    total_time=$(echo "$end_time - $start_time" | bc)
    
    # Process CPU usage data
    avg_cpu=$(awk '{sum+=$8} END {print sum/NR}' ${order_size}_cpu_usage.txt)
    
    # Get peak memory usage (in KB)
    peak_memory=$(cat ${order_size}_peak_memory.txt)
    
    # Output results to the dynamically named CSV file
    echo "$order_size,$total_time,$avg_cpu,$peak_memory" >> $output_file
}

# Main execution
echo "Order Size,Total Time,Avg CPU Usage (%),Peak Memory (KB)" > $output_file

# Loop to run simulations for the fixed order size
for ((i = 1; i <= foo; i++)); do
    run_simulation $ORDER_SIZE $i
done

./network.sh down

echo "Testing complete. Results saved in $output_file"
