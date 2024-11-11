/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	twitchapi "chat-embedder/twitch-api"
	"chat-embedder/utils"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
)

// rankCmd represents the rank command
var rankCmd = &cobra.Command{
	Use:   "rank",
	Short: "rank command ranks the viewers of a streamer on twitch for a given time period",
	Long: `rank command uses ttv command to fetch the chat logs of a twitch streamer,
	processes the json files created by it and ranks the viewers of the said streamer ordered by their comment count
	for a said time period (day/week/month/all->100days)
	`,
	Run: func(cmd *cobra.Command, args []string) {
		rateIt(cmd, args)
		Rank(cmd, args)
	},
}

func Rank(cmd *cobra.Command, args []string) {

	rankedUsers := utils.RankViewers(rankLimit, filepath.Dir(fileOutput))
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
		wg.Add(1)
		go utils.CreateUserProfiles(user, &wg, profilesOutput)
	}
	wg.Wait()
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
