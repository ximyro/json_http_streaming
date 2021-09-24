package adapters

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/textproto"
	cr "streaming/pkg/chunked_reader"
	"time"
)

type EventsClient struct {
	conn         net.Conn
	hostWithPort string
	responseBuf  *bufio.Reader
}

func NewEventsClient(hostWithPort string) *EventsClient {
	return &EventsClient{
		hostWithPort: hostWithPort,
	}
}

var requestData []string = []string{
	"GET /api/v1/events HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\nTransfer-Encoding: chunked\r\n\r\n1\r\na\r\n",
	"1\r\nb\r\n0\r\n\r\n",
}

func (e *EventsClient) createRequest() (err error) {
	e.conn, err = net.DialTimeout("tcp", e.hostWithPort, 50*time.Millisecond)
	if err != nil {
		return fmt.Errorf("can't connect %w", err)
	}

	for _, frame := range requestData {
		if n, err := e.conn.Write([]byte(frame)); err != nil {
			return fmt.Errorf("can't make a request %w", err)
		} else if n != len(frame) {
			return fmt.Errorf("short write")
		}
	}
	e.responseBuf = bufio.NewReader(e.conn)
	tr := textproto.NewReader(e.responseBuf)
	// We need to read the status first
	_, err = tr.ReadLine()
	if err != nil {
		panic(err)
	}

	// Read HTTP sheaders
	_, err = tr.ReadMIMEHeader()
	if err != nil {
		panic(err)
	}

	return nil
}

func (e *EventsClient) close() error {
	return e.conn.Close()
}

func (e *EventsClient) ReadEvents(ctx context.Context, handleFunc func([]byte) error) error {
	err := e.createRequest()
	if err != nil {
		return fmt.Errorf("can't create a request: %v", err)
	}
	defer e.close()
	reader := cr.NewChunkedReader(e.responseBuf)

	for {
		select {
		case <-ctx.Done():
			break
		default:
			data, err := reader.ReadChunk()
			if err == nil {
				err = handleFunc(data)
				if err != nil {
					return fmt.Errorf("can't process data: %v", data)
				}
				continue
			}
			switch err {
			case cr.ErrEmpyLine:
				continue
			case io.EOF, io.ErrUnexpectedEOF:
				return nil
			}
		}
	}
}
