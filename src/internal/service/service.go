package service

import (
	"encoding/json"
	"fmt"
	"log"
	"music/internal/base"
	"music/internal/model"
	"net/http"

	"github.com/gorilla/mux"
)

type Service interface {
	Run() error
	Close() error
}

type service struct {
	router *mux.Router
	repo   base.Repository
}

func NewService(r base.Repository) Service {
	router := mux.NewRouter()

	s := service{
		router: router,
		repo:   r,
	}

	s.setupRoutes()

	return &s
}

func (s *service) Run() error {
	log.Println("Server running on port :8888")
	if err := http.ListenAndServe(":8888", s.router); err != nil{
		return fmt.Errorf("Failed to listen and serve. Error: %s", err.Error())
	}

	return nil
}

func (s *service) setupRoutes() {
	s.router.HandleFunc("/music", s.Library).Methods("GET")
	s.router.HandleFunc("/music/{group}/{song}/lyrics", s.Lyrics).Methods("GET")
	s.router.HandleFunc("/music/{group}/{song}", s.Delete).Methods("DELETE")
	s.router.HandleFunc("/music/{group}/{song}", s.Update).Methods("PUT")
	s.router.HandleFunc("/music", s.Add).Methods("POST")
}

func (s *service) Library(w http.ResponseWriter, r *http.Request) {
	lib, err := s.repo.GetLibrary()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to fetch library data", http.StatusInternalServerError)
		return
	}
	log.Println("Successfully fetched song library data")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lib)
}

func (s *service) Lyrics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	group, song := params["group"], params["song"]
	text, err := s.repo.GetLyrics(group, song)
	if err != nil {
		log.Println(err)
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}
	log.Printf("Getting lyrics group: %s, song: %s completed", group, song)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(text)
}

func (s *service) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	group, song := params["group"], params["song"]
	err := s.repo.DeleteSong(group,song)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to delete song", http.StatusInternalServerError)
		return
	}
	log.Printf("Group: %s, song:%s successfully removed from the library", group, song)

	w.WriteHeader(http.StatusNoContent)
}

func (s *service) Update(w http.ResponseWriter, r *http.Request) {
	var updatedSong model.Song
	if err := json.NewDecoder(r.Body).Decode(&updatedSong); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	group, song := params["group"], params["song"]
	err := s.repo.UpdateSong(group, song, updatedSong)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to update song", http.StatusInternalServerError)
		return
	}
	log.Printf("Group: %s, song: %s successfully updated")

	w.Header().Set("Content-Type", "application/json")
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
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (s *service) Close() error{
	if err := s.repo.Close(); err != nil{
		return fmt.Errorf("Failed to close server. Error: %s", err.Error())
	}

	return nil
}
