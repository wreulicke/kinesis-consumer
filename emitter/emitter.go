package emitter

import (
	"database/sql"
	"fmt"
	"time"
)

// RedshiftEmitter is an implementation of Emitter that buffered batches of records into Redshift one by one.
// It first emits records into S3 and then perfors the Redshift JSON COPY command. S3 storage of buffered
// data achieved using the S3Emitter. A link to jsonpaths must be provided when configuring the struct.
type RedshiftEmitter struct {
	AwsAccessKey       string
	AwsSecretAccessKey string
	Delimiter          string
	Format             string
	Jsonpaths          string
	S3Bucket           string
	S3Prefix           string
	TableName          string
	Db                 *sql.DB
}

// Key create s3 key structure
func Key(prefix, firstSeq, lastSeq string) string {
	date := time.Now().UTC().Format("2006/01/02")

	if prefix == "" {
		return fmt.Sprintf("%v/%v-%v", date, firstSeq, lastSeq)
	}
	return fmt.Sprintf("%v/%v/%v-%v", prefix, date, firstSeq, lastSeq)
}
