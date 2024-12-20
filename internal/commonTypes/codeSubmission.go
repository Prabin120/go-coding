package commontypes

// CodeRunnerType represents the structure of code runner input.
type CodeRunnerType struct {
	UserId	   string `json:"userId"`
	Language   string `json:"language"`
	Code       string `json:"code"`
	QuestionId string `json:"questionId"`
}

// TestResult represents the result of executing a test case.
type TestResult struct {
	TestCaseNumber int     `json:"testCaseNumber"`
	Input          string  `json:"input"`
	ExpectedOutput string  `json:"expectedOutput"`
	ActualOutput   string  `json:"actualOutput"`
	Passed         bool    `json:"passed"`
	TimeTaken      float64 `json:"timeTaken"`
	MemoryUsed     float64 `json:"memoryUsed"`
}

// InputOutput represents the input and output for test cases.
type InputOutput struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}
