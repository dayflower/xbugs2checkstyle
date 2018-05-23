package main

// BugCollection is <BugCollection> item in FindBugs report XML
type BugCollection struct {
	Project      Project       `xml:"Project"`
	BugInstances []BugInstance `xml:"BugInstance"`
}

// Project is <Project> item in FindBugs report XML
type Project struct {
	SrcDirs []string `xml:"SrcDir"`
}

// BugInstance is <BugInstance> item in FindBugs report XML
type BugInstance struct {
	Type              string       `xml:"type,attr"`
	ClassSourceLines  []SourceLine `xml:"Class>SourceLine"`
	MethodSourceLines []SourceLine `xml:"Method>SourceLine"`
	FieldSourceLines  []SourceLine `xml:"Field>SourceLine"`
	SourceLines       []SourceLine `xml:"SourceLine"`
}

// SourceLine is <SourceLine> item in FindBugs report XML
type SourceLine struct {
	ClassName     string `xml:"classname,attr"`
	Start         int    `xml:"start,attr"`
	End           int    `xml:"end,attr"`
	StartBytecode int    `xml:"startBytecode,attr"`
	EndBytecode   int    `xml:"endBytecode,attr"`
	SourceFile    string `xml:"sourcefile,attr"`
	SourcePath    string `xml:"sourcepath,attr"`
}
