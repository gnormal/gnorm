+++
title= "Template Functions"
weight=2
alwaysopen=true
+++

<!-- {{{gocog
package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	stringrepl = regexp.MustCompile(`(strings\.(.*))`)
	kacerepl   = regexp.MustCompile(`(kace\.(.*))`)
)

func main() {
	c := exec.Command("go", "doc", "-u", "gnorm.org/gnorm/environ.FuncMap")
	b, err := c.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s := string(b)
	// strip off the first line about funcmap itself
	lines := strings.Split(s, "\n")[1:]
	for _, s := range lines {
		if len(s) == 0 {
			continue
		}
		if len(s) == 1 {
			break
		}
		// trim the trailing comma and indentation
		s := strings.TrimSpace(s[:len(s)-1])
		switch {
		case stringrepl.MatchString(s):
			s = stringrepl.ReplaceAllString(s, `[strings.$2](https://golang.org/pkg/strings/#$2)`)

		case kacerepl.MatchString(s):
			s = kacerepl.ReplaceAllString(s, `[kace.$2](https://godoc.org/github.com/codemodus/kace#$2)`)
		}
		fmt.Println("-", s)
	}
	fmt.Println()

	for _, s := range []string{"sliceString", "makeTable", "makeSlice"} {
		fmt.Println("##", s)
		c := exec.Command("go", "doc", "-u", "gnorm.org/gnorm/environ."+s)
		b, err := c.CombinedOutput()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		s := strings.TrimSpace(string(b))
		vals := strings.Split(s, "\n")
		fmt.Println("`", vals[0], "`")
		fmt.Println()
		for _, s := range vals[1:] {
			fmt.Println(strings.TrimSpace(s))
		}
	}
}
gocog}}} -->
- "compare":      [strings.Compare](https://golang.org/pkg/strings/#Compare)
- "contains":     [strings.Contains](https://golang.org/pkg/strings/#Contains)
- "containsAny":  [strings.ContainsAny](https://golang.org/pkg/strings/#ContainsAny)
- "count":        [strings.Count](https://golang.org/pkg/strings/#Count)
- "equalFold":    [strings.EqualFold](https://golang.org/pkg/strings/#EqualFold)
- "fields":       [strings.Fields](https://golang.org/pkg/strings/#Fields)
- "hasPrefix":    [strings.HasPrefix](https://golang.org/pkg/strings/#HasPrefix)
- "hasSuffix":    [strings.HasPrefix](https://golang.org/pkg/strings/#HasPrefix)
- "index":        [strings.Index](https://golang.org/pkg/strings/#Index)
- "indexAny":     [strings.IndexAny](https://golang.org/pkg/strings/#IndexAny)
- "join":         [strings.Join](https://golang.org/pkg/strings/#Join)
- "lastIndex":    [strings.LastIndex](https://golang.org/pkg/strings/#LastIndex)
- "lastIndexAny": [strings.LastIndexAny](https://golang.org/pkg/strings/#LastIndexAny)
- "repeat":       [strings.Repeat](https://golang.org/pkg/strings/#Repeat)
- "replace":      [strings.Replace](https://golang.org/pkg/strings/#Replace)
- "split":        [strings.Split](https://golang.org/pkg/strings/#Split)
- "splitAfter":   [strings.SplitAfter](https://golang.org/pkg/strings/#SplitAfter)
- "splitAfterN":  [strings.SplitAfterN](https://golang.org/pkg/strings/#SplitAfterN)
- "splitN":       [strings.SplitN](https://golang.org/pkg/strings/#SplitN)
- "title":        [strings.Title](https://golang.org/pkg/strings/#Title)
- "toLower":      [strings.ToLower](https://golang.org/pkg/strings/#ToLower)
- "toTitle":      [strings.ToTitle](https://golang.org/pkg/strings/#ToTitle)
- "toUpper":      [strings.ToUpper](https://golang.org/pkg/strings/#ToUpper)
- "trim":         [strings.Trim](https://golang.org/pkg/strings/#Trim)
- "trimLeft":     [strings.TrimLeft](https://golang.org/pkg/strings/#TrimLeft)
- "trimPrefix":   [strings.TrimPrefix](https://golang.org/pkg/strings/#TrimPrefix)
- "trimRight":    [strings.TrimRight](https://golang.org/pkg/strings/#TrimRight)
- "trimSpace":    [strings.TrimSpace](https://golang.org/pkg/strings/#TrimSpace)
- "trimSuffix":   [strings.TrimSuffix](https://golang.org/pkg/strings/#TrimSuffix)
- "camel":      [kace.Camel](https://godoc.org/github.com/codemodus/kace#Camel)
- "kebab":      [kace.Kebab](https://godoc.org/github.com/codemodus/kace#Kebab)
- "kebabUpper": [kace.KebabUpper](https://godoc.org/github.com/codemodus/kace#KebabUpper)
- "pascal":     [kace.Pascal](https://godoc.org/github.com/codemodus/kace#Pascal)
- "snake":      [kace.Snake](https://godoc.org/github.com/codemodus/kace#Snake)
- "snakeUpper": [kace.SnakeUpper](https://godoc.org/github.com/codemodus/kace#SnakeUpper)
- "sliceString": sliceString
- "makeTable":   makeTable
- "makeSlice":   makeSlice

## sliceString
` func sliceString(s string, start, end int) string `

sliceString returns a slice of s from index start to end.
## makeTable
` func makeTable(data interface{}, templateStr string, columnTitles ...string) (string, error) `

makeTable makes a nice-looking textual table from the given data using the
given template as the rendering for each line. Columns in the template
should be separated by a pipe '|'. Column titles are prepended to the table
if they exist.

For example where here people is a slice of structs with a Name and Age
fields:

```
makeTable(people, "{{.Name}}|{{.Age}}", "Name", "Age")

+----------+-----+
|   Name   | Age |
+----------+-----+
| Bob      |  30 |
| Samantha |   3 |
+----------+-----+
```
## makeSlice
` func makeSlice(vals ...interface{}) interface{} `

makeSlice returns the arguments as a single slice. If all the arguments are
strings, they are returned as a []string, otherwise they're returned as
[]interface{}.
<!-- {{{end}}} -->
