package emitter

import (
	"fmt"
	"time"
)

// Key create s3 key structure
func Key(prefix, firstSeq, lastSeq string) string {
	date := time.Now().UTC().Format("2006/01/02")

	if prefix == "" {
		return fmt.Sprintf("%v/%v-%v", date, firstSeq, lastSeq)
	}
	return fmt.Sprintf("%v/%v/%v-%v", prefix, date, firstSeq, lastSeq)
}
