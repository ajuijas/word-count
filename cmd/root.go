/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)


var (
	lineFlag bool
	wordFlag bool
	charectorFlag bool
)


func readFile(fileName string) (int, int, int) {
	var (
	lines, words, charectors int
)
	data, _ := os.ReadFile(fileName)
	if len(data) != 0 {
		lines += 1
		words += 1
	}
	for _, i := range string(data){
		charectors += 1
		if i == '\n'{
			lines += 1
		}
		if i == ' '{
			words += 1
		}
	}
	return lines, words, charectors
}

func formatIntRightAligned(num int, width int) string {
	var tem string
	strNum := strconv.Itoa(num)
	padding := width - len(strNum)
	if padding > 0 {
		tem = strings.Repeat(" ", padding) + strNum + " "
		return string(tem)
	}
	return strNum + " "
}

func printResultLine(fileName string, wg *sync.WaitGroup) {
	lines, words, charectors := readFile(fileName)
	var result string
	if !lineFlag && !wordFlag && !charectorFlag {
		lineFlag, wordFlag, charectorFlag = true, true, true
	}
	if lineFlag {
		result += formatIntRightAligned(lines, 8)
	}
	if wordFlag {
		result += formatIntRightAligned(words, 8)
	}
	if charectorFlag {
		result += formatIntRightAligned(charectors, 8)
	}
	fmt.Print(result, fileName, "\n")
	defer wg.Done()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wc",
	Short: "",
	Long: ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("please provide a file name")
		}
		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			return fmt.Errorf("%v: open: No such file or directory", args[0])
		}
		if _, err := os.Stat(args[0]); os.IsPermission(err) {
			return fmt.Errorf("%v: open: Permission denied", args[0])
		}
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		for _, fileName := range args {
			wg.Add(1)
			go printResultLine(fileName,  &wg)
		}
		wg.Wait()
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


