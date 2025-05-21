package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestChangeFileExtensions_Success(t *testing.T) {
	// 1. Setup
	tempDir, err := ioutil.TempDir("", "test_change_ext_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	originalExt := ".txt"
	newExt := ".log"

	filesToCreate := map[string]string{
		"file1.txt":    "content1",
		"file2.txt":    "content2",
		"another.txt":  "content_another",
		"image.png":    "image_data", // Should not be renamed
		"subfolder":    "",           // This is a directory, should be ignored
		"file3.md":     "markdown",   // Different extension, should be ignored
	}

	expectedRenamedBaseNames := []string{"file1.log", "file2.log", "another.log"}
	expectedRemainingBaseNames := []string{"image.png", "subfolder", "file3.md"}

	// Create files and one subdirectory
	err = os.Mkdir(filepath.Join(tempDir, "subfolder"), 0755)
	if err != nil {
		t.Fatalf("Failed to create subfolder: %v", err)
	}

	for name, content := range filesToCreate {
		if name == "subfolder" {
			continue // Already created
		}
		filePath := filepath.Join(tempDir, name)
		err := ioutil.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", name, err)
		}
	}

	// 2. Execution
	renamedFiles, errs := changeFileExtensions("txt", "log", tempDir) // Use "txt" without dot

	// 3. Verification
	if len(errs) != 0 {
		var errorMessages []string
		for _, e := range errs {
			errorMessages = append(errorMessages, e.Error())
		}
		t.Fatalf("Expected no errors, but got: %s", strings.Join(errorMessages, "; "))
	}

	if len(renamedFiles) != len(expectedRenamedBaseNames) {
		t.Errorf("Expected %d files to be renamed, but got %d. Renamed files: %v", len(expectedRenamedBaseNames), len(renamedFiles), renamedFiles)
	}

	// Check if the correct files were reported as renamed
	// Extract basenames from renamedFiles for easier comparison
	var actualRenamedBaseNames []string
	for _, rf := range renamedFiles {
		actualRenamedBaseNames = append(actualRenamedBaseNames, filepath.Base(rf))
	}
	sort.Strings(actualRenamedBaseNames)
	sort.Strings(expectedRenamedBaseNames)

	if !reflect.DeepEqual(actualRenamedBaseNames, expectedRenamedBaseNames) {
		t.Errorf("Mismatch in renamed files list. Expected: %v, Got: %v", expectedRenamedBaseNames, actualRenamedBaseNames)
	}

	// Verify file system state
	// Check that original *.txt files are gone and *.log files exist
	for _, baseName := range expectedRenamedBaseNames {
		originalFileName := strings.TrimSuffix(baseName, newExt) + originalExt
		originalFilePath := filepath.Join(tempDir, originalFileName)
		newFilePath := filepath.Join(tempDir, baseName)

		if _, err := os.Stat(originalFilePath); !os.IsNotExist(err) {
			t.Errorf("Original file %s was expected to be removed, but it still exists.", originalFilePath)
		}
		if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
			t.Errorf("New file %s was expected to be created, but it does not exist.", newFilePath)
		} else if err != nil {
			t.Errorf("Error stating new file %s: %v", newFilePath, err)
		}
	}

	// Check that other files/dirs still exist
	for _, baseName := range expectedRemainingBaseNames {
		path := filepath.Join(tempDir, baseName)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("File/dir %s was expected to remain, but it does not exist.", path)
		} else if err != nil {
			t.Errorf("Error stating file/dir %s: %v", path, err)
		}
	}
}

func TestChangeFileExtensions_NoMatchingFiles(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_no_match_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	filesToCreate := []string{"file1.md", "image.jpeg"}
	for _, name := range filesToCreate {
		filePath := filepath.Join(tempDir, name)
		if err := ioutil.WriteFile(filePath, []byte("dummy"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", name, err)
		}
	}

	renamedFiles, errs := changeFileExtensions(".txt", ".log", tempDir)

	if len(errs) != 0 {
		t.Errorf("Expected no errors, got: %v", errs)
	}
	if len(renamedFiles) != 0 {
		t.Errorf("Expected 0 renamed files, got %d: %v", len(renamedFiles), renamedFiles)
	}
}

func TestChangeFileExtensions_ErrorReadingDir(t *testing.T) {
	// Using a non-existent directory to trigger ioutil.ReadDir error
	nonExistentDir := filepath.Join(os.TempDir(), "non_existent_dir_for_test")
	// Ensure it really doesn't exist, in case a previous failed test left it
	os.RemoveAll(nonExistentDir)


	renamedFiles, errs := changeFileExtensions(".txt", ".log", nonExistentDir)

	if len(errs) == 0 {
		t.Fatalf("Expected errors when reading a non-existent directory, but got none.")
	} else {
		// Check if the error message contains the path of the directory
		expectedErrorMsgPart := fmt.Sprintf("Error reading directory %s", nonExistentDir)
		if !strings.Contains(errs[0].Error(),expectedErrorMsgPart) {
			t.Errorf("Expected error message to contain '%s', but got '%s'", expectedErrorMsgPart, errs[0].Error())
		}
	}

	if len(renamedFiles) != 0 {
		t.Errorf("Expected 0 renamed files when directory read fails, got %d", len(renamedFiles))
	}
}

// Test for when os.Rename fails (e.g. new file name is invalid or already exists as a directory)
// This is harder to reliably test without more complex setup (like permissions or specific OS states)
// For now, we assume that if ReadDir works and files match, os.Rename errors are correctly propagated.
// A more advanced test might involve trying to rename a file to a name that is an existing directory.
func TestChangeFileExtensions_RenameError(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_rename_error_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	// Create a file to be renamed
	fileToRename := "file.txt"
	filePath := filepath.Join(tempDir, fileToRename)
	if err := ioutil.WriteFile(filePath, []byte("data"), 0644); err != nil {
		t.Fatalf("Failed to create file %s: %v", filePath, err)
	}

	// Create a directory with the target new name, which should cause os.Rename to fail
	targetNameAsDir := "file.log"
	if err := os.Mkdir(filepath.Join(tempDir, targetNameAsDir), 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", targetNameAsDir, err)
	}

	renamedFiles, errs := changeFileExtensions(".txt", ".log", tempDir)

	if len(errs) == 0 {
		t.Fatalf("Expected errors when rename fails, but got none.")
	} else {
		foundRenameError := false
		for _, e := range errs {
			// Error message from os.Rename on Linux for "is a directory"
			// On Windows it might be "Access is denied." or similar if target is a dir
			// This check is OS-dependent and might need adjustment
			if strings.Contains(e.Error(), "Failed to rename") && (strings.Contains(e.Error(), "is a directory") || strings.Contains(e.Error(), "Access is denied")) {
				foundRenameError = true
				break
			}
		}
		if !foundRenameError {
			t.Errorf("Expected a specific rename error, but got: %v", errs)
		}
	}

	if len(renamedFiles) != 0 {
		t.Errorf("Expected 0 renamed files when rename fails, got %d", len(renamedFiles))
	}

	// Check that the original file still exists (since rename failed)
	if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
		t.Errorf("Original file %s should still exist after a rename failure, but it's gone.", filePath)
	}
}
