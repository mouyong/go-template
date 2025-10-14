package initialization

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	HttpHost       string   `json:"httpHost"`
	HttpPort       int      `json:"httpPort"`
	TrustedProxies []string `json:"trustedProxies"`
	DbHost         string   `json:"dbHost"`
	DbPort         int      `json:"dbPort"`
	DbDatabase     string   `json:"dbDatabase"`
	DbUsername     string   `json:"dbUsername"`
	DbPassword     string   `json:"dbPassword"`
	MqHost         string   `json:"mqHost"`
	MqPort         int      `json:"mqPort"`
}

var AppConfig Config

func LoadConfig(cfg string) Config {
	viper.SetConfigFile(cfg)
	viper.ReadInConfig()
	viper.AutomaticEnv()

	AppConfig = Config{
		HttpHost:       getViperStringValue("HTTP_Host", "0.0.0.0"),
		HttpPort:       getViperIntValue("HTTP_PORT", 3000),
		TrustedProxies: getViperStringArray("TRUSTED_PROXIES", nil),
		DbHost:         viper.GetString("DB_HOST"), // 不使用默认值，保持空字符串
		DbPort:         getViperIntValue("DB_PORT", 3306),
		DbDatabase:     viper.GetString("DB_DATABASE"), // 不使用默认值，保持空字符串
		DbUsername:     getViperStringValue("DB_USERNAME", "root"),
		DbPassword:     getViperStringValue("DB_PASSWORD", "root"),
		MqHost:         viper.GetString("MQ_HOST"), // 不使用默认值，保持空字符串
		MqPort:         getViperIntValue("MQ_PORT", 5672),
	}
	fmt.Println("读取到的配置信息：", AppConfig)

	return AppConfig
}

func getViperStringValue(key string, defaultValue string) string {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getViperStringArray(key string, defaultValue []string) []string {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	raw := strings.Split(value, ",")
	return raw
}

func getViperIntValue(key string, defaultValue int) int {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("Invalid value for %s, using default value %d\n", key, defaultValue)
		return defaultValue
	}
	return intValue
}

func getViperBoolValue(key string, defaultValue bool) bool {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		fmt.Printf("Invalid value for %s, using default value %v\n", key, defaultValue)
		return defaultValue
	}
	return boolValue
}
