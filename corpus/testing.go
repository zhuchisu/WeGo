// Copyright © 2017 Makoto Ito
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package corpus

import (
	"bytes"
	"io"
	"io/ioutil"
)

type fakeNopSeeker struct{ io.ReadCloser }

func (fake fakeNopSeeker) Seek(offset int64, whence int) (int64, error) { return 0, nil }

var (
	text       = "a b b c c c c"
	fakeSeeker = fakeNopSeeker{ReadCloser: ioutil.NopCloser(bytes.NewReader([]byte(text)))}
	// TestWord2vecCorpus is mock for test.
	TestWord2vecCorpus, _ = NewWord2vecCorpus(fakeSeeker, true, 0)
)
