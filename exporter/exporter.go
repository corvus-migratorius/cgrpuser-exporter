package exporter

import (
	"fmt"
	"time"

	"cgrpuser-exporter/utils"
)

type cgrpUserExporter struct {
	Hostname string
	Slices  []utils.UserSlice
	Timeout int
}

func (self *cgrpUserExporter) setSliceNames() {
	go func() {
		for {
			fmt.Println("Exporter is alive")
			time.Sleep(time.Duration(self.Timeout) * time.Second)
		}
	}()
}

func CgroupUserExporter(userSlicePath string, timeout int) *cgrpUserExporter {
	exporter := cgrpUserExporter{
		Hostname: utils.GetHostname(),
		Slices: utils.GetUserSlices(userSlicePath),
		Timeout: timeout,
	}

	return &exporter
}


