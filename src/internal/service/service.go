package service

import (
	"encoding/json"
	"fmt"
	"log"
	"music/internal/base"
	"music/internal/config"
	"music/internal/model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Service interface {
	Run() error
	Router() *mux.Router
	Close() error
}

type service struct {
	router *mux.Router
	repo   base.Repository
	cfg    config.Config
}

func (s *service) Router() *mux.Router {
	return s.router
}

func NewService(c config.Config, r base.Repository) Service {
	router := mux.NewRouter()

	s := service{
		router: router,
		repo:   r,
		cfg: c,
	}

	s.setupRoutes()

	return &s
}

func (s *service) Run() error {
	port := s.cfg.GetPort()
	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(port, s.router); err != nil {
		return fmt.Errorf("Failed to listen and serve. Error: %s", err.Error())
	}

	return nil
}

func (s *service) setupRoutes() {
	s.router.HandleFunc("/music", s.Library).Methods("GET")
	s.router.HandleFunc("/music", s.Add).Methods("POST")
	s.router.HandleFunc("/music/filter", s.Filter).Methods("GET")
	s.router.HandleFunc("/music/{group}/{song}", s.Update).Methods("PUT")
	s.router.HandleFunc("/music/{group}/{song}", s.Delete).Methods("DELETE")
	s.router.HandleFunc("/music/{group}/{song}/lyrics", s.Lyrics).Methods("GET")
	s.router.HandleFunc("/music/{page}/{size}", s.LibraryWithPagination).Methods("GET")
	s.router.HandleFunc("/music/filter/{page}/{size}", s.FilterWithPagination).Methods("GET")
	s.router.HandleFunc("/music/{group}/{song}/lyrics/{page}/{size}", s.LyricsWithPagination).Methods("GET")
}

// @Summary Получить библиотеку песен с пагинацией
// @Description Возвращает список песен с пагинацией
// @Tags music
// @Produce json
// @Param page path int true "Номер страницы"
// @Param size path int true "Размер страницы"
// @Success 200 {array} model.Song
// @Failure 400 {string} string "Invalid page number or size"
// @Failure 500 {string} string "Failed to fetch library data"
// @Router /music/{page}/{size} [get]
func (s *service) LibraryWithPagination(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	page, err := strconv.Atoi(params["page"])
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}
	size, err := strconv.Atoi(params["size"])
	if err != nil {
		http.Error(w, "Invalid page size", http.StatusBadRequest)
		return
	}

	lib, err := s.repo.GetLibraryWithPagination(page, size)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to fetch library data", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully fetched song library data with page: %d, size: %d", page, size)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lib)

}

// @Summary Получить библиотеку песен c фильтром и пагинацией
// @Description Возвращает список песен c фильтром и пагинацией
// @Tags music
// @Produce json
// @Param page path int true "Номер страницы"
// @Param size path int true "Размер страницы"
// @Success 200 {array} model.Song
// @Failure 400 {string} string "Invalid page number or size"
// @Failure 404 {string} string "Song not found"
// @Router /music/filter/{page}/{size} [get]
func (s *service) FilterWithPagination(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	page, err := strconv.Atoi(params["page"])
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}
	size, err := strconv.Atoi(params["size"])
	if err != nil {
		http.Error(w, "Invalid page size", http.StatusBadRequest)
		return
	}

	filter := r.URL.Query().Encode()
	target, err := s.repo.FindWithFilterAndPagination(filter, page, size)
	if err != nil {
		log.Println(err)
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	log.Printf("Successfully finded with filter: %s, page: %d, size: %d", filter, page, size)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(target)
}

// LyricsWithPagination gets lyrics with pagination
// @Summary Get lyrics with pagination
// @Description Returns lyrics for a given group and title with pagination support
// @Tags music
// @Produce json
// @Param group path string true "Group name"
// @Param song path string true "Song title"
// @Param page path int true "Page number"
// @Param size path int true "Page size"
// @Success 200 {array} string "Lyrics"
// @Failure 400 {string} string "Invalid page number or size"
// @Failure 404 {string} string "Song not found"
// @Router /music/{group}/{song}/lyrics/{page}/{size} [get]
func (s *service) LyricsWithPagination(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	page, err := strconv.Atoi(params["page"])
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}
	size, err := strconv.Atoi(params["size"])
	if err != nil {
		http.Error(w, "Invalid page size", http.StatusBadRequest)
		return
	}

	group := params["group"]
	song := params["song"]
	lyrics, err := s.repo.GetLyricsWithPagination(group, song, page, size)
	if err != nil {
		log.Println(err)
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}
	log.Printf("Successfully getting lyrics group: %s, song: %s, page: %d, size: %d", group, song, page, size)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lyrics)
}

// Filter фильтрует песни по заданным критериям
// @Summary Фильтрация песен
// @Description Возвращает список песен, отфильтрованных по заданным критериям
// @Tags music
// @Produce json
// @Param filter query string false "Критерии фильтрации в формате ключ=значение"
// @Success 200 {array} model.Song "Список отфильтрованных песен"
// @Failure 404 {string} string "Song not found"
// @Router /music/filter [get]
func (s *service) Filter(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Encode()
	target, err := s.repo.FindWithFilter(filter)
	if err != nil {
		log.Println(err)
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}
	log.Printf("Successfully finded with filter: %s", filter)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(target)
}

// Library возвращает библиотеку песен
// @Summary Получить библиотеку песен
// @Description Возвращает полную библиотеку песен
// @Tags music
// @Produce json
// @Success 200 {array} model.Song "Полный список песен"
// @Failure 500 {string} string "Failed to fetch library data"
// @Router /music/library [get]
func (s *service) Library(w http.ResponseWriter, r *http.Request) {
	lib, err := s.repo.GetLibrary()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to fetch library data", http.StatusInternalServerError)
		return
	}
	log.Println("Successfully fetched song library data")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lib)
}

// Lyrics возвращает текст песни по указанной группе и названию
// @Summary Получить текст песни
// @Description Возвращает текст песни на основе имени группы и названия песни
// @Tags music
// @Produce json
// @Param group path string true "Имя группы"
// @Param song path string true "Название песни"
// @Success 200 {string} string "Текст песни"
// @Failure 404 {string} string "Song not found"
// @Router /music/{group}/{song}/lyrics [get]
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

// Delete удаляет песню из библиотеки по указанной группе и названию
// @Summary Удалить песню
// @Description Удаляет песню на основе имени группы и названия песни
// @Tags music
// @Param group path string true "Имя группы"
// @Param song path string true "Название песни"
// @Success 204 "Песня успешно удалена"
// @Failure 500 {string} string "Failed to delete song"
// @Router /music/{group}/{song} [delete]
func (s *service) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	group, song := params["group"], params["song"]
	err := s.repo.DeleteSong(group, song)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to delete song", http.StatusInternalServerError)
		return
	}
	log.Printf("Group: %s, song:%s successfully removed from the library", group, song)

	w.WriteHeader(http.StatusNoContent)
}

// Update обновляет информацию о песне в библиотеке
// @Summary Обновить песню
// @Description Обновляет информацию о песне на основе имени группы и названия песни
// @Tags music
// @Accept json
// @Produce json
// @Param group path string true "Имя группы"
// @Param song path string true "Название песни"
// @Param song body model.Song true "Обновленная информация о песне"
// @Success 200 "Песня успешно обновлена"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Failed to update song"
// @Router /music/{group}/{song} [put]
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
		log.Println(err)
		http.Error(w, "Failed to update song", http.StatusInternalServerError)
		return
	}
	log.Printf("Group: %s, song: %s successfully updated", group, song)

	w.Header().Set("Content-Type", "application/json")
}

// Add добавляет новую песню в библиотеку
// @Summary Добавить песню
// @Description Добавляет новую песню в библиотеку, если она еще не существует
// @Tags music
// @Accept json
// @Produce json
// @Param song body model.Song true "Информация о новой песне"
// @Success 201 "Песня успешно добавлена"
// @Failure 400 {string} string "Ошибка декодирования JSON"
// @Failure 500 {string} string "Что-то пошло не так"
// @Router /music [post]
func (s *service) Add(w http.ResponseWriter, r *http.Request) {
	var newSong model.Song
	if err := json.NewDecoder(r.Body).Decode(&newSong); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Searching group: %s, song: %s", newSong.Group_name, newSong.Song)
	status, err := s.repo.Find(newSong.Group_name, newSong.Song)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if !status {
		err := s.repo.AddSong(newSong)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		log.Printf("Group: %s, song: %s is already in the library", newSong.Group_name, newSong.Song)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (s *service) Close() error {
	if err := s.repo.Close(); err != nil {
		return fmt.Errorf("Failed to close server. Error: %s", err.Error())
	}

	return nil
}
