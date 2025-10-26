package retriever

import (
	"database/sql"
	"encoding/json"
	"math"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Initialize schema
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS vectors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT,
		vector TEXT,
		repo_id TEXT,
		file_path TEXT,
		func_name TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_vectors_repo_id ON vectors (repo_id);
	`)
	return &SQLiteStore{db: db}, err
}

func (s *SQLiteStore) Add(content string, vector []float32, repoId, filePath, funcName string) error {
	vecJSON, _ := json.Marshal(vector)
	_, err := s.db.Exec(`INSERT INTO vectors (content, vector, repo_id, file_path, func_name) VALUES (?, ?, ?, ?, ?)`, content, string(vecJSON), repoId, filePath, funcName)
	return err
}

func cosineSim(a, b []float32) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func (s *SQLiteStore) Search(vector []float32, topK int, repoId string) ([]string, error) {
	rows, err := s.db.Query(`SELECT content, vector FROM vectors WHERE repo_id = ?`, repoId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type scored struct {
		text  string
		score float64
	}
	var results []scored

	for rows.Next() {
		var text, vecJSON string
		rows.Scan(&text, &vecJSON)

		var storedVec []float32
		json.Unmarshal([]byte(vecJSON), &storedVec)
		score := cosineSim(vector, storedVec)
		results = append(results, scored{text, score})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	top := make([]string, 0, topK)
	for i := 0; i < len(results) && i < topK; i++ {
		top = append(top, results[i].text)
	}
	return top, nil
}
