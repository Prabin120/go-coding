package usecases

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// Function to compile and run code based on language
func compileAndRun(language, codePath string) (*exec.Cmd, error) {
	var cmd *exec.Cmd

	switch language {
	case "python":
		// For Python, no compilation is needed, return cmd to run the script
		cmd = exec.Command("python3", codePath)
	case "cpp":
		// Compile the C++ program first
		compileCmd := exec.Command("g++", codePath, "-o", "a.out")
		if err := compileCmd.Run(); err != nil {
			return nil, fmt.Errorf("compilation error: %v", err)
		}
		// After successful compilation, set cmd to run the compiled binary
		cmd = exec.Command("./a.out")
		// Add more languages if necessary...
	default:
		return nil, fmt.Errorf("language %s not supported", language)
	}

	return cmd, nil
}

// Function to run the compiled code with test case input
func runTestCase(cmd *exec.Cmd, testCase string, timeout time.Duration) (string, error) {
	// Set up a context with timeout to limit execution time
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Use a buffer to capture the program's output
	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer
	cmd.Stdin = bytes.NewBuffer([]byte(testCase)) // Pass the test case input

	// Run the program with the provided test case input
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("execution start error: %v", err)
	}

	err := cmd.Wait()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("execution timed out")
	}
	if err != nil {
		return "", fmt.Errorf("execution error: %v", err)
	}

	return outputBuffer.String(), nil
}

// Function to run multiple test cases
func runAllTestCases(language, codePath string, testCases []string, expectedOutputs []string) {
	// Compile and run the code
	cmd, err := compileAndRun(language, codePath)
	if err != nil {
		fmt.Println("Error during compilation:", err)
		return
	}

	// Create a channel to capture the results from each test case
	type result struct {
		testCase       string
		output         string
		expectedOutput string
		passed         bool
	}

	results := make(chan result, len(testCases))

	// Run each test case concurrently
	for i, testCase := range testCases {
		go func(i int, testCase string) {
			output, err := runTestCase(cmd, testCase, 2*time.Second) // Set max execution time per test case
			if err != nil {
				fmt.Printf("Test case %d failed with error: %v\n", i+1, err)
				results <- result{testCase, output, expectedOutputs[i], false}
				return
			}
			// Compare output with expected output
			passed := output == expectedOutputs[i]
			results <- result{testCase, output, expectedOutputs[i], passed}
		}(i, testCase)
	}

	// Collect and print results
	for i := 0; i < len(testCases); i++ {
		res := <-results
		fmt.Printf("Test case %d: Passed: %v\n", i+1, res.passed)
		fmt.Printf("Expected Output: %s\nGot Output: %s\n", res.expectedOutput, res.output)
	}
}

func submit() {
	// Example test cases and expected outputs
	testCases := []string{
		"input1",
		"input2",
		"input3",
	}
	expectedOutputs := []string{
		"output1",
		"output2",
		"output3",
	}

	// Run the test cases on the code
	runAllTestCases("python", "code.py", testCases, expectedOutputs)
}
