package logview

import (
	"encoding/json"
	"fmt"
)

func (a *LogViewModel) GetLines() error {
	a.Lines = make([]*LogLineEcho, 0)
	lines, err := GetLastLines(a.FileName, a.Start, a.PerPage)
	if err != nil {
		return fmt.Errorf("get lines error %w", err)
	}
	for _, line := range lines {
		lline := &LogLineEcho{}
		err := json.Unmarshal([]byte(line), lline)
		if err != nil {
			return fmt.Errorf("get lines error %w", err)
		}
		a.Lines = append(a.Lines, lline)
	}
	return nil
}
