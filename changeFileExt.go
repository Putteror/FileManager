package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Function to change file extensions
func changeFileExtensions(oldExt string, newExt string, folderPath string) string {

	if !strings.Contains(oldExt, ".") {
		oldExt = "." + oldExt
	}

	if !strings.Contains(newExt, ".") {
		newExt = "." + newExt
	}

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Error:", err)
		return err.Error()
	}
	for _, file := range files {

		if strings.HasSuffix(file.Name(), oldExt) {

			oldName := folderPath + "/" + file.Name()
			newName := strings.TrimSuffix(oldName, oldExt) + newExt

			err := os.Rename(oldName, newName)
			if err != nil {
				fmt.Printf("Failed to rename %s to %s: %v\n", oldName, newName, err)
			} else {
				fmt.Printf("Renamed: %s -> %s\n", oldName, newName)
			}
		}
	}

	return "Change File Extension"
}

func main() {

	var oldExt, newExt string
	var folderPath string

	fmt.Println("Enter folder path ( . If this file in path )")
	fmt.Scan(&folderPath)

	fmt.Println("Enter original extension (ex=>jpg)")
	fmt.Scan(&oldExt)

	fmt.Println("Enter new extension (ex=>jpeg)")
	fmt.Scan(&newExt)

	changeFileExtensions(oldExt, newExt, folderPath)

}
