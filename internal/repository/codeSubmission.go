package repository

import (
	"bufio"
	"bytes"
	commontypes "code-compiler/internal/commonTypes"
	"code-compiler/internal/models"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// CodeRunner struct to execute code
type CodeRunner struct {
	Question *Question // Add a reference to Question
}

// SupportedLanguages maps languages to their compilation commands
type SupportedLanguages map[string]string

// GetSupportedLanguages returns a map of supported languages and their compilation commands
func getSupportedLanguages() SupportedLanguages {
	return SupportedLanguages{
		"python": "python3",
		"cpp":    "g++",
		"c":      "gcc",
		"java":   "javac",
		"go":     "go build -o", // Build command for Go
		"js":     "node",        // Execution command for JavaScript
	}
}

// CompileCode compiles code based on the language
func compileCode(codePath, language string) (string, error) {
	fileBaseName := filepath.Base(codePath)
	fileBaseNameWithoutExt := fileBaseName[:len(fileBaseName)-len(filepath.Ext(fileBaseName))]

	var cmd *exec.Cmd
	var outputFileName string
	switch language {
	case "python", "js": // Python and JS don't need compilation
		return codePath, nil
	case "c":
		outputFileName = fileBaseNameWithoutExt + ".out"
		cmd = exec.Command("gcc", codePath, "-o", outputFileName)
	case "cpp":
		outputFileName = fileBaseNameWithoutExt + ".out"
		cmd = exec.Command("g++", codePath, "-o", outputFileName)
	case "java":
		outputFileName = fileBaseNameWithoutExt + ".class"
		cmd = exec.Command("javac", codePath)
	case "go":
		outputFileName = fileBaseNameWithoutExt + ".out"
		cmd = exec.Command("go", "build", "-o", outputFileName, codePath)
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("compilation failed: %v", err)
	}
	return outputFileName, nil
}

// RunTestCases executes the compiled code with the provided test cases
func runTestCases(compiledFilePath string, testCases []models.InputOutput, language string) ([]commontypes.TestResult, error) {
	var results []commontypes.TestResult

	for i, testCase := range testCases {
		fmt.Printf("Executing Test Case %d: Input: %s, Expected Output: %s\n", i+1, testCase.Input, testCase.Output)

		cmd, err := getCommandForLanguage(compiledFilePath, language)
		if err != nil {
			return results, err
		}

		cmd.Stdin = bytes.NewBufferString(testCase.Input)

		var outputBytes []byte
		errChan := make(chan error, 1)
		go func() {
			output, err := cmd.CombinedOutput()
			outputBytes = output
			errChan <- err
		}()

		select {
		case err := <-errChan:
			if err != nil {
				return results, fmt.Errorf("failed to execute test case %d: %v", i+1, err)
			}
		case <-time.After(2 * time.Second):
			return results, fmt.Errorf("test case %d timed out", i+1)
		}

		actualOutput := string(bytes.TrimSpace(outputBytes))
		expectedOutput := testCase.Output
		passed := actualOutput == expectedOutput

		results = append(results, commontypes.TestResult{
			TestCaseNumber: i + 1,
			Input:          testCase.Input,
			ExpectedOutput: expectedOutput,
			ActualOutput:   actualOutput,
			Passed:         passed,
		})
	}

	return results, nil
}

// RunAllTestCases runs test cases and stops on the first failure
func runAllTestCases(compiledFilePath string, testCases []models.InputOutput, language string) (*commontypes.TestResult, error, int) {
	numberOfPassedTests := 0

	for i, testCase := range testCases {
		fmt.Printf("Executing Test Case %d: Input: %s, Expected Output: %s\n", i+1, testCase.Input, testCase.Output)

		cmd, err := getCommandForLanguage(compiledFilePath, language)
		if err != nil {
			return nil, err, numberOfPassedTests
		}

		cmd.Stdin = bytes.NewBufferString(testCase.Input)

		var outputBytes []byte
		errChan := make(chan error, 1)
		go func() {
			output, err := cmd.CombinedOutput()
			outputBytes = output
			errChan <- err
		}()

		select {
		case err := <-errChan:
			if err != nil {
				return nil, fmt.Errorf("failed to execute test case %d: %v", i+1, err), numberOfPassedTests
			}
		case <-time.After(2 * time.Second):
			return nil, fmt.Errorf("test case %d timed out", i+1), numberOfPassedTests
		}

		actualOutput := string(bytes.TrimSpace(outputBytes))
		expectedOutput := testCase.Output
		passed := actualOutput == expectedOutput

		if !passed {
			return &commontypes.TestResult{
				TestCaseNumber: i + 1,
				Input:          testCase.Input,
				ExpectedOutput: expectedOutput,
				ActualOutput:   actualOutput,
				Passed:         false,
			}, nil, numberOfPassedTests
		}

		numberOfPassedTests++
	}

	return nil, nil, numberOfPassedTests
}

// FileWriter writes the code to a file
func fileWriter(code string, language string) string {
	filename := fmt.Sprintf("/codeFiles/%s%s", language, uuid.NewString())
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return ""
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.WriteString(code); err != nil {
		fmt.Println("Error writing to file:", err)
		return ""
	}
	writer.Flush()

	fmt.Println("File writing successful")
	return filename
}

// FileRemoving removes the file after execution
func fileRemoving(filename string) bool {
	if err := os.Remove(filename); err != nil {
		fmt.Println("Error removing file:", err)
		return false
	}
	return true
}

// GetCommandForLanguage returns the exec command based on the language
func getCommandForLanguage(compiledFilePath, language string) (*exec.Cmd, error) {
	switch language {
	case "python":
		return exec.Command("python3", compiledFilePath), nil
	case "cpp", "c":
		return exec.Command(compiledFilePath), nil
	case "java":
		return exec.Command("java", compiledFilePath), nil
	case "go", "js":
		return exec.Command(compiledFilePath), nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
}

// Execute runs the code for either testing or submission
func (r *CodeRunner) ExecuteTest(data commontypes.CodeRunnerType) ([]commontypes.TestResult, error) {
	compiledFilePath := fileWriter(data.Code, data.Language)
	if compiledFilePath == "" {
		return nil, fmt.Errorf("file creation failed")
	}
	question, err := r.Question.GetQuestionById(data.QuestionId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve question: %v", err)
	}

	results, err := runTestCases(compiledFilePath, question.SampleTestCases, data.Language)
	if err != nil {
		return nil, err
	}

	if !fileRemoving(compiledFilePath) {
		fmt.Println("File not deleted")
	}
	return results, nil
}

func (r *CodeRunner) ExecuteSubmit(data commontypes.CodeRunnerType) (*commontypes.TestResult, error, int, int) {
	compiledFilePath := fileWriter(data.Code, data.Language)
	if compiledFilePath == "" {
		return nil, fmt.Errorf("file creation failed"), 0, 0
	}
	question, err := r.Question.GetQuestionById(data.QuestionId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve question: %v", err), 0, 0
	}

	// var passed *commontypes.TestResult
	numberOfPassedTests := 0
	totalTestCases := 0
	var failedCase *commontypes.TestResult

	testCases, err := r.Question.GetTestCases(question.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve test cases: %v", err), 0, totalTestCases
	}
	totalTestCases = len(testCases.IOPairs)
	failedCase, err, numberOfPassedTests = runAllTestCases(compiledFilePath, testCases.IOPairs, data.Language)
	if err != nil {
		return failedCase, err, numberOfPassedTests, totalTestCases
	}

	if !fileRemoving(compiledFilePath) {
		fmt.Println("File not deleted")
	}

	return failedCase, nil, numberOfPassedTests, totalTestCases
}
