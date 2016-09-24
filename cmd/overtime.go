package cmd

import (
	"fmt"
	"time"

	"github.com/Luzifer/worktime/schema"
	"github.com/spf13/cobra"
)

// overtimeCmd represents the overtime command
var overtimeCmd = &cobra.Command{
	Use:   "overtime",
	Short: "Shows total overtime over all time",
	RunE: func(cmd *cobra.Command, args []string) error {
		overtime, err := schema.GetOvertime(cfg.couchDB, time.Time{})
		if err != nil {
			return err
		}

		fmt.Printf("Total overtime: %s\n", time.Duration(overtime.Value*float64(time.Hour)))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(overtimeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// overtimeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// overtimeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
