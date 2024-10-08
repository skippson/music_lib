package config

import (
	"fmt"
	"os"
)

type Config interface {
	GetConfig() string
}

type data_struct struct {
	db_host     string
	db_port     string
	db_user     string
	db_name     string
	db_password string
	db_sslmode  string
}

func NewConfig() (Config, error) {
	cfg := data_struct{}
	file, err := os.Open("../config/config.env")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fmt.Fscanf(file, "DB_HOST=%s\nDB_PORT=%s\nDB_USER=%s\nDB_NAME=%s\nDB_PASSWORD=%s\nDB_SSLMODE=%s", &cfg.db_host, &cfg.db_port, &cfg.db_user, &cfg.db_name, &cfg.db_password, &cfg.db_sslmode)

	return data_struct{
		db_host:     cfg.db_host,
		db_port:     cfg.db_port,
		db_user:     cfg.db_user,
		db_name:     cfg.db_name,
		db_password: cfg.db_password,
		db_sslmode:  cfg.db_sslmode,
	}, nil
}

func (d data_struct) GetConfig() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", d.db_host, d.db_port, d.db_user, d.db_name, d.db_password, d.db_sslmode)
}
