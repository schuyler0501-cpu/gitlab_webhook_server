package utils

import (
	"path/filepath"
	"strings"
)

// DetectLanguage 根据文件扩展名检测编程语言
func DetectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" {
		return "unknown"
	}

	// 移除点号
	ext = strings.TrimPrefix(ext, ".")

	// 语言映射表
	languageMap := map[string]string{
		// Go
		"go": "go",
		// Java
		"java": "java",
		// JavaScript/TypeScript
		"js":   "javascript",
		"jsx":  "javascript",
		"ts":   "typescript",
		"tsx":  "typescript",
		"vue":  "vue",
		// Python
		"py":   "python",
		"pyw":  "python",
		"pyc":  "python",
		// C/C++
		"c":    "c",
		"cpp":  "cpp",
		"cc":   "cpp",
		"cxx":  "cpp",
		"h":    "c",
		"hpp":  "cpp",
		// C#
		"cs":   "csharp",
		// PHP
		"php":  "php",
		// Ruby
		"rb":   "ruby",
		// Swift
		"swift": "swift",
		// Kotlin
		"kt":   "kotlin",
		// Rust
		"rs":   "rust",
		// Shell
		"sh":   "shell",
		"bash": "shell",
		"zsh":  "shell",
		// SQL
		"sql":  "sql",
		// HTML/CSS
		"html": "html",
		"htm":  "html",
		"css":  "css",
		"scss": "css",
		"sass": "css",
		// JSON/YAML
		"json": "json",
		"yaml": "yaml",
		"yml":  "yaml",
		// Markdown
		"md":   "markdown",
		// Docker
		"dockerfile": "dockerfile",
		// Makefile
		"makefile": "makefile",
		"mk":       "makefile",
	}

	if lang, ok := languageMap[ext]; ok {
		return lang
	}

	return "unknown"
}

// GetFileExtension 获取文件扩展名（不含点号）
func GetFileExtension(filePath string) string {
	ext := filepath.Ext(filePath)
	return strings.TrimPrefix(ext, ".")
}

// GetFileName 获取文件名
func GetFileName(filePath string) string {
	return filepath.Base(filePath)
}

