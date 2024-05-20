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

def run_load_tests(test_function, num_cores, bucket_type, duration, refill_rate, capacity):
    results = []
    for _ in range(3):
        expected_tokens, actual_tokens = test_function(num_cores, bucket_type, duration, refill_rate, capacity)
        percentage_excess = (actual_tokens / expected_tokens) * 100
        results.append(percentage_excess)
    mean_percentage = np.mean(results)
    std_deviation = np.std(results)
    return mean_percentage, std_deviation

def run_load_stepwell(num_cores, bucket_type, duration, refill_rate, capacity):
    args = [directory_path + '/../exec', "TestStepWellLoad", str(num_cores), str(bucket_type), str(duration), str(refill_rate), str(capacity)]
    result = subprocess.run(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True)
    output = result.stdout
    expected_tokens, actual_tokens = parse_output(output)
    return expected_tokens, actual_tokens

def run_load_tokenbucket(num_cores, bucket_type, duration, refill_rate, capacity):
    args = [directory_path + '/../exec', "TestTokenBucketLoad", str(num_cores), str(bucket_type), str(duration), str(refill_rate), str(capacity)]
    result = subprocess.run(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True)
    output = result.stdout
    expected_tokens, actual_tokens = parse_output(output)
    return expected_tokens, actual_tokens

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
    bucket_type = 1
    duration = 10 #number of seconds in this test
    refill_rate = 10
    capacity = 10
    results_tokenbucket = []
    errors_tokenbucket = []
    results_stepwell = []
    errors_stepwell = []

    for num_cores in cores:
        avg_percentage, std_dev = run_load_tests(run_load_tokenbucket, num_cores, bucket_type, duration, refill_rate, capacity)
        results_tokenbucket.append(avg_percentage)
        errors_tokenbucket.append(std_dev)
        print(f"TokenBucket - Cores: {num_cores}, Avg Percentage: {avg_percentage}%, Std Dev: {std_dev}%")

    for num_cores in cores:
        avg_percentage, std_dev = run_load_tests(run_load_stepwell, num_cores, bucket_type, duration, refill_rate, capacity)
        results_stepwell.append(avg_percentage)
        errors_stepwell.append(std_dev)
        print(f"StepWell - Cores: {num_cores}, Avg Percentage: {avg_percentage}%, Std Dev: {std_dev}%")

    plt.figure(figsize=(10, 5))
    plt.errorbar(cores, results_tokenbucket, yerr=errors_tokenbucket, label='TokenBucket', marker='o', capsize=5, color='blue')
    plt.errorbar(cores, results_stepwell, yerr=errors_stepwell, label='StepWell', marker='x', capsize=5, color='green')
    plt.xlabel('Number of Cores')
    plt.ylabel('Percentage of the Max Amount of Tokens Issued')
    plt.ylim(0, None)
    plt.title('High Load Analysis with Varying Cores')
    plt.figtext(0.5, 0.007, f'Runtime: {duration}, Refill Rate: {refill_rate}, Token Bucket Capacity: {capacity}', ha="center", fontsize=9, style='italic')
    plt.xticks(cores)
    plt.legend()
    plt.grid(True)
    file_name = "high_load_analysis_combined.png"
    file_path = os.path.join(directory_path, file_name)
    plt.savefig(file_path, format='png', dpi=300)
    print(f"Combined plot with error bars saved to {file_path}")

if __name__ == "__main__":
    main()