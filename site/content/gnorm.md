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
"/gnorm/cli",
"/gnorm/database",
"/gnorm/database/drivers",
"/gnorm/database/drivers/mysql",
"/gnorm/database/drivers/mysql/db",
"/gnorm/database/drivers/mysql/db/charactersets",
"/gnorm/database/drivers/mysql/db/collationcharactersetapplicability",
"/gnorm/database/drivers/mysql/db/collations",
"/gnorm/database/drivers/mysql/db/columnprivileges",
"/gnorm/database/drivers/mysql/db/columns",
"/gnorm/database/drivers/mysql/db/engines",
"/gnorm/database/drivers/mysql/db/events",
"/gnorm/database/drivers/mysql/db/files",
"/gnorm/database/drivers/mysql/db/globalstatus",
"/gnorm/database/drivers/mysql/db/globalvariables",
"/gnorm/database/drivers/mysql/db/innodbbufferpage",
"/gnorm/database/drivers/mysql/db/innodbbufferpagelru",
"/gnorm/database/drivers/mysql/db/innodbbufferpoolstats",
"/gnorm/database/drivers/mysql/db/innodbcmp",
"/gnorm/database/drivers/mysql/db/innodbcmpmem",
"/gnorm/database/drivers/mysql/db/innodbcmpmemreset",
"/gnorm/database/drivers/mysql/db/innodbcmpperindex",
"/gnorm/database/drivers/mysql/db/innodbcmpperindexreset",
"/gnorm/database/drivers/mysql/db/innodbcmpreset",
"/gnorm/database/drivers/mysql/db/innodbftbeingdeleted",
"/gnorm/database/drivers/mysql/db/innodbftconfig",
"/gnorm/database/drivers/mysql/db/innodbftdefaultstopword",
"/gnorm/database/drivers/mysql/db/innodbftdeleted",
"/gnorm/database/drivers/mysql/db/innodbftindexcache",
"/gnorm/database/drivers/mysql/db/innodbftindextable",
"/gnorm/database/drivers/mysql/db/innodblocks",
"/gnorm/database/drivers/mysql/db/innodblockwaits",
"/gnorm/database/drivers/mysql/db/innodbmetrics",
"/gnorm/database/drivers/mysql/db/innodbsyscolumns",
"/gnorm/database/drivers/mysql/db/innodbsysdatafiles",
"/gnorm/database/drivers/mysql/db/innodbsysfields",
"/gnorm/database/drivers/mysql/db/innodbsysforeign",
"/gnorm/database/drivers/mysql/db/innodbsysforeigncols",
"/gnorm/database/drivers/mysql/db/innodbsysindexes",
"/gnorm/database/drivers/mysql/db/innodbsystables",
"/gnorm/database/drivers/mysql/db/innodbsystablespaces",
"/gnorm/database/drivers/mysql/db/innodbsystablestats",
"/gnorm/database/drivers/mysql/db/innodbsysvirtual",
"/gnorm/database/drivers/mysql/db/innodbtemptableinfo",
"/gnorm/database/drivers/mysql/db/innodbtrx",
"/gnorm/database/drivers/mysql/db/keycolumnusage",
"/gnorm/database/drivers/mysql/db/optimizertrace",
"/gnorm/database/drivers/mysql/db/parameters",
"/gnorm/database/drivers/mysql/db/partitions",
"/gnorm/database/drivers/mysql/db/plugins",
"/gnorm/database/drivers/mysql/db/processlist",
"/gnorm/database/drivers/mysql/db/profiling",
"/gnorm/database/drivers/mysql/db/referentialconstraints",
"/gnorm/database/drivers/mysql/db/routines",
"/gnorm/database/drivers/mysql/db/schemaprivileges",
"/gnorm/database/drivers/mysql/db/schemata",
"/gnorm/database/drivers/mysql/db/sessionstatus",
"/gnorm/database/drivers/mysql/db/sessionvariables",
"/gnorm/database/drivers/mysql/db/statistics",
"/gnorm/database/drivers/mysql/db/tableconstraints",
"/gnorm/database/drivers/mysql/db/tableprivileges",
"/gnorm/database/drivers/mysql/db/tables",
"/gnorm/database/drivers/mysql/db/tablespaces",
"/gnorm/database/drivers/mysql/db/triggers",
"/gnorm/database/drivers/mysql/db/userprivileges",
"/gnorm/database/drivers/mysql/db/views",
"/gnorm/database/drivers/mysql/pg",
"/gnorm/database/drivers/mysql/templates",
"/gnorm/database/drivers/postgres",
"/gnorm/database/drivers/postgres/pg",
"/gnorm/database/drivers/postgres/templates",
"/gnorm/environ",
"/gnorm/run",
# [[[end]]]
]
+++
