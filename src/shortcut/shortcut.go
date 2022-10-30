package shortcut

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	SHORTCUT_TOKEN = "SHORTCUT_TOKEN"
)

type ShortcutAPI struct {
	Token           string
	StoriesEndpoint string
}

type Story struct {
	Name      string `json:"name"`
	Number    string
	StoryType string   `json:"story_type"`
	Branches  []Branch `json:"branches"`
}

type Branch struct {
	Id            int    `json:"branch_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	RepositoryId  int    `json:"repository_id"`
	Repository    string
	Url           string `json:"url"`
	MergedBranchs []int  `json:"merged_branch_ids"`
}

type ShortcutInterface interface {
	GetStory(number string) Story
}

type BranchInterface interface {
	CompleteInfo() Branch
}

type StoryInterface interface {
	GetStoryAbbreviationForGitType() string
	GetOtherRelatedRepos(currentRepo string) string
}

func (story Story) GetOtherRelatedRepos(currentRepo string) string {
	otherRepos := make([]string, 0)
	for _, b := range story.Branches {
		if b.Repository != currentRepo {
			otherRepos = append(otherRepos, b.Repository)
		}
	}
	return strings.Join(otherRepos, ",")
}

func (story Story) GetStoryAbbreviationForGitType() string {
	switch strings.ToLower(story.StoryType) {
	case "feature":
		return "feat"
	case "bug":
		return "fix"
	}
	return ""
}

func (branch Branch) CompleteInfo() Branch {
	urlParts := strings.Split(branch.Url, "/")
	branch.Repository = urlParts[4]
	return branch
}

func BuildShortcutClient() ShortcutInterface {
	return ShortcutAPI{Token: os.Getenv(SHORTCUT_TOKEN), StoriesEndpoint: "https://api.app.shortcut.com/api/v3/stories/"}
}

func (client ShortcutAPI) GetStory(number string) Story {
	fmt.Println("Token", client.Token)
	httpClient := http.Client{}
	url := client.StoriesEndpoint + number
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Shortcut-Token", client.Token)
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	resBody, _ := ioutil.ReadAll(res.Body)
	story := Story{}
	json.Unmarshal([]byte(resBody), &story)
	story.Number = number
	for i, b := range story.Branches {
		story.Branches[i] = b.CompleteInfo()
	}
	return story
}
