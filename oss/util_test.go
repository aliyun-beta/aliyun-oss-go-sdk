package oss

import (
	"bytes"
	"fmt"
	"net"
	"sync"
)

const expectBut = "\n--- EXPECT ---\n%s\n--- BUT GOT ---\n%s"

type RequestRecorder struct {
	listener net.Listener
	Request  string
	Err      error
	wg       sync.WaitGroup
}

func NewRequestRecorder() (*RequestRecorder, error) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	rec := &RequestRecorder{listener: lis}
	rec.wg.Add(1)
	go rec.listen()
	return rec, nil
}

func (r *RequestRecorder) URL() string {
	return r.listener.Addr().String()
}

func (r *RequestRecorder) listen() {
	defer r.wg.Done()
	conn, err := r.listener.Accept()
	if err != nil {
		r.Err = err
		return
	}
	request := make([]byte, 1024)
	n, err := conn.Read(request)
	if err != nil {
		r.Err = err
		return
	}
	req := request[:n]
	req = bytes.Replace(req, []byte{'\r'}, nil, -1)
	req = bytes.TrimSpace(req)
	r.Request = string(req)
	conn.Close()
}

func (r *RequestRecorder) Wait() {
	r.wg.Wait()
}

func (r *RequestRecorder) Close() {
	r.listener.Close()
}

func p(v ...interface{}) {
	fmt.Println(v...)
}
