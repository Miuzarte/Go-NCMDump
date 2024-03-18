package main

import (
	"bytes"
	"crypto/aes"
	_ "embed"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	CoreKey             string `yaml:"coreKey"`
	MetaKey             string `yaml:"metaKey"`
	InputDir            string `yaml:"inputDir"`
	OutputDir           string `yaml:"outputDir"`
	CoverOutput         bool   `yaml:"coverOutput"`
	CoverEmbed          bool   `yaml:"coverEmbed"`
	HighDefinitionCover bool   `yaml:"highDefinitionCover"`
	MultiThread         bool   `yaml:"multiThread"`
}

//go:embed config_default.yml
var defaultConfig []byte

func (c *Config) init() *Config {
	var configPath string

	_, err1 := os.Stat(filepath.Join(exeDir, "config.yml"))
	_, err2 := os.Stat(filepath.Join(exeDir, "config.yaml"))
	if os.IsNotExist(err1) && os.IsNotExist(err2) {
		createConfigAndExit()
	} else if err1 == nil {
		configPath = filepath.Join(exeDir, "config.yml")
	} else if err2 == nil {
		configPath = filepath.Join(exeDir, "config.yaml")
	} else {
		log.Fatal("failed to get config path: ", err1, err2)
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("failed to read config: ", err)
	}

	err = yaml.Unmarshal(configBytes, &c)
	if err != nil {
		log.Fatal("failed to unmarshal config: ", err)
	}

	if c.CoreKey == "" || c.MetaKey == "" {
		log.Fatal("coreKey and metaKey must be set")
	}

	if len(c.CoreKey) != 32 || len(c.MetaKey) != 32 {
		log.Fatal("coreKey and metaKey must be 32 bytes long strings")
	}

	switch mode {
	case DIRECTMODE:
		c.InputDir = ""
		c.OutputDir = ""

	case CONFIGMODE:
		if c.InputDir == "" {
			c.InputDir = exeDir
		}
		if c.OutputDir == "" {
			c.OutputDir = c.InputDir
		}
	}

	return c
}

func createConfigAndExit() {
	fs, err := os.Create(filepath.Join(exeDir, "config.yml"))
	if err != nil {
		log.Fatal("failed to create config: ", err)
	}
	_, err = fs.Write(defaultConfig)
	if err != nil {
		log.Fatal("failed to write config: ", err)
	}
	log.Info("config created, please edit it and restart the program")
	os.Exit(0)
}

func checkKeys(coreKey, metaKey string) (metaOk, coreOk bool) {
	bs := 16

	cBlock, _ := aes.NewCipher([]byte(coreKey))
	coreCtext := make([]byte, len(corePlaintext))
	for i := 0; i < len(corePlaintext); i += bs {
		cBlock.Encrypt(coreCtext[i:i+bs], corePlaintext[i:i+bs])
	}

	mBlock, _ := aes.NewCipher([]byte(metaKey))
	metaCtext := make([]byte, len(metaPlaintext))
	for i := 0; i < len(metaPlaintext); i += bs {
		mBlock.Encrypt(metaCtext[i:i+bs], metaPlaintext[i:i+bs])
	}

	coreOk = bytes.Equal(coreCtext, coreCiphertext)
	metaOk = bytes.Equal(metaCtext, metaCiphertext)
	return
}
