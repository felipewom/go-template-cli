package cmd

import (
	"bufio"
	"fmt"
	"github.com/karrick/godirwalk"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		// This is just for better user interaction
		fmt.Println("Project Scaffolder CLI Tool")
		fmt.Println("")

		// Extract the repository name from the URL
		repoURL := promptUser("Enter GitHub repository URL or name (user/repo-name or https://github.com/user/repo-name.git): ")
		repoName := extractRepoName(repoURL)

		destination := promptUser("Enter the destination folder for the cloned project: ")

		// Ensure the repository name is in the correct format (user/repo-name)
		parts := strings.Split(repoName, "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			fmt.Println("Invalid repository name. Please enter in the format user/repo-name.")
			os.Exit(1)
		}

		repoURL = constructRepoURL(repoName)

		fmt.Printf("Cloning repository: %s\n", repoURL)
		fmt.Printf("Destination folder: %s\n", destination)
		fmt.Println("Note: Make sure the repository is public or you have SSH key authentication configured.")

		cloneRepo(repoURL, destination)
	}
}

// isPublicRepo checks if a GitHub repository is public
func isPublicRepo(repoName string) bool {
	// You may need to enhance this logic based on your requirements
	// This is a simplified check assuming public repositories end with .git
	return !strings.HasSuffix(repoName, ".git")
}

// extractRepoName extracts the repository name from the URL
func extractRepoName(repoURL string) string {
	for repoURL == "" {
		fmt.Println("Please enter a non-empty repository URL.")
		return extractRepoName(promptUser("Enter GitHub repository URL or name (user/repo-name or https://github.com/user/repo-name.git): "))
	}
	// Try to extract the repository name from the URL
	parts := strings.Split(repoURL, "/")
	if len(parts) > 1 && strings.HasSuffix(parts[1], ".git") {
		return strings.TrimSuffix(parts[1], ".git")
	} else if len(parts) > 2 {
		return parts[len(parts)-2] + "/" + parts[len(parts)-1]
	}
	return repoURL
}

// constructRepoURL constructs the full repository URL
func constructRepoURL(repoName string) string {
	// Check if the repository is public (use HTTPS) or private (use SSH)
	isPublic := isPublicRepo(repoName)
	if isPublic {
		return fmt.Sprintf("https://github.com/%s.git", repoName)
	}
	return fmt.Sprintf("git@github.com:%s.git", repoName)
}

// promptUser is a helper function to prompt the user for input
func promptUser(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// cloneRepo clones the repository, prints the directory structure, and performs word substitution if needed
func cloneRepo(repoURL, destination string) {
	gitBinary := viper.GetString("gitBinary")
	cmd := exec.Command(gitBinary, "clone", repoURL, destination)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error cloning repository:", err)
		os.Exit(1)
	}

	// Change into the cloned folder
	if err := os.Chdir(destination); err != nil {
		fmt.Println("Error changing into the cloned folder:", err)
		os.Exit(1)
	}

	// Print the directory structure using a recursive function
	fmt.Println("Directory structure:")
	printDirectoryStructure(".", 0)

	// Ask the user if they want to perform word substitution
	substitute := promptUser("Do you want to perform word substitution in files? (yes/no): ")
	if strings.ToLower(substitute) == "yes" || strings.ToLower(substitute) == "y" {
		oldWord := promptUser("Enter the word to replace: ")
		newWord := promptUser("Enter the new word: ")

		// Perform word substitution
		modifiedFiles, err := substituteWord(".", oldWord, newWord)
		if err != nil {
			fmt.Println("Error performing word substitution:", err)
			os.Exit(1)
		}

		// Print the list of modified files
		fmt.Println("\nModified files:")
		for _, file := range modifiedFiles {
			fmt.Println(file)
		}

		// Ask if the user wants to perform substitution again or exit
		for {
			performAgain := promptUser("Do you want to perform another substitution? (yes/no): ")
			if strings.ToLower(performAgain) == "yes" || strings.ToLower(performAgain) == "y" {
				oldWord = promptUser("Enter the word to replace: ")
				newWord = promptUser("Enter the new word: ")

				// Perform word substitution again
				modifiedFiles, err = substituteWord(".", oldWord, newWord)
				if err != nil {
					fmt.Println("Error performing word substitution:", err)
					os.Exit(1)
				}

				// Print the list of modified files
				fmt.Println("\nModified files:")
				for _, file := range modifiedFiles {
					fmt.Println(file)
				}
			} else {
				break
			}
		}
	}
}

// printDirectoryStructure recursively prints the directory structure
func printDirectoryStructure(path string, indentation int) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var folders, files []os.FileInfo

	for _, entry := range entries {
		if shouldExcludeFolder(entry) {
			continue
		}

		if entry.IsDir() {
			folders = append(folders, entry)
		} else {
			files = append(files, entry)
		}
	}

	// Print folders first
	for i, entry := range folders {
		isLast := i == len(folders)-1
		prefix := getPrefix(indentation, isLast)
		fmt.Printf("%s%s/\n", prefix, entry.Name())

		printDirectoryStructure(filepath.Join(path, entry.Name()), indentation+1)
	}

	// Print files
	for i, entry := range files {
		isLast := i == len(files)-1
		prefix := getPrefix(indentation, isLast)
		fmt.Printf("%s%s\n", prefix, entry.Name())
	}
}

// getPrefix generates the prefix for each entry in the directory structure
func getPrefix(indentation int, isLast bool) string {
	if indentation == 0 {
		return ""
	}

	var prefix string
	for i := 0; i < indentation-1; i++ {
		if i%2 == 0 {
			prefix += "│  "
		} else {
			prefix += "   "
		}
	}

	if isLast {
		prefix += "└──"
	} else {
		prefix += "├──"
	}

	return prefix
}

// shouldExcludeFolder checks if a folder should be excluded from the directory structure
func shouldExcludeFolder(de fs.FileInfo) bool {
	// Exclude hidden folders and specific folders like .git
	return de.IsDir() && (strings.HasPrefix(de.Name(), ".") || de.Name() == "node_modules" || de.Name() == "vendor")
}

// substituteWord recursively performs word substitution in specified file types
func substituteWord(rootPath, oldWord, newWord string) ([]string, error) {
	var modifiedFiles []string

	err := godirwalk.Walk(rootPath, &godirwalk.Options{
		Unsorted: true,
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.IsDir() && strings.HasPrefix(de.Name(), ".") {
				return filepath.SkipDir
			}

			if !de.IsDir() && isValidFileType(path) {
				// Read the content of the file
				content, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				// Perform word substitution
				newContent := strings.Replace(string(content), oldWord, newWord, -1)

				// Get file information to retrieve the mode
				fileInfo, err := os.Stat(path)
				if err != nil {
					return err
				}

				// Debug output
				fmt.Printf("Substituting '%s' with '%s' in file: %s\n", oldWord, newWord, path)

				// Write the modified content back to the file
				err = ioutil.WriteFile(path, []byte(newContent), fileInfo.Mode())
				if err != nil {
					return err
				}

				modifiedFiles = append(modifiedFiles, path)
			}

			return nil
		},
	})

	return modifiedFiles, err
}

// isValidFileType checks if the file has a valid type for word substitution
func isValidFileType(filePath string) bool {
	validFileTypes := map[string]bool{
		".ts":   true,
		".js":   true,
		".tsx":  true,
		".go":   true,
		".json": true,
		".yaml": true,
		".yml":  true,
		".xml":  true,
		".md":   true,
		".html": true,
		".css":  true,
		".scss": true,
		// Add more file types as needed
	}

	ext := filepath.Ext(filePath)
	return validFileTypes[ext]
}
