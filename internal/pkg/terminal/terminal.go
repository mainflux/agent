package terminal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"

	"github.com/mainflux/mainflux/errors"
	"github.com/mainflux/mainflux/logger"
)

var (
	errTerminalSessionStart = errors.New("failed to start terminal session")
)

type term struct {
	ptmx   *os.File
	writer io.Writer
	logger logger.Logger
}

type Session interface {
	Send(p []byte) errors.Error
}

func NewSession(w io.Writer, logger logger.Logger) (Session, errors.Error) {
	t := &term{
		writer: w,
		logger: logger,
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
		n, err := io.Copy(t.writer, t.ptmx)
		if err != nil {
			t.logger.Error(fmt.Sprintf("Error sending data: %s", err))
		}
		t.logger.Debug(fmt.Sprintf("Data being sent: %d", n))
	}()

	return t, nil
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
