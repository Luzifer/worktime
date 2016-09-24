package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Luzifer/worktime/schema"
	"github.com/spf13/cobra"
)

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use:   "patch <day> <time id> <start> <end> [[+/-]tag [[+/-]tag]]",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 4 {
			return fmt.Errorf("Please supply required arguments")
		}

		day, err := time.Parse("2006-01-02", args[0])
		if err != nil {
			return fmt.Errorf("'day' parameter seems to have a wrong format: %s", err)
		}

		doc, err := schema.LoadDay(cfg.couchDB, day, false)
		if err != nil {
			return err
		}

		for _, t := range doc.Times {
			if strings.HasPrefix(t.ID, args[1]) {
				if args[2] != "=" {
					t.Start = args[2]
				}
				if args[3] != "=" {
					t.End = args[3]
				}
				for _, tag := range args[4:] {
					t.Tag(tag)
				}
				return doc.Save(cfg.couchDB)
			}
		}

		return errors.New("No time frame with the given ID was found")
	},
}

func init() {
	timeCmd.AddCommand(patchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// patchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// patchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
