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

def run_load_stepwell(num_cores):
    result = subprocess.run([directory_path + '/../exec', "testStepWellLoad", str(num_cores)], capture_output=True, text=True)
    output = result.stdout
    expected_tokens, actual_tokens = parse_output(output)
    return expected_tokens, actual_tokens

def run_load_tokenbucket(num_cores):
    result = subprocess.run([directory_path + '/../exec', "TestTokenBucketLoad", str(num_cores)], capture_output=True, text=True)
    output = result.stdout
    expected_tokens, actual_tokens = parse_output(output)
    return expected_tokens, actual_tokens

def parse_output(output):
    # Split the output by lines and parse tokens information
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
    
    cores = [4, 8, 16, 32, 64]
    results_tokenbucket = []
    results_stepwell = []

    # Process TokenBucket results
    for num_cores in cores:
        expected_tokens, actual_tokens = run_load_tokenbucket(num_cores)
        print(f"TokenBucket - Expected Number of Tokens: {expected_tokens}, Actual Number of Tokens: {actual_tokens}")
        percentage_excess = (actual_tokens / expected_tokens) * 100
        results_tokenbucket.append(percentage_excess)
    
    # Plotting for TokenBucket
    plt.figure(figsize=(10, 5))
    plt.plot(cores, results_tokenbucket, marker='o')
    plt.xlabel('Number of Cores')
    plt.ylabel('Percentage of the Expected Tokens Issued')
    plt.title('TokenBucket Performance Analysis with Varying Cores')
    plt.grid(True)
    file_name = "performance_analysis_tokenbucket.png"
    file_path = os.path.join(directory_path, file_name)
    plt.savefig(file_path, format='png', dpi=300)
    print(f"TokenBucket plot saved to {file_path}")
    
    # Process StepWell results
    for num_cores in cores:
        expected_tokens, actual_tokens = run_load_stepwell(num_cores)
        print(f"StepWell - Expected Number of Tokens: {expected_tokens}, Actual Number of Tokens: {actual_tokens}")
        percentage_excess = (actual_tokens / expected_tokens) * 100
        results_stepwell.append(percentage_excess)

    # Plotting for StepWell
    plt.figure(figsize=(10, 5))
    plt.plot(cores, results_stepwell, marker='o')
    plt.xlabel('Number of Cores')
    plt.ylabel('Percentage of the Expected Tokens Issued')
    plt.title('StepWell Performance Analysis with Varying Cores')
    plt.grid(True)
    file_name = "performance_analysis_stepwell.png"
    file_path = os.path.join(directory_path, file_name)
    plt.savefig(file_path, format='png', dpi=300)
    print(f"StepWell plot saved to {file_path}")

if __name__ == "__main__":
    main()