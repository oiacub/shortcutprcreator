package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"shortcutcreator/src/shortcut"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

const (
	SHOW_CURRENT_BRANCH = "git branch --show-current"
	GH_TOKEN            = "GH_TOKEN"
	SHORTCUT_TOKEN      = "SHORTCUT_TOKEN"
)

type BranchInfo struct {
	Creator        string
	Story          string
	Name           string
	FullBranchName string
}

type BranchInterface interface {
	ParseBranchInfo()
}

func (branch *BranchInfo) ParseBranchInfo() {
	branchData := strings.Split(branch.FullBranchName, "/")
	branchData = cleanArrayNames(branchData)
	branch.Creator = branchData[0]
	branch.Story = branchData[1]
	branch.Name = branchData[2]
}

func main() {
	godotenv.Load()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the story number: ")
	numberStory, _ := reader.ReadString('\n')
	numberStory = cleanName(numberStory)
	shortcutClient := shortcut.BuildShortcutClient()
	story := shortcutClient.GetStory(numberStory)
	createPRPerBranch(story)
	if false {
		//octavioiacub/sc-67587/feature-add-program-assessment-link-to-home
		currentBranch := execCommand(SHOW_CURRENT_BRANCH)
		branchData := BranchInfo{FullBranchName: currentBranch}
		branchData.ParseBranchInfo()
		fmt.Println(branchData)
		prSubject := "Test PR"
		commitBranch := "main"
		prDescription := "Test PRName"
		prRepoOwner := "oiacub"
		prRepo := "shortcutprcreator"
		prFullBranchName := "octavioiacub/sc-67587/feature-add-program-assessment-link-to-home2"
		res := sendRequestForPRCreation(&prSubject, &prFullBranchName, &commitBranch, &prDescription, &prRepoOwner, &prRepo)
		fmt.Println(res)
	}
}

func createPRPerBranch(story shortcut.Story) {
	for _, b := range story.Branches {
		fmt.Println(b)
		fmt.Println(b.Name)
		fmt.Println("Merged Branchs", b.MergedBranchs)
		prSubject := story.GetStoryAbbreviationForGitType() + "(sc-" + story.Number + ") : " + story.Name
		relatedRepos := story.GetOtherRelatedRepos(b.Repository)
		if relatedRepos != "" {
			prSubject = prSubject + " [dep " + relatedRepos + "]"
		}
		commitBranch := "main"
		prDescription := story.Name
		prRepoOwner := "EpisourceLLC"
		prRepo := b.Repository
		prFullBranchName := b.Name
		fmt.Println("Branch data to create", prSubject, commitBranch, prDescription, prRepoOwner, prRepo, prFullBranchName)
		if len(b.MergedBranchs) == 0 {
			fmt.Println("Can create!")
		}
		//res := sendRequestForPRCreation(&prSubject, &prFullBranchName, &commitBranch, &prDescription, &prRepoOwner, &prRepo)
	}
}

func execCommand(command string) string {
	cmd := exec.Command("bash", "-c", SHOW_CURRENT_BRANCH)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return cleanName(string(out))
}

func cleanArrayNames(arr []string) []string {
	nd := make([]string, 0)
	for _, v := range arr {
		nd = append(nd, cleanName(v))
	}
	return nd
}

func cleanName(data string) string {
	re := regexp.MustCompile(`\r?\n`)
	data = re.ReplaceAllString(data, "")
	return data
}

func sendRequestForPRCreation(prSubject *string, commitBranch *string, prBranch *string, prDescription *string,
	prRepoOwner *string, prRepo *string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(GH_TOKEN)},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	newPR := &github.NewPullRequest{
		Title:               prSubject,
		Head:                commitBranch,
		Base:                prBranch,
		Body:                prDescription,
		MaintainerCanModify: github.Bool(true),
	}
	pr, _, err := client.PullRequests.Create(ctx, *prRepoOwner, *prRepo, newPR)
	if err != nil {
		return err
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	return nil
}
