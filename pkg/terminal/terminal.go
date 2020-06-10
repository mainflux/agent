package terminal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/creack/pty"

	"github.com/mainflux/agent/pkg/encoder"
	"github.com/mainflux/mainflux/errors"
	"github.com/mainflux/mainflux/logger"
)

const (
	terminal = "term"
	second   = time.Duration(1 * time.Second)
)

var (
	errTerminalSessionStart = errors.New("failed to start terminal session")
)

type term struct {
	uuid         string
	ptmx         *os.File
	writer       io.Writer
	done         chan bool
	topic        string
	timeout      time.Duration
	resetTimeout time.Duration
	timer        *time.Ticker
	publish      func(channel, payload string) error
	logger       logger.Logger
	mu           sync.Mutex
}

type Session interface {
	Send(p []byte) error
	IsDone() chan bool
	io.Writer
}

func NewSession(uuid string, timeout time.Duration, publish func(channel, payload string) error, logger logger.Logger) (Session, error) {
	t := &term{
		logger:       logger,
		uuid:         uuid,
		publish:      publish,
		timeout:      timeout,
		resetTimeout: timeout,
		topic:        fmt.Sprintf("term/%s", uuid),
		done:         make(chan bool),
	}

	c := exec.Command("bash")
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

	t.timer = time.NewTicker(1 * time.Second)

	go func() {
		for range t.timer.C {
			t.decrementCounter()
		}
		t.logger.Debug("exiting timer routine")
	}()

	return t, nil
}

func (t *term) resetCounter(timeout time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if timeout > 0 {
		t.timeout = timeout
		return
	}
}
func (t *term) decrementCounter() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.timeout = t.timeout - second
	if t.timeout == 0 {
		t.done <- true
		t.timer.Stop()
	}
}

func (t *term) IsDone() chan bool {
	return t.done
}

func (t *term) Write(p []byte) (int, error) {
	t.resetCounter(t.resetTimeout)
	n := len(p)
	payload, err := encoder.EncodeSenML(t.uuid, terminal, string(p))
	if err != nil {
		return n, err
	}

	if err := t.publish(t.topic, string(payload)); err != nil {
		return n, err
	}
	return n, nil
}

func (t *term) Send(p []byte) error {
	in := bytes.NewReader(p)
	nr, err := io.Copy(t.ptmx, in)
	t.logger.Debug(fmt.Sprintf("Written to ptmx: %d", nr))
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}
