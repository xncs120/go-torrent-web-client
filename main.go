package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"go-torrent-web-client/database"
	"go-torrent-web-client/templates"
	"go-torrent-web-client/torrent"
	"go-torrent-web-client/websocket"
)

func main() {
	dataDir := "./downloads"

	jd, err := database.NewJsonData("./GoTorrentWebClient.json")
	if err != nil {
		fmt.Println("Error initializing json file: %v", err)
	}

	tm, err := torrent.NewTorrentManager(jd, dataDir)
	if err != nil {
		fmt.Println("Error creating torrent manager:", err)
		os.Exit(1)
	}
	defer tm.Client.Close()

	wsm := websocket.NewWebSocketManager(tm)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(templates.IndexHTML)
	})

	http.HandleFunc("/media/", func(w http.ResponseWriter, r *http.Request) {
		filePath := r.URL.Path[len("/media/"):]

		playerHTML := string(templates.PlayerHTML)
		playerHTML = strings.ReplaceAll(playerHTML, "{{.FilePath}}", "/stream/"+filePath)
		playerHTML = strings.ReplaceAll(playerHTML, "{{.MediaTitle}}", filePath)

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(playerHTML))
	})

	http.HandleFunc("/stream/", func(w http.ResponseWriter, r *http.Request) {
		filePath := r.URL.Path[len("/stream/"):]
		file := "./downloads/" + filePath

		f, err := os.Open(file)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer f.Close()

		fileStat, err := f.Stat()
		if err != nil {
			http.Error(w, "File not accessible", http.StatusInternalServerError)
			return
		}

		rangeHeader := r.Header.Get("Range")
		if rangeHeader != "" {
			var start, end int64
			fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)

			if end == 0 || end >= fileStat.Size() {
				end = fileStat.Size() - 1
			}

			if start > end {
				http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
				return
			}

			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileStat.Size()))
			w.WriteHeader(http.StatusPartialContent)

			f.Seek(start, 0)
			buf := make([]byte, end-start+1)
			f.Read(buf)
			w.Write(buf)
			return
		}

		http.ServeContent(w, r, filePath, fileStat.ModTime(), f)
	})

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		source := r.FormValue("source")
		if source == "" {
			http.Error(w, "Source is required", http.StatusBadRequest)
			return
		}

		name, err := tm.AddDownload(source)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error starting download: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Download started for: %s", name)
	})

	http.HandleFunc("/histories", func(w http.ResponseWriter, r *http.Request) {
		data, err := jd.SelectAll()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving data: %v", err), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

	http.HandleFunc("/progresses", wsm.SendProgresses)

	fmt.Println("Starting server on http://localhost:60000...")
	if err := http.ListenAndServe(":60000", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
