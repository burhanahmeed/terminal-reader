package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"

	"github.com/burhanahmeed/terminal-reader/internal/embed"
	"github.com/burhanahmeed/terminal-reader/internal/llm"
	"github.com/burhanahmeed/terminal-reader/internal/repo"
	"github.com/burhanahmeed/terminal-reader/internal/retriever"
	"github.com/burhanahmeed/terminal-reader/internal/session"
	"github.com/burhanahmeed/terminal-reader/pkg/cache"
)

func hashKey(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(strings.TrimSpace(p)))
	}
	return hex.EncodeToString(h.Sum(nil))

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	localPath := flag.String("path", "", "path to local repository")
	githubURL := flag.String("github", "", "GitHub public repository URL")
	flag.Parse()

	var repoPath *string
	if *localPath != "" {
		repoPath = localPath
	} else if *githubURL != "" {
		repoPath, err = repo.CloneGithub(githubURL)
		if err != nil {
			log.Fatal("Failed to clone GitHub repository:", err)
		}
	} else {
		log.Fatal("Either --path or --github is required")
	}

	fmt.Println("📦 Indexing repo:", *localPath, *githubURL)

	cacheLayer, err := cache.NewFileCache("data/cache.json")
	if err != nil {
		log.Fatal("Failed to create cache:", err)
	}

	embedder, err := embed.NewGeminiEmbedder()
	if err != nil {
		log.Fatal("Failed to create embedder:", err)
	}

	llmClient, err := llm.NewGeminiClient()
	if err != nil {
		log.Fatal("Failed to create LLM client:", err)
	}

	store, err := retriever.NewSQLiteStore("data/vectors.db")
	if err != nil {
		log.Fatal("Failed to create vector store:", err)
	}

	loader := repo.Loader{}
	files, err := loader.LoadRepo(*repoPath)
	if err != nil {
		log.Fatal(err)
	}

	repoHash := hashKey(*repoPath)
	var allChunks []repo.Chunk
	for _, f := range files {
		chunk := repo.ChunkFile(f, 100)
		allChunks = append(allChunks, chunk...)
	}

	embed.AsyncEmbed(embedder, store, cacheLayer, allChunks, repoHash, 3)

	fmt.Println("✅ Repo indexed. Starting chat (type 'exit' to quit).")

	s := session.Session{}
	s.PromptLoop(func(query string) string {
		qKey := hashKey(repoHash, query)
		if cached, ok := cacheLayer.Get(qKey); ok {
			return "(cached)\n" + cached
		}

		qVec, err := embedder.EmbedText(query)
		if err != nil {
			return "Error embedding query: " + err.Error()
		}
		topDocs, err := store.Search(qVec, 5, repoHash)
		if err != nil {
			return "Error searching vectors: " + err.Error()
		}

		prompt := fmt.Sprintf(
			"You are a code assistant. Repo path: %s\nRelevant code:\n%s\n\nUser: %s",
			filepath.Base(*repoPath),
			strings.Join(topDocs, "\n\n---\n\n"),
			query,
		)

		fmt.Printf("Thinking...")

		resp, err := llmClient.Generate(prompt)
		if err != nil {
			return "Error: " + err.Error()
		}
		fmt.Printf("\r")
		cacheLayer.Set(qKey, resp)
		return resp
	})
}
