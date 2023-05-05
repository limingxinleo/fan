package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/logrusorgru/aurora"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "t:fmt {time}",
	Short: "格式化时间格式",
	Long:  `格式化时间格式`,
	Run: func(cmd *cobra.Command, args []string) {
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		if len(args) != 0 {
			t := args[0]
			var dt time.Time
			if ts, err := strconv.ParseInt(t, 10, 64); err == nil {
				dt = time.Unix(ts, 0)
			} else {
				dt, err = time.ParseInLocation(time.DateTime, t, zone)
				if err != nil {
					fmt.Println(aurora.Red(err))
					os.Exit(1)
				}
			}

			tbl := table.New("TimeStamp", "DateTime")
			tbl.WithHeaderFormatter(headerFmt)
			tbl.AddRow(dt.Unix(), dt.Format(time.DateTime))
			tbl.Print()
		}
	},
}

var zone = time.FixedZone("CST", 8*3600)

func init() {
	rootCmd.AddCommand(timeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// timeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// timeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
