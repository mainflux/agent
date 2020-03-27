package terminal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"

	"github.com/mainflux/agent/internal/pkg/util"
	"github.com/mainflux/mainflux/errors"
	"github.com/mainflux/mainflux/logger"
)

const terminal = "term"

var (
	errTerminalSessionStart = errors.New("failed to start terminal session")
)

type term struct {
	uuid    string
	ptmx    *os.File
	writer  io.Writer
	publish func(channel, payload string) errors.Error
	logger  logger.Logger
}

type Session interface {
	Send(p []byte) errors.Error
	io.Writer
}

func NewSession(uuid string, publish func(channel, payload string) errors.Error, logger logger.Logger) (Session, errors.Error) {
	t := &term{
		logger:  logger,
		uuid:    uuid,
		publish: publish,
	}
	// Prepare the command to execute.
	c := exec.Command("bash")
	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return t, errors.New(err.Error())
	}

	t.ptmx = ptmx
	// Copy output to mqtt
	go func() {
		n, err := io.Copy(t, t.ptmx)
		if err != nil {
			t.logger.Error(fmt.Sprintf("Error sending data: %s", err))
		}
		t.logger.Debug(fmt.Sprintf("Data being sent: %d", n))
	}()

	go func() {
		stderr, err := c.StderrPipe()
		if err != nil {
			t.logger.Error(fmt.Sprintf("Error opening pipe: %s", err))
			return
		}
		n, err := io.Copy(t, stderr)
		if err != nil {
			t.logger.Error(fmt.Sprintf("Error sending data: %s", err))
			return
		}
		t.logger.Debug(fmt.Sprintf("Data being sent: %d", n))
	}()

	return t, nil
}

func (t *term) Write(p []byte) (int, error) {
	n := len(p)
	payload, err := util.EncodeSenML(t.uuid, terminal, string(p))
	if err != nil {
		return n, err
	}
	if err := t.publish("control", string(payload)); err != nil {
		return n, err
	}
	return n, nil
}

func (t *term) Send(p []byte) errors.Error {
	in := bytes.NewReader(p)
	nr, err := io.Copy(t.ptmx, in)
	t.logger.Debug(fmt.Sprintf("Writtern to ptmx: %d", nr))
	//
	if err != nil {
		return errors.New(err.Error())
	}
	return nil

}
