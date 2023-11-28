package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type UserSlice struct {
	Path string
	// Username      string
	UID int
	// MemoryCurrent int
}

func GetUserSlices(path string) (slices []UserSlice) {
	names := scrapeSliceNames(path)

	for _, name := range names {
		slice := UserSlice{
			Path: filepath.Join(path, name),
			UID:  extractUID(name),
		}
		slices = append(slices, slice)
	}

	return
}

func scrapeSliceNames(path string) (sliceNames []string) {
	pattern := regexp.MustCompile(`user-\d+`)

	dir, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening directory:", err)
		return
	}
	defer dir.Close()

	filenames, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal("Error reading directory:", err)
	}

	for _, filename := range filenames {
		fmt.Println(filename.Name())
		if err != nil {
			log.Fatal(err)
		}
		if pattern.MatchString(filename.Name()) {
			sliceNames = append(sliceNames, filename.Name())
		}
	}

	return
}

func extractUID(name string) (UID int) {
	pattern := regexp.MustCompile(`(user-)(\d+)(.slice)`)
	UID, err := strconv.Atoi(pattern.FindStringSubmatch(name)[2])
	if err != nil {
		log.Fatalf("Failed to extract UID from '%s'", name)
	}

	return
}
