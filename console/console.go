// Copyright 2015 Matthew Collins
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

package console

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/thinkofdeath/steven/chat"
)

var (
	// stdout and a log file combind
	w io.Writer

	// For the possible option of scrolling in the
	// future
	historyBuffer [200]chat.AnyComponent
)

func init() {
	f, err := os.Create("steven-log.txt")
	if err != nil {
		panic(err)
	}
	w = io.MultiWriter(f, os.Stdout)
}

func Text(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}
	file = file[strings.LastIndexByte(file, '/')+1:]

	msg := &chat.TextComponent{
		Text: fmt.Sprintf("[%s:%d] ", file, line),
		Component: chat.Component{
			Color: chat.Aqua,
		},
	}
	msg.Extra = append(msg.Extra, chat.Wrap(&chat.TextComponent{
		Text: fmt.Sprintf(format, args...),
		Component: chat.Component{
			Color: chat.White,
		},
	}))
	Component(chat.Wrap(msg))
}

func Component(c chat.AnyComponent) {
	io.WriteString(w, c.String()+"\n")
	copy(historyBuffer[1:], historyBuffer[:])
	historyBuffer[0] = c
}

func History(lines int) []chat.AnyComponent {
	if lines == -1 {
		return historyBuffer[:]
	}
	return historyBuffer[:lines]
}