package repo

import (
	"regexp"
	"strings"
)

type Chunk struct {
	Content  string
	FilePath string
	Language string
	FuncName string // optional
}

// Function regex patterns per language
var funcPatterns = map[string]*regexp.Regexp{
	"go": regexp.MustCompile(`(?m)^func\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(`),
	"py": regexp.MustCompile(`(?m)^def\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(`),
	"js": regexp.MustCompile(`(?m)^function\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(`),
	"ts": regexp.MustCompile(`(?m)^function\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(`),
}

func ChunkFile(file FileData, maxLines int) []Chunk {
	lines := strings.Split(file.Content, "\n")
	chunks := make([]Chunk, 0)

	// If no pattern for this language, create simple chunks
	pattern, hasPattern := funcPatterns[file.Language]
	if !hasPattern {
		// Simple chunking by maxLines for unsupported languages
		for i := 0; i < len(lines); i += maxLines {
			end := i + maxLines
			if end > len(lines) {
				end = len(lines)
			}
			content := strings.Join(lines[i:end], "\n")
			if strings.TrimSpace(content) != "" {
				chunks = append(chunks, Chunk{
					Content:  content,
					FilePath: file.Path,
					Language: file.Language,
					FuncName: "",
				})
			}
		}
		return chunks
	}

	start := 0
	funcName := ""

	for i, line := range lines {
		// Detect function start
		m := pattern.FindStringSubmatch(line)
		if len(m) == 2 {
			// Save previous chunk if it has content
			if start < i {
				content := strings.Join(lines[start:i], "\n")
				if strings.TrimSpace(content) != "" {
					chunks = append(chunks, Chunk{
						Content:  content,
						FilePath: file.Path,
						Language: file.Language,
						FuncName: funcName,
					})
				}
			}
			start = i
			funcName = m[1]
		}

		// Fallback: maxLines chunking
		if i-start+1 >= maxLines {
			content := strings.Join(lines[start:i+1], "\n")
			if strings.TrimSpace(content) != "" {
				chunks = append(chunks, Chunk{
					Content:  content,
					FilePath: file.Path,
					Language: file.Language,
					FuncName: funcName,
				})
			}
			start = i + 1
			funcName = ""
		}
	}

	// Handle remaining lines
	if start < len(lines) {
		content := strings.Join(lines[start:], "\n")
		if strings.TrimSpace(content) != "" {
			chunks = append(chunks, Chunk{
				Content:  content,
				FilePath: file.Path,
				Language: file.Language,
				FuncName: funcName,
			})
		}
	}

	return chunks
}
