package settings

import (
	// _ "embed"

	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	// "gopkg.in/yaml.v3"
)

/* //go:embed settings.yaml
var settingsFile []byte

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Settings struct {
	Port      string         `yaml:"port"`
	AdminPort string         `yaml:"admin-port"`
	DB        DatabaseConfig `yaml:"database"`
} */

type AppSettings struct {
	WebserverPort      int    `toml:"webserver_port"`
	AdminWebserverPort int    `toml:"admin_webserver_port"`
	DatabaseHost       string `toml:"database_host"`
	DatabasePort       int    `toml:"database_port"`
	DatabaseUser       string `toml:"database_user"`
	DatabasePassword   string `toml:"database_password"`
	DatabaseName       string `toml:"database_name"`
	ImageDirectory     string `toml:"image_dir"`
}

func New(filepath *string) (*AppSettings, error) {
	var s *AppSettings

	if *filepath != "" {
		log.Printf("reading config file %s", *filepath)
		settings, err := ReadConfigToml(filepath)
		if err != nil {
			return nil, err
		}

		s = settings
	} else {
		log.Println("no config file, reading environment variables")
		settings, err := LoadSettings()
		if err != nil {
			return nil, err
		}

		s = settings
	}

	return s, nil
}

func LoadSettings() (*AppSettings, error) {
	// Want to load the environment variables
	var appSettings *AppSettings

	webserver_port_str := os.Getenv("WEBSERVER_PORT")
	if webserver_port_str == "" {
		webserver_port_str = "8080"
	}

	webserver_port, err := strconv.Atoi(webserver_port_str)
	if err != nil {
		return appSettings, fmt.Errorf(
			"WEBSERVER_PORT is not valid: %v", err,
		)
	}

	admin_webserver_port_str := os.Getenv("ADMIN_WEBSERVER_PORT")
	if admin_webserver_port_str == "" {
		admin_webserver_port_str = "8081"
	}

	admin_webserver_port, err := strconv.Atoi(admin_webserver_port_str)
	if err != nil {
		return appSettings, fmt.Errorf(
			"ADMIN_WEBSERVER_PORT is not valid: %v", err,
		)
	}

	database_host := os.Getenv("DATABASE_HOST")
	if len(database_host) == 0 {
		return appSettings, fmt.Errorf("DATABASE_HOST is not defined")
	}

	database_port_str := os.Getenv("DATABASE_PORT")
	if len(database_port_str) == 0 {
		return appSettings, fmt.Errorf("DATABASE_PORT is not defined")
	}

	database_port, err := strconv.Atoi(database_port_str)
	if err != nil {
		return appSettings, fmt.Errorf("DATABASE_PORT is not a valid integer: %v", err)
	}

	database_user := os.Getenv("DATABASE_USER")
	if len(database_user) == 0 {
		return appSettings, fmt.Errorf("DATABASE_USER is not defined")
	}

	database_password := os.Getenv("DATABASE_PASSWORD")
	if len(database_password) == 0 {
		return appSettings, fmt.Errorf("DATABASE_PASSWORD is not defined")
	}

	database_name := os.Getenv("DATABASE_NAME")
	if len(database_name) == 0 {
		return appSettings, fmt.Errorf("DATABASE_NAME is not defined")
	}

	image_directory := os.Getenv("IMAGE_DIRECTORY")
	if len(image_directory) == 0 {
		return appSettings, fmt.Errorf("IMAGE_DIRECTORY is not defined")
	}

	return &AppSettings{
		WebserverPort:      webserver_port,
		AdminWebserverPort: admin_webserver_port,
		DatabaseHost:       database_host,
		DatabaseUser:       database_user,
		DatabasePassword:   database_password,
		DatabasePort:       database_port,
		DatabaseName:       database_name,
		ImageDirectory:     image_directory,
	}, nil
}

func ReadConfigToml(filepath *string) (*AppSettings, error) {
	/* var s Settings

	err := yaml.Unmarshal(settingsFile, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil */
	var s AppSettings

	_, err := toml.DecodeFile(*filepath, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
