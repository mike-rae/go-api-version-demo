package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/api/v1/users/", uriV1)
	http.HandleFunc("/api/v2/users/", uriV2)
	http.HandleFunc("/api/headers/users/", headerVersioned)
	http.HandleFunc("/api/media/users/", mediaVersioned)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func uriV1(writer http.ResponseWriter, reader *http.Request) {
	id := lastPathPart(reader.URL.Path)
	writeJSON(writer, map[string]any{
		"id":   id,
		"name": "Alice Smith",
	})
}

func uriV2(w http.ResponseWriter, r *http.Request) {
	id := lastPathPart(r.URL.Path)
	writeJSON(w, map[string]any{
		"id":        id,
		"firstName": "Alice",
		"lastName":  "Smith",
		"status":    "ACTIVE",
	})
}

func headerVersioned(writer http.ResponseWriter, reader *http.Request) {
	id := lastPathPart(reader.URL.Path)
	version := reader.Header.Get("API-Version")

	fmt.Println("Header version number specified:", version)

	switch version {
	case "1":
		writeJSON(writer, map[string]any{
			"id":   id,
			"name": "Bob Smith",
		})
	case "2":
		writeJSON(writer, map[string]any{
			"id":        id,
			"firstName": "Bob",
			"lastName":  "Smith",
			"status":    "ACTIVE",
		})
	default:
		http.Error(writer, "Missing or unsupported API-Version header", http.StatusBadRequest)
	}
}

func mediaVersioned(writer http.ResponseWriter, reader *http.Request) {
	id := lastPathPart(reader.URL.Path)
	accept := reader.Header.Get("Accept")

	switch accept {
	case "application/vnd.demo.v1+json":
		fmt.Println("v1 vendor")
		writeJSON(writer, map[string]any{
			"id":   id,
			"name": "Paaa Smith",
		})
	case "application/vnd.demo.v2+json":
		fmt.Println("v2 vendor")
		writeJSON(writer, map[string]any{
			"id":        id,
			"firstName": "Paas",
			"lastName":  "Smith",
			"status":    "ACTIVE",
		})
	default:
		http.Error(writer, "Unsupported media type version", http.StatusNotAcceptable)
	}
}

func writeJSON(w http.ResponseWriter, data map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func lastPathPart(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	return parts[len(parts)-1]
}
