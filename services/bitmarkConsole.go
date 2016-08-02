// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package services

import (
	"bufio"
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/bitmark-webgui/utils"
	"github.com/bitmark-inc/logger"
	"os/exec"
	"sync"
	"time"
)

type BitmarkConsole struct {
	sync.RWMutex
	initialised bool
	bin         string
	log         *logger.L
	cmd         *exec.Cmd
	port        string
}

func (bitmarkConsole *BitmarkConsole) Initialise(binFile string) error {
	bitmarkConsole.Lock()
	defer bitmarkConsole.Unlock()

	if bitmarkConsole.initialised {
		return fault.ErrAlreadyInitialised
	}

	bitmarkConsole.log = logger.New("service-bitmarkConsole")
	if nil == bitmarkConsole.log {
		return fault.ErrInvalidLoggerChannel
	}

	// Check bitmarkConsole bin exists
	if !utils.EnsureFileExists(binFile) {
		bitmarkConsole.log.Errorf("cannot find bitmarkConsole bin: %s", binFile)
		return fault.ErrNotFoundBinFile
	}
	bitmarkConsole.bin = binFile
	bitmarkConsole.port = "2160"

	bitmarkConsole.initialised = true

	return nil
}

func (bitmarkConsole *BitmarkConsole) Finalise() error {
	bitmarkConsole.Lock()
	defer bitmarkConsole.Unlock()

	if !bitmarkConsole.initialised {
		return fault.ErrNotInitialised
	}

	bitmarkConsole.initialised = false
	return nil
}

func (bitmarkConsole *BitmarkConsole) StartBitmarkConsole() error {

	cmd := exec.Command(bitmarkConsole.bin,
		"--address", "localhost",
		"--port", bitmarkConsole.port,
		"--permit-write",
		"--once",
		"bash")

	bitmarkConsole.cmd = cmd

	// start bitmarkConsole
	stderr, err := cmd.StderrPipe()
	if err != nil {
		bitmarkConsole.log.Errorf("Error: %v", err)
		return err
	}

	if err := cmd.Start(); nil != err {
		bitmarkConsole.log.Errorf("Start bitmarkConsole failed: %v", err)
		return err
	}
	bitmarkConsole.log.Infof("process id: %d", cmd.Process.Pid)

	go func() {
		stdeReader := bufio.NewReader(stderr)
		stderrDone := make(chan bool, 1)

		go func() {
			defer close(stderrDone)
			for {
				stde, err := stdeReader.ReadString('\n')
				bitmarkConsole.log.Infof("bitmarkConsole stdout: %q", stde)
				if nil != err {
					bitmarkConsole.log.Errorf("Error: %v", err)
					return
				}
			}
		}()

		<-stderrDone
		if err := cmd.Wait(); nil != err {
			bitmarkConsole.log.Errorf("Start bitmarkConsole failed: %v", err)
			bitmarkConsole.cmd = nil
		}
	}()

	// wait for 1 second if cmd has no error then return nil
	time.Sleep(time.Second * 1)
	if nil == bitmarkConsole.cmd || nil == bitmarkConsole.cmd.Process {
		return fault.ErrBitmarkConsoleIsNotRunning
	}

	return nil
}

func (bitmarkConsole *BitmarkConsole) StopBitmarkConsole() error {
	return bitmarkConsole.cmd.Process.Kill()
}

func (bitmarkConsole *BitmarkConsole) Port() string {
	return bitmarkConsole.port
}

func (bitmarkConsole *BitmarkConsole) IsRunning() bool {
	if nil == bitmarkConsole.cmd {
		return false
	}

	if nil == bitmarkConsole.cmd.ProcessState {
		return true
	}

	return !bitmarkConsole.cmd.ProcessState.Exited()
}
