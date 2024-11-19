package config

import (
	"bufio"
	"os"
	"strings"
)

type AppConfig struct {
	AdminUser     string
	AdminPassword string
	DbName        string
	DbUser        string
	DbPassword    string
	DbHost        string
	DbPort        string
	SQLCommands   []string
}

func New(filename string) (*AppConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var (
		cfg           AppConfig
		curSqlCommand string
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			switch key {
			case "ADMIN_USERNAME":
				cfg.AdminUser = value
			case "ADMIN_PASSWORD":
				cfg.AdminPassword = value
			case "DB_NAME":
				cfg.DbName = value
			case "DB_USER":
				cfg.DbUser = value
			case "DB_PASSWORD":
				cfg.DbPassword = value
			case "DB_HOST":
				cfg.DbHost = value
			case "DB_PORT":
				cfg.DbPort = value
			}
		} else {
			curSqlCommand += line + " "
			if strings.HasSuffix(curSqlCommand, "; ") {
				cfg.SQLCommands = append(cfg.SQLCommands, curSqlCommand)
				curSqlCommand = ""
			}
		}
	}

	return &cfg, scanner.Err()
}
