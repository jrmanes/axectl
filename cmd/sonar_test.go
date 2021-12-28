package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

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

func TestGetTokenInFile(t *testing.T) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Println("user home dir not found...")
	}
	var tests = []struct {
		filename string
		path     string
		content  string
		expected []byte
	}{
		{"test-1", dirname + "/test-1", "content-test-1", []byte("content-test-1")},
		{"test-2", dirname + "/test-2", "content-test-2", []byte("content-test-2")},
		{"test-3", dirname + "/test-3", "content-test-3", []byte("content-test-3")},
	}

	for _, tt := range tests {
		t.Run("Check content in file", func(t *testing.T) {

			expected, err := GetTokenInFile(tt.path)
			if err != nil {
				t.Fatalf("failed reading: %s, %s", tt.path, err)
			}

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
