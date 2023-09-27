package payload

type FileType uint8

const (
	Regular FileType = iota + 1
	Symlink
	Directory
	CharacterDevice
	BlockDevice
	Fifo
	Socket
)

type LayoutEntry struct {
	UID  uint32
	GID  uint32
	Mode uint32
	Tag  uint32
	//Entry Entry
}
