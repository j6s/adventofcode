package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func getInput(file string) (string, error) {
	inputFile := path.Join(path.Dir(file), "input.txt")
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return "", nil
	}

	contents, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func connectInputFileToStdin(cmd *exec.Cmd, file string) {
	input, err := getInput(file)
	if err != nil {
		log.Fatalf("error running %s: %v", file, err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("error running %s: %v", file, err)
	}

	io.WriteString(stdin, input)
	stdin.Close()

}

func run(file string) string {
	cmd := exec.Command("go", "run", file)
	connectInputFileToStdin(cmd, file)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("error running %s: %v\n%s", file, err, out)
	}

	return string(out)
}

func main() {
	files, err := filepath.Glob("./**/*/main.go")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Printf("%s: %s\n", file, run(file))
	}
}
