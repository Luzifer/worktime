package cmd

import (
	"fmt"
	"net/url"
	"os"

	couch "github.com/lancecarlson/couchgo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg = struct {
	CfgFile string

	couchDB *couch.Client
}{}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "worktime",
	Short: "Manage worktimes in CouchDB",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		if viper.GetString("couchdb") == "" {
			return fmt.Errorf("couchdb url needs to be set")
		}

		u, err := url.Parse(viper.GetString("couchdb"))
		if err != nil {
			return err
		}

		cfg.couchDB = couch.NewClient(u)

		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().String("couchdb", "", "URL to access couchdb (http://user:pass@host:port/database)")
	viper.BindPFlag("couchdb", RootCmd.PersistentFlags().Lookup("couchdb"))

	RootCmd.PersistentFlags().StringVar(&cfg.CfgFile, "config", "", "config file (default is $HOME/.worktime.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfg.CfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfg.CfgFile)
	}

	viper.SetConfigName(".worktime") // name of config file (without extension)
	viper.AddConfigPath("$HOME")     // adding home directory as first search path
	viper.AutomaticEnv()             // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
