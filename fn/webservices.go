package fn

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var fileCache = make(map[string][]byte)

func ServePage(w http.ResponseWriter, filename string) {
	if content, exists := fileCache[filename]; exists {
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	} else {
		http.NotFound(w, nil)
	}
}
func LoadFilesIntoCache() {
	dir := PanicErr(filepath.Abs("static"))
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			content := PrintErr(os.ReadFile(path))
			key := filepath.ToSlash(strings.TrimPrefix(path, dir+string(filepath.Separator)))
			fileCache[key] = content
			fmt.Printf("Cached file: %s (%d kb)\n", key, len(content)/1024)
		}
		return nil
	})
}
func CachedFileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := strings.TrimPrefix(r.URL.Path, "/static/")

	if content, exists := fileCache[filePath]; exists {
		if strings.HasSuffix(filePath, ".html") {
			w.Header().Set("Content-Type", "text/html")
		} else if strings.HasSuffix(filePath, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(filePath, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}

		w.Write(content)
	} else {
		http.NotFound(w, r)
	}
}
func GetCachedFile(filename string) ([]byte, bool) {
	content, exists := fileCache[filename]
	return content, exists
}
