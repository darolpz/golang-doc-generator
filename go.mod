module github.com/darolpz/golang-doc-generator

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/joho/godotenv v1.3.0
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/mandolyte/mdtopdf v0.0.0-20190923134258-4400ba48c487
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1 // indirect
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
