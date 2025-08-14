package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   Server
	Logger   Logger
	MemoryDB MemoryDB
}

type Logger struct {
	LogLvl   string
	BathSize int
}

type MemoryDB struct {
	TTL       time.Duration
	NumShards int
}

type Server struct {
	Port string
}

func LoadConfig() *Config {
	return &Config{
		Server: Server{
			Port: parseEnvString("PORT", "8080"),
		},
		Logger: Logger{
			LogLvl:   parseEnvString("LOG_LEVEL", "debug"),
			BathSize: parseEnvInt("BATCH_SIZE", 100),
		},
		MemoryDB: MemoryDB{
			TTL:       time.Duration(parseEnvInt("MEMORY_TTL", 30)) * time.Second,
			NumShards: parseEnvInt("NUM_SHARDS", 100),
		},
	}
}

func parseEnvString(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func parseEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return i
}
