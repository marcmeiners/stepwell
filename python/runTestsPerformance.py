import subprocess
import matplotlib.pyplot as plt
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
        
        
def run_performance_test(executable_name, test_type, num_cores):
    results = []
    for _ in range(3):
        result = subprocess.run([executable_name, test_type, str(num_cores)], stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True)
        output = result.stdout
        results.append(parse_performance_output(output))
    return sum(results) / len(results)

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
    
    cores = [4, 8, 16, 32, 64]
    results_tokenbucket = []
    results_stepwell = []

    for num_cores in cores:
        execution_time = run_performance_test(executable_name, "TestTokenBucketPerformance", num_cores)
        results_tokenbucket.append(execution_time)
        print(f"TokenBucket - Execution Time for {num_cores} cores: {execution_time} ns")
    
    for num_cores in cores:
        execution_time = run_performance_test(executable_name, "TestStepWellPerformance", num_cores)
        results_stepwell.append(execution_time)
        print(f"Stepwell - Execution Time for {num_cores} cores: {execution_time} ns")
    
    # Plotting the execution times together
    plt.figure(figsize=(10, 5))
    plt.plot(cores, results_tokenbucket, label='TokenBucket', marker='o', color='blue')
    plt.plot(cores, results_stepwell, label='StepWell', marker='x', color='green')
    plt.xlabel('Number of Cores')
    plt.ylabel('Execution Time (ns)')
    plt.title('Performance Analysis by Core Count')
    plt.xticks(cores)
    plt.grid(True, which='both', linestyle='--', linewidth=0.5)
    plt.legend()
    file_name = "performance_comparison.png"
    file_path = os.path.join(directory_path, file_name)
    plt.savefig(file_path, format='png', dpi=300)
    print(f"Performance comparison plot saved to {file_path}")

if __name__ == "__main__":
    main()