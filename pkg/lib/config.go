package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/imdario/mergo"
)

type Config struct {
	Environment string      `json:"environment"`
	LogOutput   string      `json:"log_output"`
	LogLevel    string      `json:"log_level"`
	TZ          string      `json:"tz"`
	Timer       int         `json:"timer"`
	TvTime      TvTime      `json:"tv_time"`
	Plex        Plex        `json:"plex"`
	FileStorage FileStorage `json:"file_storage"`
}

type TvTime struct {
	Token          TvTimeToken `json:"token"`
	AcceptLanguage string      `json:"accept_language"`
}
type TvTimeToken struct {
	Symfony      string `json:"symfony"`
	TvstRemember string `json:"tvst_remember"`
}
type Plex struct {
	BaseUrl      string `json:"base_url"`
	Token        string `json:"token"`
	AccountId    int    `json:"account_id"`
	InitViewedAt string `json:"init_viewed_at"`
}
type FileStorage struct {
	Filename string `json:"filename"`
}

type configs []string

var (
	listConfigPath configs
	globalConfig   *Config
)

func checkConfigPath(configsPath []string) error {
	const names = "__config.go__ : checkConfigPath"
	for _, allowedValue := range []string{"development.json", "staging.json", "production.json"} {
		r := regexp.MustCompile(fmt.Sprintf("(?mi)%s$", allowedValue))
		for _, c := range configsPath {
			if r.FindString(c) != "" {
				if err := SetConfigPath(c); err != nil {
					panic(fmt.Sprintf("%s | %s", names, err))
				}
				return nil
			}
		}
	}
	return nil
}

func SetConfigPath(configPath string) error {
	listConfigPath = append(listConfigPath, configPath)
	return nil
}

func newConfig(config *Config) error {
	for _, configPath := range listConfigPath {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return err
		}
		// var jsonData map[string]interface{}
		var jsonData Config
		if err := json.Unmarshal(data, &jsonData); err != nil {
			log.Fatalf("Failed to unmarshal JSON data from file %s: %v", configPath, err)
		}
		// Merge the data into the mergedData map, overwriting existing values
		if err := mergo.Merge(config, jsonData, mergo.WithOverride); err != nil {
			log.Fatalf("Failed to merge JSON data from file %s: %v", configPath, err)
		}
	}
	return nil
}

func SetupConfig(configsPath []string) error {
	if len(configsPath) == 0 {
		return errors.New(`pass the path list`)
	}
	if err := checkConfigPath(configsPath); err != nil {
		return err
	}
	return nil
}

func GetConfig() Config {
	const names = "__config.go__ : GetConfig"
	if globalConfig == nil {
		globalConfig = &Config{}
		if err := newConfig(globalConfig); err != nil {
			panic(fmt.Sprintf("%s | %s", names, err))
		}
	}
	if globalConfig.Plex.InitViewedAt == "" {
		globalConfig.Plex.InitViewedAt = time.Now().Format("2006-01-02 15:04:05")
	}
	return *globalConfig
}
