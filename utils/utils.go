package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/fatih/color"
)

var (
	mutex         sync.Mutex
	totalComments VodComments
)

type Pair struct {
	Key    string
	Value  uint
	UserID string
}

type CommentCount struct {
	UserID string
	Count  uint
}
type VodComments struct {
	Streamer struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"streamer"`

	Comments []struct {
		ID        string `json:"_id"`
		Commenter struct {
			DisplayName string `json:"display_name"`
			ID          string `json:"_id"`
			Name        string `json:"name"`
		} `json:"commenter"`
		Message struct {
			Body      string `json:"body"`
			Fragments []struct {
				Text     string `json:"text"`
				Emoticon any    `json:"emoticon"`
			} `json:"fragments"`
			UserBadges []UserBadge `json:"user_badges"`
			UserColor  any         `json:"user_color"`
			Emoticons  []any       `json:"emoticons"`
		} `json:"message"`
	} `json:"comments"`
}

type UserBadge struct {
	ID      string `json:"_id"`
	Version string `json:"version"`
}

func getJSONFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var jsonFiles []string
	for _, file := range files {
		if ext := filepath.Ext(file.Name()); ext == ".json" {
			jsonFiles = append(jsonFiles, filepath.Join(dir, file.Name()))
		}
	}
	return jsonFiles, nil
}

func Test(temp *VodComments) {

	fmt.Println("Size of comments: ", len(temp.Comments))
	modCount := 0
	vipCount := 0
	subCount := 0
	for _, comment := range temp.Comments {
		if len(comment.Message.UserBadges) == 0 {
			continue
		}
		if comment.Message.UserBadges[0].ID == "vip" {
			vipCount++
		}
		if comment.Message.UserBadges[0].ID == "moderator" {
			modCount++
		}
		if comment.Message.UserBadges[0].ID == "subscriber" {
			subCount++
		}

	}
	fmt.Printf("Total\nvips: %d\tmods: %d\t subs: %d\n", vipCount, modCount, subCount)
}

func Temp() {
	jsonFiles, err := getJSONFiles("./output")
	if err != nil {
		color.Red("there are no json files")
	}
	var wg sync.WaitGroup
	for _, jsonFile := range jsonFiles {
		wg.Add(1)
		readJSONComments(jsonFile, &wg)
	}

	wg.Wait()
	Test(&totalComments)

	frequencyMap := make(map[string]CommentCount)

	for _, comment := range totalComments.Comments {
		_, exists := frequencyMap[comment.Commenter.DisplayName]

		if exists {
			cmtCount := frequencyMap[comment.Commenter.DisplayName]
			cmtCount.Count++
			frequencyMap[comment.Commenter.DisplayName] = cmtCount
			continue
		}

		frequencyMap[comment.Commenter.DisplayName] = CommentCount{UserID: comment.Commenter.ID, Count: 1}
	}

	sortedCmtByCount := rankByWordCount(frequencyMap)

	fmt.Println(sortedCmtByCount[:50])
}

func readJSONComments(path string, wg *sync.WaitGroup) {
	defer wg.Done()
	body, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading json file %s", err.Error())
		return
	}

	var currentComment VodComments
	err = json.Unmarshal(body, &currentComment)
	if err != nil {
		fmt.Printf("Error unmarshalling vod comment %s", err.Error())
		return
	}
	mutex.Lock()
	if totalComments.Streamer.Name == "" {
		totalComments = currentComment
	}
	totalComments.Comments = append(totalComments.Comments, currentComment.Comments...)
	mutex.Unlock()

}

// return an array of hashmap sorted by descending order in a slice
func rankByWordCount(frequencyMap map[string]CommentCount) []Pair {
	sortedSlice := make([]Pair, len(frequencyMap))
	i := 0
	for key, val := range frequencyMap {
		sortedSlice[i] = Pair{Key: key, Value: val.Count, UserID: val.UserID}
		i++
	}

	sort.Slice(sortedSlice, func(i, j int) bool {
		return sortedSlice[i].Value > sortedSlice[j].Value
	})

	return sortedSlice
}
