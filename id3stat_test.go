/*
 * Copyright (C) 2016 Upper Stream Software.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetFileStatus(t *testing.T) {
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		if err := os.Mkdir("testdata", 0755); err != nil {
			t.Fatal(err)
		}
	}

	mp3File, err := ioutil.TempFile("testdata", "test_mp3*.mp3")
	if err != nil {
		t.Fatal(err)
	}
	mp3FilePath := mp3File.Name()
	mp3File.Close()
	defer os.Remove(mp3FilePath)

	txtFile, err := ioutil.TempFile("testdata", "test_txt*.txt")
	if err != nil {
		t.Fatal(err)
	}
	txtFilePath := txtFile.Name()
	txtFile.Close()
	defer os.Remove(txtFilePath)

	tests := []struct {
		name        string
		filePath    string
		expectError bool
	}{
		{
			name:        "MP3 file",
			filePath:    mp3FilePath,
			expectError: false,
		},
		{
			name:        "TXT file (unsupported)",
			filePath:    txtFilePath,
			expectError: true,
		},
		{
			name:        "Non-existent file",
			filePath:    "non_existent_file.mp3",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := getFileStatus(tt.filePath)
			if (err != nil) != tt.expectError {
				t.Errorf("getFileStatus() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateEncodingFlag(t *testing.T) {
	tests := []struct {
		name        string
		encoding    string
		expectError bool
	}{
		{
			name:        "UTF-8 encoding",
			encoding:    "UTF-8",
			expectError: false,
		},
		{
			name:        "Empty encoding (defaults to UTF-8)",
			encoding:    "",
			expectError: false,
		},
		{
			name:        "ShiftJIS encoding",
			encoding:    "ShiftJIS",
			expectError: false,
		},
		{
			name:        "Unsupported encoding",
			encoding:    "latin1",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEncodingFlag(tt.encoding)
			if (err != nil) != tt.expectError {
				t.Errorf("validateEncodingFlag() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestParseListFile(t *testing.T) {
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		if err := os.Mkdir("testdata", 0755); err != nil {
			t.Fatal(err)
		}
	}

	listFileContent := []byte("file1.mp3\nfile2.mp3\n\"file with spaces.mp3\"\n")
	listFile, err := ioutil.TempFile("testdata", "test_list*.txt")
	if err != nil {
		t.Fatal(err)
	}
	listFilePath := listFile.Name()
	if _, err := listFile.Write(listFileContent); err != nil {
		t.Fatal(err)
	}
	listFile.Close()
	defer os.Remove(listFilePath)

	tests := []struct {
		name           string
		listFilePath   string
		encoding       string
		expectedCount  int
		expectError    bool
	}{
		{
			name:           "Valid list file with UTF-8 encoding",
			listFilePath:   listFilePath,
			encoding:       "UTF-8",
			expectedCount:  3,
			expectError:    false,
		},
		{
			name:           "Valid list file with empty encoding (defaults to UTF-8)",
			listFilePath:   listFilePath,
			encoding:       "",
			expectedCount:  3,
			expectError:    false,
		},
		{
			name:           "Non-existent list file",
			listFilePath:   "non_existent_list.txt",
			encoding:       "UTF-8",
			expectedCount:  0,
			expectError:    true,
		},
		{
			name:           "Valid list file with unsupported encoding",
			listFilePath:   listFilePath,
			encoding:       "latin1",
			expectedCount:  0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := parseListFile(tt.listFilePath, tt.encoding)
			
			if (err != nil) != tt.expectError {
				t.Errorf("parseListFile() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			if !tt.expectError && len(files) != tt.expectedCount {
				t.Errorf("parseListFile() returned %d files, expected %d", len(files), tt.expectedCount)
			}
		})
	}
}

func TestListFilesIn(t *testing.T) {
	testDir, err := ioutil.TempDir("", "test_traverse")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	subDir := filepath.Join(testDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	mp3Files := []string{
		filepath.Join(testDir, "file1.mp3"),
		filepath.Join(testDir, "file2.mp3"),
		filepath.Join(subDir, "file3.mp3"),
	}

	for _, file := range mp3Files {
		if err := ioutil.WriteFile(file, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	nonMp3File := filepath.Join(testDir, "file.txt")
	if err := ioutil.WriteFile(nonMp3File, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	files, err := listFilesIn(testDir)
	if err != nil {
		t.Errorf("listFilesIn() error = %v", err)
		return
	}

	if len(files) != len(mp3Files) {
		t.Errorf("listFilesIn() returned %d files, expected %d", len(files), len(mp3Files))
	}

	for _, file := range files {
		if !strings.HasSuffix(strings.ToLower(file), ".mp3") {
			t.Errorf("listFilesIn() returned non-MP3 file: %s", file)
		}
	}

	_, err = listFilesIn("non_existent_directory")
	if err == nil {
		t.Errorf("listFilesIn() with non-existent directory did not return an error")
	}

	_, err = listFilesIn(mp3Files[0])
	if err == nil {
		t.Errorf("listFilesIn() with a file instead of a directory did not return an error")
	}
}

func TestTraverse(t *testing.T) {
	testDir, err := ioutil.TempDir("", "test_traverse")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	subDirs := []string{
		filepath.Join(testDir, "subdir1"),
		filepath.Join(testDir, "subdir2"),
		filepath.Join(testDir, "subdir1", "subsubdir"),
	}

	for _, dir := range subDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}

	mp3Files := []string{
		filepath.Join(testDir, "file1.mp3"),
		filepath.Join(subDirs[0], "file2.mp3"),
		filepath.Join(subDirs[1], "file3.mp3"),
		filepath.Join(subDirs[2], "file4.mp3"),
	}

	for _, file := range mp3Files {
		if err := ioutil.WriteFile(file, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	nonMp3Files := []string{
		filepath.Join(testDir, "file1.txt"),
		filepath.Join(subDirs[0], "file2.txt"),
	}

	for _, file := range nonMp3Files {
		if err := ioutil.WriteFile(file, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	dirs := []string{testDir}
	resultDirs, resultFiles, err := traverse(dirs, []string{})
	if err != nil {
		t.Errorf("traverse() error = %v", err)
		return
	}

	if len(resultDirs) != 0 {
		t.Errorf("traverse() returned %d dirs, expected 0 (all traversed)", len(resultDirs))
	}

	if len(resultFiles) != len(mp3Files) {
		t.Errorf("traverse() returned %d files, expected %d", len(resultFiles), len(mp3Files))
	}

	for _, file := range resultFiles {
		if !strings.HasSuffix(strings.ToLower(file), ".mp3") {
			t.Errorf("traverse() returned non-MP3 file: %s", file)
		}
	}
}

func TestFlagParsing(t *testing.T) {
	originalArgs := os.Args
	originalFlagCommandLine := flag.CommandLine
	
	defer func() {
		os.Args = originalArgs
		flag.CommandLine = originalFlagCommandLine
	}()

	tests := []struct {
		name         string
		args         []string
		expectPanic  bool
	}{
		{
			name:        "Valid args with MP3 file",
			args:        []string{"id3stat", "file.mp3"},
			expectPanic: false,
		},
		{
			name:        "Version flag",
			args:        []string{"id3stat", "-V"},
			expectPanic: true, // os.Exit in the function will cause panic in test
		},
		{
			name:        "License flag",
			args:        []string{"id3stat", "-L"},
			expectPanic: true, // os.Exit in the function will cause panic in test
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			versionFlag = flag.Bool("V", false, "Print the version number.")
			licenceFlag = flag.Bool("L", false, "Print the licencing notice.")
			filesFlag = flag.String("files", "", "Provides a list of files to process.")
			encodingFlag = flag.String("encoding", "UTF-8", "Encoding of a file that -files flag provides.")
			dirFlag = flag.String("dir", "", "Specifies the directory to test files in.")
			
			os.Args = tt.args
			
			if !tt.expectPanic {
				parseFlagsAndExit()
			}
		})
	}
}
