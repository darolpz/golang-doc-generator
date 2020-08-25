package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/darolpz/golang-doc-generator/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//CheckError throw panic if error is distinct of nil
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

// LogFormat is a custom logging format
func LogFormat(param gin.LogFormatterParams) string {

	// your custom format
	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

// GetEnvVariable returns value stored in .env
func GetEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	CheckError(err)

	return os.Getenv(key)
}

// FilterCommits find all commits of a type and returns them
func FilterCommits(commitType *string, commits *[]models.Commit) []models.Commit {
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

// ClasifyCommit returns the type of a commit
func ClasifyCommit(CommitTypes *[12]string, commit *models.Commit) string {
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
