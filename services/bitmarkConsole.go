// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package services

import (
	"bufio"
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/bitmark-webgui/utils"
	"github.com/bitmark-inc/logger"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type BitmarkConsole struct {
	sync.RWMutex
	initialised bool
	bin         string
	crt         string
	key         string
	log         *logger.L
	process     *os.Process
	url         string
}

func (bitmarkConsole *BitmarkConsole) Initialise(binFile string, crtFile string, keyFile string) error {
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
	bitmarkConsole.crt = crtFile
	bitmarkConsole.key = keyFile


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
		"--port", "2160",
		"--permit-write",
		"--random-url",
		"--tls",
		"--tls-crt", bitmarkConsole.crt,
		"--tls-key", bitmarkConsole.key,
		"--once",
		"bash")

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

	bitmarkConsole.process = cmd.Process
	bitmarkConsole.log.Infof("process id: %d", cmd.Process.Pid)

	go func() {

		var url string
		stdeReader := bufio.NewReader(stderr)
		// stdoReader := bufio.NewReader(stdout)
		stderrDone := make(chan bool, 1)
		// stdoutDone := make(chan bool, 1)

		go func() {
			defer close(stderrDone)
			for {
				stde, err := stdeReader.ReadString('\n')
				bitmarkConsole.log.Infof("bitmarkConsole stdout: %q", stde)
				if nil != err {
					bitmarkConsole.log.Errorf("Error: %v", err)
					return
				}
				if "" == url && strings.Contains(string(stde), "URL: http"){
					url = string(stde)
					tmpStrArr := strings.Split(url, ":")
					url = tmpStrArr[len(tmpStrArr)-1]
					bitmarkConsole.url = url
				}
			}
		}()

		<-stderrDone
		if err := cmd.Wait(); nil != err {
			bitmarkConsole.log.Errorf("Start bitmarkConsole failed: %v", err)
			bitmarkConsole.process = nil
		}
	}()

	// wait for 1 second if cmd has no error then return nil
	time.Sleep(time.Second * 1)
	if nil == bitmarkConsole.process {
		return fault.ErrBitmarkConsoleIsNotRunning
	}

	return nil
}

func (bitmarkConsole *BitmarkConsole) GetBitmarkConsoleUrl() string {
	return bitmarkConsole.url
}
