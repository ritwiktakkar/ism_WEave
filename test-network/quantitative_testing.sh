#!/bin/bash

# Array of order sizes to test
ORDER_SIZES=(200 500 1000 2000 3000 5000 10000 15000 20000)

# Function to run the simulation and collect metrics
run_simulation() {
    # bring network down to known state and restart necessary containers and channel
    ./network.sh down
    ./networkSetup.sh 1

    local order_size=$1
    echo "Running simulation for $order_size shirts"
    
    # Record start time
    start_time=$(date +%s.%N)
    
    # Runing simulation command here, using the order size in the filename
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
    
    # Output results to a CSV file
    echo "$order_size,$total_time,$avg_cpu,$peak_memory" >> results.csv
}

# Main execution
echo "Order Size,Total Time,Avg CPU Usage (%),Peak Memory (KB)" > results.csv
for size in "${ORDER_SIZES[@]}"; do
    run_simulation $size
done

./network.sh down

echo "Testing complete. Results saved in results.csv"
