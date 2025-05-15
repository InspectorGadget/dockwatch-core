package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/InspectorGadget/dockwatch-core/structs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var dockerClient *client.Client

func ConnectToDockerSock() error {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return errors.New(err.Error())
	}

	dockerClient = client

	return nil
}

func GetClient() *client.Client {
	return dockerClient
}

func FetchContainers() ([]structs.Container, error) {
	var dockerContainers []structs.Container

	containers, err := GetClient().ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return []structs.Container{}, err
	}

	for _, container := range containers {
		reader, err := GetClient().ContainerStatsOneShot(context.Background(), container.ID)
		if err != nil {
			return []structs.Container{}, err
		}

		// Decode the stat to JSON stream
		var stat structs.StatsResponse
		decoder := json.NewDecoder(reader.Body)
		if err := decoder.Decode(&stat); err != nil && err != io.EOF {
			log.Println("Failed to decode stats:", err)
			continue
		}
		reader.Body.Close()

		// Calculate CPU usage
		var cpuPercent float64
		cpuDelta := float64(stat.CPUStats.CPUUsage.TotalUsage - stat.PreCPUStats.CPUUsage.TotalUsage)
		systemDelta := float64(stat.CPUStats.SystemUsage - stat.PreCPUStats.SystemUsage)

		if systemDelta > 0 && cpuDelta > 0 {
			cpuPercent = (cpuDelta / systemDelta) * 100.0
		}

		// Fallback: Use simplified estimation (less accurate, but safer)
		if stat.CPUStats.CPUUsage.TotalUsage > 0 {
			cpuPercent = float64(stat.CPUStats.CPUUsage.TotalUsage) / float64(stat.CPUStats.SystemUsage) * 100.0
		}

		// Calculate Memory
		memUsage := float64(stat.MemoryStats.Usage) / (1024 * 1024)
		memLimit := float64(stat.MemoryStats.Limit) / (1024 * 1024)
		memPercent := (memUsage / memLimit) * 100.0

		dockerContainers = append(
			dockerContainers,
			structs.Container{
				Name:   strings.Join(container.Names, ", "),
				State:  container.State,
				Status: container.Status,
				Stat: structs.ContainerMetric{
					CPU:      fmt.Sprintf("%.1f%%", cpuPercent),
					MemUsage: fmt.Sprintf("%.1f MB", memUsage),
					MemLimit: fmt.Sprintf("%.1f MB", memLimit),
					MemPerc:  fmt.Sprintf("%.1f%%", memPercent),
				},
			},
		)
	}

	return dockerContainers, nil
}
