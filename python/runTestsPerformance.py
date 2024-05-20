import subprocess
import matplotlib.pyplot as plt
import numpy as np
import os

directory_path = os.path.dirname(os.path.abspath(__file__))

def get_go_binary_path():
    config_file = os.path.join(directory_path, "go_path.conf")
    if os.path.exists(config_file):
        with open(config_file, "r") as file:
            go_binary = file.readline().strip()
    else:
        go_binary = "go"
    return go_binary

def compile_go_executable(source_path, output_name):
    go_binary = get_go_binary_path()
    command = [go_binary, 'build', '-o', output_name, source_path]
    result = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True)
    
    if result.returncode != 0:
        print("Error compiling Go executable:")
        print(result.stderr)
        exit(1)
    else:
        print(f"Successfully compiled {output_name}.")

def run_performance_test(executable_name, test_type, num_cores, bucket_type, duration, refill_rate, capacity):
    results = []
    for _ in range(20):
        args = [executable_name, test_type, str(num_cores), str(bucket_type), str(duration), str(refill_rate), str(capacity)]
        result = subprocess.run(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True)
        output = result.stdout
        results.append(parse_performance_output(output) / float(duration))  # Normalize by duration
    return results

def parse_performance_output(output):
    lines = output.split('\n')
    for line in lines:
        if "Time" in line:
            time_ns = int(line.split(':')[-1].strip())
            return time_ns
    return 0

def main():
    go_source_path = directory_path + "/../main.go"
    executable_name = directory_path + "/../exec"
    
    compile_go_executable(go_source_path, executable_name)
    
    cores = [1, 2, 4, 8, 32, 64]
    duration = 1000000 #number of requests in this test
    refill_rate = 100
    capacity = 10
    bucket_types = [1,2,3,4]
    bucket_labels = {
        1: "Tokenbucket Trivial",
        2: "Tokenbucket Atomic with Loops",
        3: "Tokenbucket with Locks",
        4: "Tokenbucket Helia"
    }
    results_stepwell = []
    errors_stepwell = []

    # Run StepWell tests once per core count
    for num_cores in cores:
        results_sw = run_performance_test(executable_name, "TestStepWellPerformance", num_cores, 1, duration, refill_rate, capacity)
        mean_sw = np.mean(results_sw)
        std_sw = np.std(results_sw)
        results_stepwell.append(mean_sw)
        errors_stepwell.append(std_sw)
        print(f"StepWell Performance {num_cores} cores: {mean_sw:.3f} ns/request ± {std_sw:.3f}")

    # Run TokenBucket tests for each bucket type
    for bucket_type in bucket_types:
        label = bucket_labels[bucket_type]
        results_tokenbucket = []
        errors_tokenbucket = []

        for num_cores in cores:
            results_tb = run_performance_test(executable_name, "TestTokenBucketPerformance", num_cores, bucket_type, duration, refill_rate, capacity)
            mean_tb = np.mean(results_tb)
            std_tb = np.std(results_tb)
            results_tokenbucket.append(mean_tb)
            errors_tokenbucket.append(std_tb)
            print(f"{label} Performance {num_cores} cores: {mean_tb:.3f} ns/request ± {std_tb:.3f}")

        plt.figure(figsize=(10, 5))
        plt.errorbar(cores, results_tokenbucket, yerr=errors_tokenbucket, label=label, marker='o', color='blue', capsize=5)
        plt.errorbar(cores, results_stepwell, yerr=errors_stepwell, label='StepWell w/ Trivial Tokenbucket', marker='x', color='green', capsize=5)
        plt.xlabel('Number of Cores')
        plt.ylabel('Execution Time per Request (ns/request)')
        plt.title(f'Performance Analysis by Core Count - {label}')
        plt.figtext(0.5, 0.007, f'Number of Requests: {duration}, Refill Rate: {refill_rate}, Token Bucket Capacity: {capacity}', ha="center", fontsize=9, style='italic')
        plt.xticks(cores)
        plt.grid(True, which='both', linestyle='--', linewidth=0.5)
        if bucket_type == 3:  # Apply log scale only for "Tokenbucket with Locks"
            plt.yscale('log')
        plt.legend()
        file_name = f"performance_comparison_{bucket_type}.png"
        file_path = os.path.join(directory_path, file_name)
        plt.savefig(file_path, format='png', dpi=300)
        plt.close()
        print(f"Performance comparison plot for {label} saved to {file_path}")

if __name__ == "__main__":
    main()