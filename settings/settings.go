package settings

import (
	// _ "embed"

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
