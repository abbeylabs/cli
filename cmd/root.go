/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
	"os"
)

const (
	letters = `
 _______  _______  _______  _______  __   __    ___      _______  _______  _______ 
|   _   ||  _    ||  _    ||       ||  | |  |  |   |    |   _   ||  _    ||       |
|  |_|  || |_|   || |_|   ||    ___||  |_|  |  |   |    |  |_|  || |_|   ||  _____|
|       ||       ||       ||   |___ |       |  |   |    |       ||       || |_____ 
|       ||  _   | |  _   | |    ___||_     _|  |   |___ |       ||  _   | |_____  |
|   _   || |_|   || |_|   ||   |___   |   |    |       ||   _   || |_|   | _____| |
|__| |__||_______||_______||_______|  |___|    |_______||__| |__||_______||_______|
`
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "abbey",
	Short: "",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cc.Init(&cc.Config{
		RootCmd:  rootCmd,
		Headings: cc.Green + cc.Bold + cc.Underline,
		Commands: cc.Magenta + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Blue + cc.Bold,
		Flags:    cc.Yellow + cc.Bold,
	})

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.abbeycli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
