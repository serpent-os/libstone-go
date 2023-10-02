package payload

import "fmt"

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%9d   B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%9.2f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
