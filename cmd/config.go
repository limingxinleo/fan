package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gobeam/stringy"
	"github.com/limingxinleo/fan/config"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
	"reflect"
	"strings"
)

// configCmd represents the entity command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show Config",
	Long:  `展示配置信息`,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.GetConfig(cmd)

		res, err := json.Marshal(c)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(string(res))
	},
}

var configSetterCmd = &cobra.Command{
	Use:   "config:set {key} {value}",
	Short: "设置配置信息",
	Long:  `设置配置信息 config:set oss_config.endpoint test`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println(aurora.Yellow("请输入 key 和 value 值，例如 config:set oss_config.endpoint test"))
			os.Exit(1)
		}

		key := args[0]
		value := args[1]

		cf := config.GetConfig(cmd)

		keys := strings.Split(key, ".")

		ref := reflect.ValueOf(cf).Elem()

		for _, s := range keys {
			str := stringy.New(s)
			ref = ref.FieldByName(str.CamelCase())
		}

		ref.SetString(value)

		config.SetConfig(cmd, cf)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(configSetterCmd)
}
