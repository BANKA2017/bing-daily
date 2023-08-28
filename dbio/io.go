package dbio

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func dbioPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	if _, err = os.Stat(filepath.Dir(ex) + "/bing.json"); os.IsNotExist(err) {
		return "", err
	} else {
		return filepath.Dir(ex) + "/bing.json", nil
	}

}

func DbioRead[T any](template *T) error {
	path, err := dbioPath()
	if err != nil {
		return err
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, template); err != nil {
		return err
	}
	return err
}

func DbioWrite(content []byte) error {
	path, err := dbioPath()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	n, _ := file.Seek(0, io.SeekEnd)
	_, err = file.WriteAt(content, n)
	defer file.Close()
	return err
}
