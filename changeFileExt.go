package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Function to change file extensions
func changeFileExtensions(oldExt string, newExt string, folderPath string) ([]string, []error) {
	var renamedFiles []string
	var errors []error

	if !strings.Contains(oldExt, ".") {
		oldExt = "." + oldExt
	}

	if !strings.Contains(newExt, ".") {
		newExt = "." + newExt
	}

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		errors = append(errors, fmt.Errorf("Error reading directory %s: %w", folderPath, err))
		return renamedFiles, errors
	}
	for _, file := range files {

		if strings.HasSuffix(file.Name(), oldExt) {

			oldName := folderPath + "/" + file.Name()
			newName := strings.TrimSuffix(oldName, oldExt) + newExt

			err := os.Rename(oldName, newName)
			if err != nil {
				errors = append(errors, fmt.Errorf("Failed to rename %s to %s: %w", oldName, newName, err))
			} else {
				renamedFiles = append(renamedFiles, newName)
			}
		}
	}

	return renamedFiles, errors
}

func main() {

	var oldExt, newExt string
	var folderPath string

	fmt.Println("Enter the path to the folder (e.g., /path/to/your/files or . for current directory):")
	fmt.Scan(&folderPath)

	// Validate folder path
	fileInfo, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		fmt.Printf("Error: Folder path '%s' does not exist.\n", folderPath)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("Error accessing folder path '%s': %v\n", folderPath, err)
		os.Exit(1)
	}
	if !fileInfo.IsDir() {
		fmt.Printf("Error: Path '%s' is not a directory.\n", folderPath)
		os.Exit(1)
	}

	fmt.Println("Enter original extension (ex=>jpg)")
	fmt.Scan(&oldExt)

	// Validate oldExt
	if oldExt == "" {
		fmt.Println("Error: Original extension cannot be empty.")
		os.Exit(1)
	}

	fmt.Println("Enter new extension (ex=>jpeg)")
	fmt.Scan(&newExt)

	// Validate newExt
	if newExt == "" {
		fmt.Println("Error: New extension cannot be empty.")
		os.Exit(1)
	}

	renamed, errs := changeFileExtensions(oldExt, newExt, folderPath)

	if len(errs) > 0 {
		fmt.Println("Errors encountered:")
		for _, err := range errs {
			fmt.Println("- ", err)
		}
	}

	if len(renamed) > 0 {
		fmt.Println("Successfully renamed files:")
		for _, file := range renamed {
			fmt.Println("- ", file)
		}
		fmt.Printf("%d file(s) renamed successfully.\n", len(renamed))
	}

	if len(renamed) == 0 && len(errs) == 0 {
		fmt.Println("No files found with the original extension or no files needed renaming.")
	}

}
