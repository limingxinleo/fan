package cmd

import (
	"bufio"
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
		var t string
		if len(args) != 0 {
			dump(args[0])
		}

		in := bufio.NewReader(os.Stdin)
		scanner := bufio.NewScanner(in)
		for true {
			fmt.Println(aurora.Yellow("Please Input TimeStamp or Datetime:"))
			scanner.Scan()
			t = string(scanner.Bytes())
			dump(t)
		}
	},
}

func toTime(value string) time.Time {
	switch value {
	case "now":
		return time.Now()
	case "today":
		y, m, d := time.Now().Date()
		return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	case "tomorrow":
		y, m, d := time.Now().Add(time.Hour * 24).Date()
		return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	default:
		dt, err := time.ParseInLocation(time.DateTime, value, time.Local)
		if err != nil {
			fmt.Println(aurora.Red(err))
			os.Exit(1)
		}
		return dt
	}
}

func dump(val string) {
	var dt time.Time
	if ts, err := strconv.ParseInt(val, 10, 64); err == nil {
		dt = time.Unix(ts, 0)
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
