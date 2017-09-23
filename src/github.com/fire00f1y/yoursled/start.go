package main

import (
	"io/ioutil"
	"os"
	"flag"
	"strings"
	"path/filepath"
	"log"
)

var scriptLocation string

func getCurrentDir() (string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func main() {
	scriptLocation = getCurrentDir()

	suffix := flag.String("ext", ".bat", "Extension of files to run.")
	prefix := flag.String("base", "", "Base file name to run. If you have file1.bat, file2.bat you would give '-base=file' Give nothing if you want to run all files of that type")
	filePath := flag.String("path", "", "Path to script files directory. Takes current directory if not given")
	flag.Parse()

	if filePath != nil && *filePath != "" {
		scriptLocation = *filePath
	}
	if !(strings.HasSuffix(scriptLocation, "/") || strings.HasSuffix(scriptLocation, "\\")) {
		scriptLocation = scriptLocation + "/"
	}

	files,err := ioutil.ReadDir(scriptLocation)
	if err != nil {
		log.Fatalf("PROBLEMS: %v\n", err)
	}

	runFiles := make([]string, 0)
	log.Printf("Looking for files to run in %s\n", scriptLocation)
	for _,v := range files {
		if v.IsDir() {
			continue
		}
		if !strings.HasPrefix(v.Name(), *prefix) {
			continue
		}
		if !strings.HasSuffix(v.Name(), *suffix) {
			continue
		}
		runFiles = append(runFiles, v.Name())
	}
	os.Exit(runBatFiles(runFiles))
}

func runBatFiles(files []string) (int) {
	log.Printf("Files to run: %s\n", files)
	// To pause execution and see output
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	return 0
}