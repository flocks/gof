package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	FileExist = _fileExist
)

type Filematch struct {
	filePath string
	line     int
	col      int
	desc     string // could potentially hold error/warning from program like linter
}

func main() {
	stdin := ""
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		stdin = stdin + "\n" + scanner.Text()
	}

	// fmt.Println(stdin)

	if scanner.Err() != nil {
		// Handle error.
	}
	app := &cli.App{
		Name:  "gof",
		Usage: "extract files from stdin",
		Action: func(cCtx *cli.Context) error {
			files := FindFiles(stdin)
			printFiles(files)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func FindFiles(input string) []Filematch {
	var result []Filematch
	files := make(map[string]Filematch)

	lines := strings.Split(string(input), "\n")

	for _, val := range lines {
		words := strings.Split(val, " ")
		for _, w := range words {

			filePath := extractFilePath(w)
			_, exist := files[filePath]
			if FileExist(filePath) && !exist {
				files[filePath] = Filematch{filePath: filePath}
			}
		}
	}

	for _, val := range files {
		result = append(result, val)
	}

	return result
}

func printFiles(files []Filematch) {
	for _, val := range files {
		fmt.Println(val.filePath)
	}
}

func _fileExist(filePath string) bool {
	if stat, err := os.Stat(filePath); err == nil {
		return !stat.IsDir()
	} else {
		return false
	}
}

func extractFilePath(_filePath string) string {
	filePath := _filePath
	if !strings.HasPrefix(filePath, "/") {
		pwd, _ := os.LookupEnv("PWD")
		filePath = pwd + "/" + filePath
	}
	r, _ := regexp.Compile("[0-9]+:[0-9]+")
	index := r.FindStringIndex(_filePath)

	if index != nil && index[0] != 0 {
		filePath = _filePath[:index[0]-1]
	}
	return filePath
}
