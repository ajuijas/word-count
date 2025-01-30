/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var (
	lineFlag, wordFlag, charectorFlag bool
)

const (
	maxOpenFileCount = 2
	fileBufferSize = 1024 * 1024
)

type result struct {
	lineCount int
	wordCount int
	charCount int
	filename  string
	err       error
}

func countLine(line string, res *result) {

	res.lineCount++

	// Count words in the current line
	words := strings.Fields(line) // Fields splits the line by whitespace
	res.wordCount += len(words)

	// Use len(line) for total character count (including spaces, newlines, etc.)
	res.charCount += len(line)
}

func readFromInput(resultChan chan result){

	var res result
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan(){
		line := scanner.Text()
		countLine(line, &res)
	}
	fmt.Println()
	resultChan <- res
}

func readFromFile(fileName string, resultChan chan result, wg *sync.WaitGroup){

	var res result

	defer wg.Done()

	res.filename = fileName

	file, err := os.Open(fileName)
	if err != nil {
		res.err = err
		resultChan <- res
		return
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, fileBufferSize)

	for {
		// Read a line from the file
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				res.err = err
			}
			break // End of file
		}
		countLine(line, &res)
	}
	
	resultChan <- res

}

func printError(resu result) {
	if os.IsNotExist(resu.err) {
		fmt.Printf("./wc: '%s': open: No such file or directory\n", resu.filename)
	} else if os.IsPermission(resu.err) {
		fmt.Printf("./wc: '%s': open: Permission denied\n", resu.filename)
	} else if strings.Contains(resu.err.Error(), "is a directory") {
		fmt.Printf("./wc: '%s': read: Is a directory", resu.filename)
	} else {
		fmt.Printf("./wc: '%s': error: Unhandled error", resu.filename)
	}
	fmt.Println()
}

func printResult(resu result){
	if !lineFlag && !wordFlag &&!charectorFlag {
		// All flags are set to true, while no flags are present in the command
		lineFlag = true
		wordFlag = true
		charectorFlag = true
	}

	if resu.err != nil {
		printError(resu)
		return
	}
	if lineFlag{fmt.Printf("%8d ", resu.lineCount)}
	if wordFlag{fmt.Printf("%8d ", resu.wordCount)}
	if charectorFlag{fmt.Printf("%8d ", resu.charCount)}
	fmt.Println(resu.filename)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wc",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) { 
		resultChan := make(chan result, maxOpenFileCount)
		if len(args) == 0 {
			readFromInput(resultChan)
		}else{
			var wg sync.WaitGroup
			for _, fileName := range args {
				// Read the content from file
				wg.Add(1)
				go readFromFile(fileName, resultChan, &wg)
			}

			wg.Wait()
		}
		close(resultChan)
		
		for _res := range resultChan {
			printResult(_res)  // TODO: Try formating the result
		}
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolVarP(&lineFlag, "line", "l", false, "Enable the lines count")
	rootCmd.PersistentFlags().BoolVarP(&wordFlag, "word", "w", false, "Enable the words count")
	rootCmd.PersistentFlags().BoolVarP(&charectorFlag, "charector", "c", false, "Enable the charectors count")
}


