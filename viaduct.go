package main

import (
	"log"
	"os"
	"path/filepath"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func linkUp(targetDir string, sourceDir string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		relSourcePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			log.Fatal(err)
		}

		if relSourcePath == "." {
			return nil
		}

		targetPath := filepath.Join(targetDir, relSourcePath)

		if !exists(targetPath) {
			err := os.Symlink(path, targetPath)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Linked: " + path + " ---> " + targetPath)
		} else {
			log.Println("Exists: " + path + " ---> " + targetPath)
		}

		return nil
	}
}

func main() {
	sourceDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	sourceDir = filepath.Clean(sourceDir)
	log.Println("Sourcedir: " + sourceDir)

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	targetDir := filepath.Join(currentDir, "..")
	targetDir = filepath.Clean(targetDir)
	log.Println("Targetdir: " + targetDir)

	err = os.Chdir(targetDir)
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(sourceDir, linkUp(targetDir, sourceDir))
	if err != nil {
		log.Fatal(err)
	}
}