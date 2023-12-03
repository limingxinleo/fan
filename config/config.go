package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type Config struct {
	OssConfig OssConfig `json:"oss_config"`
	PdfConfig PdfConfig `json:"pdf_config"`
}

type OssConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Bucket          string `json:"bucket"`
	BaseUri         string `json:"base_uri"`
}

type PdfConfig struct {
	MeteredKey string `json:"metered_key"`
}

func SetConfig(cmd *cobra.Command, config *Config) {
	path, err := cmd.Flags().GetString("config")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	InitConfig(path, config)
}

func GetConfig(cmd *cobra.Command) *Config {
	path, err := cmd.Flags().GetString("config")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bf, err := os.ReadFile(path)
	if err != nil {
		return InitConfig(path, nil)
	}

	c := &Config{}
	err = json.Unmarshal(bf, c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return c
}

func InitConfig(path string, config *Config) *Config {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	w := bufio.NewWriter(f)

	if config == nil {
		config = &Config{}
	}

	bt, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = w.Write(bt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = w.Flush()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return config
}
