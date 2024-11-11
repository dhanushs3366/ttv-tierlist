/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	twitchapi "chat-embedder/twitch-api"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	username   string
	period     string
	ttvClient  *twitchapi.TwitchApi
	fileOutput string
	fileCount  uint = 0
	rankLimit  uint
	mutex      sync.Mutex
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chat-embedder",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args[:])
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	godotenv.Load()
	ttvClient = twitchapi.Init()

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(ttvCmd)
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringVarP(&username, "username", "u", "minicake", "uhmm pass in ur twitch username it will print out their top viewers")

	ttvCmd.Flags().StringVarP(&username, "username", "u", "minicake", "uhmm pass in ur twitch username it will print out their top viewers")
	ttvCmd.Flags().StringVarP(&period, "period", "p", "week", "uhmmm this for the # of week u wanna analayze")
	ttvCmd.Flags().StringVarP(&fileOutput, "output", "o", fmt.Sprintf("%s.json", period), "uhm for the output ig")
	// add a rank command

	// rankCmd.Flags().StringVarP(&username, "username", "u", "minicake", "uhmm pass in ur twitch username it will print out their top viewers")
	// rankCmd.Flags().StringVarP(&period, "period", "p", "week", "uhmmm this for the # of week u wanna analayze")
	// rankCmd.Flags().StringVarP(&fileOutput, "output", "o", fmt.Sprintf("%s.json", period), "uhm for the output ig")
	rankCmd.Flags().UintVarP(&rankLimit, "limit", "l", 25, "set the limit for viewers default is 25")
}
