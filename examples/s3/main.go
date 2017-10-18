package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/telenor-digital-asia/kinesis-connectors"
	"github.com/telenor-digital-asia/kinesis-connectors/emitter/s3"
)

func main() {

	var (
		app    = flag.String("a", "", "App name")
		bucket = flag.String("b", "", "Bucket name")
		stream = flag.String("s", "", "Stream name")
	)
	flag.Parse()

	e := &s3.Emitter{
		Bucket: *bucket,
		Region: "us-west-1",
	}

	c := connector.NewConsumer(connector.Config{
		AppName:    *app,
		StreamName: *stream,
	})

	c.Start(connector.HandlerFunc(func(b connector.Buffer) {
		body := new(bytes.Buffer)

		for _, r := range b.GetRecords() {
			body.Write(r.Data)
		}

		err := e.Emit(
			s3.Key("", b.FirstSeq(), b.LastSeq()),
			bytes.NewReader(body.Bytes()),
			"",
		)

		if err != nil {
			fmt.Printf("error %s\n", err)
			os.Exit(1)
		}
	}))

	select {} // run forever
}
