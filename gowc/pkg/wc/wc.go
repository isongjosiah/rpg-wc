package wc

import (
	"bufio"
	"errors"
	"fmt"
	"gowc/pkg/wcflag"
	"io/fs"
	"log"
	"os"
	"strings"
)

// Input would contain statistic about a file
type Input struct {
	FileName  string
	Content   string
	ByteCount int
}

// wcEngine ...
type wcEngine struct {
	*wcflag.Options
	Input []Input
}

// InitEngine ...
func InitEngine(options *wcflag.Options) wcEngine {

	// We need to identify command line arguments that are not flags. The working assumption is that flags
	// will always start with a `-`. If a command line argument doesn't start with the `-` then it isn't a flag
	inputFiles := handleCommandLineInput(os.Args[1:])
	return wcEngine{
		Options: options,
		Input:   inputFiles,
	}

}

// todo: explain tomorrow
func handleCommandLineInput(cmdArg []string) []Input {

	// If we are in pipe mode. Handle text from standard input
	stdInF, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error reading stdIn " + err.Error())
	}

	if stdInF.Mode()&os.ModeCharDevice == 0 { // if true, we are in pipe mode

		scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
		text := ""
		for scanner.Scan() {
			text += scanner.Text()
		}

		return []Input{
			{
				FileName: stdInF.Name(),
				Content:  text,
			},
		}

	}

	inputList := make([]Input, 0)

	// loop through from the end and store everything that isn't a flag in inputList
	for i := len(cmdArg) - 1; i >= 0; i-- {

		arg := cmdArg[i]
		if strings.HasPrefix(arg, "-") {
			break
		}

		// check if it is a file
		f, err := os.Stat(arg)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				log.Fatal("Fatal Error: file " + arg + " does not exist.")
			}
		}

		// if a directory is provided error out
		if f.IsDir() {
			log.Fatal("Fatal Error: " + arg + " is a directory.")
		}

		// get the content of the file
		content, err := os.ReadFile(arg)
		if err != nil {
			log.Fatal("Fatal Error: file " + arg + ".")
		}

		inputList = append(inputList, Input{
			FileName: arg,
			Content:  string(content),
		})
	}

	// If there is no input provided stop the program
	if len(inputList) == 0 {
		log.Fatalf("Fatal Error: no input files specified!")
	}

	return inputList
}

// Count is the main function of wcEngine. It does a count based
// on the option defined and returns the output to the command line
func (we *wcEngine) Count() {

	switch {

	case *we.Options.CountBytes:
		countByte(we.Input)

	default:
		log.Fatal("Option is not supported!")
	}

	we.printResult()

}

// countByte counts the number byte in a given string
func countByte(input []Input) {
	for i := range input {
		input[i].ByteCount = len(input[i].Content)
	}
}

func countLine(text string) {

}

func countWord(text string) {

}

func (we *wcEngine) printResult() {

	totalByteCount := 0

	for index := range we.Input {
		fmt.Printf("%v %s\n", we.Input[index].ByteCount, we.Input[index].FileName)
		totalByteCount += we.Input[index].ByteCount
	}

	if totalByteCount > we.Input[0].ByteCount {
		fmt.Printf("%v Total\n", totalByteCount)
	}

}
