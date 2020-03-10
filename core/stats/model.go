package stats

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/host"
)

// Statistic is the struct used to parse each statistic.
type Statistic struct {
	// PiWorker stats
	ActiveTasks          uint16        `json:"activeTasks"`
	InactiveTasks        uint16        `json:"inactiveTasks"`
	OnExecutionTasks     uint8         `json:"onExecutionTasks"`
	FailedTasks          uint8         `json:"failedTasks"`
	AverageExecutionTime time.Duration `json:"averageExecutionTime"`
	BackupLoopState      bool          `json:"backupLoopState"`

	// Host (probably a RPi) stats
	RaspberryStats RaspberryStats `json:"raspberryStats"`

	sumExecTime time.Duration
	obs         uint64

	sync.RWMutex
}

// RaspberryStats is the struct that contains the statistics related with the Host (generally it will be a Raspberry Pi).
type RaspberryStats struct {
	Host      HostStats    `json:"temperature"`
	CPULoad   float64      `json:"cpuLoad"`
	Storage   StorageStats `json:"storage"`
	RAM       RAMStats     `json:"ram"`
	Timestamp time.Time    `json:"timestamp"`
}

// RAMStats is the struct used to parse the statistics related with the RAM of the Host.
type RAMStats struct {
	Total     uint64 `json:"total"`
	Available uint64 `json:"available"`
	Used      uint64 `json:"used"`
}

// StorageStats is the struct used to parse storage stats of the Host.
type StorageStats struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

// HostStats is the struct used to parse some additional statistics about the Host.
type HostStats struct {
	BootTime     uint64 `json:"bootTime"`
	UpTime       uint64 `json:"uptime"`
	Temperatures []host.TemperatureStat
}

// NewAvgObs is a method with the purpose of add new data to be calculated into the
// `Statistic.AverageExecutionTime` field.
func (s *Statistic) NewAvgObs(duration time.Duration) time.Duration {
	s.Lock()
	defer s.Unlock()

	s.sumExecTime += duration
	s.obs++
	s.AverageExecutionTime = time.Duration(float32(s.sumExecTime) / float32(s.obs))

	return s.AverageExecutionTime
}
