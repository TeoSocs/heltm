/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "heltm",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		scanAndParse("templates", "out")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.heltm.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var wg sync.WaitGroup

func scanAndParse(basePath string, outPath string) {

	// A sample config
	config := map[string]string{ // TODO: switch to actual config
		"textColor":      "#abcdef",
		"linkColorHover": "#ffaacc",
	}

	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error { // TODO: check return errors
		log.Println(path)
		if !info.IsDir() {
			log.Println("^ is a file")
			wg.Add(1)
			go parse(path, config, basePath, outPath)
		}

		return nil
	})

	wg.Wait()
}

func parse(path string, config map[string]string, basePath string, outPath string) { // function too big: refactor
	defer wg.Done()

	t, err := template.ParseFiles(path)
	if err != nil {
		log.Print(err)
		return
	}

	writePath := strings.Replace(path, basePath, outPath, 1)

	err = os.MkdirAll(filepath.Dir(writePath), os.ModePerm)
	if err != nil {
		log.Println("create directory: ", err)
		return
	}

	f, err := os.Create(writePath)
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = t.Execute(f, config)

	if err != nil {
		log.Print("execute: ", err)
		return
	}
}
