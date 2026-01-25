package crud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func A(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	fmt.Println(r.URL.Path)
	if r.Method == "GET" {
		selectByID(w, r)
	} else if r.Method == "PUT" {
		updateCategory(w, r)
	} else if r.Method == "DELETE" {
		deleteCategory(w, r)
	}
}
func selectByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/") // ambil ID dari URL path
	id, e := strconv.Atoi(idStr)                                // ganti ke int
	if e != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	for _, c := range Category1 {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}
	http.Error(w, "Category belum ada", http.StatusNotFound)
}
func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, e := strconv.Atoi(idStr)
	if e != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var updateCategory Category
	e = json.NewDecoder(r.Body).Decode(&updateCategory)
	if e != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	for i := range Category1 {
		if Category1[i].ID == id {
			updateCategory.ID = id
			Category1[i] = updateCategory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategory)
			return
		}
	}
	http.Error(w, "Category belum ada", http.StatusNotFound)
}
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, e := strconv.Atoi(idStr)
	if e != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	for i, c := range Category1 {
		if c.ID == id {
			// bikin slice baru dengan data sebelum dan sesudah index
			Category1 = append(Category1[:i], Category1[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Sukses delete",
			})
			return
		}
	}
	http.Error(w, "Category belum ada", http.StatusNotFound)
}
