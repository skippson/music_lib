package service

import (
	"encoding/json"
	"log"
	"music/internal/base"
	"music/internal/model"
	"net/http"

	"github.com/gorilla/mux"
)

type Service interface {
	Run() error
}

type service struct {
	router *mux.Router
	repo   base.Repository
}

func NewService(r base.Repository) (Service){
	router := mux.NewRouter()

	s := service{
		router: router,
		repo: r,
	}

	s.setupRoutes()

	return &s
}

func (s *service) Run() error{
	log.Println("Server running on 8080")
	http.ListenAndServe(":8080", s.router)
	return nil
}

func (s *service) setupRoutes() {
	s.router.HandleFunc("/songs", s.Library).Methods("GET")
	s.router.HandleFunc("/songs/{song}", s.Lyrics).Methods("GET")
	s.router.HandleFunc("/songs/{song}", s.Delete).Methods("DELETE")
	s.router.HandleFunc("/songs/{song}", s.Update).Methods("PUT")
	s.router.HandleFunc("/songs", s.Add).Methods("POST")
}

func (s *service) Library(w http.ResponseWriter, r *http.Request) {
	lib, err := s.repo.GetLibrary()
	if err != nil {
		log.Println("Error fetching library data:", err)
		http.Error(w, "Failed to fetch library data", http.StatusInternalServerError)
		return
	}
	log.Println("success")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lib)
}

func (s *service) Lyrics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	songName := params["song"]

	text, err := s.repo.GetLyrics(songName)
	if err != nil {
		log.Println(err)
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}
	log.Println("success")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(text)
}

func (s *service) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	songName := params["song"]

	err := s.repo.DeleteSong(songName)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to delete song", http.StatusInternalServerError)
		return
	}
	log.Println("success")

	w.WriteHeader(http.StatusNoContent)
}

func (s *service) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	songName := params["song"]

	var updatedSong model.Song
	if err := json.NewDecoder(r.Body).Decode(&updatedSong); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := s.repo.UpdateSong(songName, updatedSong)
	if err != nil {
		log.Println("Error updating song:", err)
		http.Error(w, "Failed to update song", http.StatusInternalServerError)
		return
	}
	log.Println("success")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSong)
}

func (s *service) Add(w http.ResponseWriter, r *http.Request) {
	var newSong model.Song
	if err := json.NewDecoder(r.Body).Decode(&newSong); err != nil {
		log.Println("Error decoding JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := s.repo.AddSong(newSong)
	if err != nil {
		log.Println("Error adding song", nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println("success")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(newSong)
}
