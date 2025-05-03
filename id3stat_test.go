package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func createTestFileWithID3v1Tag(t *testing.T, path string) {
	
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

func createTestFileWithoutID3v1Tag(t *testing.T, path string) {
	frameData := []byte{0xFF, 0xFB, 0x90, 0x44, 0x00} // Just some random bytes
	
	if err := ioutil.WriteFile(path, frameData, 0644); err != nil {
		t.Fatalf("Failed to create mock MP3 file without ID3v1 tag: %v", err)
	}
}

func TestGetFileStatus(t *testing.T) {
	testDir := "testdata"
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		if err := os.Mkdir(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	withTagPath := filepath.Join(testDir, "with_id3v1.mp3")
	createTestFileWithID3v1Tag(t, withTagPath)
	defer os.Remove(withTagPath)

	withoutTagPath := filepath.Join(testDir, "without_id3v1.mp3")
	createTestFileWithoutID3v1Tag(t, withoutTagPath)
	defer os.Remove(withoutTagPath)

	noPermissionPath := filepath.Join(testDir, "no_permission.mp3")
	createTestFileWithoutID3v1Tag(t, noPermissionPath)
	if err := os.Chmod(noPermissionPath, 0000); err != nil {
		t.Fatalf("Failed to change file permissions: %v", err)
	}
	defer func() {
		os.Chmod(noPermissionPath, 0644)
		os.Remove(noPermissionPath)
	}()

	tests := []struct {
		name     string
		path     string
		wantErr  bool
		errType  string
	}{
		{
			name:     "MP3 file with ID3v1 tag",
			path:     filepath.Join("testdata", "with_id3v1.mp3"),
			wantErr:  false,
			errType:  "",
		},
		{
			name:     "MP3 file without ID3v1 tag",
			path:     filepath.Join("testdata", "without_id3v1.mp3"),
			wantErr:  false, // No error, but output will be printed
			errType:  "",
		},
		{
			name:     "Non-existent file",
			path:     filepath.Join("testdata", "non_existent.mp3"),
			wantErr:  true,
			errType:  "os.PathError",
		},
		{
			name:     "File with no read permissions",
			path:     filepath.Join("testdata", "no_permission.mp3"),
			wantErr:  true,
			errType:  "os.PathError",
		},
		{
			name:     "Unsupported file type",
			path:     "test.txt",
			wantErr:  true,
			errType:  "id3Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := getFileStatus(tt.path)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errType != "" {
				switch tt.errType {
				case "id3Error":
					if _, ok := err.(id3Error); !ok {
						t.Errorf("getFileStatus() error type = %T, want %s", err, tt.errType)
					}
				case "os.PathError":
					if _, ok := err.(*os.PathError); !ok {
						t.Errorf("getFileStatus() error type = %T, want %s", err, tt.errType)
					}
				}
			}
		})
	}
}

func TestGetFileStatuses(t *testing.T) {
	testDir := "testdata"
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		if err := os.Mkdir(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	withTagPath := filepath.Join(testDir, "with_id3v1.mp3")
	createTestFileWithID3v1Tag(t, withTagPath)
	defer os.Remove(withTagPath)

	withoutTagPath := filepath.Join(testDir, "without_id3v1.mp3")
	createTestFileWithoutID3v1Tag(t, withoutTagPath)
	defer os.Remove(withoutTagPath)

	tests := []struct {
		name      string
		pathnames []string
		wantSuccess int
		wantError  int
	}{
		{
			name:      "All files exist",
			pathnames: []string{
				withTagPath,
				withoutTagPath,
			},
			wantSuccess: 2,
			wantError:   0,
		},
		{
			name:      "Some files don't exist",
			pathnames: []string{
				withTagPath,
				filepath.Join(testDir, "non_existent.mp3"),
			},
			wantSuccess: 1,
			wantError:   1,
		},
		{
			name:      "Unsupported file types",
			pathnames: []string{
				"test.txt",
				"test.png",
			},
			wantSuccess: 0,
			wantError:   2,
		},
		{
			name:      "Mixed file types and existence",
			pathnames: []string{
				withTagPath,
				"test.txt",
				filepath.Join(testDir, "non_existent.mp3"),
			},
			wantSuccess: 1,
			wantError:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSuccess, gotError := getFileStatuses(tt.pathnames)
			
			if gotSuccess != tt.wantSuccess {
				t.Errorf("getFileStatuses() success count = %v, want %v", gotSuccess, tt.wantSuccess)
			}
			
			if gotError != tt.wantError {
				t.Errorf("getFileStatuses() error count = %v, want %v", gotError, tt.wantError)
			}
		})
	}
}
