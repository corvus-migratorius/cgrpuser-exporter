package utils

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type UserSlice struct {
	Path          string
	UID           string
	Username      string
	MemoryCurrent uint64
	// MemoryCurrentHuman string
}

func GetUserSlices(path string) (slices []UserSlice) {
	names := scrapeSliceNames(path)

	for _, name := range names {
		slice := UserSlice{
			Path: filepath.Join(path, name),
			UID:  extractUID(name),
		}
		slice.Username = getUsername(slice.UID)
		slice.MemoryCurrent = getMemoryCurrent(slice.Path)
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

func extractUID(name string) (uid string) {
	pattern := regexp.MustCompile(`(user-)(\d+)(.slice)`)
	uid = string(pattern.FindStringSubmatch(name)[2])
	return
}

func getUsername(uid string) (username string) {
	user, err := user.LookupId(uid)
	if err != nil {
		log.Printf("ERROR: Could not find a user with UID %s!", uid)
	}
	return user.Username
}

func getMemoryCurrent(slicePath string) (memCurrent uint64) {
	memFilePath := filepath.Join(slicePath, "memory.current")
	data, err := os.ReadFile(memFilePath)
	if err != nil {
		log.Printf("ERROR: Could not read '%s'!", memFilePath)
	}

	memCurrent, err = strconv.ParseUint(strings.Trim(string(data), "\n"), 10, 64)
	if err != nil {
		log.Printf("ERROR: Failed to parse '%s': %s", memFilePath, err)
	}

	return
}
