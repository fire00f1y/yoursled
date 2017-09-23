package main

import (
	"io/ioutil"
	"os"
	"flag"
	"strings"
	"path/filepath"
	"log"
	"os/exec"
)

const simultaneousRuns = 2
var scriptLocation string

func runBatFiles(files []string, suffix string) (int) {
	log.Printf("Files to run: %s\n", files)
	// To pause execution and see output
	//bufio.NewReader(os.Stdin).ReadBytes('\n')

	cmdStack := make([]exec.Cmd, 0)
	for _,fileName := range files {
		cmd := exec.Cmd{Path:fileName, Dir: scriptLocation, Stdout: os.Stdout, Stderr: os.Stderr}
		cmdStack = append(cmdStack, cmd)
	}

	log.Printf("Size of cmdStack: %d\n", len(cmdStack))
	ch := make(chan error, len(cmdStack))
	cmdChan := make(chan exec.Cmd, len(cmdStack))

	for r := 1; r < simultaneousRuns; r++ {
		go runner(r, cmdChan, ch)
	}
	for _,c := range cmdStack {
		cmdChan <- c
	}
	close(cmdChan)
	for i := 0; i < len(cmdStack); i++ {
		err := <- ch
		if err != nil {
			log.Fatalln(err)
		}
	}
	close(ch)

	return 0
}

func runner(id int, jobs <-chan exec.Cmd, result chan<- error) {
	for job := range jobs {
		log.Printf("[%d] Started Running %s\n", id, job.Path)
		result <- job.Run()
		log.Printf("[%d] Finished Running %s\n", id, job.Path)
	}
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

	log.Printf("Looking for files to run in %s\n", scriptLocation)
	runFiles := getFiles(*prefix, *suffix)
	os.Exit(runBatFiles(runFiles, *suffix))
}

func getFiles(prefix, suffix string) ([]string) {
	runFiles := make([]string, 0)
	files,err := ioutil.ReadDir(scriptLocation)
	if err != nil {
		log.Fatalf("YOU GOT PROBLEMS: %v\n", err)
	}

	for _,v := range files {
		if v.IsDir() {
			continue
		}
		if !strings.HasPrefix(v.Name(), prefix) {
			continue
		}
		if !strings.HasSuffix(v.Name(), suffix) {
			continue
		}
		runFiles = append(runFiles, v.Name())
	}
	return runFiles
}

func getCurrentDir() (string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}