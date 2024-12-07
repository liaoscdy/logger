package logger

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileExecutor struct {
	fileAbsPath  string
	fileWriter   *os.File
	fileOpenTime time.Time
	// rotate file by day, eg: xxx.log.2006-01-02
	rotateMaxDay    int
	rotateCleanWg   sync.WaitGroup
	isRotateLogFile bool
}

func NewFileExecutor(fileAbsPath string) *FileExecutor {
	return &FileExecutor{fileAbsPath: fileAbsPath}
}

func (f *FileExecutor) openLogFile() error {
	fd, err := os.OpenFile(f.fileAbsPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0660))
	if err != nil {
		return err
	}

	f.fileWriter = fd
	f.fileOpenTime = time.Now()
	return nil
}

func (f *FileExecutor) isNeedRotate() bool {
	if !f.isRotateLogFile || f.rotateMaxDay <= 0 || f.fileWriter == nil {
		return false
	}

	return time.Now().Day() != f.fileOpenTime.Day()
}

func (f *FileExecutor) rotateLogFile() error {
	if err := f.fileWriter.Close(); err != nil {
		return err
	}

	// save xxx.log to xxx.log.2006-01-02
	prevOpenTime := f.fileOpenTime
	saveFileName := fmt.Sprintf("%s.%s", f.fileAbsPath, prevOpenTime.Format(time.DateOnly))
	if err := os.Rename(f.fileAbsPath, saveFileName); err != nil {
		return err
	}

	// open new logFile
	if err := f.openLogFile(); err != nil {
		return err
	}

	if f.rotateMaxDay > 0 {
		f.rotateCleanWg.Add(1)
		go func() {
			defer f.rotateCleanWg.Done()
			logDir := filepath.Dir(f.fileAbsPath)
			current := time.Now()
			_ = filepath.Walk(logDir, func(path string, info fs.FileInfo, err error) error {
				if info == nil || info.IsDir() {
					return nil
				}

				timestampExt := filepath.Ext(path)
				if len(timestampExt) == 0 {
					return nil
				}

				fileTime, err := time.ParseInLocation(time.DateOnly, timestampExt[1:], prevOpenTime.Location())
				if err != nil {
					return nil
				}

				deleteTime := fileTime.Add(24 * time.Hour * time.Duration(f.rotateMaxDay))
				if deleteTime.Before(current) {
					_ = os.Remove(path)
				}
				return nil
			})
		}()
	}
	return nil
}

func (f *FileExecutor) SetRotateMaxDays(maxDay int) {
	f.rotateMaxDay = maxDay
}

func (f *FileExecutor) EnableFileRotate() {
	f.isRotateLogFile = true
}

func (f *FileExecutor) DisableFileRotate() {
	f.isRotateLogFile = false
}

func (f *FileExecutor) WriteMsg(msg []byte) error {
	if len(msg) == 0 {
		return nil
	}

	if f.fileWriter == nil {
		if err := f.openLogFile(); err != nil {
			return err
		}
	}

	if f.isNeedRotate() {
		if err := f.rotateLogFile(); err != nil {
			return err
		}
	}

	if f.fileWriter != nil {
		// requires '\n' as the end of the line in the file
		if msg[len(msg)-1] != '\n' {
			msg = append(msg, '\n')
		}

		_, err := f.fileWriter.Write(msg)
		return err
	}
	return nil
}

func (f *FileExecutor) Flush() {
	if f.fileWriter != nil {
		_ = f.fileWriter.Sync()
	}
}

func (f *FileExecutor) Close() {
	if f.fileWriter != nil {
		_ = f.fileWriter.Close()
	}

	f.fileWriter = nil
	f.rotateCleanWg.Wait()
}
