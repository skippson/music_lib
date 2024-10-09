package base

import (
	"database/sql"
	"fmt"
	"log"
	"music/internal/config"
	"music/internal/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pressly/goose/v3"
	_ "github.com/lib/pq"
)

type Repository interface {
	AddSong(newSong model.Song) error
	GetLibrary() ([]model.Song, error)
	GetLyrics(group, song string) (string, error)
	DeleteSong(group, song string) error
	UpdateSong(group, song string, updateSong model.Song) error
	Close() error
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
	b, err := gorm.Open("postgres", c)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database. Error: %s", err.Error())
	}

	sqlDB := b.DB()
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping database. Error: %s", err.Error())
	}
	log.Println("Connecting to database: success")

	log.Print("Running migrations...")
	if err := applyMigrations(sqlDB); err != nil {
		return nil, fmt.Errorf("Failed to make migrations. Error: %s", err.Error())
	}
	log.Println("Database migrations completed")

	return &repository{
		base: b,
	}, nil
}

func (r *repository) AddSong(newSong model.Song) error {
	log.Printf("Trying to add group: %s, song: %s", newSong.Group_name, newSong.Song)
	r.base.Create(&newSong)
	if r.base.Error != nil {
		return fmt.Errorf("Failed to add group: %s, song: %s. Error:%s", newSong.Group_name, newSong.Song, r.base.Error.Error())
	}

	log.Printf("Group: %s, song: %s added with ID:%d",newSong.Group_name, newSong.Song, newSong.ID)
	return nil
}

func (r *repository) GetLibrary() ([]model.Song, error) {
	log.Print("Trying to fetching library data...")
	data := make([]model.Song, 0)

	if err := r.base.Find(&data).Error; err != nil {
		return nil, fmt.Errorf("Failed to fetch library data. Error: %s ", err.Error())
	}

	return data, nil
}

func (r *repository) GetLyrics(group, song string) (string, error) {
	log.Printf("Trying to get lyrics of group: %s, song: %s", group, song)
	var target model.Song
	if err := r.base.Where("group_name = ? and song = ?", group, song).First(&target).Error; err != nil {
		return "", fmt.Errorf("Failed to get lyrics of group: %s, song: %s. Error: %s", group, song, err.Error())
	}

	return target.Lyrics, nil
}

func (r *repository) DeleteSong(group, song string) error {
	log.Printf("Trying to delete group: %s, song: %s", group, song)
	if err := r.base.Where("group_name = ? and song = ?", group, song).Delete(&model.Song{}).Error; err != nil {
		return fmt.Errorf("Failed to delete group: %s, song: %s. Error: %s ", group, song, err.Error())
	}

	return nil
}

func (r *repository) UpdateSong(group, song string, updateSong model.Song) error {
	log.Printf("Trying to update group: %s, song: %s")
	if err := r.base.Model(&model.Song{}).Where("group_name = ? and song = ?", group, song).Update(&updateSong).Error; err != nil {
		return fmt.Errorf("Failed to update group: %s, song: %s. Error: %s", group, song, err.Error())
	}

	return nil
}

func (r * repository) Close() error{
	if err := r.base.Close(); err != nil{
		return fmt.Errorf("Failed to close database. Error: %s", err.Error())
	}

	return nil
}