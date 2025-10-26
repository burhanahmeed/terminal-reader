package embed

import (
	"fmt"
	"strings"
	"sync"

	"github.com/burhanahmeed/terminal-reader/internal/repo"
	"github.com/burhanahmeed/terminal-reader/internal/retriever"
	"github.com/burhanahmeed/terminal-reader/pkg/cache"
)

func AsyncEmbed(
	embedder *GeminiEmbedder,
	store *retriever.SQLiteStore,
	cacheLayer *cache.FileCache,
	chunks []repo.Chunk,
	repoHash string,
	concurrency int,
) {
	taskChan := make(chan repo.Chunk, len(chunks))
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for chunk := range taskChan {
			// Skip empty chunks
			if strings.TrimSpace(chunk.Content) == "" {
				continue
			}

			chunkKey := repoHash + "|" + chunk.FilePath + "|" + chunk.FuncName
			if _, ok := cacheLayer.Get(chunkKey); ok {
				continue // already embedded
			}

			vec, err := embedder.EmbedText(chunk.Content)
			if err != nil {
				fmt.Printf("Failed to embed chunk from %s: %v\n", chunk.FilePath, err)
				continue
			}

			err = store.Add(chunk.FilePath+"|"+chunk.Language+"|"+chunk.FuncName, vec, repoHash, chunk.FilePath, chunk.FuncName)
			if err != nil {
				fmt.Printf("Failed to store chunk from %s: %v\n", chunk.FilePath, err)
				continue
			}

			cacheLayer.Set(chunkKey, "done")
		}
	}

	// start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker()
	}

	// send tasks to workers
	for _, chunk := range chunks {
		taskChan <- chunk
	}
	close(taskChan)

	wg.Wait()
}
