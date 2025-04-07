package main

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"log"
	"time"
)

type SystemStats struct {
	CPUUsage    float64              `json:"cpu_usage"`
	MemoryUsage float64              `json:"memory_usage"`
	NetIO       []net.IOCountersStat `json:"net_io"`
	DiskUsage   disk.UsageStat       `json:"disk_usage"`
}

func (c *connMap) handleConnections() {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			stats := getStats()
			data, err := json.Marshal(stats)
			if err != nil {
				log.Println(fmt.Sprintf("msg json.Marshal failed. msg:%v, err:%v", stats, err))
				continue
			}
			var allConn = c.getAllConn()
			for _, conn := range allConn {
				if err := conn.SendMsg(SERVER_RESOURCE_REPORTING, data); err != nil {
					log.Println(fmt.Sprintf("msg send failed. msg:%s, err:%v", string(data), err))
					return
				}
			}
		}
	}
}

func getStats() SystemStats {
	stats := SystemStats{}
	cpuUsage, _ := cpu.Percent(0, false)
	if len(cpuUsage) > 0 {
		stats.CPUUsage = cpuUsage[0]
	}
	memInfo, _ := mem.VirtualMemory()
	stats.MemoryUsage = memInfo.UsedPercent
	netIO, _ := net.IOCounters(false)
	if len(netIO) > 0 {
		stats.NetIO = netIO
	}
	diskUsage, _ := disk.Usage("/")
	stats.DiskUsage = *diskUsage
	return stats
}
