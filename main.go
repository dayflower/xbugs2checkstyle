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
	BaseDir string `short:"b" long:"basedir" description:"base directory"`
}

func addCheckstyleError(files map[string]*checkstyle.File, bug BugInstance, source SourceLine) {
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

	file.AddError(checkstyle.NewError(source.Start, 1, checkstyle.SeverityWarning, bug.Type, bug.Type))
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	if opts.BaseDir == "" {
		opts.BaseDir, _ = os.Getwd()
	}

	body, _ := ioutil.ReadAll(os.Stdin)

	document := BugCollection{}
	xml.Unmarshal(body, &document)

	checkstyleFiles := make(map[string]*checkstyle.File)

	for _, bug := range document.BugInstances {
		for _, source := range bug.ClassSourceLines {
			addCheckstyleError(checkstyleFiles, bug, source)
		}
		for _, source := range bug.MethodSourceLines {
			addCheckstyleError(checkstyleFiles, bug, source)
		}
		for _, source := range bug.FieldSourceLines {
			addCheckstyleError(checkstyleFiles, bug, source)
		}
		for _, source := range bug.SourceLines {
			addCheckstyleError(checkstyleFiles, bug, source)
		}
	}

	checkstyleResult := checkstyle.New()
	for _, file := range checkstyleFiles {
		checkstyleResult.AddFile(file)
	}

	buf, _ := xml.MarshalIndent(checkstyleResult, "", "\t")
	fmt.Println(string(buf))
}
