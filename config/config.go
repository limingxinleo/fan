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
}

type OssConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

func GetConfig(cmd *cobra.Command) *Config {
	path, err := cmd.Flags().GetString("config")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bf, err := os.ReadFile(path)
	if err != nil {
		return InitConfig(path)
	}

	c := &Config{}
	err = json.Unmarshal(bf, c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return c
}

func InitConfig(path string) *Config {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(123123)
		fmt.Println(err)
		os.Exit(1)
	}

	w := bufio.NewWriter(f)

	config := &Config{}

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
