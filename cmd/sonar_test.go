package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

// TestCreateFileWithContent check write/read content to files
func TestStartSonar(t *testing.T) {
	t.Parallel()
	// create a cobra instance
	var sonarCmd = &cobra.Command{}

	var detectOk = []struct {
		name    string
		command []string
		args    []string
	}{
		{"check Install", []string{"install"}, []string{"install"}},
		{"check Status", []string{"status"}, []string{"status"}},
		{"check Run", []string{"run"}, []string{"run"}},
		{"check Run, status, stop", []string{"run", "status", "stop"}, []string{"run", "status", "stop"}},
		{"check Run, status, stop", []string{"test", "status", "stop"}, []string{"run", "status", "stop"}},
	}

	t.Run("Check arguments in command", func(t *testing.T) {
		for _, tt := range detectOk {
			// set the values from our data
			sonarCmd.SetArgs(tt.command)

			// verify if the flags provided are the same as the array expected
			if reflect.DeepEqual(tt.args, sonarCmd.Flags()) {
				t.Errorf("ERROR: \n name: %s, \n command: %s, \n args: %s, \n cmd: %t",
					tt.name,
					tt.command,
					tt.args[0],
					sonarCmd.Flags().Changed(tt.command[0]),
				)
			}
		}
	})
}

// TestCreateFileWithContent check write/read content to files
func TestCreateFileWithContent(t *testing.T) {
	var tests = []struct {
		path     string
		content  string
		expected []byte
	}{
		{"/tmp/test-1", "content-test-1", []byte("content-test-1")},
		{"/tmp/test-2", "content-test-2", []byte("content-test-2")},
		{"/tmp/test-3", "content-test-3", []byte("content-test-3")},
	}

	for _, tt := range tests {
		t.Run("Check content in file", func(t *testing.T) {
			expected := CreateFileWithContent(tt.path, tt.content)
			fmt.Println("expected", expected)

			g, err := ioutil.ReadFile(tt.path)
			if err != nil {
				t.Fatalf("failed reading: %s, %s", tt.path, err)
			}

			if !bytes.Equal(tt.expected, g) {
				t.Fatalf("failed reading: %s, %s", tt.path, err)
			}

			defer func() {
				err = os.Remove(tt.path)
				if err != nil {
					fmt.Println(err)
				}
			}()
		})
	}
}

// TestCommandExists test the func CommandExists which returns true/false depending on the command exists depending on
func TestCommandExists(t *testing.T) {
	// Create data for existent commands
	var exists = []struct {
		got  string
		want bool
	}{
		{"ls", true},
		{"cat", true},
		{"mv", true},
		{"cp", true},
	}

	// Create data for not existent commands
	var notExists = []struct {
		got  string
		want bool
	}{
		{"testing", true},
		{"abc", true},
		{"jui", true},
		{"mmmm", true},
	}

	// loop inside our data and validates the tests cases
	t.Run("Command exists", func(t *testing.T) {
		for _, tt := range exists {
			e := CommandExists(tt.got)
			if e != tt.want {
				t.Errorf("ERROR: exists: %t, want: %t", e, tt.want)
			}
		}
	})
	// loop inside our data and validates the tests cases
	t.Run("Command DOES NOT exists", func(t *testing.T) {
		for _, tt := range notExists {
			ne := CommandExists(tt.got)
			if ne == tt.want {
				t.Errorf("ERROR: notExists: %t, want: %t", ne, tt.want)
			}
		}
	})
}
