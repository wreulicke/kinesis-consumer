package connector

// Checkpoint interface for functions that checkpoints need to
// implement in order to track consumer progress.
type Checkpoint interface {
	Get(streamName, shardID string) (string, error)
	Set(streamName, shardID, sequenceNumber string) error
}
