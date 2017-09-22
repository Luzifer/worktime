package cmd

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/Luzifer/worktime/schema"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <day> <start> <end> [tag [tag...]]",
	Short: "Add a time frame to the given day",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return fmt.Errorf("Please supply required arguments")
		}

		day, err := parseTime("2006-01-02", args[0])
		if err != nil {
			return fmt.Errorf("'day' parameter seems to have a wrong format: %s", err)
		}

		doc, err := schema.LoadDay(cfg.couchDB, day, true)
		if err != nil {
			return err
		}

		doc.Times = append(doc.Times, &schema.Time{
			ID:    fmt.Sprintf("%x", sha1.Sum([]byte(time.Now().String()))),
			Start: args[1],
			End:   args[2],
			Tags:  args[3:],
		})

		return doc.Save(cfg.couchDB)
	},
}

func init() {
	timeCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
