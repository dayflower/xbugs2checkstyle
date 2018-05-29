package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"

	"github.com/phayes/checkstyle"
)

var opts struct {
	Language string `short:"l" long:"lang" description:"language (default: en)"`
}

func toAbsPath(srcDirs []string, source string) string {
	minseps := math.MaxInt32
	var res string

	for _, dir := range srcDirs {
		if strings.HasSuffix(dir, string(os.PathSeparator)+source) {
			seps := strings.Count(dir[0:len(dir)-len(source)-1], string(os.PathSeparator))
			if seps < minseps {
				minseps = seps
				res = dir
			}
		}
	}

	if res != "" {
		return res
	} else {
		return source
	}
}

func addCheckstyleError(bugDescriptions map[string]string, srcDirs []string, files map[string]*checkstyle.File, bug BugInstance, source SourceLine) {
	sourcePath := toAbsPath(srcDirs, source.SourcePath)

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

	srcDirs := document.Project.SrcDirs
	checkstyleFiles := make(map[string]*checkstyle.File)

	for _, bug := range document.BugInstances {
		for _, source := range bug.ClassSourceLines {
			addCheckstyleError(bugDescriptions, srcDirs, checkstyleFiles, bug, source)
		}
		for _, source := range bug.MethodSourceLines {
			addCheckstyleError(bugDescriptions, srcDirs, checkstyleFiles, bug, source)
		}
		for _, source := range bug.FieldSourceLines {
			addCheckstyleError(bugDescriptions, srcDirs, checkstyleFiles, bug, source)
		}
		for _, source := range bug.SourceLines {
			addCheckstyleError(bugDescriptions, srcDirs, checkstyleFiles, bug, source)
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
