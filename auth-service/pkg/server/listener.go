package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type listener struct {
	Addr     string `json:"addr"`
	Fd       int    `json:"fd"`
	Filename string `json:"filename"`
}

var ListenerEnvKey string = "LISTENER_CONFIG"

func importListener(addr string) (net.Listener, error) {
	listenerConfig := os.Getenv(ListenerEnvKey)
	if listenerConfig == "" {
		return nil, fmt.Errorf("unable to find listener config")
	}

	var l listener

	err := json.NewDecoder(strings.NewReader(listenerConfig)).Decode(&l)
	if err != nil {
		return nil, err
	}

	if l.Addr != addr {
		return nil, fmt.Errorf("unable to find address listener, got %s", l.Addr)
	}

	listenerFile := os.NewFile(uintptr(l.Fd), l.Filename)
	if listenerFile == nil {
		return nil, fmt.Errorf("unable to create file listener")
	}

	defer listenerFile.Close()

	ln, err := net.FileListener(listenerFile)
	if err != nil {
		return nil, err
	}

	return ln, nil

}

func getOrNewListener(addr string) (net.Listener, error) {
	ln, err := importListener(addr)
	if err != nil {
		return nil, err
	}

	if ln == nil {
		ln, err = net.Listen("tcp", addr)
	}

	if err != nil {
		return nil, err
	}

	return ln, nil

}

func getListenerFile(ln net.Listener) (*os.File, error) {
	switch t := ln.(type) {
	case *net.TCPListener:
		return t.File()
	case *net.UnixListener:
		return t.File()
	}
	return nil, fmt.Errorf("unsupported listener: %T", ln)
}

func forkChild(addr string, ln net.Listener) (*os.Process, error) {
	lnFile, err := getListenerFile(ln)
	if err != nil {
		return nil, err
	}

	defer lnFile.Close()

	listenerConfig := &listener{
		Addr:     addr,
		Fd:       3,
		Filename: lnFile.Name(),
	}
	var b bytes.Buffer

	err = json.NewEncoder(&b).Encode(&listenerConfig)
	if err != nil {
		return nil, err
	}
	files := []*os.File{os.Stdin, os.Stdout, os.Stderr, lnFile}
	environment := append(os.Environ(), fmt.Sprintf("%s=%s", ListenerEnvKey, b.String()))

	execName, err := os.Executable()
	if err != nil {
		return nil, err
	}

	execDir := filepath.Dir(execName)

	p, err := os.StartProcess(execName, []string{execName}, &os.ProcAttr{
		Dir:   execDir,
		Env:   environment,
		Files: files,
		Sys:   &syscall.SysProcAttr{},
	})

	if err != nil {
		return nil, err
	}

	return p, nil

}
