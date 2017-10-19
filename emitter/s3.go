package emitter

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/s3"
)

// NewS3Emitter create new S3Emitter
func NewS3Emitter(bucket, region string) *S3Emitter {
	svc := s3.New(
		session.New(aws.NewConfig().WithMaxRetries(10)),
		aws.NewConfig().WithRegion(region),
	)
	return &S3Emitter{
		S3Svc:  svc,
		Bucket: bucket,
	}
}

// S3Emitter stores data in S3 bucket.
//
// The use of  this struct requires the configuration of an S3 bucket/endpoint. When the buffer is full, this
// struct's Emit method adds the contents of the buffer to S3 as one file. The filename is generated
// from the first and last sequence numbers of the records contained in that file separated by a
// dash. This struct requires the configuration of an S3 bucket and endpoint.
type S3Emitter struct {
	S3Svc  *s3.S3
	Bucket string
}

// EmitWithACL This method allows you to specify ACL
func (e S3Emitter) EmitWithACL(ACL, s3Key string, b io.ReadSeeker) error {
	params := &s3.PutObjectInput{
		ACL:         aws.String(ACL),
		Body:        b,
		Bucket:      aws.String(e.Bucket),
		ContentType: aws.String("text/plain"),
		Key:         aws.String(s3Key),
	}
	_, err := e.S3Svc.PutObject(params)
	return err
}

// Emit is invoked when the buffer is full. This method emits the set of filtered records.
func (e S3Emitter) Emit(s3Key string, b io.ReadSeeker) error {
	params := &s3.PutObjectInput{
		Body:        b,
		Bucket:      aws.String(e.Bucket),
		ContentType: aws.String("text/plain"),
		Key:         aws.String(s3Key),
	}
	_, err := e.S3Svc.PutObject(params)
	return err
}

// NewManifestS3Emitter new manifest emitter
func NewManifestS3Emitter(bucket, region, outputStream string) *ManifestEmitter {
	return &ManifestEmitter{
		emitter:      NewS3Emitter(bucket, region),
		OutputStream: outputStream,
	}
}

// ManifestEmitter An implementation of Emitter that puts event data on S3 file, and then puts the
// S3 file path onto the output stream for processing by manifest application.
type ManifestEmitter struct {
	emitter      *S3Emitter
	OutputStream string
}

// Emit is invoked when the buffer is full. This method emits the set of filtered records.
func (e ManifestEmitter) Emit(s3Key string, b io.ReadSeeker) error {
	// put contents to S3 Bucket
	e.emitter.Emit(s3Key, b)

	// put file path on Kinesis output stream
	params := &kinesis.PutRecordInput{
		Data:         []byte(s3Key),
		PartitionKey: aws.String(s3Key),
		StreamName:   aws.String(e.OutputStream),
	}

	svc := kinesis.New(session.New())
	_, err := svc.PutRecord(params)

	if err != nil {
		return err
	}

	return nil
}
