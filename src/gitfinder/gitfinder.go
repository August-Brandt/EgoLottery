package gitfinder

import (
	"log"
	"os"
	"path/filepath"
)

// Takes a slice of directory paths (dirs) and looks for a .git folder.
//
// Will look into a sup-directories until reaching a level of depth.
//
// A negative depth value will be seen as a value of 0.
func FindGitFolders(dirs []string, depth int) []string {
	gitDirs := []string{}
	var absPath string
	for _, dir := range dirs {
		absPath = dir
		if !pathIsDir(dir) {
			log.Printf("'%s' is not a directory\n", dir)
			continue
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, entry := range entries {
			if entry.Name() == ".git" {
				gitDirs = append(gitDirs, filepath.Join(absPath, entry.Name()))
			}
		}
	}
	return gitDirs
}

func pathIsDir(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return false
	}
	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(err)
		return false
	}
	return fileInfo.IsDir()
}