// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package services

import (
	"bufio"
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/logger"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Bitcoind struct {
	sync.RWMutex
	initialised bool
	log         *logger.L
	process     *os.Process
	running     bool
	ModeStart   chan bool
}

func (bitcoind *Bitcoind) Initialise() error {
	bitcoind.Lock()
	defer bitcoind.Unlock()

	if bitcoind.initialised {
		return fault.ErrAlreadyInitialised
	}

	bitcoind.log = logger.New("service-bitcoind")
	if nil == bitcoind.log {
		return fault.ErrInvalidLoggerChannel
	}

	bitcoind.running = false
	bitcoind.ModeStart = make(chan bool, 1)

	// all data initialised
	bitcoind.initialised = true
	return nil
}

func (bitcoind *Bitcoind) Finalise() error {
	bitcoind.Lock()
	defer bitcoind.Unlock()

	if !bitcoind.initialised {
		return fault.ErrNotInitialised
	}

	bitcoind.initialised = false
	return nil
}

func (bitcoind *Bitcoind) IsRunning() bool {
	return bitcoind.running
}

func (bitcoind *Bitcoind) Run(args interface{}, shutdown <-chan struct{}) {
loop:
	for {
		select {

		case <-shutdown:
			break loop
		case start := <-bitcoind.ModeStart:
			if start {
				bitcoind.startBitcoind()

			} else {
				bitcoind.stopBitcoind()
			}
		}

	}
	close(bitcoind.ModeStart)
}

func (bitcoind *Bitcoind) startBitcoind() error {
	if bitcoind.running {
		bitcoind.log.Errorf("Start bitcoind failed: %v", fault.ErrBitcoindIsRunning)
		return fault.ErrBitcoindIsRunning
	}

	// start bitcoind as sub process
	cmd := exec.Command("bitcoind")
	// start bitcoind as sub process
	stderr, err := cmd.StderrPipe()
	if err != nil {
		bitcoind.log.Errorf("Error: %v", err)
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		bitcoind.log.Errorf("Error: %v", err)
		return err
	}
	if err := cmd.Start(); nil != err {
		bitcoind.log.Errorf("Start bitcoind failed: %v", err)
		return err
	}

	bitcoind.running = true
	bitcoind.process = cmd.Process
	bitcoind.log.Infof("process id: %d", cmd.Process.Pid)

	go func() {

		stdeReader := bufio.NewReader(stderr)
		stdoReader := bufio.NewReader(stdout)
		stderrDone := make(chan bool, 1)
		stdoutDone := make(chan bool, 1)

		go func() {
			defer close(stderrDone)
			for {
				stde, err := stdeReader.ReadString('\n')
				bitcoind.log.Errorf("bitcoind stderr: %q", stde)
				if nil != err {
					bitcoind.log.Errorf("Error: %v", err)
					return
				}
			}
		}()

		go func() {
			defer close(stdoutDone)
			for {
				stdo, err := stdoReader.ReadString('\n')
				bitcoind.log.Infof("bitcoind stdout: %q", stdo)
				if nil != err {
					bitcoind.log.Errorf("Error: %v", err)
					return
				}
			}
		}()

		<-stderrDone
		<-stdoutDone
		if err := cmd.Wait(); nil != err {
			bitcoind.log.Errorf("Start bitcoind failed: %v", err)
			bitcoind.running = false
			bitcoind.process = nil
		}
	}()

	// wait for 1 second if cmd has no error then return nil
	time.Sleep(time.Second * 1)
	return nil
}

func (bitcoind *Bitcoind) stopBitcoind() error {
	if !bitcoind.running {
		bitcoind.log.Errorf("Stop bitcoind failed: %v", fault.ErrBitcoindIsNotRunning)
		return fault.ErrBitcoindIsNotRunning
	}

	if err := bitcoind.process.Signal(os.Interrupt); nil != err {
		bitcoind.log.Errorf("Send interrupt to bitcoind failed: %v", err)
		if err := bitcoind.process.Signal(os.Kill); nil != err {
			bitcoind.log.Errorf("Send kill to bitcoind failed: %v", err)
			return err
		}
	}

	bitcoind.log.Infof("Stop bitcoind. PID: %d", bitcoind.process.Pid)
	bitcoind.running = false
	bitcoind.process = nil
	return nil
}

func (bitcoind *Bitcoind) GetInfo() ([]byte, error) {
	out, err := exec.Command("bitcoin-cli", "getinfo").Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}
