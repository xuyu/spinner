package sensor

const (
	PROC_STAT        = "/proc/stat"
	PROC_CPUINFO     = "/proc/cpuinfo"
	PROC_MEMINFO     = "/proc/meminfo"
	PROC_NET_DEV     = "/pro/net/dev"
	PROC_PARTITIONS  = "/proc/partitions"
	PROC_DISKSTATS   = "/proc/diskstats"
	PROC_FILESYSTEMS = "/proc/filesystems"
	ETC_MTAB         = "/etc/mtab"
)

var (
	NEWLINE     = []byte("\n")
	NODEV       = []byte("nodev")
	PHYSICAL_ID = []byte("physical id")
	MEM_KB      = []byte("kB")
	BTIME       = []byte("btime")
)
