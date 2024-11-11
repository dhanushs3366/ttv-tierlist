/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	twitchapi "chat-embedder/twitch-api"
	"chat-embedder/utils"
	"fmt"
	"sync"

	"github.com/spf13/cobra"
)

// rankCmd represents the rank command
var rankCmd = &cobra.Command{
	Use:   "rank",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: Test,
}

func Test(cmd *cobra.Command, args []string) {
	// err := ttvCmd.Execute()
	// if err != nil {
	// 	color.Red("Error trying to execute ttv cmd %s", err.Error())
	// }

	rankedUsers := utils.RankViewers(rankLimit)
	userDetails := make(chan twitchapi.TwitchUser, len(rankedUsers))
	fmt.Println(len(rankedUsers))
	var wg sync.WaitGroup

	for _, user := range rankedUsers {
		wg.Add(1)
		go ttvClient.GetUserDetails(user.UserID, userDetails, &wg)
	}
	wg.Wait()
	close(userDetails)
	for user := range userDetails {
		fmt.Println(user.ProfileImageURL, user.DisplayName)
	}
}

func init() {
	rootCmd.AddCommand(rankCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rankCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rankCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
