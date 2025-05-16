package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"

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

func calculateCPUPercent(containerID string, ch chan<- float64, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("docker", "stats", "--no-stream", "--format", "{{.CPUPerc}}", containerID)
	output, err := cmd.Output()
	if err != nil {
		log.Println("Failed to get CPU usage:", err)
		return
	}
	cpuPercentStr := strings.TrimSpace(string(output))
	cpuPercentStr = strings.TrimSuffix(cpuPercentStr, "%")
	cpuPercent, _ := strconv.ParseFloat(cpuPercentStr, 64)

	ch <- cpuPercent
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

		// Calculate CPU usage - using a separate channel for each container
		ch := make(chan float64, 1)
		var wg sync.WaitGroup

		wg.Add(1)
		go calculateCPUPercent(container.ID, ch, &wg)

		// Wait for the goroutine to finish before closing the channel
		wg.Wait()
		close(ch)

		var cpuPercent float64
		for percent := range ch {
			cpuPercent = percent
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
