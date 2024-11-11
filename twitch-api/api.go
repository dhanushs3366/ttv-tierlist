package twitchapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type TwitchApi struct {
	client    *http.Client
	authToken string
	EndPoint  string
	clienID   string
}

type TIME_PERIOD string
type FILE_TYPE string

const (
	ALL   TIME_PERIOD = "all"
	DAY   TIME_PERIOD = "day"
	MONTH TIME_PERIOD = "month"
	WEEK  TIME_PERIOD = "week"
)

const (
	TEXT_FILE FILE_TYPE = "txt"
	JSON_FILE FILE_TYPE = "json"
)

type TwitchVOD struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	URL         string `json:"url"`
}

type TwitchUser struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
}

type TwitchResponse[T interface{}] struct {
	Data []T `json:"data"`
}

func Init() *TwitchApi {
	twitchURL := os.Getenv("TWITCH_API")
	authToken := os.Getenv("TWITCH_OAUTH_TOKEN")
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	return &TwitchApi{
		client:    &http.Client{},
		EndPoint:  twitchURL,
		authToken: authToken,
		clienID:   clientID,
	}

}

func getFileType(filePath string) FILE_TYPE {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".txt":
		return TEXT_FILE
	case ".json":
		return JSON_FILE
	default:
		return TEXT_FILE
	}
}

func getTimePeriod(period TIME_PERIOD) uint {
	switch period {
	case "week":
		return 7
	case "month":
		return 28
	case "day":
		return 1
	default:
		return 100
	}
}

func (t *TwitchApi) GetVideoIDs(userID string, timePeriod TIME_PERIOD) ([]TwitchVOD, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/videos?user_id=%s&period=%s&first=%d", t.EndPoint, userID, timePeriod, getTimePeriod(timePeriod)), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.authToken))
	req.Header.Set("Client-Id", t.clienID)

	resp, err := t.client.Do(req)
	if err != nil {
		log.Printf("Couldnt create a client")
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("cant process resp.Body")
		return nil, err
	}
	var VODs TwitchResponse[TwitchVOD]
	err = json.Unmarshal(body, &VODs)

	if err != nil {
		log.Printf("cant unmarshall")
		return nil, err
	}

	return VODs.Data, nil
}

// we marshalling and unmarshalling one more time for no reason but i wanna follow some nice code practice
func (t *TwitchApi) WriteVideoID(userID string, timePeriod string, filePath string) error {
	ext := getFileType(filePath)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	if ext != JSON_FILE {
		vods, err := t.GetVideoIDs(userID, TIME_PERIOD(timePeriod))
		if err != nil {
			return fmt.Errorf("failed to get video IDs: %v", err)
		}
		data := []byte("Video IDs: " + fmt.Sprintf("%v", vods))
		_, err = file.Write(data)
		if err != nil {
			return fmt.Errorf("failed to write data to file: %v", err)
		}
		log.Printf("Successfully wrote VOD data to %s", filePath)
		return nil
	}

	vods, err := t.GetVideoIDs(userID, TIME_PERIOD(timePeriod))
	if err != nil {
		return fmt.Errorf("failed to get video IDs: %v", err)
	}

	data, err := json.MarshalIndent(vods, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	log.Printf("Successfully wrote VOD data to %s", filePath)
	return nil
}

func (t *TwitchApi) GetUserID(username string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users?login=%s", t.EndPoint, username), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.authToken))
	req.Header.Set("Client-Id", t.clienID)

	res, err := t.client.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var TwitchUser TwitchResponse[TwitchUser]

	err = json.Unmarshal(body, &TwitchUser)
	if err != nil {
		return "", err
	}

	return TwitchUser.Data[0].ID, nil
}

func (t *TwitchApi) GetUserDetails(userID string, userChannel chan TwitchUser, wg *sync.WaitGroup) {
	defer wg.Done()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users?id=%s", t.EndPoint, userID), nil)
	if err != nil {
		fmt.Printf("Error creating a GET req %s\n", err.Error())
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.authToken))
	req.Header.Set("Client-Id", t.clienID)
	resp, err := t.client.Do(req)
	if err != nil {
		fmt.Printf("Error executing the GET req %s\n", err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading the response body %s\n", err.Error())
		return
	}
	var userDetail TwitchResponse[TwitchUser]
	err = json.Unmarshal(body, &userDetail)
	if err != nil {
		fmt.Printf("Error unmarshalling the body %s\n", err.Error())
		return
	}

	userChannel <- userDetail.Data[0]
}
