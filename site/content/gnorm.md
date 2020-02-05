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
# 	exclude := map[string]bool{".git": true, "site": true, "vendor": true, "_testdata": true}
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
"/gnorm/.idea",
"/gnorm/.idea/dataSources",
"/gnorm/cli",
"/gnorm/cli/testdata",
"/gnorm/database",
"/gnorm/database/drivers",
"/gnorm/database/drivers/mysql",
"/gnorm/database/drivers/mysql/gnorm",
"/gnorm/database/drivers/mysql/gnorm/columns",
"/gnorm/database/drivers/mysql/gnorm/statistics",
"/gnorm/database/drivers/mysql/gnorm/tables",
"/gnorm/database/drivers/mysql/templates",
"/gnorm/database/drivers/postgres",
"/gnorm/database/drivers/postgres/_static",
"/gnorm/database/drivers/postgres/gnorm",
"/gnorm/database/drivers/postgres/gnorm/columns",
"/gnorm/database/drivers/postgres/gnorm/tables",
"/gnorm/database/drivers/postgres/templates",
"/gnorm/database/drivers/sqlite",
"/gnorm/database/drivers/sqlite/gnorm",
"/gnorm/database/drivers/sqlite/gnorm/columns",
"/gnorm/database/drivers/sqlite/gnorm/statistics",
"/gnorm/database/drivers/sqlite/gnorm/tables",
"/gnorm/environ",
"/gnorm/output",
"/gnorm/output/Main",
"/gnorm/output/Main/tables",
"/gnorm/run",
"/gnorm/run/data",
"/gnorm/run/testdata",
"/gnorm/run/testdata/base",
"/gnorm/run/testdata/base/level_one",
"/gnorm/run/testdata/base/level_two",
"/gnorm/static",
"/gnorm/templates",
# [[[end]]]
]
+++
