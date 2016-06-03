package services

import (
	"fmt"
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/logger"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
	"time"
)

var (
	cmd     *exec.Cmd
	cmdType string
)

func generateBitmarkPay(cmdStr string, s ...string) BitmarkPay {
	bitmarkPay := BitmarkPay{
		initialised: true,
		log:         logger.New("service-bitmarkPay"),
	}
	if nil == bitmarkPay.log {
		fmt.Println(fault.ErrInvalidLoggerChannel)
	}

	var cs []string
	cs = append(cs, s...)
	cmd = exec.Command(cmdStr, cs...)
	cmdType = "test"

	return bitmarkPay
}

func TestLocalStatus(t *testing.T) {
	localBitmarkPay := generateBitmarkPay("pwd", "-L")
	localBitmarkPay.asyncJob.hash = "12345678"
	localBitmarkPay.asyncJob.cmdType = "test"

	// cmd is empty, status should return stopped
	assert.Equal(t, "stopped", localBitmarkPay.status())

	// cmd is not empty, should return running
	localCmd := exec.Command("pwd", "-L")
	localBitmarkPay.asyncJob.command = localCmd
	assert.Equal(t, "stopped", localBitmarkPay.status())

	// cmd is complete
	if err := localCmd.Run(); nil != err {
		assert.Fail(t, "fail to run cmd: pwd -L")
	} else {
		assert.Equal(t, "success", localBitmarkPay.status())
	}

	// cmd is comptet with error, should return fail
	localCmd = exec.Command("cat", "/notExist")
	localBitmarkPay.asyncJob.command = localCmd
	if err := localCmd.Run(); nil != err {
		assert.Equal(t, "fail", localBitmarkPay.status())
	} else {
		assert.Fail(t, "should not pass to run cmd: cat /notExist")
	}

	// cmd is invalid, should return stopped
	localCmd = exec.Command("pwdsss")
	localBitmarkPay.asyncJob.command = localCmd
	if err := localCmd.Run(); nil != err {
		assert.Equal(t, "stopped", localBitmarkPay.status())
	} else {
		assert.Fail(t, "should not pass to run cmd: pwdsss")
	}
}

func TestRunBitmarkPayJob(t *testing.T) {
	bitmarkPay := generateBitmarkPay("pwd", "-L")
	if err := bitmarkPay.runBitmarkPayJob(cmd, cmdType); nil != err {
		t.Errorf("runBitmarkPayJob get hash error: %v", err)
	}
	assert.NotEmpty(t, bitmarkPay.asyncJob.hash, "job hash is empty")
	assert.NotNil(t, bitmarkPay.asyncJob.command, "asynjob command is nil")
	assert.NotEmpty(t, bitmarkPay.asyncJob.cmdType, "asyncjob command type is empty")
	assert.Equal(t, cmdType, bitmarkPay.asyncJob.cmdType, "asyncjob command type is incorrect")
}

func TestGetBitmarkPayJobHash(t *testing.T) {
	bitmarkPay := generateBitmarkPay("pwd", "-L")
	if err := bitmarkPay.runBitmarkPayJob(cmd, cmdType); nil != err {
		t.Errorf("runBitmarkPayJob get hash error: %v", err)
	}
	assert.NotEmpty(t, bitmarkPay.GetBitmarkPayJobHash(), "bitmarkPay jobHash is empty")
}

func TestGetBitmarkPayJobType(t *testing.T) {
	bitmarkPay := generateBitmarkPay("pwd", "-L")
	if err := bitmarkPay.runBitmarkPayJob(cmd, cmdType); nil != err {
		t.Errorf("runBitmarkPayJob get hash error: %v", err)
	}
	assert.Equal(t, cmdType, bitmarkPay.GetBitmarkPayJobType(bitmarkPay.asyncJob.hash), "asyncjob command type is incorrect")
}

func TestGetBitmarkPayJobResult(t *testing.T) {
	bitmarkPay := generateBitmarkPay("pwd", "-L")
	if err := bitmarkPay.runBitmarkPayJob(cmd, cmdType); nil != err {
		t.Errorf("runBitmarkPayJob get hash error: %v", err)
	}

	// Test for empty hash
	bitmarkPayType := BitmarkPayType{}
	result, err := bitmarkPay.GetBitmarkPayJobResult(bitmarkPayType)
	assert.Nil(t, result, "in running status, result should be nil")
	assert.Equal(t, fault.ErrInvalidCommandParams, err, "not exit while jobHash is invalid")

	// Test when pass wrong hash
	bitmarkPayType = BitmarkPayType{
		JobHash: "12345678",
	}
	result, err = bitmarkPay.GetBitmarkPayJobResult(bitmarkPayType)
	assert.Nil(t, result, "in running status, result should be nil")
	assert.Equal(t, fault.ErrNotFoundBitmarkPayJob, err, "not exit while the hash is mismatched")

	// Test in running
	bitmarkPayType = BitmarkPayType{
		JobHash: bitmarkPay.asyncJob.hash,
	}
	result, err = bitmarkPay.GetBitmarkPayJobResult(bitmarkPayType)
	assert.Nil(t, result, "in running status, result should be nil")
	assert.Equal(t, fault.ErrInvalidAccessBitmarkPayJobResult, err, "not exit while test status equals running or stopped ")

loop:
	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Printf("status: %s, processState: %s\n", bitmarkPay.status(), bitmarkPay.asyncJob.command.ProcessState.String())
			if bitmarkPay.status() == "success" {
				bitmarkPayType = BitmarkPayType{
					JobHash: bitmarkPay.asyncJob.hash,
				}
				result, err = bitmarkPay.GetBitmarkPayJobResult(bitmarkPayType)
				if nil == bitmarkPay.asyncJob.result {
					assert.Nil(t, result, "in running status, result should be nil")
					assert.Equal(t, fault.ErrExecBitmarkPayJob, err, "not exit while the hash is mismatched")
				} else {
					assert.Nil(t, err, "in running status, result should be nil")
					assert.NotNil(t, result, "result is nil")
				}
				t.Logf("result: %s", string(result))
				break loop
			} else {
				t.Logf("status: %s", bitmarkPay.status())
			}
		}
	}
}
