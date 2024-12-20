package utils

import (
	"regexp"
	"strings"
)

func CheckRiskyImports(code string, language string) bool {
	switch language {
	case "go":
		return checkGoRiskyImports(code)
	case "py":
		return checkPythonRiskyImports(code)
	case "js":
		return checkJsRiskyImports(code)
	case "c":
		return checkCRiskyImports(code)
	case "cpp":
		return checkCPPRiskyImports(code)
	case "java":
		return checkJavaRiskyImports(code)
	default:
		return false
	}
}

// Function to check if Go code contains imports
func checkGoRiskyImports(code string) bool {
	// Use a regular expression to find any 'import' keyword
	match, _ := regexp.MatchString(`\bimport\b`, code)
	return match
}

// Function to check if Python code contains imports
func checkPythonRiskyImports(code string) bool {
	// Check if the code contains 'import' or 'from'
	if strings.Contains(code, "import") || strings.Contains(code, "from") {
		return true
	}
	return false
}

// Function to check if JavaScript code contains imports
func checkJsRiskyImports(code string) bool {
	// Check for 'import' keyword (ES6 modules)
	if strings.Contains(code, "import") {
		return true
	}
	// Check for 'require' keyword (CommonJS)
	if strings.Contains(code, "require") {
		return true
	}
	// No risky imports found
	return false
}

// Function to check if C code contains includes (C doesn't use `import`, it uses `#include`)
func checkCRiskyImports(code string) bool {
	// Check if the code contains '#include' directive
	if strings.Contains(code, "#include") {
		return true
	}
	return false
}

// Function to check if C++ code contains includes (C++ uses `#include` like C)
func checkCPPRiskyImports(code string) bool {
	// Check if the code contains '#include' directive
	if strings.Contains(code, "#include") {
		return true
	}
	return false
}

// Function to check if Java code contains imports
func checkJavaRiskyImports(code string) bool {
	// Check if the code contains 'import' keyword
	if strings.Contains(code, "import") {
		return true
	}
	return false
}