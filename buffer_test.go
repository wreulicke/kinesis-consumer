package connector

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/kinesis"
)

func BenchmarkBufferLifecycle(b *testing.B) {
	buf := Buffer{MaxRecordCount: 1000}
	seq := "1"
	rec := &kinesis.Record{SequenceNumber: &seq}

	for i := 0; i < b.N; i++ {
		buf.AddRecord(rec)

		if buf.ShouldFlush() {
			buf.Flush()
		}
	}
}

func Test_FirstSeq(t *testing.T) {
	b := Buffer{}
	s1, s2 := "1", "2"
	r1 := &kinesis.Record{SequenceNumber: &s1}
	r2 := &kinesis.Record{SequenceNumber: &s2}

	b.AddRecord(r1)
	if b.FirstSeq() != "1" {
		t.Fail()
	}

	b.AddRecord(r2)
	if b.FirstSeq() != "1" {
		t.Fail()
	}
}

func Test_LastSeq(t *testing.T) {
	b := Buffer{}
	s1, s2 := "1", "2"
	r1 := &kinesis.Record{SequenceNumber: &s1}
	r2 := &kinesis.Record{SequenceNumber: &s2}

	b.AddRecord(r1)
	if b.LastSeq() != "1" {
		t.Fail()
	}

	b.AddRecord(r2)
	if b.LastSeq() != "2" {
		t.Fail()
	}
}

func Test_ShouldFlush(t *testing.T) {
	b := Buffer{MaxRecordCount: 2}
	s1, s2 := "1", "2"
	r1 := &kinesis.Record{SequenceNumber: &s1}
	r2 := &kinesis.Record{SequenceNumber: &s2}

	b.AddRecord(r1)
	if b.ShouldFlush() != false {
		t.Fail()
	}
	b.AddRecord(r2)
	if b.ShouldFlush() == false {
		t.Fail()
	}
}
