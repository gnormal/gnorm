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
var FuncMap = map[string]interface{}{
	// funcs from strings package
	"compare":      strings.Compare,
	"contains":     strings.Contains,
	"containsAny":  strings.ContainsAny,
	"count":        strings.Count,
	"equalFold":    strings.EqualFold,
	"fields":       strings.Fields,
	"hasPrefix":    strings.HasPrefix,
	"hasSuffix":    strings.HasPrefix,
	"index":        strings.Index,
	"indexAny":     strings.IndexAny,
	"join":         strings.Join,
	"lastIndex":    strings.LastIndex,
	"lastIndexAny": strings.LastIndexAny,
	"repeat":       strings.Repeat,
	"replace":      strings.Replace,
	"split":        strings.Split,
	"splitAfter":   strings.SplitAfter,
	"splitAfterN":  strings.SplitAfterN,
	"splitN":       strings.SplitN,
	"title":        strings.Title,
	"toLower":      strings.ToLower,
	"toTitle":      strings.ToTitle,
	"toUpper":      strings.ToUpper,
	"trim":         strings.Trim,
	"trimLeft":     strings.TrimLeft,
	"trimPrefix":   strings.TrimPrefix,
	"trimRight":    strings.TrimRight,
	"trimSpace":    strings.TrimSpace,
	"trimSuffix":   strings.TrimSuffix,

	"sliceString": sliceString,
	"makeTable":   makeTable,
	"makeSlice":   makeSlice,
	"pascal":      kace.Pascal,
	"camel":       kace.Camel,
	"snake":       kace.Snake,
	"kebab":       kace.Kebab,
	"snakeUpper":  kace.SnakeUpper,
	"kebabUpper":  kace.KebabUpper,
}

// sliceString returns a slice of s from index start to end.
func sliceString(s string, start, end int) string {
	return s[start:end]
}

// makeTable makes a nice-looking textual table from the given data using the
// given template as the rendering for each line.  Columns in the template
// should be separated by a pipe '|'.  Column titles are prepended to the table
// if they exist.
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
