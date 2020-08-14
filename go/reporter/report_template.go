package reporter

import (
	"strings"
	"text/template"

	"github.com/tmcphillips/blazegraph-util/util"
)

type Properties struct {
	Delimiters DelimiterPair
	Funcs      template.FuncMap
	Prefixes   map[string]string
}

var GraveDelimiters DelimiterPair
var JSPDelimiters DelimiterPair
var TripleSingleQuoteDelimiters DelimiterPair
var DefaultDelimiters DelimiterPair

func init() {
	GraveDelimiters = DelimiterPair{"`", "`"}
	JSPDelimiters = DelimiterPair{Start: "<%", End: "%>"}
	TripleSingleQuoteDelimiters = DelimiterPair{Start: "'''", End: "'''"}
	DefaultDelimiters = TripleSingleQuoteDelimiters
}

// ReportTemplate is customizable to create reports from different kinds of templates.
type ReportTemplate struct {
	Text       string
	Properties Properties
}

// NewReportTemplate returns a ReportTemplate with the given customizations.
func NewReportTemplate(text string) *ReportTemplate {
	// print(text)
	rt := new(ReportTemplate)
	rt.Text = text
	rt.Properties.Delimiters = DefaultDelimiters
	rt.Properties.Prefixes = map[string]string{}
	return rt
}

func (rp *ReportTemplate) SetDelimiters(delimiters DelimiterPair) {
	rp.Properties.Delimiters = delimiters
}

func (rp *ReportTemplate) SetFuncs(funcs template.FuncMap) {
	rp.Properties.Funcs = funcs
}

func (rp *ReportTemplate) Expand(data interface{}, removeNewlines bool) (result string, err error) {

	text := util.TrimEachLine(rp.Text)
	text = EscapeRawText(rp.Properties.Delimiters, text)

	if removeNewlines {
		text = RemoveNewlines(text)
	}

	textTemplate := template.New("main")
	if rp.Properties.Funcs != nil {
		textTemplate.Funcs(rp.Properties.Funcs)
	}

	textTemplate.Parse(text)
	var buffer strings.Builder

	err = textTemplate.Execute(&buffer, data)
	if err != nil {
		return
	}
	result = buffer.String()
	result = InsertNewlines(result)
	return
}
