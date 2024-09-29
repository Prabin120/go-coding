package utils

import (
	"fmt"
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

// List of risky imports to check for
var goRiskyImports = []string{
	"os", "syscall", "net", "net/http", "os/exec", "io/ioutil", "unsafe",
	"net/rpc", "net/smtp", "crypto/tls", "bufio", "path/filepath",
	"runtime/debug", "log/syslog", "time", "database/sql",
	"encoding/gob", "encoding/json", "mime/multipart",
	"archive/zip", "archive/tar", "crypto/md5", "crypto/sha1",
	"debug/elf", "debug/pe", "debug/macho", "reflect", "strconv",
	"os/signal", "os/user", "errors", "html/template", "regexp",
	"net/mail", "net/ftp", "crypto/rand",
}

// checkRiskyImports checks if the provided code contains any risky imports
func checkGoRiskyImports(code string) bool {
	// Regex to find all imports in the code
	importRegex := regexp.MustCompile(`import\s*\(\s*([\s\S]*?)\s*\)`)
	matches := importRegex.FindStringSubmatch(code)
	if len(matches) > 1 {
		importBlock := matches[1]
		for _, riskyImport := range goRiskyImports {
			if strings.Contains(importBlock, riskyImport) {
				fmt.Printf("Risky import detected: %s\n", riskyImport)
				return true
			}
		}
	}
	fmt.Println("No risky imports detected")
	return false
}

var pythonRiskyImports = []string{
	"os", "sys", "subprocess", "socket", "threading", "multiprocessing",
	"http.server", "ftplib", "telnetlib", "popen2", "pickle",
	"ctypes", "shutil", "fileinput", "wsgiref", "xmlrpc",
	"asyncio", "requests", "paramiko", "psycopg2", "pymysql",
	"sqlite3", "cryptography", "ssl", "hashlib", "email",
	"xml.etree.ElementTree", "xml.dom", "xml.sax", "configparser",
	"json", "csv", "trace", "inspect", "re", "pandas",
	"numpy", "matplotlib", "tkinter", "sqlite", "os.path",
}

// checkPythonRiskyImports checks if the provided Python code contains any risky imports
func checkPythonRiskyImports(code string) bool {
	// Regex to find all import statements in the code
	importRegex := regexp.MustCompile(`^\s*(import|from)\s+(.*?)(\s+as\s+\w+)?\s*$`)
	lines := strings.Split(code, "\n")

	for _, line := range lines {
		if importRegex.MatchString(line) {
			matches := importRegex.FindStringSubmatch(line)
			// matches[2] contains the imported modules
			importedModules := strings.Split(matches[2], ",")
			for _, module := range importedModules {
				module = strings.TrimSpace(module)
				for _, riskyImport := range pythonRiskyImports {
					if strings.HasPrefix(module, riskyImport) {
						fmt.Printf("Risky import detected: %s\n", riskyImport)
						return true
					}
				}
			}
		}
	}
	fmt.Println("No risky imports detected")
	return false
}

var jsRiskyImports = []string{
	"fs", "child_process", "http", "https", "net", "os", "path", "dns", "process", "vm", "cluster", "socket.io", "express", "mongoose", "mysql", "pg", "mongodb", "axios", "request", "ws", "sharp", "puppeteer", "bcrypt", "jsonwebtoken",
}

// checkJsRiskyImports checks if the provided JavaScript code contains any risky imports
func checkJsRiskyImports(code string) bool {
	// Regex to find all import statements (ES6) and require statements (CommonJS)
	importRegex := regexp.MustCompile(`^\s*import\s+(.*?)(\s+from\s+['"]([^'"]+)['"])?\s*;?`)
	requireRegex := regexp.MustCompile(`^\s*const\s+\w+\s*=\s*require\(['"]([^'"]+)['"]\);?`)

	lines := strings.Split(code, "\n")

	for _, line := range lines {
		// Check for import statements
		if importRegex.MatchString(line) {
			matches := importRegex.FindStringSubmatch(line)
			importedModules := strings.Split(matches[1], ",")
			for _, module := range importedModules {
				module = strings.TrimSpace(module)
				for _, riskyImport := range jsRiskyImports {
					if strings.HasPrefix(module, riskyImport) {
						fmt.Printf("Risky import detected: %s\n", riskyImport)
						return true
					}
				}
			}
		}

		// Check for require statements
		if requireRegex.MatchString(line) {
			matches := requireRegex.FindStringSubmatch(line)
			riskyImport := matches[1] // The imported module from require
			for _, module := range jsRiskyImports {
				if strings.HasPrefix(riskyImport, module) {
					fmt.Printf("Risky import detected: %s\n", riskyImport)
					return true
				}
			}
		}
	}

	fmt.Println("No risky imports detected")
	return false
}

var cRiskyImports = []string{
	"<stdlib.h>", "<unistd.h>", "<fcntl.h>", "<sys/types.h>", "<sys/stat.h>", "<signal.h>", "<exec.h>", "<dirent.h>", "<sys/socket.h>", "<netinet/in.h>", "<arpa/inet.h>", "<errno.h>", "<pthread.h>", "<sys/time.h>",
}

// checkCRiskyImports checks if the provided C code contains any risky imports
func checkCRiskyImports(code string) bool {
	// Regex to find all #include directives in the code
	includeRegex := regexp.MustCompile(`#include\s*["<]([^">]+)[">]`)

	matches := includeRegex.FindAllStringSubmatch(code, -1)
	for _, match := range matches {
		if len(match) > 1 {
			importedModule := match[1]
			for _, riskyImport := range cRiskyImports {
				if strings.Trim(importedModule, "<>") == riskyImport {
					fmt.Printf("Risky import detected: %s\n", riskyImport)
					return true
				}
			}
		}
	}
	fmt.Println("No risky imports detected")
	return false
}

var cppRiskyImports = []string{
	"<fstream>", "<cstdlib>", "<cstdio>", "<cstdio>", "<system>", "<unistd.h>", "<thread>", "<mutex>", "<signal>", "<csignal>", "<functional>", "<locale>", "<thread>", "<chrono>", "<future>", "<atomic>", "<condition_variable>", "networking>", "<sys/socket.h>", "<arpa/inet.h>", "<netinet/in.h>", "<sys/types.h>", "<sys/stat.h>", "<dirent.h>", "<pwd.h>", "<grp.h>", "<signal.h>", "<exec.h>",
}

// checkCPPRiskyImports checks if the provided C++ code contains any risky imports
func checkCPPRiskyImports(code string) bool {
	// Regex to find all #include directives in the code
	includeRegex := regexp.MustCompile(`#include\s*["<]([^">]+)[">]`)

	matches := includeRegex.FindAllStringSubmatch(code, -1)
	for _, match := range matches {
		if len(match) > 1 {
			importedModule := match[1]
			for _, riskyImport := range cppRiskyImports {
				if strings.Trim(importedModule, "<>") == riskyImport {
					fmt.Printf("Risky import detected: %s\n", riskyImport)
					return true
				}
			}
		}
	}
	fmt.Println("No risky imports detected")
	return false
}

// List of risky imports to check for in Java
var javaRiskyImports = []string{
	"java.lang.Runtime", "java.lang.Process", "java.lang.ProcessBuilder",
	"java.io.File", "java.io.FileInputStream", "java.io.FileOutputStream",
	"java.nio.file.Files", "java.nio.file.Paths", "java.net.Socket",
	"java.net.ServerSocket", "java.net.URL", "java.net.HttpURLConnection",
	"java.util.Scanner", "java.util.zip.ZipInputStream",
	"java.util.zip.ZipFile", "java.util.concurrent.ExecutorService",
	"java.util.concurrent.Executors", "java.sql.Connection",
	"java.sql.DriverManager", "java.sql.PreparedStatement",
}

// checkJavaRiskyImports checks if the provided Java code contains any risky imports
func checkJavaRiskyImports(code string) bool {
	// Regex to find all import statements in Java
	importRegex := regexp.MustCompile(`^\s*import\s+([a-zA-Z0-9._*]+)\s*;`)

	lines := strings.Split(code, "\n")

	for _, line := range lines {
		if importRegex.MatchString(line) {
			matches := importRegex.FindStringSubmatch(line)
			if len(matches) > 1 {
				importedModule := matches[1]
				for _, riskyImport := range javaRiskyImports {
					if importedModule == riskyImport {
						fmt.Printf("Risky import detected: %s\n", riskyImport)
						return true
					}
				}
			}
		}
	}
	fmt.Println("No risky imports detected")
	return false
}
