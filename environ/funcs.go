package environ

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
	"text/template"

	"github.com/codemodus/kace"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

// FuncMap is the standard list of functions available to templates.
// Strings methods are from the strings package - https://golang.org/pkg/strings/
// kace methods are from https://github.com/codemodus/kace/
var FuncMap = map[string]interface{}{
	"compare":      strings.Compare,      // https://golang.org/pkg/strings/#Compare
	"contains":     strings.Contains,     // https://golang.org/pkg/strings/#Contains
	"containsAny":  strings.ContainsAny,  // https://golang.org/pkg/strings/#ContainsAny
	"count":        strings.Count,        // https://golang.org/pkg/strings/#Count
	"equalFold":    strings.EqualFold,    // https://golang.org/pkg/strings/#EqualFold
	"fields":       strings.Fields,       // https://golang.org/pkg/strings/#Fields
	"hasPrefix":    strings.HasPrefix,    // https://golang.org/pkg/strings/#HasPrefix
	"hasSuffix":    strings.HasPrefix,    // https://golang.org/pkg/strings/#HasPrefix
	"index":        strings.Index,        // https://golang.org/pkg/strings/#Index
	"indexAny":     strings.IndexAny,     // https://golang.org/pkg/strings/#IndexAny
	"join":         strings.Join,         // https://golang.org/pkg/strings/#Join
	"lastIndex":    strings.LastIndex,    // https://golang.org/pkg/strings/#LastIndex
	"lastIndexAny": strings.LastIndexAny, // https://golang.org/pkg/strings/#LastIndexAny
	"repeat":       strings.Repeat,       // https://golang.org/pkg/strings/#Repeat
	"replace":      strings.Replace,      // https://golang.org/pkg/strings/#Replace
	"split":        strings.Split,        // https://golang.org/pkg/strings/#Split
	"splitAfter":   strings.SplitAfter,   // https://golang.org/pkg/strings/#SplitAfter
	"splitAfterN":  strings.SplitAfterN,  // https://golang.org/pkg/strings/#SplitAfterN
	"splitN":       strings.SplitN,       // https://golang.org/pkg/strings/#SplitN
	"title":        strings.Title,        // https://golang.org/pkg/strings/#Title
	"toLower":      strings.ToLower,      // https://golang.org/pkg/strings/#ToLower
	"toTitle":      strings.ToTitle,      // https://golang.org/pkg/strings/#ToTitle
	"toUpper":      strings.ToUpper,      // https://golang.org/pkg/strings/#ToUpper
	"trim":         strings.Trim,         // https://golang.org/pkg/strings/#Trim
	"trimLeft":     strings.TrimLeft,     // https://golang.org/pkg/strings/#TrimLeft
	"trimPrefix":   strings.TrimPrefix,   // https://golang.org/pkg/strings/#TrimPrefix
	"trimRight":    strings.TrimRight,    // https://golang.org/pkg/strings/#TrimRight
	"trimSpace":    strings.TrimSpace,    // https://golang.org/pkg/strings/#TrimSpace
	"trimSuffix":   strings.TrimSuffix,   // https://golang.org/pkg/strings/#TrimSuffix

	"camel":      kace.Camel,      // https://godoc.org/github.com/codemodus/kace#Camel
	"kebab":      kace.Kebab,      // https://godoc.org/github.com/codemodus/kace#Kebab
	"kebabUpper": kace.KebabUpper, // https://godoc.org/github.com/codemodus/kace#KebabUpper
	"pascal":     kace.Pascal,     // https://godoc.org/github.com/codemodus/kace#Pascal
	"snake":      kace.Snake,      // https://godoc.org/github.com/codemodus/kace#Snake
	"snakeUpper": kace.SnakeUpper, // https://godoc.org/github.com/codemodus/kace#SnakeUpper

	"sliceString": sliceString,
	"makeTable":   makeTable,
	"makeSlice":   makeSlice,
}

// sliceString returns a slice of s from index start to end.
func sliceString(s string, start, end int) string {
	return s[start:end]
}

// makeTable makes a nice-looking textual table from the given data using the
// given template as the rendering for each line.  Columns in the template
// should be separated by a pipe '|'.  Column titles are prepended to the table
// if they exist.
//
// For example where here people is a slice of structs with a Name and Age fields:
//    ```
//    makeTable(people, "{{.Name}}|{{.Age}}", "Name", "Age")
//
//    +----------+-----+
//    |   Name   | Age |
//    +----------+-----+
//    | Bob      |  30 |
//    | Samantha |   3 |
//    +----------+-----+
//    ```
func makeTable(data interface{}, templateStr string, columnTitles ...string) (string, error) {
	t, err := template.New("table").Parse("{{range .}}" + templateStr + "\n{{end}}")
	if err != nil {
		return "", errors.WithMessage(err, "failed to parse table template")
	}
	buf := &bytes.Buffer{}
	hasHeader := len(columnTitles) > 0
	if hasHeader {
		// this can't fail so we drop the error
		_, _ = buf.WriteString(strings.Join(columnTitles, "|") + "\n")
	}
	if err := t.Execute(buf, data); err != nil {
		return "", errors.WithMessage(err, "failed to run table template")
	}
	r := csv.NewReader(buf)
	r.Comma = '|'
	output := &bytes.Buffer{}
	table, err := tablewriter.NewCSVReader(output, r, hasHeader)
	if err != nil {
		return "", errors.WithMessage(err, "failed to render from pipe delimited bytes")
	}
	table.SetAutoFormatHeaders(false)
	table.Render()
	return output.String(), nil
}

// makeSlice returns the arguments as a single slice.  If all the arguments are
// strings, they are returned as a []string, otherwise they're returned as
// []interface{}.
func makeSlice(vals ...interface{}) interface{} {
	ss := make([]string, len(vals))
	for x := range vals {
		if s, ok := vals[x].(string); ok {
			ss[x] = s
		} else {
			return vals
		}
	}
	return ss
}

// makeMap expects an even number of parameters, in order to have name:value
// pairs.  All even values must be strings as keys.  Odd values may be any
// value. This is used to make maps to pass information into sub templates,
// range statements, etc.
func makeMap(vals ...interface{}) (map[string]interface{}, error) {
	if len(vals)%2 != 0 {
		return nil, errors.New("odd number of arguments passed to makeMap")
	}
	ret := make(map[string]interface{}, len(vals)/2)
	for x := 0; x < len(vals); x += 2 {
		s, ok := vals[x].(string)
		if !ok {
			return nil, fmt.Errorf("expected key values to be string, but got %T", vals[x])
		}
		ret[s] = vals[x+1]
	}
	return ret, nil
}
