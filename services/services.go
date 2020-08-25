package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/darolpz/golang-doc-generator/models"
	"github.com/darolpz/golang-doc-generator/utils"
	"github.com/jung-kurt/gofpdf"
	"github.com/mandolyte/mdtopdf"
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
		commitClasification := utils.ClasifyCommit(&commitTypes, &commit)
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

//GeneratePdf2 generates 2 files, a md file and a pdf file generated from the md file
func GeneratePdf2(params *models.Parameter, commits *[]models.Commit) (fileName string, err error) {
	title := fmt.Sprintf("%s_%s_%s", params.Repo[8:], params.From, params.To)
	baseURL := utils.GetEnvVariable("GITLAB_URL")
	repoURL := fmt.Sprintf(`%s/%s`, baseURL, params.Repo)

	/* Getting template */
	template := "template.md"
	content, err := ioutil.ReadFile(template)
	utils.CheckError(err)

	/* Transforming to string */
	stringifiedContent := string(content)

	/* Creating new md  file */
	f, err := os.Create(fmt.Sprintf("docs/%s.md", title))
	utils.CheckError(err)
	defer f.Close()

	/* Replacing strings */
	stringifiedContent = strings.ReplaceAll(stringifiedContent, "PROJECT", params.Repo[8:])
	stringifiedContent = strings.ReplaceAll(stringifiedContent, "TAG_FROM", params.From)
	stringifiedContent = strings.ReplaceAll(stringifiedContent, "TAG_TO", params.To)
	stringifiedContent = strings.ReplaceAll(stringifiedContent, "REPO_URL", repoURL)

	/* Writing first part of md file */
	_, err = f.WriteString(stringifiedContent)
	utils.CheckError(err)

	/* Start to fullfil with commits */

	commitTypes := [12]string{"Feat", "Fix", "Chore", "Test", "Docs", "Build", "Ci", "Perf", "Refactor", "Revert", "Style", "Others"}

	//¡Map of slices!
	sliceMap := make(map[string][]models.Commit)

	//Initializing map
	for _, commitType := range commitTypes {
		sliceMap[commitType] = make([]models.Commit, 0)
	}

	//Filling map
	for _, commit := range *commits {
		commitClasification := utils.ClasifyCommit(&commitTypes, &commit)
		sliceMap[commitClasification] = append(sliceMap[commitClasification], commit)
	}

	//Printing commits according to type
	for _, commitType := range commitTypes {
		if len(sliceMap[commitType]) > 0 {
			f.WriteString(fmt.Sprintf("### %s\n\n", commitType))

			for index, commit := range sliceMap[commitType] {
				// Ignore commit with repeated title
				if index == 0 || sliceMap[commitType][index].Title != sliceMap[commitType][index-1].Title {
					f.WriteString(fmt.Sprintf("* **%s** (%s): %s\n", commit.ShortID, commit.CreatedAt.Format("2006-01-02"), commit.Title))

				}
			}
			f.WriteString("\n")
		}
	}

	f.Sync()

	/* Generating pdf */
	input := fmt.Sprintf("docs/%s.md", title)
	output := fmt.Sprintf("docs/%s.pdf", title)

	mdContent, err := ioutil.ReadFile(input)
	utils.CheckError(err)

	pf := mdtopdf.NewPdfRenderer("", "", output, "")
	pf.Pdf.SetSubject("title", true)
	pf.THeader = mdtopdf.Styler{Font: "Times", Style: "IUB", Size: 20, Spacing: 2,
		TextColor: mdtopdf.Color{Red: 0, Green: 0, Blue: 0},
		FillColor: mdtopdf.Color{Red: 179, Green: 179, Blue: 255}}
	pf.TBody = mdtopdf.Styler{Font: "Arial", Style: "", Size: 12, Spacing: 2,
		TextColor: mdtopdf.Color{Red: 0, Green: 0, Blue: 0},
		FillColor: mdtopdf.Color{Red: 255, Green: 102, Blue: 129}}

	fileName = output
	err = pf.Process(mdContent)
	return
}

//NotifyThroughtSlack notifies slack
func NotifyThroughtSlack(fileName, channel string) {
	reqURL := utils.GetEnvVariable("SLACK_URL")
	token := utils.GetEnvVariable("SLACK_APP_TOKEN")
	token = fmt.Sprintf("Bearer %s", token)
	channel = utils.GetEnvVariable(strings.ToUpper(channel))
	fileURL := utils.GetEnvVariable("HOST_URL")
	// Getting file
	fileDir, _ := os.Getwd()
	filePath := path.Join(fileDir, fileName)
	file, err := os.Open(filePath)
	utils.CheckError(err)
	fileContents, err := ioutil.ReadAll(file)
	utils.CheckError(err)

	//Setting parameters
	params := map[string]string{
		"channels":        channel,
		"initial_comment": fmt.Sprintf("New file upload at: %s/%s", fileURL, fileName),
	}
	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	part, err := writer.CreateFormFile("file", fileName)
	utils.CheckError(err)
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	utils.CheckError(err)

	//Setting http client and sending request
	request, err := http.NewRequest("POST", reqURL, reqBody)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", token)
	utils.CheckError(err)
	client := &http.Client{}
	resp, err := client.Do(request)
	utils.CheckError(err)

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	utils.CheckError(err)
}
