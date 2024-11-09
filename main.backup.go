package main

// import (
// 	twitchapi "chat-embedder/twitch-api"

// 	"github.com/joho/godotenv"
// )

// func main() {
// 	godotenv.Load()
// 	// const (
// 	// 	TTV_DOWNLOADER = "./TwitchDownloaderCLI"
// 	// 	TTV_FEATURE    = "chatdownload"
// 	// 	VOD_ID         = "2294104188"
// 	// 	EMBED_IMGS     = "--embed-images"
// 	// 	IS_BTTV        = "--bttv=true"
// 	// 	FFZ_ARG        = "--ffz=false"
// 	// 	STV_ARG        = "--stv=false"
// 	// 	OUTPUT_PATH    = "-o"
// 	// 	OUTPUT_FILE    = "chat-embed1.json"
// 	// )

// 	// cmd := exec.Command(TTV_DOWNLOADER, TTV_FEATURE, "--id", VOD_ID, EMBED_IMGS, IS_BTTV, FFZ_ARG, STV_ARG, OUTPUT_PATH, OUTPUT_FILE)
// 	// stdout, err := cmd.Output()

// 	// if err != nil {
// 	// 	fmt.Printf("Error: %v\n", err)
// 	// } else {
// 	// 	fmt.Println("Executed successfully")
// 	// 	fmt.Println(string(stdout))
// 	// }
// 	api := twitchapi.Init()
// 	err := api.WriteVideoID("405606875", "week", "output/week.json")

// 	if err != nil {
// 		panic(err)
// 	}
// }
