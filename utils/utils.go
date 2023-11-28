package utils

import (
	"fmt"
	"log"
	"os"
	"regexp"
)

func ScrapeSliceNames() (sliceNames []string) {
	pattern := regexp.MustCompile(`user-\d+`)

	dir, err := os.Open("/sys/fs/cgroup/user.slice")
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
	// output, err := exec.Command("stat", "-fc", "%T", "/sys/fs/cgroup/").Output()
	//
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	// return strings.Trim(string(output), "\n")
}
