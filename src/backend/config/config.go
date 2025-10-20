package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Security SecurityConfig `mapstructure:"security"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
}

type StorageConfig struct {
	UploadDir     string `mapstructure:"uploadDir"`
	IncomingDir   string `mapstructure:"incomingDir"`
	PrivateDir    string `mapstructure:"privateDir"`
	MaxUploadSize int64  `mapstructure:"maxUploadSize"`
}

type SecurityConfig struct {
	AllowedExtensions []string `mapstructure:"allowedExtensions"`
	BlockedPaths      []string `mapstructure:"blockedPaths"`
}

type LoggingConfig struct {
	Level   string `mapstructure:"level"`
	Format  string `mapstructure:"format"`
	ToFile  bool   `mapstructure:"toFile"`
	Enabled bool   `mapstructure:"enabled"`
	LogDir  string `mapstructure:"logDir"`
}

// LoadConfig loads the configuration file and environment variables
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Configuration file settings
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Environment variable settings
	viper.AutomaticEnv()
	viper.SetEnvPrefix("SIMPLE_SERVER")

	// Try to read the configuration file
	configFileFound := false
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Configuration file not found, checking for example config file...")
			// Config file does not exist, check for example config file
			if _, err := os.Stat("config.yaml.example"); err == nil {
				log.Println("Example config file found, copying to actual config file...")
				// Copy example config file to actual config file
				if err := copyConfigFile("config.yaml.example", "config.yaml"); err == nil {
					log.Println("Config file created successfully, re-reading...")
					// Retry reading config file
					if err := viper.ReadInConfig(); err == nil {
						configFileFound = true
						log.Printf("Config file loaded successfully: %s", viper.ConfigFileUsed())
					}
				} else {
					log.Printf("Failed to copy config file: %v", err)
				}
			} else {
				log.Println("Example config file not found, will use default config")
			}
		} else {
			return nil, fmt.Errorf("Error reading config file: %w", err)
		}
	} else {
		configFileFound = true
		log.Printf("Config file loaded successfully: %s", viper.ConfigFileUsed())
	}

	// If config file not found, set default values
	if !configFileFound {
		setDefaultValues()
	}

	// Environment variable mapping (highest priority)
	if host := os.Getenv("HOST"); host != "" {
		viper.Set("server.host", host)
	}
	if port := os.Getenv("PORT"); port != "" {
		viper.Set("server.port", port)
	}

	// Parse config
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("Error parsing config: %w", err)
	}

	return config, nil
}

// setDefaultValues sets default config values
func setDefaultValues() {
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8000)
	viper.SetDefault("server.readTimeout", "30s")
	viper.SetDefault("server.writeTimeout", "30s")

	viper.SetDefault("storage.uploadDir", "./files")
	viper.SetDefault("storage.incomingDir", "./files/incoming")
	viper.SetDefault("storage.privateDir", "./files/private-files")
	viper.SetDefault("storage.maxUploadSize", 1000*1024*1024) // 1000MB

	viper.SetDefault("security.allowedExtensions", []string{".jpg", ".png", ".pdf", ".md", ".txt", ".html", ".css", ".js"})
	viper.SetDefault("security.blockedPaths", []string{"incoming", "private-files"})

	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.toFile", false)
	viper.SetDefault("logging.enabled", true)
	viper.SetDefault("logging.logDir", "./logs")
}

// copyConfigFile copies a config file
func copyConfigFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = srcFile.Seek(0, 0) // Reset file pointer
	if err != nil {
		return err
	}

	_, err = dstFile.ReadFrom(srcFile)
	return err
}

// GetListenAddr gets the listen address
func (c *Config) GetListenAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// PrintConfig prints current config information
func (c *Config) PrintConfig() {
	log.Println("=== Current Config Information ===")
	log.Printf("Server address: %s", c.GetListenAddr())
	log.Printf("Upload directory: %s", c.Storage.UploadDir)
	log.Printf("Max upload size: %d MB", c.Storage.MaxUploadSize/(1024*1024))
	log.Printf("Logging enabled: %v", c.Logging.Enabled)
	log.Printf("Log level: %s", c.Logging.Level)
	log.Printf("Log format: %s", c.Logging.Format)
	log.Printf("Log output to file: %v", c.Logging.ToFile)
	if c.Logging.ToFile {
		log.Printf("Log directory: %s", c.Logging.LogDir)
	}
	log.Printf("Allowed extensions: %v", c.Security.AllowedExtensions)
	log.Printf("Blocked paths: %v", c.Security.BlockedPaths)
}
