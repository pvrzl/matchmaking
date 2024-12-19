package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type GeneratorInput struct {
	PackageName string
	UseCaseName string
	Monitoring  bool
	OutputPath  string
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the package name: ")
	packageName, _ := reader.ReadString('\n')
	packageName = strings.TrimSpace(packageName)

	useCaseName := cases.Title(language.English).String(packageName)

	includeMonitoring := askYesNo(reader, "Include Monitoring? (y/n): ")

	fmt.Print("Enter the output path: ")
	outputPath, _ := reader.ReadString('\n')
	outputPath = strings.TrimSpace(outputPath)

	input := GeneratorInput{
		PackageName: packageName,
		UseCaseName: useCaseName,
		Monitoring:  includeMonitoring,
		OutputPath:  outputPath,
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

)

type {{.UseCaseName}} struct {
	{{- if .Monitoring }}
	monitoring.Helper
	{{- end }}
}

type UseCaseArgs struct {
	
}

func New{{.UseCaseName}}UseCase(in UseCaseArgs) *{{.UseCaseName}} {
	return &{{.UseCaseName}}{
		
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
