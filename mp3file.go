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
	"os"

	"github.com/dhowden/tag"
)

// CheckMp3FileStatus returns true if an MP3 file has an ID3 tag, false otherwise.
func CheckMp3FileStatus(pathname string) (bool, error) {
	f, err1 := os.Open(pathname)
	if err1 != nil {
		return false, err1
	}
	defer f.Close()
	_, err2 := tag.ReadID3v1Tags(f)
	if err2 != nil {
		return false, nil
	}
	return true, nil
}
