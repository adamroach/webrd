package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	viper *viper.Viper

	BindAddresses []string    `mapstructure:"bind_addresses" yaml:"bind_addresses"`
	Video         Video       `mapstructure:"video" yaml:"video"`
	IceServers    []IceServer `mapstructure:"ice_servers" yaml:"ice_servers"`
	Tls           Tls         `mapstructure:"tls" yaml:"tls"`
	Security      Security    `mapstructure:"security" yaml:"security"`
	Auth          Auth        `mapstructure:"auth" yaml:"auth"`
}

type Auth struct {
	UseSystemAuth      bool   `mapstructure:"use_system_auth" yaml:"use_system_auth"`
	HmacKey            string `mapstructure:"hmac_key" yaml:"hmac_key"`
	TokenValidityHours int    `mapstructure:"token_validity_hours" yaml:"token_validity_hours"`
	Users              []User `mapstructure:"users" yaml:"users"`
}

type User struct {
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
}

type Video struct {
	Bitrate   int `mapstructure:"bitrate" yaml:"bitrate"`
	Framerate int `mapstructure:"framerate" yaml:"framerate"`
}

type IceServer struct {
	Username   *string  `mapstructure:"username" yaml:"username,omitempty" json:"username,omitempty"`
	Credential *string  `mapstructure:"credential" yaml:"credential,omitempty" json:"credential,omitempty"`
	Urls       []string `mapstructure:"urls" yaml:"urls" json:"urls"`
}

type Tls struct {
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled"`
	CertFile string `mapstructure:"cert_file" yaml:"cert_file"`
	KeyFile  string `mapstructure:"key_file" yaml:"key_file"`
}

type Security struct {
	CheckOrigin       bool     `mapstructure:"check_origin" yaml:"check_origin"`
	AdditionalOrigins []string `mapstructure:"additional_origins" yaml:"additional_origins"`
}

func NewConfig() *Config {
	c := &Config{}
	c.viper = viper.New()

	// Setup config file paths and formats
	c.viper.SetConfigName("config")
	c.viper.AddConfigPath("/etc/webrdd/")
	c.viper.AddConfigPath("$HOME/.webrdd")
	c.viper.AddConfigPath(".")
	err := c.viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Println("Config file not found, using defaults")
	} else if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Setup environment variable reading
	c.viper.SetEnvPrefix("WEBRDD")
	c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	c.viper.AutomaticEnv()

	// Setup default values
	c.viper.SetDefault("bind_addresses", []string{":8080"})
	c.viper.SetDefault("video.bitrate", 8_000_000)
	c.viper.SetDefault("video.framerate", 30)
	c.viper.SetDefault("tls.cert_file", "./cert.pem")
	c.viper.SetDefault("tls.key_file", "./key.pem")
	c.viper.SetDefault("security.check_origin", true)

	err = c.viper.Unmarshal(c)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}
	return c
}

func (c *Config) String() string {
	yaml, err := yaml.Marshal(c)
	if err != nil {
		return ""
	}
	return string(yaml)
}
