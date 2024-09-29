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

// CompileCode compiles code based on the language
func compileCode(codePath, language string) (string, error) {
	fileBaseName := filepath.Base(codePath)
	fileBaseNameWithoutExt := fileBaseName[:len(fileBaseName)-len(filepath.Ext(fileBaseName))]
	var cmd *exec.Cmd
	var outputFileName string

	// Determine the output filename based on the language
	switch language {
	case "py", "js": // Python and JS don't need compilation
		return codePath, nil
	case "c":
		outputFileName = filepath.Join(filepath.Dir(codePath), fileBaseNameWithoutExt+".out")
		cmd = exec.Command("gcc", codePath, "-o", outputFileName)
	case "cpp":
		outputFileName = filepath.Join(filepath.Dir(codePath), fileBaseNameWithoutExt+".out")
		cmd = exec.Command("g++", codePath, "-o", outputFileName)
	case "java":
		outputFileName = filepath.Join(filepath.Dir(codePath), fileBaseNameWithoutExt+".class")
		cmd = exec.Command("javac", codePath)
	case "go":
		outputFileName = filepath.Join(filepath.Dir(codePath), fileBaseNameWithoutExt+".out")
		cmd = exec.Command("go", "build", "-o", outputFileName, codePath)
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	// Run the command and check for errors
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("compilation failed: %v", err)
	}
	fileRemoving(codePath)
	return outputFileName, nil
}

// GetCommandForLanguage returns the exec command based on the language
func getCommandForLanguage(compiledFilePath, language string) (*exec.Cmd, error) {
	switch language {
	case "py":
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
func runAllTestCases(compiledFilePath string, testCases []models.InputOutput, language string) (*commontypes.TestResult, int, error) {
	numberOfPassedTests := 0

	for i, testCase := range testCases {
		fmt.Printf("Executing Test Case %d: Input: %s, Expected Output: %s\n", i+1, testCase.Input, testCase.Output)

		cmd, err := getCommandForLanguage(compiledFilePath, language)
		if err != nil {
			return nil, numberOfPassedTests, err
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
				return nil, numberOfPassedTests, fmt.Errorf("failed to execute test case %d: %v", i+1, err)
			}
		case <-time.After(2 * time.Second):
			return nil, numberOfPassedTests, fmt.Errorf("test case %d timed out", i+1)
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
			}, numberOfPassedTests, nil
		}

		numberOfPassedTests++
	}

	return nil, numberOfPassedTests, nil
}

// FileWriter writes the code to a file
func fileWriter(code string, language string) string {
	filename := fmt.Sprintf("./codeFiles/%s%s.%s", language, uuid.NewString(), language)
	println(filename)
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

// Execute runs the code for either testing or submission
func (r *CodeRunner) ExecuteTest(data commontypes.CodeRunnerType) ([]commontypes.TestResult, error) {
	codeFileath := fileWriter(data.Code, data.Language)
	if codeFileath == "" {
		return nil, fmt.Errorf("file creation failed")
	}
	question, err := r.Question.GetQuestionById(data.QuestionId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve question: %v", err)
	}
	compiledFilePath, err := compileCode(codeFileath, data.Language)
	if err != nil {
		return nil, fmt.Errorf("code compilation failed: %v", err)
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

func (r *CodeRunner) ExecuteSubmit(data commontypes.CodeRunnerType) (*commontypes.TestResult, int, int, error) {
	compiledFilePath := fileWriter(data.Code, data.Language)
	if compiledFilePath == "" {
		return nil, 0, 0, fmt.Errorf("file creation failed")
	}

	question, err := r.Question.GetQuestionById(data.QuestionId)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to retrieve question: %v", err)
	}
	// var passed *commontypes.TestResult
	numberOfPassedTests := 0
	totalTestCases := 0
	var failedCase *commontypes.TestResult

	testCases, err := r.Question.GetTestCases(question.ID)
	if err != nil {
		return nil, 0, totalTestCases, fmt.Errorf("failed to retrieve test cases: %v", err)
	}
	totalTestCases = len(testCases.IOPairs)
	failedCase, numberOfPassedTests, err = runAllTestCases(compiledFilePath, testCases.IOPairs, data.Language)

	if !fileRemoving(compiledFilePath) {
		fmt.Println("File not deleted")
	}

	if err != nil {
		return failedCase, numberOfPassedTests, totalTestCases, err
	}

	return failedCase, numberOfPassedTests, totalTestCases, nil
}
