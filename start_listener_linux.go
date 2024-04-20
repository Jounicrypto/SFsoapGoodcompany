// +build !windows,!wasm

package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"syscall"

	"golang.org/x/sys/unix"
)

func startListener(bindip string, port int) net.Listener {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var opErr error
			err := c.Control(func(fd uintptr) {
				opErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			})
			if err != nil {
				return err
			}
			if err = c.Control(func(fd uintptr) {
				opErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
			}); err != nil {
				return err
			}
			return opErr
		},
	}
	lis, err := lc.Listen(context.Background(), "tcp", bindip+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("could not open socket ", bindip, ":", port, " error ", err)
	}
	log.Printf("Listener on: %s\n", lis.Addr())
	return lis
}

func startUds(path string) net.Listener {
	uLis, err := net.Listen("unix", path)
	if err != nil {
		log.Fatal("could not open path", path, " error ", err)
	}
	log.Printf("Listener on: %s\n", uLis.Addr())
	return uLis
}
