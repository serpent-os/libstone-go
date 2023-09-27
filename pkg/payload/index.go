package payload

type IndexEntry struct {
	Start  uint64
	End    uint64
	digest [2]uint64
}
