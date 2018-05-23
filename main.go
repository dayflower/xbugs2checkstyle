package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/phayes/checkstyle"
)

type absolutePathCache struct {
	paths []string
	cache map[string]string
}

func newAbsolutePathCache(paths []string) *absolutePathCache {
	sorted := make([]string, len(paths))
	copy(sorted, paths)

	sort.Slice(sorted, func(i, j int) bool {
		a, b := paths[i], paths[j]
		if len(a) > len(b) {
			return true
		}
		if len(a) < len(b) {
			return false
		}
		return strings.Compare(a, b) > 0
	})

	return &absolutePathCache{sorted, make(map[string]string)}
}

func (a *absolutePathCache) find(relative string) string {
	res, exist := a.cache[relative]
	if exist {
		return res
	}

	for _, v := range a.paths {
		if strings.HasSuffix(v, "/"+relative) {
			a.cache[relative] = v
			return v
		}
	}

	a.cache[relative] = relative
	return relative
}

func addCheckstyleError(pathCache *absolutePathCache, files map[string]*checkstyle.File, bug BugInstance, source SourceLine) {
	sourcePath := pathCache.find(source.SourcePath)

	file, exist := files[sourcePath]
	if !exist {
		file = checkstyle.NewFile(sourcePath)
		files[sourcePath] = file
	}

	file.AddError(checkstyle.NewError(source.Start, 1, checkstyle.SeverityWarning, bug.Type, bug.Type))
}

func main() {
	body, _ := ioutil.ReadAll(os.Stdin)

	document := BugCollection{}
	xml.Unmarshal(body, &document)

	pathCache := newAbsolutePathCache(document.Project.SrcDirs)

	checkstyleFiles := make(map[string]*checkstyle.File)

	for _, bug := range document.BugInstances {
		for _, source := range bug.ClassSourceLines {
			addCheckstyleError(pathCache, checkstyleFiles, bug, source)
		}
		for _, source := range bug.MethodSourceLines {
			addCheckstyleError(pathCache, checkstyleFiles, bug, source)
		}
		for _, source := range bug.FieldSourceLines {
			addCheckstyleError(pathCache, checkstyleFiles, bug, source)
		}
		for _, source := range bug.SourceLines {
			addCheckstyleError(pathCache, checkstyleFiles, bug, source)
		}
	}

	checkstyleResult := checkstyle.New()
	for _, file := range checkstyleFiles {
		checkstyleResult.AddFile(file)
	}

	buf, _ := xml.MarshalIndent(checkstyleResult, "", "\t")
	fmt.Println(string(buf))
}
