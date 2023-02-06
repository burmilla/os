// Copyright 2019 Wataru Ishida. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sctp

import (
	"fmt"
	"io"
	"math/rand"
	"sync"
	"testing"
	"time"
)

const (
	STREAM_TEST_CLIENTS = 128
	STREAM_TEST_STREAMS = 11
)

func TestStreams(t *testing.T) {
	var rMu sync.Mutex
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomStr := func(strlen int) string {
		const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		result := make([]byte, strlen)
		rMu.Lock()
		for i := range result {
			result[i] = chars[r.Intn(len(chars))]
		}
		rMu.Unlock()
		return string(result)
	}

	addr, _ := ResolveSCTPAddr("sctp", "127.0.0.1:0")
	ln, err := ListenSCTPExt("sctp", addr, InitMsg{NumOstreams: STREAM_TEST_STREAMS, MaxInstreams: STREAM_TEST_STREAMS})
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	addr = ln.Addr().(*SCTPAddr)
	t.Logf("Listen on %s", ln.Addr())

	go func() {
		for {
			c, err := ln.Accept()
			sconn := c.(*SCTPConn)
			if err != nil {
				t.Errorf("failed to accept: %v", err)
				return
			}
			defer sconn.Close()

			sconn.SubscribeEvents(SCTP_EVENT_DATA_IO)
			go func() {
				totalrcvd := 0
				for {
					buf := make([]byte, 512)
					n, info, err := sconn.SCTPRead(buf)
					if err != nil {
						if err == io.EOF || err == io.ErrUnexpectedEOF {
							if n == 0 {
								break
							}
							t.Logf("EOF on server connection. Total bytes received: %d, bytes received: %d", totalrcvd, n)
						} else {
							t.Errorf("Server connection read err: %v. Total bytes received: %d, bytes received: %d", err, totalrcvd, n)
							return
						}
					}
					t.Logf("server read: info: %+v, payload: %s", info, string(buf[:n]))
					n, err = sconn.SCTPWrite(buf[:n], info)
					if err != nil {
						t.Error(err)
						return
					}
				}
			}()
		}
	}()

	wait := make(chan struct{})
	i := 0
	for ; i < STREAM_TEST_CLIENTS; i++ {
		go func(test int) {
			defer func() { wait <- struct{}{} }()
			conn, err := DialSCTPExt(
				"sctp", nil, addr, InitMsg{NumOstreams: STREAM_TEST_STREAMS, MaxInstreams: STREAM_TEST_STREAMS})
			if err != nil {
				t.Errorf("failed to dial address %s, test #%d: %v", addr.String(), test, err)
				return
			}
			defer conn.Close()
			conn.SubscribeEvents(SCTP_EVENT_DATA_IO)
			for ppid := uint16(0); ppid < STREAM_TEST_STREAMS; ppid++ {
				info := &SndRcvInfo{
					Stream: uint16(ppid),
					PPID:   uint32(ppid),
				}
				rMu.Lock()
				randomLen := r.Intn(255)
				rMu.Unlock()
				text := fmt.Sprintf("Test %s ***\n\t\t%d %d ***", randomStr(randomLen), test, ppid)
				n, err := conn.SCTPWrite([]byte(text), info)
				if err != nil {
					t.Errorf("failed to write %s, len: %d, err: %v, bytes written: %d", text, len(text), err, n)
					return
				}
				rn := 0
				cn := 0
				buf := make([]byte, 512)
				for {
					cn, info, err = conn.SCTPRead(buf[rn:])
					if err != nil {
						if err == io.EOF || err == io.ErrUnexpectedEOF {
							rn += cn
							break
						}
						t.Errorf("failed to read: %v", err)
						return
					}
					if info.Stream != ppid {
						t.Errorf("Mismatched PPIDs: %d != %d", info.Stream, ppid)
						return
					}
					rn += cn
					if rn >= n {
						break
					}
				}
				rtext := string(buf[:rn])
				if rtext != text {
					t.Fatalf("Mismatched payload: %s != %s", rtext, text)
				}
			}
		}(i)
	}
	for ; i > 0; i-- {
		select {
		case <-wait:
		case <-time.After(time.Second * 30):
			close(wait)
			t.Fatal("timed out")
		}
	}
}
