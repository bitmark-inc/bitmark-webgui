// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package services

import (
	"crypto/tls"
	"github.com/bitmark-inc/bitmark-mgmt/fault"
	"github.com/bitmark-inc/bitmark-mgmt/utils"
	"github.com/bitmark-inc/bitmarkd/rpc"
	"github.com/bitmark-inc/logger"
	"io/ioutil"
	"net"
	netrpc "net/rpc"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Bitmarkd struct {
	sync.RWMutex
	initialised bool
	log         *logger.L
	configFile  string
	process     *os.Process
	running     bool
	ModeStart   chan bool
}

func (bitmarkd *Bitmarkd) Initialise(configFile string) error {
	bitmarkd.Lock()
	defer bitmarkd.Unlock()

	if bitmarkd.initialised {
		return fault.ErrAlreadyInitialised
	}

	bitmarkd.configFile = configFile

	bitmarkd.log = logger.New("service-bitmarkd")
	if nil == bitmarkd.log {
		return fault.ErrInvalidLoggerChannel
	}

	bitmarkd.running = false
	bitmarkd.ModeStart = make(chan bool, 1)

	// all data initialised
	bitmarkd.initialised = true
	return nil
}

func (bitmarkd *Bitmarkd) Finalise() error {
	bitmarkd.Lock()
	defer bitmarkd.Unlock()

	if !bitmarkd.initialised {
		return fault.ErrNotInitialised
	}

	bitmarkd.initialised = false
	return nil
}

func (bitmarkd *Bitmarkd) IsRunning() bool {
	return bitmarkd.running
}

func (bitmarkd *Bitmarkd) BitmarkdBackground(args interface{}, shutdown <-chan bool, finished chan<- bool) {
loop:
	for {
		select {

		case <-shutdown:
			break loop
		case start := <-bitmarkd.ModeStart:
			if start {
				if err := bitmarkd.startBitmarkd(); nil != err {
					bitmarkd.log.Errorf("Start bitmarkd failed: %v", err)
				}

			} else {
				if err := bitmarkd.stopBitmarkd(); nil != err {
					bitmarkd.log.Errorf("Stop bitmarkd failed: %v", err)
				}
			}
		}

	}
	close(bitmarkd.ModeStart)
	close(finished)
}

func (bitmarkd *Bitmarkd) startBitmarkd() error {
	if bitmarkd.running {
		return nil
	}

	// Check bitmarkConfigFile exists
	if !utils.EnsureFileExists(bitmarkd.configFile) {
		return fault.ErrNotFoundConfigFile
	}

	// start bitmarkd as sub process
	cmd := exec.Command("bitmarkd", "--config-file="+bitmarkd.configFile)
	// start bitmarkd as sub process
	stderr, err := cmd.StderrPipe()
	if err != nil {
		bitmarkd.log.Errorf("Error: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		bitmarkd.log.Errorf("Error: %v", err)
	}
	if err := cmd.Start(); nil != err {
		return err
	}

	bitmarkd.running = true
	bitmarkd.process = cmd.Process
	bitmarkd.log.Infof("process id: %d", cmd.Process.Pid)
	go func() {
		stde, err := ioutil.ReadAll(stderr)
		if nil != err {
			bitmarkd.log.Errorf("Error: %v", err)
		}

		stdo, err := ioutil.ReadAll(stdout)
		if nil != err {
			bitmarkd.log.Errorf("Error: %v", err)
		}

		bitmarkd.log.Errorf("bitmarkd stderr: %s", stde)
		bitmarkd.log.Infof("bitmarkd stdout: %s", stdo)

		if err := cmd.Wait(); nil != err {
			bitmarkd.log.Errorf("Start bitmarkd failed: %v", err)
			bitmarkd.running = false
			bitmarkd.process = nil
		}
	}()

	// wait for 1 second if cmd has no error then return nil
	time.Sleep(time.Second * 1)
	return nil

}

func (bitmarkd *Bitmarkd) stopBitmarkd() error {
	if !bitmarkd.running {
		return nil
	}

	if err := bitmarkd.process.Signal(os.Interrupt); nil != err {
		bitmarkd.log.Errorf("Send interrupt to bitmarkd failed: %v", err)
		if err := bitmarkd.process.Signal(os.Kill); nil != err {
			bitmarkd.log.Errorf("Send kill to bitmarkd failed: %v", err)
			return err
		}
	}

	bitmarkd.log.Infof("Stop bitmarkd. PID: %d", bitmarkd.process.Pid)
	bitmarkd.running = false
	bitmarkd.process = nil
	return nil
}

func (bitmarkd *Bitmarkd) GetInfo(client *netrpc.Client) (*rpc.InfoReply, error) {

	var reply rpc.InfoReply
	if err := client.Call("Node.Info", rpc.InfoArguments{}, &reply); err != nil {
		bitmarkd.log.Errorf("Node.Info error: %v\n", err)
		return nil, fault.ErrNodeInfoRequestFail
	}

	return &reply, nil
}

// connect to bitmarkd RPC
func (bitmarkd *Bitmarkd) Connect(connect string) (net.Conn, error) {

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", connect, tlsConfig)
	if nil != err {
		return nil, err
	}

	return conn, nil
}
