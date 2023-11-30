/* Package utils provides utility functions for the cgruser-exporter project.*/
package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

type UserSlice struct {
	Path               string
	UID                string
	Username           string
	MemoryCurrent      uint64
	MemoryCurrentHuman string
	MemoryAnon         uint64
	MemoryAnonHuman    string
	SwapCurrent        uint64
	SwapCurrentHuman   string
}

func GetUserSlices(path string) (slices []UserSlice) {
	names := scrapeSliceNames(path)

	for _, name := range names {
		slice := UserSlice{
			Path: filepath.Join(path, name),
			UID:  extractUID(name),
		}
		slice.Username = getUsername(slice.UID)
		slice.MemoryCurrent = getNumericFileContents(filepath.Join(slice.Path, "memory.current"))
		slice.MemoryCurrentHuman = humanize.IBytes(slice.MemoryCurrent)
		slice.MemoryAnon = getNumValueFromFile(filepath.Join(slice.Path, "memory.stat"), "anon")
		slice.MemoryAnonHuman = humanize.IBytes(slice.MemoryAnon)
		slice.SwapCurrent = getNumericFileContents(filepath.Join(slice.Path, "memory.swap.current"))
		slice.SwapCurrentHuman = humanize.IBytes(slice.SwapCurrent)
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
		if err != nil {
			log.Fatal(err)
		}
		if pattern.MatchString(filename.Name()) {
			sliceNames = append(sliceNames, filename.Name())
			fmt.Printf("Detected a user slice: %s\n", filename.Name())
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

func getNumericFileContents(path string) (value uint64) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("ERROR: Could not read '%s'!", path)
	}

	value, err = strconv.ParseUint(strings.Trim(string(data), "\n"), 10, 64)
	if err != nil {
		log.Printf("ERROR: Failed to parse '%s': %s", path, err)
	}

	return
}

func getNumValueFromFile(path, key string) (valueUint64 uint64) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, key) {
			value := strings.Split(line, " ")[1]
			valueUint64, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				log.Printf("ERROR: Failed to parse '%s': %s", path, err)
			}
			return valueUint64
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Key '%s' not found in '%s'\n", key, path)
	return
}

// GetHostname returns the hostname of the current node
func GetHostname() (hostname string) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get the hostname from the OS")
	}
	return hostname
}
