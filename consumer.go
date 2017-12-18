package connector

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

// New creates a new consumer with initialied kinesis connection
func New(config Config) *Consumer {
	config.setDefaults()

	svc := kinesis.New(
		session.New(
			aws.NewConfig().WithMaxRetries(10).WithRegion(config.StreamRegion),
		),
	)

	return &Consumer{
		svc:    svc,
		Config: config,
	}
}

// Consumer wraps the interaction with the Kinesis stream
type Consumer struct {
	svc *kinesis.Kinesis
	Config
}

// Start takes a handler and then loops over each of the shards
// processing each one with the handler.
func (c *Consumer) Start(handler Handler) {
	resp, err := c.svc.DescribeStream(
		&kinesis.DescribeStreamInput{
			StreamName: aws.String(c.StreamName),
		},
	)

	if err != nil {
		c.Logger.Panicf("Error DescribeStream %v", err.Error())
	}

	for _, shard := range resp.StreamDescription.Shards {
		go c.handlerLoop(*shard.ShardId, handler)
	}
}

func (c *Consumer) handlerLoop(shardID string, handler Handler) {
	buf := &Buffer{
		MaxRecordCount: c.BufferSize,
		shardID:        shardID,
	}
	c.Logger.Printf("Processing, %v", shardID)
	shardIterator := c.getShardIterator(shardID)
	for {
		resp, err := c.svc.GetRecords(
			&kinesis.GetRecordsInput{
				ShardIterator: shardIterator,
			},
		)
		if err != nil {
			c.Logger.Printf("GetRecords %v", err.Error())
		} else {
			if len(resp.Records) > 0 {
				for _, r := range resp.Records {
					buf.AddRecord(r)
					if buf.ShouldFlush() {
						handler.HandleRecords(*buf)
						c.Logger.Printf("Count %v flushed", buf.RecordCount())
						c.Checkpoint.Set(c.StreamName, shardID, buf.LastSeq())
						buf.Flush()
					}
				}
			}
		}
		if resp == nil || resp.NextShardIterator == nil || shardIterator == resp.NextShardIterator {
			shardIterator = c.getShardIterator(shardID)
		} else {
			shardIterator = resp.NextShardIterator
		}
	}
}

func (c *Consumer) getShardIterator(shardID string) *string {
	params := &kinesis.GetShardIteratorInput{
		StreamName: aws.String(c.StreamName),
		ShardId:    aws.String(shardID),
	}

	if lastSeq, err := c.Checkpoint.Get(c.StreamName, shardID); err == nil {
		params.ShardIteratorType = aws.String("AFTER_SEQUENCE_NUMBER")
		params.StartingSequenceNumber = aws.String(lastSeq)
	} else {
		params.ShardIteratorType = aws.String("TRIM_HORIZON") //Read from beginning of the stream
	}

	resp, err := c.svc.GetShardIterator(params)
	if err != nil {
		c.Logger.Panicf("Error GetShardIterator %v", err.Error())
	}

	return resp.ShardIterator
}
