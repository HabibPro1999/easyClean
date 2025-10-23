package models

// ReferenceType represents the type of code reference
type ReferenceType int

const (
	RefTypeImport ReferenceType = iota
	RefTypeStringLiteral
	RefTypeTemplateLiteral
	RefTypeCSSUrl
	RefTypeHTMLAttribute
	RefTypeConstant
	RefTypeFunctionCall
	RefTypeConfig // For YAML/config file declarations
)

// String returns the string representation of ReferenceType
func (rt ReferenceType) String() string {
	return [...]string{
		"Import",
		"StringLiteral",
		"TemplateLiteral",
		"CSSUrl",
		"HTMLAttribute",
		"Constant",
		"FunctionCall",
		"Config",
	}[rt]
}

// Reference represents a single location in code where an asset is referenced
type Reference struct {
	// Location
	SourceFile string `json:"source_file"`
	LineNumber int    `json:"line_number"`
	Column     int    `json:"column,omitempty"`

	// Content
	MatchedText string `json:"matched_text"`
	Context     string `json:"context,omitempty"`

	// Classification
	Type       ReferenceType `json:"type"`
	Confidence float32       `json:"confidence"`

	// Flags
	IsComment  bool `json:"is_comment"`
	IsDynamic  bool `json:"is_dynamic"`
	IsDeadCode bool `json:"is_dead_code,omitempty"`
}
