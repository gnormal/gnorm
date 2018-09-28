+++
title= "Template Functions"
weight=2
alwaysopen=true
+++

Note that template functions are not available for external template engines.

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
	stringsPkg = regexp.MustCompile(`"(.*?)".*?(strings\.(.*))$`)
	kacePkg    = regexp.MustCompile(`"(.*?)".*?(kace\.(.*))`)
	localFunc  = regexp.MustCompile(`"(.*?)".*?[a-zA-Z0-9_]+$`)
)

func main() {
	c := exec.Command("go", "doc", "-u", "gnorm.org/gnorm/environ.FuncMap")
	b, err := c.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("<table>")

	var locals []string
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
		case stringsPkg.MatchString(s):
			s = stringsPkg.ReplaceAllString(s, `<tr><td>$1</td><td>[https://golang.org/pkg/strings/#$3](https://golang.org/pkg/strings/#$3)</td></tr>`)

		case kacePkg.MatchString(s):
			s = kacePkg.ReplaceAllString(s, `<tr><td>$1</td><td>[https://godoc.org/github.com/codemodus/kace#$3](https://godoc.org/github.com/codemodus/kace#$3)</td></tr>`)
		case localFunc.MatchString(s):
			name := localFunc.ReplaceAllString(s, "$1")
			locals = append(locals, name)
			lowername := strings.ToLower(name)
			s = fmt.Sprintf("<tr><td>%s</td><td>[%s (see below)](/templates/functions/#%s)</td></tr>", name, name, lowername)
		}
		fmt.Println(s)
	}
	fmt.Println("</table>")

	for _, s := range locals {
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
<table>
<tr><td>camel</td><td>[https://godoc.org/github.com/codemodus/kace#Camel](https://godoc.org/github.com/codemodus/kace#Camel)</td></tr>
<tr><td>compare</td><td>[https://golang.org/pkg/strings/#Compare](https://golang.org/pkg/strings/#Compare)</td></tr>
<tr><td>contains</td><td>[https://golang.org/pkg/strings/#Contains](https://golang.org/pkg/strings/#Contains)</td></tr>
<tr><td>containsAny</td><td>[https://golang.org/pkg/strings/#ContainsAny](https://golang.org/pkg/strings/#ContainsAny)</td></tr>
<tr><td>count</td><td>[https://golang.org/pkg/strings/#Count](https://golang.org/pkg/strings/#Count)</td></tr>
<tr><td>dec</td><td>[dec (see below)](/templates/functions/#dec)</td></tr>
<tr><td>equalFold</td><td>[https://golang.org/pkg/strings/#EqualFold](https://golang.org/pkg/strings/#EqualFold)</td></tr>
<tr><td>fields</td><td>[https://golang.org/pkg/strings/#Fields](https://golang.org/pkg/strings/#Fields)</td></tr>
<tr><td>hasPrefix</td><td>[https://golang.org/pkg/strings/#HasPrefix](https://golang.org/pkg/strings/#HasPrefix)</td></tr>
<tr><td>hasSuffix</td><td>[https://golang.org/pkg/strings/#HasPrefix](https://golang.org/pkg/strings/#HasPrefix)</td></tr>
<tr><td>inc</td><td>[inc (see below)](/templates/functions/#inc)</td></tr>
<tr><td>strIndex</td><td>[https://golang.org/pkg/strings/#Index](https://golang.org/pkg/strings/#Index)</td></tr>
<tr><td>indexAny</td><td>[https://golang.org/pkg/strings/#IndexAny](https://golang.org/pkg/strings/#IndexAny)</td></tr>
<tr><td>join</td><td>[https://golang.org/pkg/strings/#Join](https://golang.org/pkg/strings/#Join)</td></tr>
<tr><td>kebab</td><td>[https://godoc.org/github.com/codemodus/kace#Kebab](https://godoc.org/github.com/codemodus/kace#Kebab)</td></tr>
<tr><td>kebabUpper</td><td>[https://godoc.org/github.com/codemodus/kace#KebabUpper](https://godoc.org/github.com/codemodus/kace#KebabUpper)</td></tr>
<tr><td>lastIndex</td><td>[https://golang.org/pkg/strings/#LastIndex](https://golang.org/pkg/strings/#LastIndex)</td></tr>
<tr><td>lastIndexAny</td><td>[https://golang.org/pkg/strings/#LastIndexAny](https://golang.org/pkg/strings/#LastIndexAny)</td></tr>
<tr><td>makeMap</td><td>[makeMap (see below)](/templates/functions/#makemap)</td></tr>
<tr><td>makeSlice</td><td>[makeSlice (see below)](/templates/functions/#makeslice)</td></tr>
<tr><td>numbers</td><td>[numbers (see below)](/templates/functions/#numbers)</td></tr>
<tr><td>pascal</td><td>[https://godoc.org/github.com/codemodus/kace#Pascal](https://godoc.org/github.com/codemodus/kace#Pascal)</td></tr>
<tr><td>repeat</td><td>[https://golang.org/pkg/strings/#Repeat](https://golang.org/pkg/strings/#Repeat)</td></tr>
<tr><td>replace</td><td>[https://golang.org/pkg/strings/#Replace](https://golang.org/pkg/strings/#Replace)</td></tr>
<tr><td>sliceString</td><td>[sliceString (see below)](/templates/functions/#slicestring)</td></tr>
<tr><td>snake</td><td>[https://godoc.org/github.com/codemodus/kace#Snake](https://godoc.org/github.com/codemodus/kace#Snake)</td></tr>
<tr><td>snakeUpper</td><td>[https://godoc.org/github.com/codemodus/kace#SnakeUpper](https://godoc.org/github.com/codemodus/kace#SnakeUpper)</td></tr>
<tr><td>split</td><td>[https://golang.org/pkg/strings/#Split](https://golang.org/pkg/strings/#Split)</td></tr>
<tr><td>splitAfter</td><td>[https://golang.org/pkg/strings/#SplitAfter](https://golang.org/pkg/strings/#SplitAfter)</td></tr>
<tr><td>splitAfterN</td><td>[https://golang.org/pkg/strings/#SplitAfterN](https://golang.org/pkg/strings/#SplitAfterN)</td></tr>
<tr><td>splitN</td><td>[https://golang.org/pkg/strings/#SplitN](https://golang.org/pkg/strings/#SplitN)</td></tr>
<tr><td>sub</td><td>[sub (see below)](/templates/functions/#sub)</td></tr>
<tr><td>sum</td><td>[sum (see below)](/templates/functions/#sum)</td></tr>
<tr><td>title</td><td>[https://golang.org/pkg/strings/#Title](https://golang.org/pkg/strings/#Title)</td></tr>
<tr><td>toLower</td><td>[https://golang.org/pkg/strings/#ToLower](https://golang.org/pkg/strings/#ToLower)</td></tr>
<tr><td>toTitle</td><td>[https://golang.org/pkg/strings/#ToTitle](https://golang.org/pkg/strings/#ToTitle)</td></tr>
<tr><td>toUpper</td><td>[https://golang.org/pkg/strings/#ToUpper](https://golang.org/pkg/strings/#ToUpper)</td></tr>
<tr><td>trim</td><td>[https://golang.org/pkg/strings/#Trim](https://golang.org/pkg/strings/#Trim)</td></tr>
<tr><td>trimLeft</td><td>[https://golang.org/pkg/strings/#TrimLeft](https://golang.org/pkg/strings/#TrimLeft)</td></tr>
<tr><td>trimPrefix</td><td>[https://golang.org/pkg/strings/#TrimPrefix](https://golang.org/pkg/strings/#TrimPrefix)</td></tr>
<tr><td>trimRight</td><td>[https://golang.org/pkg/strings/#TrimRight](https://golang.org/pkg/strings/#TrimRight)</td></tr>
<tr><td>trimSpace</td><td>[https://golang.org/pkg/strings/#TrimSpace](https://golang.org/pkg/strings/#TrimSpace)</td></tr>
<tr><td>trimSuffix</td><td>[https://golang.org/pkg/strings/#TrimSuffix](https://golang.org/pkg/strings/#TrimSuffix)</td></tr>
</table>
## dec
` func dec(x int) int `

dec decrements the argument's value by 1.
## inc
` func inc(x int) int `

inc increments the argument's value by 1.
## makeMap
` func makeMap(vals ...interface{}) (map[string]interface{}, error) `

makeMap expects an even number of parameters, in order to have name:value
pairs. All even values must be strings as keys. Odd values may be any value.
This is used to make maps to pass information into sub templates, range
statements, etc.
## makeSlice
` func makeSlice(vals ...interface{}) interface{} `

makeSlice returns the arguments as a single slice. If all the arguments are
strings, they are returned as a []string, otherwise they're returned as
[]interface{}.
## numbers
` func numbers(start, end int) data.Strings `

numbers returns a slice of strings of the numbers start to end (inclusive).
## sliceString
` func sliceString(s string, start, end int) string `

sliceString returns a slice of s from index start to end.
## sub
` func sub(x int, vals ...int) int `

sub subtracts the second and following values from the first argument.
## sum
` func sum(vals ...int) int `

sum returns the sum of its arguments.
<!-- {{{end}}} -->
