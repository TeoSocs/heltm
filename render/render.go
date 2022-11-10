package render

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/magiconair/properties"
)

var wg sync.WaitGroup

var ConcurrencyEnabled = true

func ProcessTemplatesIn(basePath string, outPath string, props string) {

	if !ConcurrencyEnabled {
		log.Println("Warning: the concurrent rendering of template is disabled, the ETA will be longer")
	}

	config := properties.MustLoadFile(props, properties.UTF8)

	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error { // TODO: check return errors
		if info.IsDir() {
			log.Println("  scanning", path)
		} else {
			log.Println("processing", path)
			wg.Add(1)
			if ConcurrencyEnabled {
				go render(path, config.Map(), basePath, outPath)
			} else {
				render(path, config.Map(), basePath, outPath)
			}
		}

		return nil
	})

	wg.Wait()
}

func render(path string, config map[string]string, basePath string, outPath string) {
	defer wg.Done()

	t := template.Must(template.New(filepath.Base(path)).Funcs(sprig.FuncMap()).ParseFiles(path))

	writePath := strings.Replace(path, basePath, outPath, 1)

	err := os.MkdirAll(filepath.Dir(writePath), os.ModePerm)
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
}
