package cmd

import (
	"fmt"
	"strings"

	"github.com/Luzifer/worktime/schema"
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag <day> [time-id] <[+/-]tag>",
	Short: "Adds or removes a tag from the day or time entry inside the day",
	RunE: func(cmd *cobra.Command, args []string) error {
		var timeId, tag string
		switch len(args) {
		case 2:
			tag = args[1]
		case 3:
			timeId = args[1]
			tag = args[2]
		default:
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

		if timeId != "" {
			for _, t := range doc.Times {
				if strings.HasPrefix(t.ID, timeId) {
					t.Tag(tag)
					return doc.Save(cfg.couchDB)
				}
			}
			return fmt.Errorf("Could not find time with ID '%s'", timeId)
		}

		doc.Tag(tag)
		return doc.Save(cfg.couchDB)
	},
}

func init() {
	RootCmd.AddCommand(tagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
