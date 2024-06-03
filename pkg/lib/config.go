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
	Environment string  `json:"environment"`
	Storage     Storage `json:"storage"`
	LogLevel    string  `json:"log_level"`
	LogOutput   string  `json:"log_output"`
	Plex        Plex    `json:"plex"`
	Timer       int     `json:"timer"`
	TvTime      TvTime  `json:"tv_time"`
	TZ          string  `json:"tz"`
	Mailer      Mailer  `json:"mailer"`
}

type TvTime struct {
	AcceptLanguage     string `json:"accept_language"`
	BaseUrl            string `json:"base_url"`
	CountryCode        string `json:"country_code"`
	Locale             string `json:"locale"`
	SearchUrl          string `json:"search_url"`
	SeriesUrl          string `json:"series_url"`
	SeasonsUrl         string `json:"seasons_url"`
	WatchedEpisodesUrl string `json:"watched_episodes_url"`
	Token              string `json:"token"`
	TokenType          string `json:"token_type"`
}
type Plex struct {
	AccountId    int    `json:"account_id"`
	BaseUrl      string `json:"base_url"`
	InitViewedAt string `json:"init_viewed_at"`
	ViewedAtSort string `json:"viewed_at_sort"`
	Token        string `json:"token"`
}
type Storage struct {
	Filename string `json:"filename"`
}

type Mailer struct {
	Recipient   string `json:"recipient"`
	TemplateDir string `json:"template_dir"`
	SMTP        SMTP   `json:"smtp"`
}

type SMTP struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
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
