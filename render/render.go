package render

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v3"
)

var wg sync.WaitGroup

var ConcurrencyEnabled = true

func mustReadYamlFile(path string) map[string]interface{} {
	valueMap := make(map[string]interface{})
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("error reading file:", err)
	}
	err = yaml.Unmarshal(f, &valueMap)
	if err != nil {
		log.Fatalln("error unmarshaling yaml:", err)
	}
	return valueMap
}

func ProcessTemplatesIn(basePath string, outPath string, valuesPath string) {

	if !ConcurrencyEnabled {
		log.Println("Warning: the concurrent rendering of template is disabled, the ETA will be longer")
	}

	config := mustReadYamlFile(valuesPath)

	processNode := func(path string, info os.FileInfo, err error) error { // ignored errors in render function: I want the script to continue anyway
		if info.IsDir() {
			log.Println("  scanning", path)
		} else {
			log.Println("processing", path)
			wg.Add(1)
			if ConcurrencyEnabled {
				go render(path, config, basePath, outPath)
			} else {
				render(path, config, basePath, outPath)
			}
		}
		return nil
	}

	filepath.Walk(basePath, processNode)

	wg.Wait()
}

func render(path string, config map[string]interface{}, basePath string, outPath string) (err error) {
	defer wg.Done()

	t := template.Must(template.New(filepath.Base(path)).Funcs(sprig.FuncMap()).ParseFiles(path))

	writePath := strings.Replace(path, basePath, outPath, 1)

	err = os.MkdirAll(filepath.Dir(writePath), os.ModePerm)
	if err != nil {
		log.Println(path, " error creating directory: ", err)
		return
	}

	f, err := os.Create(writePath)
	if err != nil {
		log.Println(path, " error creating file: ", err)
		return
	}

	err = t.Execute(f, config)

	if err != nil {
		log.Println(path, " error executing: ", err)
		return
	}

	return nil
}
