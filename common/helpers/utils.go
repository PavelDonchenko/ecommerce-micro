package helpers

import (
	"encoding/json"
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
