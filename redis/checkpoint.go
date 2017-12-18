package redis

import (
	"fmt"

	"gopkg.in/redis.v5"
)

// New RedisCheckpoint returns a checkpoint that uses Redis for underlying storage
func New(appName, addr string) (*Checkpoint, error) {
	client := redis.NewClient(&redis.Options{Addr: addr})

	// verify we can ping server
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &Checkpoint{
		AppName: appName,
		Client:  client,
	}, nil
}

// Checkpoint implements the Checkpont interface.
// Used to enable the Pipeline.ProcessShard to checkpoint it's progress
// while reading records from Kinesis stream.
type Checkpoint struct {
	AppName string
	Client  *redis.Client
}

// Get last seq.
func (c *Checkpoint) Get(streamName, shardID string) (string, error) {
	return c.Client.Get(c.key(streamName, shardID)).Result()
}

// Set stores a checkpoint for a shard (e.g. sequence number of last record processed by application).
// Upon failover, record processing is resumed from this point.
func (c *Checkpoint) Set(streamName, shardID, sequenceNumber string) error {
	if sequenceNumber == "" {
		return fmt.Errorf("sequence number should not be empty")
	}
	return c.Client.Set(c.key(streamName, shardID), sequenceNumber, 0).Err()
}

// key generates a unique Redis key for storage of Checkpoint.
func (c *Checkpoint) key(streamName, shardID string) string {
	return fmt.Sprintf("%v:checkpoint:%v:%v", c.AppName, streamName, shardID)
}
