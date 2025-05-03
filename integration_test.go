// +build integtest

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func createIntegTestFileWithID3v1Tag(t *testing.T, path string) {
	frameData := []byte{0xFF, 0xFB, 0x90, 0x44, 0x00} // Just some random bytes
	
	id3v1Tag := make([]byte, 128)
	copy(id3v1Tag[0:3], []byte("TAG"))                      // ID3v1 signature
	copy(id3v1Tag[3:33], []byte("Test Title               ")) // Title (30 bytes)
	copy(id3v1Tag[33:63], []byte("Test Artist              ")) // Artist (30 bytes)
	copy(id3v1Tag[63:93], []byte("Test Album               ")) // Album (30 bytes)
	copy(id3v1Tag[93:97], []byte("2020"))                    // Year (4 bytes)
	copy(id3v1Tag[97:125], []byte("Test Comment                   ")) // Comment (28 bytes)
	id3v1Tag[125] = 0                                        // Zero byte
	id3v1Tag[126] = 1                                        // Track number
	id3v1Tag[127] = 0                                        // Genre
	
	fileData := append(frameData, id3v1Tag...)
	
	if err := ioutil.WriteFile(path, fileData, 0644); err != nil {
		t.Fatalf("Failed to create mock MP3 file with ID3v1 tag: %v", err)
	}
}

func createIntegTestFileWithoutID3v1Tag(t *testing.T, path string) {
	frameData := []byte{0xFF, 0xFB, 0x90, 0x44, 0x00} // Just some random bytes
	
	if err := ioutil.WriteFile(path, frameData, 0644); err != nil {
		t.Fatalf("Failed to create mock MP3 file without ID3v1 tag: %v", err)
	}
}

// TestDirectCommandLineArguments tests the application with direct file paths as command-line arguments
func TestDirectCommandLineArguments(t *testing.T) {
	// Create a temporary test directory
	testDir := "testdata"
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		if err := os.Mkdir(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}
	defer os.RemoveAll(testDir)

	// Create test files
	withTagPath := filepath.Join(testDir, "with_id3v1.mp3")
	createIntegTestFileWithID3v1Tag(t, withTagPath)
	defer os.Remove(withTagPath)

	withoutTagPath := filepath.Join(testDir, "without_id3v1.mp3")
	createIntegTestFileWithoutID3v1Tag(t, withoutTagPath)
	defer os.Remove(withoutTagPath)

	// Build the application
	cmd := exec.Command("go", "build", "-o", "id3stat")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build the application: %v", err)
	}
	defer os.Remove("id3stat")

	// Test with a file that has ID3v1 tag
	cmd = exec.Command("./id3stat", withTagPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v, output: %s", err, output)
	}
	if strings.Contains(string(output), withTagPath) {
		t.Errorf("File with ID3v1 tag was reported as not having a tag: %s", output)
	}

	// Test with a file that doesn't have ID3v1 tag
	cmd = exec.Command("./id3stat", withoutTagPath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v, output: %s", err, output)
	}
	if !strings.Contains(string(output), withoutTagPath) {
		t.Errorf("File without ID3v1 tag was not reported: %s", output)
	}

	// Test with multiple files
	cmd = exec.Command("./id3stat", withTagPath, withoutTagPath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v, output: %s", err, output)
	}
	if strings.Contains(string(output), withTagPath) {
		t.Errorf("File with ID3v1 tag was reported as not having a tag: %s", output)
	}
	if !strings.Contains(string(output), withoutTagPath) {
		t.Errorf("File without ID3v1 tag was not reported: %s", output)
	}
}

// TestListFileInput tests the application with a list file input
func TestListFileInput(t *testing.T) {
	// Create a temporary test directory
	testDir := "testdata"
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		if err := os.Mkdir(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}
	defer os.RemoveAll(testDir)

	// Create test files
	withTagPath := filepath.Join(testDir, "with_id3v1.mp3")
	createIntegTestFileWithID3v1Tag(t, withTagPath)
	defer os.Remove(withTagPath)

	withoutTagPath := filepath.Join(testDir, "without_id3v1.mp3")
	createIntegTestFileWithoutID3v1Tag(t, withoutTagPath)
	defer os.Remove(withoutTagPath)

	// Create list file
	listFilePath := filepath.Join(testDir, "files.txt")
	listContent := withTagPath + "\n" + withoutTagPath + "\n"
	if err := ioutil.WriteFile(listFilePath, []byte(listContent), 0644); err != nil {
		t.Fatalf("Failed to create list file: %v", err)
	}
	defer os.Remove(listFilePath)

	// Build the application
	cmd := exec.Command("go", "build", "-o", "id3stat")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build the application: %v", err)
	}
	defer os.Remove("id3stat")

	// Test with list file
	cmd = exec.Command("./id3stat", "--files="+listFilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v, output: %s", err, output)
	}
	
	if strings.Contains(string(output), withTagPath) {
		t.Errorf("File with ID3v1 tag was reported as not having a tag: %s", output)
	}
	if !strings.Contains(string(output), withoutTagPath) {
		t.Errorf("File without ID3v1 tag was not reported: %s", output)
	}
}

// TestDirectoryInput tests the application with a directory input
func TestDirectoryInput(t *testing.T) {
	// Create a temporary test directory
	testDir := "testdata"
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		if err := os.Mkdir(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}
	defer os.RemoveAll(testDir)

	// Create subdirectory
	subDir := filepath.Join(testDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create test files in main directory and subdirectory
	withTagPath := filepath.Join(testDir, "with_id3v1.mp3")
	createIntegTestFileWithID3v1Tag(t, withTagPath)

	withoutTagPath := filepath.Join(testDir, "without_id3v1.mp3")
	createIntegTestFileWithoutID3v1Tag(t, withoutTagPath)

	withTagSubPath := filepath.Join(subDir, "with_id3v1_sub.mp3")
	createIntegTestFileWithID3v1Tag(t, withTagSubPath)

	withoutTagSubPath := filepath.Join(subDir, "without_id3v1_sub.mp3")
	createIntegTestFileWithoutID3v1Tag(t, withoutTagSubPath)

	// Build the application
	cmd := exec.Command("go", "build", "-o", "id3stat")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build the application: %v", err)
	}
	defer os.Remove("id3stat")

	// Test with directory
	cmd = exec.Command("./id3stat", "--dir="+testDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v, output: %s", err, output)
	}
	
	if strings.Contains(string(output), withTagPath) {
		t.Errorf("File with ID3v1 tag was reported as not having a tag: %s", output)
	}
	if !strings.Contains(string(output), withoutTagPath) {
		t.Errorf("File without ID3v1 tag was not reported: %s", output)
	}
	if strings.Contains(string(output), withTagSubPath) {
		t.Errorf("File with ID3v1 tag in subdirectory was reported as not having a tag: %s", output)
	}
	if !strings.Contains(string(output), withoutTagSubPath) {
		t.Errorf("File without ID3v1 tag in subdirectory was not reported: %s", output)
	}
}
