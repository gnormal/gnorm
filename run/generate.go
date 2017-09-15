package run // import "gnorm.org/gnorm/run"

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/environ"
)

// Generate reads your database, gets the schema for it, and then generates
// files based on your templates and your configuration.
func Generate(env environ.Values, cfg *Config) error {
	info, err := cfg.Driver.Parse(env.Log, cfg.ConnStr, cfg.Schemas, makeFilter(cfg.IncludeTables, cfg.ExcludeTables))
	if err != nil {
		return err
	}
	if err := convertNames(env.Log, info, cfg); err != nil {
		return err
	}
	if len(cfg.SchemaPaths) == 0 {
		env.Log.Println("No SchemaPaths specified, skipping schemas.")
	} else {
		if err := generateSchemas(env, cfg, info); err != nil {
			return err
		}
	}
	if len(cfg.EnumPaths) == 0 {
		env.Log.Println("No EnumPath specified, skipping enums.")
	} else {
		if err := generateEnums(env, cfg, info); err != nil {
			return err
		}
	}
	if len(cfg.TablePaths) == 0 {
		env.Log.Println("No table path specified, skipping tables.")
	} else {
		if err := generateTables(env, cfg, info); err != nil {
			return err
		}
	}
	return nil
}

func generateSchemas(env environ.Values, cfg *Config, info *database.Info) error {
	for _, schema := range info.Schemas {
		for _, target := range cfg.SchemaPaths {
			if err := generateSchema(env, schema, target, cfg.PostRun); err != nil {
				return err
			}
		}
	}
	return nil
}

func generateSchema(env environ.Values, schema *database.Schema, target OutputTarget, postrun []string) error {
	env.Log.Printf("Generating output for schema %v", schema.Name)
	buf := &bytes.Buffer{}
	err := target.Filename.Execute(buf, struct{ Schema string }{Schema: schema.Name})
	if err != nil {
		return errors.WithMessage(err, "failed to run SchemaPath Filename template with schema "+schema.Name)
	}
	outputPath := buf.String()
	if err := os.MkdirAll(filepath.Dir(outputPath), 0700); err != nil {
		return errors.WithMessage(err, "error creating output directory for schema "+schema.Name)
	}
	f, err := os.OpenFile(outputPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return errors.WithMessage(err, "failed to create output file for schema "+schema.Name)
	}
	defer f.Close()
	if err := target.Contents.Execute(f, schema); err != nil {
		return errors.WithMessage(err, "failed to run schema template over schema "+schema.Name)
	}
	if err := f.Close(); err != nil {
		return errors.Wrapf(err, "error closing generated file %q", outputPath)
	}
	if len(postrun) > 0 {
		return doPostRun(env, outputPath, postrun)
	}
	return nil
}

func generateEnums(env environ.Values, cfg *Config, info *database.Info) error {
	for _, schema := range info.Schemas {
		for _, enum := range schema.Enums {
			for _, target := range cfg.EnumPaths {
				if err := generateEnum(env, enum, target, cfg.PostRun); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func generateEnum(env environ.Values, enum *database.Enum, target OutputTarget, postrun []string) error {
	env.Log.Printf("Generating output for enum %v", enum.Name)
	buf := &bytes.Buffer{}
	err := target.Filename.Execute(buf, struct{ Schema, Enum string }{Schema: enum.Schema, Enum: enum.Name})
	if err != nil {
		return errors.Wrapf(err, "failed to run enum filename template with enum %v.%v"+enum.Schema, enum.Name)
	}
	outputPath := buf.String()
	if err := os.MkdirAll(filepath.Dir(outputPath), 0700); err != nil {
		return errors.Wrapf(err, "error creating output directory for enum %v.%v "+enum.Schema, enum.Name)
	}
	f, err := os.OpenFile(outputPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrapf(err, "failed to create output file for enum %v.%v"+enum.Schema, enum.Name)
	}
	defer f.Close()
	if err := target.Contents.Execute(f, enum); err != nil {
		return errors.Wrapf(err, "failed to run enum contents template over enum %v.%v"+enum.Schema, enum.Name)
	}
	if err := f.Close(); err != nil {
		return errors.Wrapf(err, "error closing generated file %q", outputPath)
	}
	if len(postrun) > 0 {
		return doPostRun(env, outputPath, postrun)
	}
	return nil
}

func generateTables(env environ.Values, cfg *Config, info *database.Info) error {
	for _, schema := range info.Schemas {
		for _, table := range schema.Tables {
			for _, target := range cfg.TablePaths {
				if err := generateTable(env, table, target, cfg.PostRun); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func generateTable(env environ.Values, table *database.Table, target OutputTarget, postrun []string) error {
	env.Log.Printf("Generating output for table %v", table.Name)
	buf := &bytes.Buffer{}
	err := target.Filename.Execute(buf, struct{ Schema, Table string }{Schema: table.Schema, Table: table.Name})
	if err != nil {
		return errors.Wrapf(err, "failed to run table filename template with table %v.%v"+table.Schema, table.Name)
	}
	outputPath := buf.String()
	if err := os.MkdirAll(filepath.Dir(outputPath), 0700); err != nil {
		return errors.Wrapf(err, "error creating output directory for table %v.%v "+table.Schema, table.Name)
	}
	f, err := os.OpenFile(outputPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrapf(err, "failed to create output file for table %v.%v"+table.Schema, table.Name)
	}
	defer f.Close()
	if err := target.Contents.Execute(f, table); err != nil {
		return errors.Wrapf(err, "failed to run table contents template over table %v.%v"+table.Schema, table.Name)
	}
	if err := f.Close(); err != nil {
		return errors.Wrapf(err, "error closing generated file %q", outputPath)
	}
	if len(postrun) > 0 {
		return doPostRun(env, outputPath, postrun)
	}
	return nil
}

func doPostRun(env environ.Values, file string, postrun []string) error {
	newenv := make(map[string]string, len(env.Env)+1)
	for k := range env.Env {
		newenv[k] = env.Env[k]
	}
	run := make([]string, len(postrun))
	newenv["GNORMFILE"] = file
	conv := func(s string) string { return newenv[s] }
	for x, s := range postrun {
		run[x] = os.Expand(s, conv)
	}
	var cmd *exec.Cmd
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if len(run) > 1 {
		cmd = exec.CommandContext(ctx, run[0], run[1:]...)
	} else {
		cmd = exec.CommandContext(ctx, run[0])
	}
	cmd.Stderr = env.Stderr
	cmd.Stdout = env.Stdout
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error running postrun command %q", run)
	}
	return nil
}
