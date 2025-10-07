package logview

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// получить count строк начиная с конца файла
func GetLastLines(filepath string, start, count int) ([]string, error) {
	fileHandle, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %w", err)
	}
	defer fileHandle.Close()
	allLine, err := lineCounter(fileHandle)
	if err != nil {
		return nil, fmt.Errorf("count lines in file error %w", err)
	}
	if start > allLine {
		return nil, fmt.Errorf("in file %d lines but %d greate this", allLine, start)
	}
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()
	out := make([]string, 0)
	if start > 0 {
		for {
			// cursor -= 1
			err := skipNext(fileHandle, filesize, &cursor)
			if err != nil {
				return nil, fmt.Errorf("skip next line error %w", err)
			}
			start -= 1
			if start <= 0 { // stop if count zero
				break
			}
			if cursor == -filesize { // stop if we are at the begining
				break
			}
		}
	}
	if cursor != -filesize {
		for {
			// cursor -= 1
			nextLine, err := lineNext(fileHandle, filesize, &cursor)
			if err != nil {
				return nil, fmt.Errorf("read next line error %w", err)
			}
			out = append(out, nextLine)
			count -= 1
			if count <= 0 { // stop if count zero
				break
			}
			if cursor == -filesize { // stop if we are at the begining
				break
			}
		}
	}
	return out, nil
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}
	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)
		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}

func lineNext(r *os.File, filesize int64, cursor *int64) (string, error) {
	line := ""
	// var cursor int64 = 0
	for {
		*cursor -= 1
		r.Seek(*cursor, io.SeekEnd)
		char := make([]byte, 1)
		r.Read(char)
		if *cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
			break
		}
		line = fmt.Sprintf("%s%s", string(char), line) // there is more efficient way
		if *cursor == -filesize {                      // stop if we are at the begining
			break
		}
	}
	return line, nil
}

func skipNext(r *os.File, filesize int64, cursor *int64) error {
	for {
		*cursor -= 1
		r.Seek(*cursor, io.SeekEnd)
		char := make([]byte, 1)
		r.Read(char)
		if *cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
			break
		}
		if *cursor == -filesize { // stop if we are at the begining
			break
		}
	}
	return nil
}
