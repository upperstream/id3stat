package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckMp3FileStatus(t *testing.T) {
	testDir := "testdata"
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		if err := os.Mkdir(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	withTagPath := filepath.Join(testDir, "with_id3v1.mp3")
	createMockMp3FileWithID3v1Tag(t, withTagPath)
	defer os.Remove(withTagPath)

	withoutTagPath := filepath.Join(testDir, "without_id3v1.mp3")
	createMockMp3FileWithoutID3v1Tag(t, withoutTagPath)
	defer os.Remove(withoutTagPath)

	nonExistentPath := filepath.Join(testDir, "non_existent.mp3")

	noPermissionPath := filepath.Join(testDir, "no_permission.mp3")
	createMockMp3FileWithoutID3v1Tag(t, noPermissionPath)
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
		expected bool
		wantErr  bool
	}{
		{
			name:     "File with ID3v1 tag",
			path:     withTagPath,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "File without ID3v1 tag",
			path:     withoutTagPath,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "Non-existent file",
			path:     nonExistentPath,
			expected: false,
			wantErr:  true,
		},
		{
			name:     "File with no read permissions",
			path:     noPermissionPath,
			expected: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckMp3FileStatus(tt.path)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckMp3FileStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if got != tt.expected {
				t.Errorf("CheckMp3FileStatus() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func createMockMp3FileWithID3v1Tag(t *testing.T, path string) {
	
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

func createMockMp3FileWithoutID3v1Tag(t *testing.T, path string) {
	frameData := []byte{0xFF, 0xFB, 0x90, 0x44, 0x00} // Just some random bytes
	
	if err := ioutil.WriteFile(path, frameData, 0644); err != nil {
		t.Fatalf("Failed to create mock MP3 file without ID3v1 tag: %v", err)
	}
}
