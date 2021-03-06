package cmd

import (
	"crypto/sha1"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/Luzifer/worktime/schema"
)

// trackCmd represents the track command
var trackCmd = &cobra.Command{
	Use:   "track [tag [tag]]",
	Short: "Track a time frame from the command start until interruption",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			c     = make(chan os.Signal, 1)
			start = time.Now()
			run   = true
		)

		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		fmt.Println("Press Ctrl+C to stop time tracking...")

		for run {
			select {
			case <-c: // Blocks until a signal arrives
				run = false

			case <-time.Tick(time.Second):
				fmt.Printf("\r\x1b[K")
				fmt.Printf("Tracking since %s (%s)...", start.Format(time.RFC1123), time.Since(start))
			}
		}

		end := time.Now()

		doc, err := schema.LoadDay(cfg.couchDB, start, true)
		if err != nil {
			return err
		}

		doc.Times = append(doc.Times, &schema.Time{
			ID:    fmt.Sprintf("%x", sha1.Sum([]byte(time.Now().String()))),
			Start: start.Format("15:04:05"),
			End:   end.Format("15:04:05"),
			Tags:  args,
		})

		return doc.Save(cfg.couchDB)
	},
}

func init() {
	timeCmd.AddCommand(trackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// trackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// trackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
