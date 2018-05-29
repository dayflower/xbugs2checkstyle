package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"

	"github.com/phayes/checkstyle"
)

var opts struct {
	BaseDir  string `short:"b" long:"basedir" description:"base directory"`
	Language string `short:"l" long:"lang" description:"language (default: en)"`
}

func addCheckstyleError(bugDescriptions map[string]string, files map[string]*checkstyle.File, bug BugInstance, source SourceLine) {
	var sourcePath string
	if filepath.IsAbs(source.SourcePath) {
		sourcePath = source.SourcePath
	} else {
		sourcePath = filepath.Join(opts.BaseDir, source.SourcePath)
	}

	file, exist := files[sourcePath]
	if !exist {
		file = checkstyle.NewFile(sourcePath)
		files[sourcePath] = file
	}

	message := bug.Type + ": " + bugDescriptions[bug.Type]

	file.AddError(checkstyle.NewError(source.Start, 1, checkstyle.SeverityWarning, message, bug.Type))
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	if opts.BaseDir == "" {
		opts.BaseDir, _ = os.Getwd()
	}

	if opts.Language == "" {
		opts.Language = "en"
	}

	var bugDescriptions map[string]string
	switch opts.Language {
	case "en":
		bugDescriptions = BugDescriptionEn
	case "ja":
		bugDescriptions = BugDescriptionJa
	case "fr":
		bugDescriptions = BugDescriptionFr
	default:
		panic("Unsupported language: " + opts.Language)
	}

	body, _ := ioutil.ReadAll(os.Stdin)

	document := BugCollection{}
	xml.Unmarshal(body, &document)

	checkstyleFiles := make(map[string]*checkstyle.File)

	for _, bug := range document.BugInstances {
		for _, source := range bug.ClassSourceLines {
			addCheckstyleError(bugDescriptions, checkstyleFiles, bug, source)
		}
		for _, source := range bug.MethodSourceLines {
			addCheckstyleError(bugDescriptions, checkstyleFiles, bug, source)
		}
		for _, source := range bug.FieldSourceLines {
			addCheckstyleError(bugDescriptions, checkstyleFiles, bug, source)
		}
		for _, source := range bug.SourceLines {
			addCheckstyleError(bugDescriptions, checkstyleFiles, bug, source)
		}
	}

	checkstyleResult := checkstyle.New()
	checkstyleResult.Version = "8.9"
	for _, file := range checkstyleFiles {
		checkstyleResult.AddFile(file)
	}

	buf, _ := xml.MarshalIndent(checkstyleResult, "", "\t")
	fmt.Println(string(buf))
}
