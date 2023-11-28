package exporter

import (
	"fmt"
	"time"
)

type cgrpUserExporter struct {
	// Slices  []UserSlice
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

func CgroupUserExporter(timeout int) *cgrpUserExporter {
	exporter := cgrpUserExporter{
		Timeout: timeout,
	}

	return &exporter
}


