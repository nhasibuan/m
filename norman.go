package main

import (
	"encoding/json"
	"errors"
	"log"
	"m/crud"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}
	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}
	db, err := crud.Connect(config.DBConn)
	if err != nil {
		log.Fatal("", err)
	}
	defer db.Close()
	//
	repository := crud.NewRepository(db)
	service := crud.NewService(repository)
	//
	handler := crud.NewHandler(service)
	http.HandleFunc("/test", handler.HandleHealth)
	http.HandleFunc("/api/categories", handler.Handle)
	http.HandleFunc("/api/categories/", handler.HandleByID)
	//
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"Method":  r.Method,
			"Content": w.Header().Get("Content-Type"),
		})
	})
	// curl http://localhost:8080/health
	http.HandleFunc("/api/categories1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" { // select all
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(crud.Category1)
		} else if r.Method == "POST" { // insert
			var categoryBaru crud.Category // Category
			err := json.NewDecoder(r.Body).Decode(&categoryBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			categoryBaru.ID = len(crud.Category1) + 1
			crud.Category1 = append(crud.Category1, categoryBaru)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(categoryBaru)
		}
	})
	// curl http://localhost:8080/api/categories
	http.HandleFunc("/api/categories1/", crud.A)

	e := http.ListenAndServe(config.Port, nil)
	if e != nil && !errors.Is(e, http.ErrServerClosed) {

		log.Fatal(e)
	}
}

// https://m-production-f151.up.railway.app/health
// https://railway.com?referralCode=rbJYUL
// https://m-nhasibuan5181-xe4oymdo.leapcell.dev/health
