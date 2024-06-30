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
	"unicode/utf8"
)

// Input would contain statistic about a file
type Input struct {
	FileName       string
	Content        string
	ByteCount      int
	LineCount      int
	WordCount      int
	CharacterCount int
}

// wcEngine ...
type WcEngine struct {
	*wcflag.Options
	Input []Input
}

// InitEngine ...
func InitEngine(options *wcflag.Options) WcEngine {

	// We need to identify command line arguments that are not flags. The working assumption is that flags
	// will always start with a `-`. If a command line argument doesn't start with the `-` then it isn't a flag
	inputFiles := handleCommandLineInput(os.Args[1:])
	return WcEngine{
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
			text += scanner.Text() + "\n"
		}

		return []Input{
			{
				FileName: stdInF.Name(),
				Content:  strings.Trim(text, "\n"),
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
			Content:  strings.Trim(string(content), "\n "),
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
func (we *WcEngine) Count() {

	switch {
	case *we.Options.CountLines:
		countLine(we.Input)

	case *we.Options.CountBytes:
		countByte(we.Input)

	case *we.Options.CountWords:
		countWord(we.Input)

	case *we.Options.CountCharacters:
		countCharacter(we.Input)

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

func countLine(input []Input) {
	for i := range input {
		input[i].LineCount = len(strings.Split(input[i].Content, "\n"))
	}
}

func countWord(input []Input) {
	for i := range input {
		input[i].WordCount = len(strings.Split(input[i].Content, " "))
	}
}

func countCharacter(input []Input) {
	for i := range input {
		input[i].CharacterCount = utf8.RuneCount([]byte(input[i].Content))
	}
}

func (we *WcEngine) printResult() {

	totalByteCount := 0
	totalLineCount := 0

	for index := range we.Input {
		fmt.Printf("%v %v %v %v %s\n", we.Input[index].ByteCount, we.Input[index].LineCount, we.Input[index].WordCount, we.Input[index].CharacterCount, we.Input[index].FileName)
		totalByteCount += we.Input[index].ByteCount
		totalLineCount += we.Input[index].LineCount
	}

	if totalByteCount > we.Input[0].ByteCount {
		fmt.Printf("%v %v Total\n", totalByteCount, totalLineCount)
	}

}
