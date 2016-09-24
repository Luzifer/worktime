package cmd

import (
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/Luzifer/worktime/schema"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [day]",
	Short: "Display a summary of the given / current day",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			args = []string{time.Now().Format("2006-01-02")}
		}

		day, err := time.Parse("2006-01-02", args[0])
		if err != nil {
			return fmt.Errorf("'day' parameter seems to have a wrong format: %s", err)
		}

		doc, err := schema.LoadDay(cfg.couchDB, day, false)
		if err != nil {
			return err
		}

		overtime, err := schema.GetOvertime(cfg.couchDB, day)
		if err != nil {
			return err
		}

		tplSrc, err := Asset("templates/show.tpl")
		if err != nil {
			return err
		}

		tpl, err := template.New("show").Parse(string(tplSrc))
		if err != nil {
			return err
		}

		return tpl.Execute(os.Stdout, map[string]interface{}{
			"day":      doc,
			"overtime": time.Duration(overtime.Value * float64(time.Hour)),
		})

	},
}

func init() {
	RootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
