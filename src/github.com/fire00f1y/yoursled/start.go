package main

import (
	"io/ioutil"
	"os"
	"flag"
	"strings"
	"path/filepath"
	"log"
	"os/exec"
	"github.com/fire00f1y/yoursled/constants"
)

var simultaneousRuns = 2
var scriptLocation string

func runBatFiles(files []string, suffix string) (int) {
	log.Printf("Files to run: %s\n", files)
	// To pause execution and see output
	//bufio.NewReader(os.Stdin).ReadBytes('\n')

	cmdStack := make([]exec.Cmd, 0)
	for _, fileName := range files {
		cmd := exec.Cmd{Path: fileName, Dir: scriptLocation, Stdout: os.Stdout, Stderr: os.Stderr}
		cmdStack = append(cmdStack, cmd)
	}

	ch := make(chan error, len(cmdStack))
	cmdChan := make(chan exec.Cmd, len(cmdStack))

	for r := 1; r <= simultaneousRuns; r++ {
		go runner(r, cmdChan, ch)
	}
	for _, c := range cmdStack {
		cmdChan <- c
	}
	close(cmdChan)
	for i := 0; i < len(cmdStack); i++ {
		err := <-ch
		if err != nil {
			log.Fatalln(err)
		}
	}
	close(ch)

	return len(cmdStack)
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
	workers := flag.Int("runners", simultaneousRuns, "Number of scripts to execute simultaneously")
	job := flag.String("job", constants.BatchScriptsJob, "Which utility job to run")
	flag.Parse()

	simultaneousRuns = *workers
	if filePath != nil && *filePath != "" {
		scriptLocation = *filePath
	}
	if !(strings.HasSuffix(scriptLocation, "/") || strings.HasSuffix(scriptLocation, "\\")) {
		scriptLocation = scriptLocation + "/"
	}

	switch *job {
	case constants.BatchScriptsJob:
		{
			log.Printf("Looking for files to run in %s\n", scriptLocation)
			runFiles := getFiles(*prefix, *suffix)
			log.Printf("Finished processing %d files\n", runBatFiles(runFiles, *suffix))
		}
	default:
		{
			log.Fatalf("Unknown job: %s\n", *job)
		}
	}

	os.Exit(0)
}

func getFiles(prefix, suffix string) ([]string) {
	runFiles := make([]string, 0)
	files, err := ioutil.ReadDir(scriptLocation)
	if err != nil {
		log.Fatalf("YOU GOT PROBLEMS: %v\n", err)
	}

	for _, v := range files {
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
	lstat, err := os.Lstat(".")
	folder, err := filepath.Abs(filepath.Dir(lstat.Name()))
	if err != nil {
		log.Fatalf("Error getting current dir: %v\n", err)
	}
	return folder
}
