package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/darolpz/golang-doc-generator/models"
	"github.com/darolpz/golang-doc-generator/utils"
	"github.com/jung-kurt/gofpdf"
)

// GetCommits return commit titles
func GetCommits(params *models.Parameter) []models.Commit {
	baseURL := utils.GetEnvVariable("GITLAB_URL")

	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/compare?from=%s&to=%s", baseURL, url.QueryEscape(params.Repo), params.From, params.To)
	resp, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	bodyStringified := string(body)
	var gitResponse models.GitResponse
	errUnmarshal := json.Unmarshal([]byte(bodyStringified), &gitResponse)
	if errUnmarshal != nil {
		fmt.Println(errUnmarshal)
	}
	return gitResponse.Commits
}

// GeneratePdf generates a pdf for release notes
func GeneratePdf(params *models.Parameter, commits *[]models.Commit) error {
	baseURL := utils.GetEnvVariable("GITLAB_URL")
	//Setting up document
	title := fmt.Sprintf("%s_%s_%s", params.Repo[8:], params.From, params.To)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle(title, false)
	pdf.SetTopMargin(30)
	pdf.SetHeaderFuncMode(func() {
		pdf.Image("img/logo.png", 100, 6, 10, 10, false, "", 0, "")
		pdf.SetY(10)
		pdf.SetFont("Arial", "B", 15)
		pdf.SetTextColor(100, 100, 100)
		pdf.WriteAligned(50, 0, "Backend Team", "L")
		pdf.WriteAligned(0, 0, "Flow Factory", "R")
	}, true)

	pdf.AddPage()
	pdf.Ln(0)
	pdf.SetFont("Arial", "B", 30)
	pdf.CellFormat(0, 6, title, "", 1, "C", false, 0, "")
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 25)
	pdf.Bookmark("Repo", 0, 0)
	pdf.CellFormat(0, 6, "Repo", "", 1, "L", false, 0, "")
	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 15)
	_, lineHt := pdf.GetFontSize()
	htmlStr := fmt.Sprintf(`<a href="%s/%s">%s/%s</a>`, baseURL, params.Repo, baseURL, params.Repo)
	html := pdf.HTMLBasicNew()
	html.Write(lineHt, htmlStr)
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 25)
	pdf.Bookmark("Changes", 0, 0)
	pdf.CellFormat(0, 6, "Changes", "", 1, "L", false, 0, "")
	pdf.Ln(10)

	commitTypes := [12]string{"Feat", "Fix", "Chore", "Test", "Docs", "Build", "Ci", "Perf", "Refactor", "Revert", "Style", "Others"}

	//¡Map of arrays!
	bigMap := make(map[string][]models.Commit)

	//Initializing map
	for _, commitType := range commitTypes {
		bigMap[commitType] = make([]models.Commit, 0)
	}

	//Filling map
	for _, commit := range *commits {
		commitClasification := clasifyCommit(&commitTypes, &commit)
		bigMap[commitClasification] = append(bigMap[commitClasification], commit)
	}

	//Printing commits according to type
	for _, commitType := range commitTypes {
		if len(bigMap[commitType]) > 0 {
			pdf.SetFont("Arial", "B", 20)
			pdf.Bookmark(commitType, 1, 0)
			pdf.CellFormat(0, 6, commitType, "", 1, "L", false, 0, "")
			pdf.Ln(6)

			for index, commit := range bigMap[commitType] {
				// Ignore commit with repeated title
				if index == 0 || bigMap[commitType][index].Title != bigMap[commitType][index-1].Title {
					pdf.SetFont("Arial", "B", 15)
					pdf.CellFormat(0, 6, commit.ShortID, "", 1, "L", false, 0, "")
					pdf.SetFont("Arial", "I", 10)
					pdf.CellFormat(0, 6, commit.Title, "", 1, "", false, 0, "")
					pdf.Ln(2)
				}
			}
			pdf.Ln(4)
		}
	}

	return pdf.OutputFileAndClose(fmt.Sprintf("docs/%s.pdf", title))
}

// filterCommits look for all commits from one type and returns them
func filterCommits(commitType *string, commits *[]models.Commit) []models.Commit {
	slice := make([]models.Commit, 0)
	for _, commit := range *commits {
		//Compare commmit type with first characters of commit title
		if strings.ToLower(*commitType) == commit.Title[0:len(*commitType)] {
			colonIndex := strings.Index(commit.Title, ":")
			commit.Title = strings.TrimSpace(commit.Title[(colonIndex + 1):])
			slice = append(slice, commit)
		}
	}

	return slice

}

// clasifyCommit returns the type of a commit
func clasifyCommit(CommitTypes *[12]string, commit *models.Commit) string {
	for _, commitType := range CommitTypes {
		//Compare commmit type with first characters of commit title
		if strings.ToLower(commitType) == commit.Title[0:len(commitType)] {
			colonIndex := strings.Index(commit.Title, ":")
			commit.Title = strings.TrimSpace(commit.Title[(colonIndex + 1):])
			return commitType
		}
	}
	return "Others"
}
