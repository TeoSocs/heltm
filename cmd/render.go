/*
Copyright © 2022 MATTEO SOVILLA teo.sovi@gmail.com
*/
package cmd

import (
	"fmt"
	"heltm/render"

	"github.com/spf13/cobra"
)

var fromParam = "from"
var outputParam = "out"
var valuesParam = "values"

var defaultFrom = "templates"
var defaultOut = "out"
var defaultProps = "values.properties"

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "locally render templates",
	Long: fmt.Sprintf(`Render templates locally and save the output.

The default folder for the templates is "%s", and the default output
location is "%s". The "out" folder will be created if nonexistent, and its
content will be overwrited if there is any.
The default values location is the file "%s".
Any rendered file will maintain the same name and path relative to the output
folder.

The expected format is the go template format, with the functions from Sprig.
No validation is performed over the output files.`, defaultFrom, defaultOut, defaultProps),
	Run: func(cmd *cobra.Command, args []string) {
		c, _ := cmd.Flags().GetBool(disableConcurrencyParam)
		if c {
			render.ConcurrencyEnabled = false
		}
		from, _ := cmd.Flags().GetString(fromParam)
		out, _ := cmd.Flags().GetString(outputParam)
		props, _ := cmd.Flags().GetString(valuesParam)
		render.ProcessTemplatesIn(from, out, props)
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// renderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	renderCmd.Flags().StringP(fromParam, "f", defaultFrom, "Templates directory")
	renderCmd.Flags().StringP(outputParam, "o", defaultOut, "Output directory")
	renderCmd.Flags().StringP(valuesParam, "p", defaultProps, "Template values file")
}
