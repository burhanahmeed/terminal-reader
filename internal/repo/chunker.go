package repo

import (
	"regexp"
	"strings"
)

// Chunk represents a meaningful section of a file.
type Chunk struct {
	Content  string
	FilePath string
	Language string
	FuncName string // unused for markdown, kept for consistency
}

// Regex to detect Markdown headings (#, ##, ###, etc.)
var headingPattern = regexp.MustCompile(`^(#{1,6})\s+(.*)`)

// ChunkFile splits a Markdown file into semantic chunks based on headings.
func ChunkFile(file FileData, maxLines int) []Chunk {
	lines := strings.Split(file.Content, "\n")
	chunks := make([]Chunk, 0)

	start := 0
	currentHeading := ""

	for i, line := range lines {
		if headingPattern.MatchString(line) {
			// Save previous chunk before starting a new section
			if start < i {
				content := strings.Join(lines[start:i], "\n")
				if strings.TrimSpace(content) != "" {
					chunks = append(chunks, Chunk{
						Content:  content,
						FilePath: file.Path,
						Language: "md",
						FuncName: currentHeading,
					})
				}
			}

			// Start new section at this heading
			start = i
			match := headingPattern.FindStringSubmatch(line)
			currentHeading = match[2]
		}
	}

	// Handle the last chunk (if any)
	if start < len(lines) {
		content := strings.Join(lines[start:], "\n")
		if strings.TrimSpace(content) != "" {
			chunks = append(chunks, Chunk{
				Content:  content,
				FilePath: file.Path,
				Language: "md",
				FuncName: currentHeading,
			})
		}
	}

	// Optionally, apply maxLines as a fallback for very large sections
	if maxLines > 0 {
		chunks = splitLargeChunks(chunks, maxLines)
	}

	return chunks
}

// splitLargeChunks breaks overly large chunks by maxLines.
func splitLargeChunks(chunks []Chunk, maxLines int) []Chunk {
	result := make([]Chunk, 0)
	for _, c := range chunks {
		lines := strings.Split(c.Content, "\n")
		if len(lines) <= maxLines {
			result = append(result, c)
			continue
		}

		for i := 0; i < len(lines); i += maxLines {
			end := i + maxLines
			if end > len(lines) {
				end = len(lines)
			}
			subContent := strings.Join(lines[i:end], "\n")
			if strings.TrimSpace(subContent) != "" {
				result = append(result, Chunk{
					Content:  subContent,
					FilePath: c.FilePath,
					Language: c.Language,
					FuncName: c.FuncName,
				})
			}
		}
	}
	return result
}
