+++
date = 2016-10-02T23:00:00Z
title = "Gnorm"
vanity = "https://github.com/gnormal/gnorm"
type="code"
aliases = [
# [[[gocog
# package main
# 
# import (
# 	"fmt"
# 	"os"
# 	"path/filepath"
# )
# 
# func main() {
# 	exclude := map[string]bool{".git": true, "site": true, "vendor": true}
# 	// this runs in the root dir of the repo, thanks to go:generate, so
# 	// this is actually getting the top level subdirectories
# 	err := filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
# 		if exclude[path] {
# 			return filepath.SkipDir
# 		}
# 		if f.IsDir() && path != "." {
# 			fmt.Printf("\"/gnorm/%s\",\n", path)
# 		}
# 		return nil
# 	})
# 	if err != nil {
# 		fmt.Fprintf(os.Stderr, "%v\n", err)
# 		os.Exit(1)
# 	}
# }
# gocog]]]
"/gnorm/_testdata",
"/gnorm/_testdata/templates",
"/gnorm/cli",
"/gnorm/database",
"/gnorm/database/drivers",
"/gnorm/database/drivers/postgres",
"/gnorm/database/drivers/postgres/pg",
"/gnorm/database/drivers/postgres/templates",
"/gnorm/environ",
"/gnorm/run",
# [[[end]]]
]
+++
