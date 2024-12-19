package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type GeneratorInput struct {
	PackageName  string
	RepoName     string
	IncludeDBW   bool
	IncludeDBR   bool
	IncludeCache bool
	Monitoring   bool
	OutputPath   string
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the package name: ")
	packageName, _ := reader.ReadString('\n')
	packageName = strings.TrimSpace(packageName)

	repoName := strings.Title(packageName)

	includeDBW := askYesNo(reader, "Include DB writer (DBW)? (y/n): ")
	includeDBR := askYesNo(reader, "Include DB reader (DBR)? (y/n): ")
	includeCache := askYesNo(reader, "Include Redis? (y/n): ")
	includeMonitoring := askYesNo(reader, "Include Monitoring? (y/n): ")

	fmt.Print("Enter the output path: ")
	outputPath, _ := reader.ReadString('\n')
	outputPath = strings.TrimSpace(outputPath)

	input := GeneratorInput{
		PackageName:  packageName,
		RepoName:     repoName,
		IncludeDBW:   includeDBW,
		IncludeDBR:   includeDBR,
		IncludeCache: includeCache,
		Monitoring:   includeMonitoring,
		OutputPath:   outputPath,
	}

	err := generateCode(input)
	if err != nil {
		fmt.Printf("Error generating code: %v\n", err)
	}
}

func askYesNo(reader *bufio.Reader, question string) bool {
	fmt.Print(question)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}

func generateCode(input GeneratorInput) error {
	const repoTemplate = `package {{.PackageName}}

import (
	{{- if .Monitoring }}
	"app/pkg/monitoring"
	{{- end }}

	{{- if .IncludeCache }}
	"github.com/go-redis/redis/v8"
	{{- end }}
	{{- if .IncludeDBW }}
	"github.com/jmoiron/sqlx"
	{{- end }}
)

type {{.RepoName}} struct {
	{{- if .IncludeDBW }}
	dbW   *sqlx.DB
	{{- end }}
	{{- if .IncludeDBR }}
	dbR   *sqlx.DB
	{{- end }}
	{{- if .IncludeCache }}
	redisClient *redis.Client
	{{- end }}
	{{- if .Monitoring }}
	monitoring.Helper
	{{- end }}
}

type RepoArgs struct {
	{{- if .IncludeDBW }}
	DBW   *sqlx.DB
	{{- end }}
	{{- if .IncludeDBR }}
	DBR   *sqlx.DB
	{{- end }}
	{{- if .IncludeCache }}
	RedisClient *redis.Client
	{{- end }}
}

func New{{.RepoName}}Repo(in RepoArgs) *{{.RepoName}} {
	return &{{.RepoName}}{
		{{- if .IncludeDBW }}
		dbW:   in.DBW,
		{{- end }}
		{{- if .IncludeDBR }}
		dbR:   in.DBR,
		{{- end }}
		{{- if .IncludeCache }}
		redisClient: in.RedisClient,
		{{- end }}
	}
}
`

	tmpl, err := template.New("repo").Parse(repoTemplate)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(input.OutputPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	outputFile, err := os.Create(filepath.Join(input.OutputPath, "init.go"))
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return tmpl.Execute(outputFile, input)
}
