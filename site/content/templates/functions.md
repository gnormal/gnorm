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
)

func main() {
	fmt.Println("```")
	for _, s := range []string{"FuncMap", "sliceString", "makeTable", "makeSlice"} {
		c := exec.Command("go", "doc","-u", "gnorm.org/gnorm/environ."+s)
		b, err := c.CombinedOutput()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(b))
	}
	fmt.Println("```")
}
gocog}}} -->
```
var FuncMap = map[string]interface{}{
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

	"pascal":     kace.Pascal,
	"camel":      kace.Camel,
	"snake":      kace.Snake,
	"kebab":      kace.Kebab,
	"snakeUpper": kace.SnakeUpper,
	"kebabUpper": kace.KebabUpper,

	"sliceString": sliceString,
	"makeTable":   makeTable,
	"makeSlice":   makeSlice,
}
    FuncMap is the standard list of functions available to templates. Strings
    methods are from the strings package - https://golang.org/pkg/strings/ kace
    methods are from https://github.com/codemodus/kace/


func sliceString(s string, start, end int) string
    sliceString returns a slice of s from index start to end.


func makeTable(data interface{}, templateStr string, columnTitles ...string) (string, error)
    makeTable makes a nice-looking textual table from the given data using the
    given template as the rendering for each line. Columns in the template
    should be separated by a pipe '|'. Column titles are prepended to the table
    if they exist.


func makeSlice(vals ...interface{}) interface{}
    makeSlice returns the arguments as a single slice. If all the arguments are
    strings, they are returned as a []string, otherwise they're returned as
    []interface{}.


```
<!-- {{{end}}} -->
