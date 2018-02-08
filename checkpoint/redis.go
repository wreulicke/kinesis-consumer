package checkpoint

import (
	"fmt"
	"os"

	redis "gopkg.in/redis.v5"
)

// NewWithLocalRedis returns a checkpoint that uses Redis in local for underlying storage
func NewWithLocalRedis(appName string) (Checkpoint, error) {
	return NewRedis("127.0.0.1:6379", appName)
}

// NewRedis returns a checkpoint that uses Redis for underlying storage
func NewRedis(addr, appName string) (Checkpoint, error) {
	if addr == "" {
		addr = os.Getenv("REDIS_URL")
	}

	client := redis.NewClient(&redis.Options{Addr: addr})

	// verify we can ping server
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &checkpoint{
		appName,
		client,
	}, nil
}

// checkpoint stores and retreives the last evaluated key from a DDB scan
type checkpoint struct {
	appName string
	client  *redis.Client
}

// Get fetches the checkpoint for a particular Shard.
func (c *checkpoint) Get(streamName, shardID string) (string, error) {
	val, _ := c.client.Get(c.key(streamName, shardID)).Result()
	return val, nil
}

// Set stores a checkpoint for a shard (e.g. sequence number of last record processed by application).
// Upon failover, record processing is resumed from this point.
func (c *checkpoint) Set(streamName, shardID, sequenceNumber string) error {
	if sequenceNumber == "" {
		return fmt.Errorf("sequence number should not be empty")
	}
	err := c.client.Set(c.key(streamName, shardID), sequenceNumber, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// key generates a unique Redis key for storage of Checkpoint.
func (c *checkpoint) key(streamName, shardID string) string {
	return fmt.Sprintf("%v:checkpoint:%v:%v", c.appName, streamName, shardID)
}
