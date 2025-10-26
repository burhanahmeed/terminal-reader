package embed

import (
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
	type task struct {
		c repo.Chunk
	}

	tasks := make([]task, len(chunks))
	results := make(chan struct{}, len(chunks))

	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for t := range tasks {
			task := tasks[t]
			chunkKey := repoHash + "|" + task.c.FilePath + "|" + task.c.FuncName
			if _, ok := cacheLayer.Get(chunkKey); ok {
				results <- struct{}{}
				continue
			}
			vec, err := embedder.EmbedText(task.c.Content)
			if err == nil {
				store.Add(task.c.FilePath+"|"+task.c.Language+"|"+task.c.FuncName, vec, repoHash, task.c.FilePath, task.c.FuncName)
				cacheLayer.Set(chunkKey, "done")
			}
			results <- struct{}{}
		}
	}

	// start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker()
	}

	// enqueue tasks
	for _, ch := range chunks {
		tasks = append(tasks, task{c: ch})
	}
	wg.Wait()
}
