package repo

import (
	"io/fs"
	"os"
	"path/filepath"
)

type FileData struct {
	Path      string
	Language  string
	Content   string
	Extension string
}

type Loader struct{}

// supported file extensions
var allowedExt = map[string]string{
	".md": "markdown",
}

func (l *Loader) LoadRepo(root string) ([]FileData, error) {
	var files []FileData

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		lang, ok := allowedExt[ext]
		if !ok {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files = append(files, FileData{
			Path:      path,
			Language:  lang,
			Content:   string(data),
			Extension: ext,
		})
		return nil
	})

	return files, err
}
