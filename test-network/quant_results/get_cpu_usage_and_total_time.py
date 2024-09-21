import pandas as pd
import sys
from datetime import datetime


def get_cpu_usage_and_total_time(file_path):
    # Load the file and skip the first few rows to extract the data
    data = pd.read_csv(file_path, delim_whitespace=True, skiprows=3, header=None)

    # Rename columns to match expected structure
    column_names = [
        "Time",
        "Period",
        "UID",
        "PID",
        "%usr",
        "%system",
        "%guest",
        "%wait",
        "%CPU",
        "CPU",
        "Command",
    ]
    data.columns = column_names

    # Extract %CPU column as float
    cpu_usage = data["%CPU"].astype(float)

    # Calculate average CPU usage including 0
    average_cpu_including_zeros = cpu_usage.mean()

    # Calculate average CPU usage excluding 0
    average_cpu_excluding_zeros = cpu_usage[cpu_usage > 0].mean()

    # Convert time column to datetime (handle AM/PM)
    start_time = datetime.strptime(
        data["Time"].iloc[0] + " " + data["Period"].iloc[0], "%I:%M:%S %p"
    )
    end_time = datetime.strptime(
        data["Time"].iloc[-1] + " " + data["Period"].iloc[-1], "%I:%M:%S %p"
    )

    total_time = (
        end_time - start_time
    ).total_seconds() + 1  # Adding 1 second to include both times

    # Print the results
    # print(f"Start Time: {start_time}")
    # print(f"End Time: {end_time}")
    print(f"Total time: {total_time} seconds")
    print(f"Average CPU usage (including zeros): {average_cpu_including_zeros:.2f}%")
    print(f"Average CPU usage (excluding zeros): {average_cpu_excluding_zeros:.2f}%")


if __name__ == "__main__":
    # Check if the file path is provided as an argument
    if len(sys.argv) != 2:
        print("Usage: python script_name.py <file_path>")
    else:
        # Pass the file path argument to the function
        get_cpu_usage_and_total_time(sys.argv[1])
