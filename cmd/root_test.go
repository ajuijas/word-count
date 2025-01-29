package cmd

import (
	// "reflect"
	"bytes"
	"fmt"
	"os"
	"testing"
)

func Test_no_such_file(t *testing.T){
    rootCmd.SetArgs([]string{})
    err := rootCmd.Args(rootCmd, []string{"foo.txt"})

    if err == nil || err.Error() != "foo.txt: open: No such file or directory" {
		t.Errorf("Expected error 'foo.txt: open: No such file or directory', got: %v", err)
	}
}
func Test_no_argument(t *testing.T){
    rootCmd.SetArgs([]string{})
    err := rootCmd.Args(rootCmd, []string{})

    if err == nil || err.Error() != "please provide a file name" {
		t.Errorf("Expected error 'please provide a file name', got: %v", err)
	}
}

// func TestRootCommandPermissionDenied(t *testing.T){
//     rootCmd.SetArgs([]string{})
//     err := rootCmd.Args(rootCmd, []string{"locked.txt"})

//     if err == nil || err.Error() != "foo.txt: open: No such file or directory" {
// 		t.Errorf("Expected error 'foo.txt: open: No such file or directory', got: %v", err)
// 	}
// }


func Test_lines_count(t *testing.T){

    originalStdout := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	
    rootCmd.SetArgs([]string{})
    err := rootCmd.Args(rootCmd, []string{"test.txt", "-l"})
    rootCmd.Execute()
    fmt.Print(err)

	w.Close()

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	os.Stdout = originalStdout

	expected := "Hello, world!\n"
	if buf.String() != expected {
		t.Errorf("expected %q but got %q", expected, buf.String())
	}
}