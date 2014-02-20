package sensor

const (
	PROC_PARTITIONS  = "/proc/partitions"
	PROC_DISKSTATS   = "/proc/diskstats"
	PROC_FILESYSTEMS = "/proc/filesystems"
	ETC_MTAB         = "/etc/mtab"
)

var (
	NEWLINE = []byte("\n")
	NODEV   = []byte("nodev")
)
