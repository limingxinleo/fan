/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/limingxinleo/fan/config"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

// cosUploaderCmd represents the cosUploader command
var cosUploaderCmd = &cobra.Command{
	Use:   "cos:upload",
	Short: "上传文件到腾讯云 COS",
	Long:  `上传文件到腾讯云 COS`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(aurora.Yellow("请输入待上传文件或目录地址."))
			os.Exit(1)
		}

		file := args[0]
		cf := config.GetConfig(cmd)
		bucket, err := cmd.Flags().GetString("bucket")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if bucket != "" {
			cf.CosConfig.Bucket = bucket
		}

		if cf.CosConfig.Bucket == "" {
			fmt.Println(errors.New("请配置 cos_config.bucket"))
			os.Exit(1)
		}

		if cf.CosConfig.Region == "" {
			fmt.Println(errors.New("请配置 cos_config.region"))
			os.Exit(1)
		}

		if cf.CosConfig.BaseUri == "" {
			cf.CosConfig.BaseUri = "https://" + cf.CosConfig.Bucket + ".cos." + cf.CosConfig.Region + ".myqcloud.com"
		}

		stat, err := os.Stat(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		target, err := cmd.Flags().GetString("target")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if target == "" && !stat.IsDir() {
			target = generateCosTarget(file)
		}

		client, err := newCosClient(&cf.CosConfig)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if stat.IsDir() {
			prefix := strings.Trim(target, "/")
			if prefix == "" {
				fmt.Println("请使用 -t 输入上传的目标目录")
				os.Exit(1)
			}

			err := filepath.Walk(file, func(p string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				rel, err := filepath.Rel(file, p)
				if err != nil {
					return err
				}

				key := prefix + "/" + filepath.ToSlash(rel)
				err = uploadToCos(client, p, key)
				if err != nil {
					return err
				}

				fmt.Println(strings.TrimSuffix(cf.CosConfig.CdnUri, "/") + "/" + key)
				return nil
			})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			err := uploadToCos(client, file, target)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println(strings.TrimSuffix(cf.CosConfig.CdnUri, "/") + "/" + target)
		}
	},
}

func generateCosTarget(file string) string {
	ext := strings.TrimPrefix(filepath.Ext(file), ".")
	return fmt.Sprintf("%s/%s.%s", time.Now().Format("2006/01/02"), uuid.NewString(), ext)
}

func newCosClient(cf *config.CosConfig) (*cos.Client, error) {
	u, err := url.Parse(cf.BaseUri)
	if err != nil {
		return nil, err
	}

	return cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cf.SecretId,
			SecretKey: cf.SecretKey,
		},
	}), nil
}

func uploadToCos(client *cos.Client, file string, target string) error {
	_, err := client.Object.PutFromFile(context.Background(), target, file, nil)
	return err
}

func init() {
	rootCmd.AddCommand(cosUploaderCmd)

	cosUploaderCmd.Flags().StringP("target", "t", "", "上传路径")
	cosUploaderCmd.Flags().StringP("bucket", "b", "", "COS Bucket")
}
