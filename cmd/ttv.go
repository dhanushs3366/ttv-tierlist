/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	twitchapi "chat-embedder/twitch-api"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ttvCmd represents the ttv command
var ttvCmd = &cobra.Command{
	Use:   "ttv",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: rateIt,
}

func rateIt(cmd *cobra.Command, args []string) {
	userID, err := ttvClient.GetUserID(username)
	if err != nil {
		color.Red("username not found")
	}

	vods, err := ttvClient.GetVideoIDs(userID, twitchapi.TIME_PERIOD(period))
	if err != nil {
		color.Red("couldnt get VOD IDs :(")
	}
	var vodIDs []string
	var wg sync.WaitGroup
	for _, vod := range vods {
		vodIDs = append(vodIDs, vod.ID)
	}

	for _, vodID := range vodIDs {
		wg.Add(1)
		go writeToJSON(vodID, &wg)
	}
	wg.Wait()
	fmt.Println("Its done check the output file")

}

func init() {
	rootCmd.AddCommand(ttvCmd)

}

func writeToJSON(vodID string, wg *sync.WaitGroup) {
	defer wg.Done()

	mutex.Lock()
	fileCount++
	randomNumber := fileCount
	mutex.Unlock()

	const (
		TTV_DOWNLOADER = "./TwitchDownloaderCLI"
		TTV_FEATURE    = "chatdownload"
		EMBED_IMGS     = "--embed-images"
		IS_BTTV        = "--bttv=true"
		FFZ_ARG        = "--ffz=false"
		STV_ARG        = "--stv=true"
		OUTPUT_PATH    = "-o"
	)
	ext := filepath.Ext(fileOutput)
	rawFileName := strings.TrimSuffix(fileOutput, ext)
	cmd := exec.Command(TTV_DOWNLOADER, TTV_FEATURE, "--id", vodID, OUTPUT_PATH, fmt.Sprintf("%s-%d%s", rawFileName, randomNumber, ext), "--collision", "Overwrite")

	stdout, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Execution error: %v\nOutput: %s\n", err, string(stdout))
		return
	}

	fmt.Println("Executed successfully")
	fmt.Println(string(stdout))
}
