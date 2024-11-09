package twitchapi

import (
	"encoding/json"
	"io"
	"os"
)

// this should fetch the weekly week.json and execute TwitchDownloaderCLI

const VIDEO_ID_PATH = "output/week.json"

func (t *TwitchApi) GetIDs() ([]string, error) {
	f, err := os.Open(VIDEO_ID_PATH)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)

	if err != nil {
		return nil, err
	}

	var (
		IDs  []string
		vods []TwitchVOD
	)

	err = json.Unmarshal(data, &vods)
	if err != nil {
		return nil, err
	}
	for _, vod := range vods {
		IDs = append(IDs, vod.ID)
	}

	return IDs, nil
}
