package sensor

const (
	PROC            = "/proc"
	ProcStat        = "/proc/stat"
	ProcCPUInfo     = "/proc/cpuinfo"
	ProcMemInfo     = "/proc/meminfo"
	ProcNetDev      = "/proc/net/dev"
	ProcPartitions  = "/proc/partitions"
	ProcDiskStats   = "/proc/diskstats"
	ProcFileSystems = "/proc/filesystems"
	EtcMtab         = "/etc/mtab"
)

var (
	Newline    = []byte("\n")
	NoDev      = []byte("nodev")
	PhysicalID = []byte("physical id")
	MemKB      = []byte("kB")
	BTime      = []byte("btime")
)
