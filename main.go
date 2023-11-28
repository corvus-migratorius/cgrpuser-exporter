package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"cgrpuser-exporter/exporter"
	"cgrpuser-exporter/utils"

	"github.com/akamensky/argparse"
)

func parseArgs(APPNAME string, VERSION string) (*int, *int) {
	parser := argparse.NewParser(APPNAME, "Publishes Systemd's user 'slice' metrics for Prometheus")
	port := parser.Int("p", "port", &argparse.Options{
		Help:    "Port for publishing the Prometheus exporter metrics",
		Default: 9201,
	})
	timeout := parser.Int("t", "timeout", &argparse.Options{
		Help:    "Timeout for gather user slice scraping (seconds)",
		Default: 15,
	})
	version := parser.Flag("V", "version", &argparse.Options{
		Help: "Display version number and quit",
	})

	// Parse input and display usage on error (equivalent to `--help`)
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	if *version {
		fmt.Printf("%s: version %s\n", APPNAME, VERSION)
		os.Exit(0)
	}

	return port, timeout
}

func getCgroupsVersion() string {
	version, err := exec.Command("stat", "-fc", "%T", "/sys/fs/cgroup/").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Trim(string(version), "\n")
}

func main() {
	APPNAME := "cgrpuser-exporter"
	VERSION := "0.1.0"

	port, timeout := parseArgs(APPNAME, VERSION)

	log.Printf("Using port %d to publish /metrics\n", *port)
	log.Printf("Setting user.slice scraping timeout to %d seconds\n", *timeout)

	cgrpVersion := getCgroupsVersion()

	if cgrpVersion == "tmpfs" {
		log.Fatal("Detected cgroups v1: we only support v2.")
	} else if cgrpVersion == "cgroup2fs" {
		log.Println("Detected cgroups v2 - good to go!")
	} else {
		log.Fatalf("Could not determine cgroups version (got '%s' for '/sys/fs/cgroup/' fs type)", cgrpVersion)
	}

	fmt.Println(utils.ScrapeSliceNames())

	exporter := exporter.CgroupUserExporter(*timeout)
	fmt.Printf("%#v\n", exporter)

	// cnexporter.RecordCounts()
	// cnexporter.RecordMetadata()

	// http.Handle("/metrics", promhttp.Handler())
	// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
