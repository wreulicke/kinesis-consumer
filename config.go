package connector

import (
	"io/ioutil"
	"time"

	"log"
)

// Config vars for the application
type Config struct {
	// AppName is the application name and checkpoint namespace.
	AppName string

	// StreamName is the Kinesis stream.
	StreamName string

	// StreamRegion is the Kinesis stream.
	StreamRegion string

	// FlushInterval is a regular interval for flushing the buffer. Defaults to 1s.
	FlushInterval time.Duration

	// BufferSize determines the batch request size. Must not exceed 500. Defaults to 500.
	BufferSize int

	// Logger is the logger used. Defaults to log.Log.
	Logger *log.Logger

	// Checkpoint for tracking progress of consumer.
	Checkpoint Checkpoint
}

type noopCheckpoint struct{}

func (n noopCheckpoint) Set(string, string, string) error   { return nil }
func (n noopCheckpoint) Get(string, string) (string, error) { return "", nil }

// defaults for configuration.
func (c *Config) setDefaults() {
	if c.Logger == nil {
		c.Logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}
	c.Logger.Println("kinesis-connectors")

	if c.AppName == "" {
		c.Logger.Panicf("AppName required")
	}

	if c.StreamName == "" {
		c.Logger.Panicf("StreamName required")
	}

	if c.StreamRegion == "" {
		c.Logger.Panicf("StreamRegion required")
	}

	if c.BufferSize == 0 {
		c.BufferSize = 500
	}

	if c.FlushInterval == 0 {
		c.FlushInterval = time.Second
	}

	if c.Checkpoint == nil {
		c.Checkpoint = &noopCheckpoint{}
	}
}
