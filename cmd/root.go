/*
Copyright © 2022 MATTEO SOVILLA teo.sovi@gmail.com
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "heltm",
	Short: "a go template engine based on Helm",
	Long: `A command line utility to locally render multiple template files using the go
template language and the functions provided by Sprig library.

No validation is performed over the output files, thus making the tool agnostic
respect to the intended use of the templated files.

The tool makes heavy use of goroutines, a parameter is provided to disable the
concurrency if needed.

REFERENCE DOCUMENTATION:
Golang template language:   https://pkg.go.dev/text/template
Sprig template functions:   https://pkg.go.dev/github.com/Masterminds/sprig
Heltm repository:           https://github.com/TeoSocs/heltm
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var disableConcurrencyParam = "non-parallel"

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.heltm.yaml)")
	rootCmd.PersistentFlags().BoolP(disableConcurrencyParam, "", false, "Force non-parallel execution")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
