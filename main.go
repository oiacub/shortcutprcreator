package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

const (
	SHOW_CURRENT_BRANCH = "git branch --show-current"
)

type BranchInfo struct {
	Creator string
	Story   string
	Name    string
}

func main() {
	//octavioiacub/sc-67587/feature-add-program-assessment-link-to-home
	currentBranch := execCommand(SHOW_CURRENT_BRANCH)
	branchData := branchParser(currentBranch)
	fmt.Println(branchData)
}

func execCommand(command string) string {
	cmd := exec.Command("bash", "-c", SHOW_CURRENT_BRANCH)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func branchParser(info string) BranchInfo {
	branchData := strings.Split(info, "/")
	branchData = cleanArrayNames(branchData)
	return BranchInfo{Creator: branchData[0], Story: branchData[1], Name: branchData[2]}
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
