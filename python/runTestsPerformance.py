import subprocess
import matplotlib.pyplot as plt
import os

directory_path = os.path.dirname(os.path.abspath(__file__))

def compile_go_executable(source_path, output_name):
    command = ['go', 'build', '-o', output_name, source_path]
    result = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    
    if result.returncode != 0:
        print("Error compiling Go executable:")
        print(result.stderr)
        exit(1)
    else:
        print(f"Successfully compiled {output_name}.")
        
        
def run_performance_test(executable_name, test_type, num_cores):
    results = []
    for _ in range(3):
        result = subprocess.run([executable_name, test_type, str(num_cores)], capture_output=True, text=True)
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
    
    # Plotting the execution times
    plt.figure(figsize=(10, 5))
    plt.bar(cores, results_tokenbucket, color='blue')
    plt.xlabel('Number of Cores')
    plt.ylabel('Execution Time (ns)')
    plt.title('TokenBucket Performance Analysis by Core Count')
    plt.xticks(cores)  # Ensure we have a tick for every core count
    plt.grid(True, which='both', linestyle='--', linewidth=0.5)
    file_name = "tokenbucket_performance_times.png"
    file_path = os.path.join(directory_path, file_name)
    plt.savefig(file_path, format='png', dpi=300)
    print(f"Performance plot saved to {file_path}")
    
    for num_cores in cores:
        execution_time = run_performance_test(executable_name, "TestStepWellPerformance", num_cores)
        results_stepwell.append(execution_time)
        print(f"Stepwell - Execution Time for {num_cores} cores: {execution_time} ns")
    
    # Plotting the execution times
    plt.figure(figsize=(10, 5))
    plt.bar(cores, results_stepwell, color='blue')
    plt.xlabel('Number of Cores')
    plt.ylabel('Execution Time (ns)')
    plt.title('Stepwell Performance Analysis by Core Count')
    plt.xticks(cores)  # Ensure we have a tick for every core count
    plt.grid(True, which='both', linestyle='--', linewidth=0.5)
    file_name = "stepwell_performance_times.png"
    file_path = os.path.join(directory_path, file_name)
    plt.savefig(file_path, format='png', dpi=300)
    print(f"Performance plot saved to {file_path}")

if __name__ == "__main__":
    main()