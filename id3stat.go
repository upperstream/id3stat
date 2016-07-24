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
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const appVersion = "0.0.1"

var versionFlag = flag.Bool("V", false, "Print the version number.")
var licenceFlag = flag.Bool("L", false, "Print the licencing notice.")
var filesFlag = flag.String("files", "", "Provides a list of files to process.")
var encodingFlag = flag.String("encoding", "UTF-8",
	"Encoding of a file that -files flag provides.")
var dirFlag = flag.String("dir", "", "Specifies the directory to test files in.")

type id3Error struct {
	Path string
	What string
}

//go:generate go run tools/files2go.go -o notice.go NOTICE.txt

func (e id3Error) Error() string {
	return fmt.Sprintf("%s: %s", e.What, e.Path)
}

func main() {
	parseFlagsAndExit()

	var files []string
	var err error
	if len(*dirFlag) > 0 {
		files, err = listFilesIn(*dirFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else if len(*filesFlag) > 0 {
		files, err = parseListFile(*filesFlag, *encodingFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else {
		files = os.Args[len(os.Args)-flag.NArg() : len(os.Args)]
	}
	nSuccess, _ := getFileStatuses(files)
	if nSuccess == 0 {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func listFilesIn(dirname string) ([]string, error) {
	stat1, err1 := os.Stat(dirname)
	if err1 != nil {
		return nil, err1
	}
	if !stat1.IsDir() {
		return nil, fmt.Errorf("Not a directory: %s", dirname)
	}
	_, files, err2 := traverse(append(make([]string, 0, 128), *dirFlag), make([]string, 0, 128))
	if err2 != nil {
		return nil, err2
	}
	return files, nil
}

func traverse(directories []string, acc []string) (dirs []string, files []string, err error) {
	var dirs2 []string
	var files2 []string
	dirs = make([]string, 0, 128)
	files = acc
	for _, dir := range directories {
		dirs2, files2, _ = readdir(dir)
		dirs = append(dirs, dirs2...)
		files = append(files, files2...)
	}
	if len(dirs) > 0 {
		return traverse(dirs, files)
	}
	return dirs, files, nil
}

func readdir(dirname string) (dirs []string, files []string, err error) {
	var d *os.File
	d, err = os.Open(dirname)
	if err != nil {
		return nil, nil, err
	}
	defer func(name string) {
		d.Close()
	}(dirname)
	dirs = make([]string, 0, 128)
	files = make([]string, 0, 128)
	names, err2 := d.Readdirnames(0)
	if err2 != nil {
		return nil, nil, err2
	}
	for _, name := range names {
		var stat os.FileInfo
		var err3 error
		path := filepath.Join(dirname, name)
		stat, err3 = os.Stat(path)
		if err3 != nil {
			return dirs, files, err3
		}
		if stat.IsDir() {
			dirs = append(dirs, path)
		} else {
			if strings.EqualFold(filepath.Ext(path), ".mp3") == true {
				files = append(files, path)
			}
		}
	}
	return dirs, files, nil
}

func parseFlagsAndExit() {
	flag.Usage = printUsage
	flag.Parse()

	if *versionFlag {
		fmt.Println("Version:", appVersion)
		os.Exit(0)
	}

	if *licenceFlag {
		printLicence()
		os.Exit(0)
	}

	if err := validateEncodingFlag(*encodingFlag); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	if len(*filesFlag) > 0 && len(*dirFlag) > 0 {
		fmt.Fprintf(os.Stderr, "You cannot specify --files and --dir at the same time\n\n")
		printUsage()
		os.Exit(2)
	}

	if len(*dirFlag) > 0 && flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr,
			"You cannot specify command line arguments when you provide --dir option.\n\n")
		printUsage()
		os.Exit(2)
	}

	if len(*filesFlag) > 0 && flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr,
			"You cannot specify command line arguments when you provide --files option.\n\n")
		printUsage()
		os.Exit(2)
	}

	if len(*dirFlag) == 0 && len(*filesFlag) == 0 && flag.NArg() == 0 {
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	executable := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", executable)
	fmt.Fprintln(os.Stderr, executable, "mp3file [...]")
	fmt.Fprintln(os.Stderr, executable, "--files=<list> --encoding=<encoding>")
	fmt.Fprintln(os.Stderr, executable, "-H | -L | -V")
	flag.PrintDefaults()
}

func validateEncodingFlag(encoding string) error {
	switch encoding {
	case "ShiftJIS":
		return nil
	case "", "UTF-8":
		return nil
	default:
		return fmt.Errorf("Unsupported encoding: %s", encoding)
	}
}

func getFileStatuses(pathnames []string) (nSuccess int, nError int) {
	nSuccess = 0
	nError = 0
	for _, pathname := range pathnames {
		if err := getFileStatus(pathname); err == nil {
			nSuccess++
		} else {
			fmt.Fprintln(os.Stderr, err.Error())
			nError++
		}
	}
	return nSuccess, nError
}

func getFileStatus(pathname string) error {
	switch strings.ToLower(filepath.Ext(pathname)) {
	case ".mp3":
		result, err := CheckMp3FileStatus(pathname)
		if err != nil {
			return err
		}
		if !result {
			fmt.Println(pathname)
		}
	default:
		return id3Error{
			pathname,
			"Unsupported file type",
		}
	}
	return nil
}

func newReader(reader io.Reader, encoding string) (io.Reader, error) {
	switch encoding {
	case "ShiftJIS":
		return transform.NewReader(reader, japanese.ShiftJIS.NewDecoder()), nil
	case "", "UTF-8":
		return reader, nil
	default:
		return nil, fmt.Errorf("Unsupported encoding: %s", encoding)
	}
}

func parseListFile(listfile string, encoding string) (files []string, err error) {
	f, _ := os.Open(listfile)
	defer f.Close()
	var reader io.Reader
	reader, err = newReader(f, encoding)
	if err != nil {
		return nil, err
	}
	files = make([]string, 0, 256)
	s := bufio.NewScanner(reader)
	for s.Scan() {
		quoted := s.Text()
		if filename, err1 := strconv.Unquote(quoted); err1 == nil {
			files = append(files, filename)
		} else {
			files = append(files, quoted)
		}
	}
	return files, nil
}

func printLicence() {
	fmt.Println(NOTICE)
}
