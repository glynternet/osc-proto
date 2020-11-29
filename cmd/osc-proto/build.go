package main

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/glynternet/osc-proto/pkg/generate"
	"github.com/glynternet/osc-proto/pkg/generate/golang"
	"github.com/glynternet/osc-proto/pkg/types"
	"github.com/glynternet/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func buildCmdTree(logger log.Logger, _ io.Writer, rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{
		Use:  "generate TYPES_FILE",
		Args: cobra.ExactArgs(1),
		RunE: generateTypesRun(logger),
	})
}

func generateTypesRun(logger log.Logger) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		typesPath := args[0]
		typesSerialised, err := ioutil.ReadFile(typesPath)
		if err != nil {
			return errors.Wrap(err, "reading types file")
		}
		var typesDeserialised types.Types
		if err := yaml.Unmarshal(typesSerialised, &typesDeserialised); err != nil {
			return errors.Wrap(err, "deserialising types definition")
		}

		allFiles := make(map[string][]byte)
		for _, generator := range []generate.Generator{
			golang.Generator{Package: "types"},
		} {
			outFiles, err := generator.Generate(typesDeserialised)
			if err != nil {
				return errors.Wrap(err, "generating files content for types")
			}
			for path, content := range outFiles {
				if _, exists := allFiles[path]; exists {
					return fmt.Errorf("multiple generators would be writing content over each other at %s", path)
				}
				allFiles[path] = content
			}
		}
		for path, content := range allFiles {
			if err := ioutil.WriteFile(path, content, 0640); err != nil {
				return errors.Wrap(err, "writing content to file")
			}
			if err := logger.Log(
				log.Message("Generated file written."),
				log.KV{
					K: "path",
					V: path,
				}); err != nil {
				return errors.Wrap(err, "logging message")
			}
		}
		return nil
	}
}
