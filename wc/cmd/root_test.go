package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func resetFlags() {
    lineFlag = false
    wordFlag = false
    charectorFlag = false
}

func Test_count(t *testing.T) {
	var line string = "12 3456 789012"
	var res result = result{filename: "test", lineCount: 10, wordCount: 10, charCount: 10}
	var expected result = result{filename: "test", lineCount: 11, wordCount: 13, charCount: 22}

	count(line, &res)

	switch{
	case res.filename != expected.filename,
	 res.lineCount != expected.lineCount,
	 res.wordCount != expected.wordCount,
	 res.charCount != expected.charCount:	
		t.Errorf("Expected %+v, got %+v", expected, res) 	
	}
	
}

func Test_readFromInput(t *testing.T) {

	var input string = "1234 5678 \n901234 56789\n0"
	reader := strings.NewReader(input)

	var expected result = result{lineCount: 3, wordCount: 5, charCount: 20}

	resultChan := make(chan result, 1)

	readFromInput(reader, resultChan)
	close(resultChan)

	res := <-resultChan

	switch{
	case res.filename != expected.filename,
	 res.lineCount != expected.lineCount,
	 res.wordCount != expected.wordCount,
	 res.charCount != expected.charCount:	
		t.Errorf("Expected %+v, got %+v", expected, res) 	
	}
}

func Test_readFromFile(t *testing.T) {

	tmpFile, err := os.CreateTemp("", "testfile.txt")
	if err!=nil {
		t.Fatal("1: Failed to create a temp file for the test")
	}
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.WriteString("1234 5678 \n901234 56789\n089"); err != nil{
		t.Fatal("2: Failed to create a temp file for the test")
	}

	var expected result = result{lineCount: 3, wordCount: 5, charCount: 22}

	resultChan := make(chan result, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	readFromFile(tmpFile.Name(), resultChan, &wg)

	close(resultChan)

	res := <- resultChan

		switch{
	case res.filename != tmpFile.Name(),
	 res.lineCount != expected.lineCount,
	 res.wordCount != expected.wordCount,
	 res.charCount != expected.charCount:	
		t.Errorf("Expected %+v, got %+v", expected, res) 	
	}
}

func Test_RootCmd_FileInput(t *testing.T){

	defer resetFlags()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() { os.Stdout = oldStdout }()

	tmpFile, err := os.CreateTemp("", "testfile.txt")
	if err!=nil {
		t.Fatal("1: Failed to create a temp file for the test")
	}
	defer os.Remove(tmpFile.Name())


	var expected string = "       3        5       22 " + tmpFile.Name()

	if _, err = tmpFile.WriteString("1234 5678 \n901234 56789\n089"); err != nil{
		t.Fatal("2: Failed to create a temp file for the test")
	}

	rootCmd.SetArgs([]string{tmpFile.Name()})
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("rootCmd.Execute() returned an error: %v", err)
	}

	w.Close()

	var out bytes.Buffer
	io.Copy(&out, r)

	expected = strings.TrimSpace(expected)
	outString := strings.TrimSpace(out.String())
	
	if outString != expected{
		t.Errorf("Expected %+v, got %+v", expected, outString)
	}
	
}

func Test_RootCmd_FileInput_LineCountFlag(t *testing.T){

	defer resetFlags()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() { os.Stdout = oldStdout }()

	tmpFile, err := os.CreateTemp("", "testfile.txt")
	if err!=nil {
		t.Fatal("1: Failed to create a temp file for the test")
	}
	defer os.Remove(tmpFile.Name())


	var expected string = "       3 " + tmpFile.Name()

	if _, err = tmpFile.WriteString("1234 5678 \n901234 56789\n089"); err != nil{
		t.Fatal("2: Failed to create a temp file for the test")
	}

	rootCmd.SetArgs([]string{"-l", tmpFile.Name()})
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("rootCmd.Execute() returned an error: %v", err)
	}

	w.Close()

	var out bytes.Buffer
	io.Copy(&out, r)

	expected = strings.TrimSpace(expected)
	outString := strings.TrimSpace(out.String())
	
	if outString != expected{
		t.Errorf("Expected %+v, got %+v", expected, outString)
	}
	
}

func Test_RootCmd_FileInput_WordCountFlag(t *testing.T){

	defer resetFlags()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() { os.Stdout = oldStdout }()

	tmpFile, err := os.CreateTemp("", "testfile.txt")
	if err!=nil {
		t.Fatal("1: Failed to create a temp file for the test")
	}
	defer os.Remove(tmpFile.Name())


	var expected string = "       5 " + tmpFile.Name()

	if _, err = tmpFile.WriteString("1234 5678 \n901234 56789\n089"); err != nil{
		t.Fatal("2: Failed to create a temp file for the test")
	}

	rootCmd.SetArgs([]string{"-w", tmpFile.Name()})
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("rootCmd.Execute() returned an error: %v", err)
	}

	w.Close()

	var out bytes.Buffer
	io.Copy(&out, r)

	expected = strings.TrimSpace(expected)
	outString := strings.TrimSpace(out.String())
	
	if outString != expected{
		t.Errorf("Expected %+v, got %+v", expected, outString)
	}
	
}

func Test_RootCmd_FileInput_CharectorCountFlag(t *testing.T){

	defer resetFlags()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() { os.Stdout = oldStdout }()

	tmpFile, err := os.CreateTemp("", "testfile.txt")
	if err!=nil {
		t.Fatal("1: Failed to create a temp file for the test")
	}
	defer os.Remove(tmpFile.Name())


	var expected string = "       22 " + tmpFile.Name()

	if _, err = tmpFile.WriteString("1234 5678 \n901234 56789\n089"); err != nil{
		t.Fatal("2: Failed to create a temp file for the test")
	}

	rootCmd.SetArgs([]string{"-c", tmpFile.Name()})
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("rootCmd.Execute() returned an error: %v", err)
	}

	w.Close()

	var out bytes.Buffer
	io.Copy(&out, r)

	expected = strings.TrimSpace(expected)
	outString := strings.TrimSpace(out.String())
	
	if outString != expected{
		t.Errorf("Expected %+v, got %+v", expected, outString)
	}
	
}
