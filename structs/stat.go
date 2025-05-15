package structs

import "time"

// Struct data retrieved from https://pkg.go.dev/github.com/docker/docker@v28.1.1+incompatible/api/types/container

type StatsResponse struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`

	// Common stats
	Read    time.Time `json:"read"`
	PreRead time.Time `json:"preread"`

	PidsStats struct {
		Current uint64 `json:"current,omitempty"`
		Limit   uint64 `json:"limit,omitempty"`
	} `json:"pids_stats"`
	BlkioStats struct {
		IoServiceBytesRecursive []BlkioStatEntry `json:"io_service_bytes_recursive"`
		IoServicedRecursive     []BlkioStatEntry `json:"io_serviced_recursive"`
		IoQueuedRecursive       []BlkioStatEntry `json:"io_queue_recursive"`
		IoServiceTimeRecursive  []BlkioStatEntry `json:"io_service_time_recursive"`
		IoWaitTimeRecursive     []BlkioStatEntry `json:"io_wait_time_recursive"`
		IoMergedRecursive       []BlkioStatEntry `json:"io_merged_recursive"`
		IoTimeRecursive         []BlkioStatEntry `json:"io_time_recursive"`
		SectorsRecursive        []BlkioStatEntry `json:"sectors_recursive"`
	} `json:"blkio_stats"`

	NumProcs     uint32 `json:"num_procs"`
	StorageStats struct {
		ReadCountNormalized  uint64 `json:"read_count_normalized,omitempty"`
		ReadSizeBytes        uint64 `json:"read_size_bytes,omitempty"`
		WriteCountNormalized uint64 `json:"write_count_normalized,omitempty"`
		WriteSizeBytes       uint64 `json:"write_size_bytes,omitempty"`
	} `json:"storage_stats"`

	// Shared stats
	CPUStats struct {
		CPUUsage       CPUUsage       `json:"cpu_usage"`
		SystemUsage    uint64         `json:"system_cpu_usage,omitempty"`
		OnlineCPUs     uint32         `json:"online_cpus,omitempty"`
		ThrottlingData ThrottlingData `json:"throttling_data,omitempty"`
	} `json:"cpu_stats"`

	PreCPUStats struct {
		// CPU Usage. Linux and Windows.
		CPUUsage CPUUsage `json:"cpu_usage"`

		// System Usage. Linux only.
		SystemUsage uint64 `json:"system_cpu_usage,omitempty"`

		// Online CPUs. Linux only.
		OnlineCPUs uint32 `json:"online_cpus,omitempty"`

		// Throttling Data. Linux only.
		ThrottlingData ThrottlingData `json:"throttling_data,omitempty"`
	} `json:"precpu_stats"` // "Pre"="Previous"

	MemoryStats struct {

		// current res_counter usage for memory
		Usage uint64 `json:"usage,omitempty"`
		// maximum usage ever recorded.
		MaxUsage uint64 `json:"max_usage,omitempty"`
		// TODO(vishh): Export these as stronger types.
		// all the stats exported via memory.stat.
		Stats map[string]uint64 `json:"stats,omitempty"`
		// number of times memory usage hits limits.
		Failcnt uint64 `json:"failcnt,omitempty"`
		Limit   uint64 `json:"limit,omitempty"`

		// committed bytes
		Commit uint64 `json:"commitbytes,omitempty"`
		// peak committed bytes
		CommitPeak uint64 `json:"commitpeakbytes,omitempty"`
		// private working set
		PrivateWorkingSet uint64 `json:"privateworkingset,omitempty"`
	} `json:"memory_stats,omitempty"`
	Networks map[string]struct {
		// Bytes received. Windows and Linux.
		RxBytes uint64 `json:"rx_bytes"`
		// Packets received. Windows and Linux.
		RxPackets uint64 `json:"rx_packets"`
		// Received errors. Not used on Windows. Note that we don't `omitempty` this
		// field as it is expected in the >=v1.21 API stats structure.
		RxErrors uint64 `json:"rx_errors"`
		// Incoming packets dropped. Windows and Linux.
		RxDropped uint64 `json:"rx_dropped"`
		// Bytes sent. Windows and Linux.
		TxBytes uint64 `json:"tx_bytes"`
		// Packets sent. Windows and Linux.
		TxPackets uint64 `json:"tx_packets"`
		// Sent errors. Not used on Windows. Note that we don't `omitempty` this
		// field as it is expected in the >=v1.21 API stats structure.
		TxErrors uint64 `json:"tx_errors"`
		// Outgoing packets dropped. Windows and Linux.
		TxDropped uint64 `json:"tx_dropped"`
		// Endpoint ID. Not used on Linux.
		EndpointID string `json:"endpoint_id,omitempty"`
		// Instance ID. Not used on Linux.
		InstanceID string `json:"instance_id,omitempty"`
	} `json:"networks,omitempty"`
}

type BlkioStatEntry struct {
	Major uint64 `json:"major"`
	Minor uint64 `json:"minor"`
	Op    string `json:"op"`
	Value uint64 `json:"value"`
}

type CPUUsage struct {
	// Total CPU time consumed.
	// Units: nanoseconds (Linux)
	// Units: 100's of nanoseconds (Windows)
	TotalUsage uint64 `json:"total_usage"`

	// Total CPU time consumed per core (Linux). Not used on Windows.
	// Units: nanoseconds.
	PercpuUsage []uint64 `json:"percpu_usage,omitempty"`

	// Time spent by tasks of the cgroup in kernel mode (Linux).
	// Time spent by all container processes in kernel mode (Windows).
	// Units: nanoseconds (Linux).
	// Units: 100's of nanoseconds (Windows). Not populated for Hyper-V Containers.
	UsageInKernelmode uint64 `json:"usage_in_kernelmode"`

	// Time spent by tasks of the cgroup in user mode (Linux).
	// Time spent by all container processes in user mode (Windows).
	// Units: nanoseconds (Linux).
	// Units: 100's of nanoseconds (Windows). Not populated for Hyper-V Containers
	UsageInUsermode uint64 `json:"usage_in_usermode"`
}

type ThrottlingData struct {
	// Number of periods with throttling active
	Periods uint64 `json:"periods"`
	// Number of periods when the container hits its throttling limit.
	ThrottledPeriods uint64 `json:"throttled_periods"`
	// Aggregate time the container was throttled for in nanoseconds.
	ThrottledTime uint64 `json:"throttled_time"`
}
