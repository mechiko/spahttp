package logview

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testsLastLines = []struct {
	name  string
	err   bool
	file  string
	lines int
	start int
}{
	// the table itself
	{"test 0", true, "1236", 0, 0},
	{"test 1", false, "../../../cmd/.spahttp/echo", 10, 0},
	{"test 2", false, "../../../cmd/.spahttp/echo", 10, 1100},
}

func TestGetLastLines(t *testing.T) {
	wd, _ := os.Getwd()
	fmt.Printf("wd=%s\n", wd)
	for _, tt := range testsLastLines {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			lines, err := GetLastLines(tt.file, tt.start, tt.lines)
			if tt.err {
				assert.NotNil(t, err, "ожидаем ошибку")
			} else {
				// ожидаем отсутствие ошибки
				assert.NoError(t, err)
				assert.Equal(t, len(lines), tt.lines, "ожидаемое значение")
			}
		})
	}
}

func TestGetLines(t *testing.T) {
	wd, _ := os.Getwd()
	fmt.Printf("wd=%s\n", wd)
	for _, tt := range testsLastLines {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			lvm := &LogViewModel{
				Start:    tt.start,
				PerPage:  tt.lines,
				FileName: tt.file,
			}
			err := lvm.GetLines()
			if tt.err {
				assert.NotNil(t, err, "ожидаем ошибку")
			} else {
				// ожидаем отсутствие ошибки
				assert.NoError(t, err)
				assert.Equal(t, len(lvm.Lines), tt.lines, "ожидаемое значение")
			}
		})
	}
}
