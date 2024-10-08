package base

import (
	"database/sql"
	"fmt"
	"log"
	"music/internal/config"
	"music/internal/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Импортируем драйвер PostgreSQL
	"github.com/pressly/goose/v3"
	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL для sql.DB
)

type Repository interface {
	AddSong(newSong model.Song) error
	GetLibrary() ([]model.Song, error)
	GetLyrics(name string) (string, error)
	DeleteSong(name string) error
	UpdateSong(oldNameSong string, updateSong model.Song) error
}

type repository struct {
	base *gorm.DB
}

func applyMigrations(db *sql.DB) error {
	goose.SetDialect("postgres")
	if err := goose.Up(db, "internal/base/migrations"); err != nil {
		return err
	}
	return nil
}

func NewRepository(cfg config.Config) (Repository, error) {
	c := cfg.GetConfig()
	log.Print("Connecting to database...")
	
	fmt.Println(c)
	b, err := gorm.Open("postgres", c)
	if err != nil {
		return nil, err
	}

	sqlDB := b.DB()
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	log.Println("success")

	log.Print("Running migrations...")
	if err := applyMigrations(sqlDB); err != nil {
		return nil, err
	}

	// if err := goose.Up(chance, "./migrations"); err != nil {
	// 	return nil, err
	// }
	log.Println("database migrations completed")

	return &repository{
		base: b,
	}, nil
}

func (r *repository) AddSong(newSong model.Song) error {
	fmt.Println(newSong)
	r.base.Create(&newSong)
	if r.base.Error != nil {
		return fmt.Errorf("Failed to add song:%w", r.base.Error)
	}

	log.Println("Song added with ID:", newSong.ID)
	return nil
}

func (r *repository) GetLibrary() ([]model.Song, error) {
	log.Print("Trying to get library data...")
	data := make([]model.Song, 0)

	if err := r.base.Find(&data).Error; err != nil {
		return nil, fmt.Errorf("error:", err)
	}

	return data, nil
}

func (r *repository) GetLyrics(name string) (string, error) {
	log.Print("Trying to get lyrics...")
	var song model.Song
	if err := r.base.Where("song = ?", name).First(&song).Error; err != nil {
		return "", fmt.Errorf("error:", err)
	}

	return song.Group_name, nil
}

func (r *repository) DeleteSong(name string) error {
	log.Print("Trying to delete song...")
	if err := r.base.Where("song = ?", name).Delete(&model.Song{}).Error; err != nil {
		return fmt.Errorf("error:", err)
	}

	return nil
}

func (r *repository) UpdateSong(oldNameSong string, updateSong model.Song) error {
	log.Print("Trying to update song...")
	if err := r.base.Model(&model.Song{}).Where("song = ?", oldNameSong).Update(&updateSong).Error; err != nil {
		return fmt.Errorf("error:", err)
	}

	return nil
}
