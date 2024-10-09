package config

import (
	"fmt"
	"os"
)

type Config interface {
	GetConfigSQL() string
	GetPort() string
}

type config struct {
	port        string
	host        string
	db_port     string
	db_user     string
	db_name     string
	db_password string
	db_sslmode  string
}

func NewConfig() (Config, error) {
	cfg := config{}
	file, err := os.Open("../config/config.env")
	if err != nil {
		return nil, fmt.Errorf("Failed to read configuration. Error:%s", err.Error())
	}
	defer file.Close()

	fmt.Fscanf(file, "HOST=%s\nPORT=%s\nDB_PORT=%s\nDB_USER=%s\nDB_NAME=%s\nDB_PASSWORD=%s\nDB_SSLMODE=%s", &cfg.host, &cfg.port, &cfg.db_port, &cfg.db_user, &cfg.db_name, &cfg.db_password, &cfg.db_sslmode)

	return config{
		port:        cfg.port,
		host:        cfg.host,
		db_port:     cfg.db_port,
		db_user:     cfg.db_user,
		db_name:     cfg.db_name,
		db_password: cfg.db_password,
		db_sslmode:  cfg.db_sslmode,
	}, nil
}

func (c config) GetConfigSQL() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", c.host, c.db_port, c.db_user, c.db_name, c.db_password, c.db_sslmode)
}

func (c config) GetPort() string {
	return fmt.Sprintf(":%s", c.port)
}
