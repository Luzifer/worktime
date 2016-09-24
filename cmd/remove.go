package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/Luzifer/worktime/schema"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove [day] <id>",
	Short:   "Deletes a time frame from the given day",
	Aliases: []string{"rm"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var inDay, inID string

		switch len(args) {
		case 2:
			inDay = args[0]
			inID = args[1]
		case 1:
			inDay = time.Now().Format("2006-01-02")
			inID = args[0]
		default:
			return fmt.Errorf("Please supply required arguments")
		}

		day, err := time.Parse("2006-01-02", inDay)
		if err != nil {
			return fmt.Errorf("'day' parameter seems to have a wrong format: %s", err)
		}

		doc, err := schema.LoadDay(cfg.couchDB, day, false)
		if err != nil {
			return err
		}

		nt := []*schema.Time{}
		for i := range doc.Times {
			if !strings.HasPrefix(doc.Times[i].ID, inID) {
				nt = append(nt, doc.Times[i])
			}
		}
		doc.Times = nt

		return doc.Save(cfg.couchDB)
	},
}

func init() {
	timeCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
