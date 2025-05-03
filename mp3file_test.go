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
	"io/ioutil"
	"os"
	"testing"
)

func createTestMP3WithID3v1Tag(t *testing.T) string {
	tmpfile, err := ioutil.TempFile("testdata", "test_with_id3v1*.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpfile.Close()

	header := []byte("ID3")
	if _, err := tmpfile.Write(header); err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Seek(128, os.SEEK_SET); err != nil {
		t.Fatal(err)
	}

	tagData := []byte("TAG" + "Title                          " +
		"Artist                         " +
		"Album                          " +
		"2023" + "Comment                         " + "\x00")
	if _, err := tmpfile.Write(tagData); err != nil {
		t.Fatal(err)
	}

	return tmpfile.Name()
}

func createTestMP3WithoutID3v1Tag(t *testing.T) string {
	tmpfile, err := ioutil.TempFile("testdata", "test_without_id3v1*.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpfile.Close()

	data := []byte("This is a test MP3 file without an ID3v1 tag")
	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}

	return tmpfile.Name()
}

func createInvalidFile(t *testing.T) string {
	tmpfile, err := ioutil.TempFile("testdata", "test_invalid*.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpfile.Close()

	data := []byte("Invalid MP3 file")
	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}

	return tmpfile.Name()
}

func TestCheckMp3FileStatus(t *testing.T) {
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		if err := os.Mkdir("testdata", 0755); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name          string
		createTestFile func(*testing.T) string
		expectedResult bool
		expectError    bool
	}{
		{
			name:          "MP3 file with ID3v1 tag",
			createTestFile: createTestMP3WithID3v1Tag,
			expectedResult: true,
			expectError:    false,
		},
		{
			name:          "MP3 file without ID3v1 tag",
			createTestFile: createTestMP3WithoutID3v1Tag,
			expectedResult: false,
			expectError:    false,
		},
		{
			name:          "Invalid MP3 file",
			createTestFile: createInvalidFile,
			expectedResult: false,
			expectError:    false, // The function handles errors from tag.ReadID3v1Tags
		},
		{
			name:          "Non-existent file",
			createTestFile: func(*testing.T) string { return "non_existent_file.mp3" },
			expectedResult: false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.createTestFile(t)
			
			if tt.name != "Non-existent file" {
				defer os.Remove(filePath)
			}
			
			result, err := CheckMp3FileStatus(filePath)
			
			if (err != nil) != tt.expectError {
				t.Errorf("CheckMp3FileStatus() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			if result != tt.expectedResult {
				t.Errorf("CheckMp3FileStatus() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}
