package cmd

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/golang-module/carbon/v2"
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
		var t string
		if len(args) != 0 {
			dumpTime(args[0])
		}

		in := bufio.NewReader(os.Stdin)
		scanner := bufio.NewScanner(in)
		for true {
			fmt.Println(aurora.Yellow("Please Input TimeStamp or Datetime:"))
			scanner.Scan()
			t = string(scanner.Bytes())
			dumpTime(t)
		}
	},
}

func toTime(value string) time.Time {
	switch value {
	case "now":
		return carbon.Now().ToStdTime()
	case "today":
		return carbon.Now().StartOfDay().ToStdTime()
	case "tomorrow":
		return carbon.Tomorrow().ToStdTime()
	default:
		return carbon.Parse(value).ToStdTime()
	}
}

func dumpTime(val string) {
	var dt time.Time
	if ts, err := strconv.ParseInt(val, 10, 64); err == nil {
		dt = carbon.CreateFromTimestamp(ts).ToStdTime()
	} else {
		dt = toTime(val)
	}

	tbl := table.New("TimeStamp", "DateTime")
	tbl.WithHeaderFormatter(color.New(color.FgGreen, color.Underline).SprintfFunc())
	tbl.AddRow(dt.Unix(), dt.Format(time.DateTime))
	tbl.Print()
	fmt.Println("")
}

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
