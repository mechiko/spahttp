package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mechiko/utility"
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
	inHome         bool
	cwd            string // как относительный путь для окружения, "." путь домашнего каталога или пусто
	absConfigPath  string
	absDbPath      string
	absLogPath     string
	configuration  *Configuration
	configFileName string
	warning        string
}

const modError = "pkg:config"

func New(inUserHome bool) (cfg *Config, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s panic %v", modError, r)
		}
	}()

	configName := Name
	viperOrigin := viper.GetViper()
	cfg = &Config{
		Viper:  viperOrigin,
		cwd:    ".",
		inHome: inUserHome,
	}
	if cfg.inHome {
		cfg.cwd = osUserHomeDir()
	}

	// прописываем пути каталогов absConfigPath absDbPath absLogPath
	if err := cfg.initPaths(); err != nil {
		return nil, fmt.Errorf("%s: failed to initialize paths: %w", modError, err)
	}

	cfg.configFileName = filepath.Join(cfg.absConfigPath, configName+".toml")
	viperConfigPath := configPath
	if !strings.HasPrefix(viperConfigPath, ".") {
		viperConfigPath = "." + viperConfigPath
	}
	viperConfigPath = filepath.Join(cfg.cwd, viperConfigPath)
	viperOrigin.SetConfigName(configName)
	viperOrigin.SetConfigType("toml")
	viperOrigin.AddConfigPath(viperConfigPath)

	if err := viperOrigin.MergeConfig(strings.NewReader(string(TomlConfig))); err != nil {
		return nil, fmt.Errorf("%s %w", modError, err)
	}
	if err := viperOrigin.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// другая ошибка
			return nil, fmt.Errorf("%s %w", modError, err)
		} else {
			// конфиг файл не найден
			cfg.warning = fmt.Sprintf("config file ('%s') not found\n", cfg.configFileName)
		}
	}
	cfg.configuration = &Configuration{}
	if err := viperOrigin.Unmarshal(cfg.configuration); err != nil {
		return nil, fmt.Errorf("%s %w", modError, err)
	}
	// игнорируем ошибку этот вызов для первого сохранения файла конфига когда его нет
	err = viperOrigin.SafeWriteConfig()
	if err != nil {
		cfg.warning += fmt.Sprintf("%s\n", err.Error())
	}
	return cfg, nil
}

// вроде как возвращает копию структуры через разыменование
// проверено :)
func (c *Config) Configuration() *Configuration {
	sc := *c.configuration
	return &sc
}

func osUserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to HOME environment variable
		return os.Getenv("HOME")
	}
	return home
}

// создаем папки по конфигурации
func (c *Config) initPaths() (err error) {
	if c.absConfigPath, err = c.createDirectory(configPath); err != nil {
		return err
	}
	if c.absLogPath, err = c.createDirectory(logPath); err != nil {
		return err
	}
	if c.absDbPath, err = c.createDirectory(dbPath); err != nil {
		return err
	}
	return nil
}

// проверяем путь каталога, если он абсолютный (линукс обычно) то возвращаем его если он существует
// иначе создаем пути которых нет
// если путь относительный добавляем точку и аналогично абсолютному
func (c *Config) createDirectory(path string) (string, error) {
	if filepath.IsAbs(path) {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("config create dir error %w", err)
		}
		if !utility.PathOrFileExists(abs) {
			if err := pathCreate(abs); err != nil && !errors.Is(err, fs.ErrExist) {
				return "", fmt.Errorf("create path %s error %w", abs, err)
			}
		}
		return abs, nil
	}
	if !strings.HasPrefix(path, ".") {
		path = "." + path
	}
	fullPath := filepath.Join(c.cwd, path)
	if err := pathCreate(fullPath); err != nil && !errors.Is(err, fs.ErrExist) {
		return "", fmt.Errorf("невозможно создать путь %s %w", fullPath, err)
	}
	return filepath.Abs(fullPath)
}

// создаст весь путь вложенных каталогов
func pathCreate(path string) error {
	if path != "" {
		if err := os.MkdirAll(path, os.ModePerm); err != nil { // создает весь путь
			// if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) ConfigPath() string {
	return c.absConfigPath
}

func (c *Config) DbPath() string {
	return c.absDbPath
}

func (c *Config) LogPath() string {
	return c.absLogPath
}
