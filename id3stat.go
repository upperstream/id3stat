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
 * limitations under the License.package main
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const appVersion = "0.0.1"

var versionFlag = flag.Bool("V", false, "Print the version number.")
var licenceFlag = flag.Bool("L", false, "Print the licencing notice.")
var helpFlag = flag.Bool("H", false, "Print this message.")

type id3Error struct {
	Path string
	What string
}

//go:generate go run tools/files2go.go -o notice.go NOTICE.txt

func (e id3Error) Error() string {
	return fmt.Sprintf("%s: %s", e.What, e.Path)
}

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Println("Version:", appVersion)
		return
	}
	if *licenceFlag {
		printLicence()
		return
	}
	if *helpFlag {
		flag.Usage()
		os.Exit(2)
		return
	}
	if flag.NArg() == 0 {
		flag.PrintDefaults()
		return
	}

	var nError int
	_, nError = getFileStatuses(os.Args[len(os.Args)-flag.NArg() : len(os.Args)])
	if nError == 0 {
		os.Exit(0)
	} else {
		os.Exit(1)
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

func printLicence() {
	fmt.Println(NOTICE)
}
