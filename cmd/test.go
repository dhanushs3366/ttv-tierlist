/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"chat-embedder/utils"
	"fmt"
	"sync"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.RankViewers(50, fileOutput)
	},
}

var counter = 0

func increment(wg *sync.WaitGroup) {
	counter++
	wg.Done()
}

func Bain() {
	var wg sync.WaitGroup
	expectedCounter := 1000

	for i := 0; i < expectedCounter; i++ {
		wg.Add(1)
		go increment(&wg)
	}

	wg.Wait()
	fmt.Println("Expected Counter:", expectedCounter)
	fmt.Println("Actual Counter:", counter)
	// Check for race condition
	if expectedCounter != counter {
		fmt.Println("Race condition detected!")
	} else {
		fmt.Println("No race condition detected.")
	}
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
