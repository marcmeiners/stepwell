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

def run_load_tests(executable_name, test_type, num_cores, bucket_type, duration, refill_rate, capacity):
    results = []
    for _ in range(3):
        args = [executable_name, test_type, str(num_cores), str(bucket_type), str(duration), str(refill_rate), str(capacity)]
        result = subprocess.run(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True)
        output = result.stdout
        expected_tokens, actual_tokens = parse_output(output)
        percentage_excess = (actual_tokens / expected_tokens) * 100
        results.append(percentage_excess)
    mean_percentage = np.mean(results)
    std_deviation = np.std(results)
    return mean_percentage, std_deviation

def parse_output(output):
    lines = output.split('\n')
    for line in lines:
        if "Expected" in line and "Actual" in line:
            parts = line.split()
            expected_tokens = float(parts[1])
            actual_tokens = int(parts[3])
            return expected_tokens, actual_tokens
    return 0.0, 0  # Return defaults if not found

def main():
    go_source_path = directory_path + "/../main.go"
    executable_name = directory_path + "/../exec"
    
    compile_go_executable(go_source_path, executable_name)
    
    cores = [1, 2, 4, 8, 32, 64]
    duration = 10 # number of seconds in this test
    refill_rate = 10
    capacity = 10
    bucket_types = [1, 2, 3, 4]
    bucket_labels = {
        1: "Tokenbucket Trivial",
        2: "Tokenbucket Atomic with Loops",
        3: "Tokenbucket with Locks",
        4: "Tokenbucket Helia"
    }
    
    # Run StepWell tests once per core count
    results_stepwell = []
    errors_stepwell = []
    for num_cores in cores:
        mean_sw, std_sw = run_load_tests(executable_name, "TestStepWellPerformance", num_cores, 1, duration, refill_rate, capacity)
        results_stepwell.append(mean_sw)
        errors_stepwell.append(std_sw)
        print(f"StepWell Performance {num_cores} cores: {mean_sw:.3f} % ± {std_sw:.3f}")

    # Run TokenBucket tests for each bucket type and generate plots
    for bucket_type in bucket_types:
        label = bucket_labels[bucket_type]
        results_tokenbucket = []
        errors_tokenbucket = []

        for num_cores in cores:
            mean_tb, std_tb = run_load_tests(executable_name, "TestTokenBucketPerformance", num_cores, bucket_type, duration, refill_rate, capacity)
            results_tokenbucket.append(mean_tb)
            errors_tokenbucket.append(std_tb)
            print(f"{label} Performance {num_cores} cores: {mean_tb:.3f} % ± {std_tb:.3f}")

        plt.figure(figsize=(10, 5))
        plt.errorbar(cores, results_tokenbucket, yerr=errors_tokenbucket, label=label, marker='o', color='blue', capsize=5)
        plt.errorbar(cores, results_stepwell, yerr=errors_stepwell, label='StepWell w/ Trivial Tokenbucket', marker='x', color='green', capsize=5)
        plt.xlabel('Number of Cores')
        plt.ylabel('Percentage of the Max Amount of Tokens Issued')
        plt.ylim(0, None)
        plt.title(f'High Load Analysis with Varying Cores - {label}')
        plt.figtext(0.5, 0.85, f'Runtime: {duration}s, Refill Rate: {refill_rate}, Token Bucket Capacity: {capacity}', ha="center", fontsize=9, style='italic')
        plt.xticks(cores)
        plt.grid(True, which='both', linestyle='--', linewidth=0.5)
        plt.legend()
        file_name = f"performance_comparison_{bucket_type}.png"
        file_path = os.path.join(directory_path, file_name)
        plt.savefig(file_path, format='png', dpi=300)
        plt.close()
        print(f"High Load Analysis plot for {label} saved to {file_path}")

if __name__ == "__main__":
    main()