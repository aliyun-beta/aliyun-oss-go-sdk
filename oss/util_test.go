package oss

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

const (
	expectBut         = "\n--- EXPECT ---\n%v\n--- BUT GOT ---\n%v"
	testcaseExpectBut = "testcase %d:\n--- EXPECT ---\n%v\n--- BUT GOT ---\n%v"
)

type response struct {
	headers []string
	body    string
}

func newResponse(resp string) *response {
	lines := strings.Split(resp, "\n")
	isHeader := true
	var headers []string
	var body string
	for _, line := range lines {
		if line == "" {
			isHeader = false
		}
		if isHeader {
			headers = append(headers, line)
		} else {
			body += line + "\n"
		}
	}
	return &response{
		headers: headers,
		body:    body,
	}
}

func (r *response) Write(w io.Writer) error {
	for _, header := range r.headers {
		w.Write([]byte(header))
		w.Write([]byte{'\r', '\n'})
	}
	w.Write([]byte(r.body))
	return nil
}

type MockServer struct {
	listener net.Listener
	conn     net.Conn
	Request  string
	Err      error
	resp     *response
	wg       sync.WaitGroup
}

func NewMockServer(resp string) (*MockServer, error) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}
	rec := &MockServer{listener: lis, resp: newResponse(resp)}
	rec.wg.Add(1)
	go rec.listen()
	return rec, nil
}

func (r *MockServer) URL() string {
	return r.listener.Addr().String()
}

func (r *MockServer) listen() {
	defer r.wg.Done()
	var err error
	r.conn, err = r.listener.Accept()
	if err != nil {
		r.Err = err
		return
	}
	request := make([]byte, 1024)
	n, err := r.conn.Read(request)
	if err != nil {
		r.Err = err
		return
	}
	req := request[:n]
	req = bytes.Replace(req, []byte{'\r'}, nil, -1)
	req = bytes.TrimSpace(req)
	r.Request = string(req)
	r.resp.Write(r.conn)
}

func (r *MockServer) Wait() {
	r.wg.Wait()
}

func (r *MockServer) Close() {
	r.conn.Close()
	r.listener.Close()
}

func p(v ...interface{}) {
	fmt.Println(v...)
}
