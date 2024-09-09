package helpers

import (
	"os"
	"testing"
)

func TestCreateFolders(t *testing.T) {
	folders := []string{"test1", "test2"}
	CreateFolders(folders)

	for _, folder := range folders {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			t.Errorf("Expected folder %s to be created, but it was not", folder)
		}
	}

	for _, folder := range folders {
		os.RemoveAll(folder)
	}
}
