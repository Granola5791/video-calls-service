package config

import (
	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.SetConfigName("constants")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}

func GetInt(key string) int {
	if viper.IsSet(key) {
		return viper.GetInt(key)
	}
	panic(viper.GetString("error.missing_config") + key)
}

func GetString(key string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}
	panic(viper.GetString("error.missing_config") + key)
}

func GetBool(key string) bool {
	if viper.IsSet(key) {
		return viper.GetBool(key)
	}
	panic(viper.GetString("error.missing_config") + key)
}

func GetFloat64(key string) float64 {
	if viper.IsSet(key) {
		return viper.GetFloat64(key)
	}
	panic(viper.GetString("error.missing_config") + key)
}

func GetUint32(key string) uint32 {
	if viper.IsSet(key) {
		return viper.GetUint32(key)
	}
	panic(viper.GetString("error.missing_config") + key)
}

func GetUint8(key string) uint8 {
	if viper.IsSet(key) {
		return viper.GetUint8(key)
	}
	panic(viper.GetString("error.missing_config") + key)
}
