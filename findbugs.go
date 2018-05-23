package main

type BugCollection struct {
	BugInstances []BugInstance `xml:"BugInstance"`
}

type BugInstance struct {
	Type              string       `xml:"type,attr"`
	ClassSourceLines  []SourceLine `xml:"Class>SourceLine"`
	MethodSourceLines []SourceLine `xml:"Method>SourceLine"`
	FieldSourceLines  []SourceLine `xml:"Field>SourceLine"`
	SourceLines       []SourceLine `xml:"SourceLine"`
}

type SourceLine struct {
	ClassName     string `xml:"classname,attr"`
	Start         int    `xml:"start,attr"`
	End           int    `xml:"end,attr"`
	StartBytecode int    `xml:"startBytecode,attr"`
	EndBytecode   int    `xml:"endBytecode,attr"`
	SourceFile    string `xml:"sourcefile,attr"`
	SourcePath    string `xml:"sourcepath,attr"`
}
