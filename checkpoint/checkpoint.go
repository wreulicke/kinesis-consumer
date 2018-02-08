package checkpoint

// Checkpoint stores and retreives the last evaluated key from a scan
type Checkpoint interface {
	Get(streamName, shardID string) (string, error)
	Set(streamName, shardID, sequenceNumber string) error
}
