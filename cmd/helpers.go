package cmd

import "time"

func parseTime(format, input string) (time.Time, error) {
	if input == "today" || input == "now" {
		input = time.Now().Format(format)
	}

	return time.Parse(format, input)
}
