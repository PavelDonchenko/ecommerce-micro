package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func CreateFolders(folders []string) {
	jsonStr, _ := json.Marshal(folders)

	var data []string
	_ = json.Unmarshal(jsonStr, &data)
	for _, name := range data {
		if _, err := os.Stat(name); err != nil {
			os.Mkdir(name, os.ModePerm)
		}
	}
}

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

func CreateFile(data []byte, pathFile string) error {
	file, err := os.Create(pathFile)
	if err != nil {
		os.Exit(1)
		return errors.New("invalid file path")
	}

	_, err = file.Write(data)
	if err != nil {
		os.Exit(1)
		return fmt.Errorf("error when write file %s: %s \n", pathFile, err)
	}

	return nil
}
