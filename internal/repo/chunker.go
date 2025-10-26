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

// Simple function regex per language
var funcPatterns = map[string]*regexp.Regexp{
	"go": regexp.MustCompile(`(?m)^func\s+([^\(]+)`),
	"py": regexp.MustCompile(`(?m)^def\s+([^\(]+)`),
	"js": regexp.MustCompile(`(?m)^function\s+([^\(]+)`),
	"ts": regexp.MustCompile(`(?m)^function\s+([^\(]+)`),
}

func ChunkFile(file FileData, maxLines int) []Chunk {
	lines := strings.Split(file.Content, "\n")
	chunks := make([]Chunk, 0)

	pattern, ok := funcPatterns[file.Language]
	if !ok {
		return nil
	}

	start := 0
	funcName := ""
	for i, line := range lines {
		if i >= maxLines {
			break
		}
		if pattern.MatchString(line) {
			chunks = append(chunks, Chunk{Content: strings.Join(lines[start:i], "\n"), FilePath: file.Path, Language: file.Language, FuncName: funcName})
			start = i
			funcName = line
		} else if funcName != "" && strings.HasPrefix(line, "}") {
			funcName = ""
		}
	}
	if start < len(lines) {
		chunks = append(chunks, Chunk{Content: strings.Join(lines[start:], "\n"), FilePath: file.Path, Language: file.Language, FuncName: funcName})
	}
	return chunks
}
