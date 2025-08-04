package smooth

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	sentinelEnvVar = "SUPERGO_SERVER"
	reloadKey      = "9959D7AD-26C9-4B02-8FB8-0EB84779198B"
)

var (
	superRestartCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "super_restart_count",
			Help: "super restart count",
		},
	)
)

type SuperServer struct {
	Executor    func(fds []*Fd)
	ListenAddrs []string
	child       *child
	files       []*os.File
	retry       int
	logger      *log.Logger

	Config     SuperServerConfig
	pipeWriter *os.File
	pipeReader *os.File
}

type SuperServerConfig struct {
	StopTimeout    time.Duration
	HttpListenAddr string
	ReloadKey      string
	MaxRetry       int
}

func init() {
	prometheus.MustRegister(superRestartCount)
}

func NewSuperServerConfig() SuperServerConfig {
	return SuperServerConfig{
		StopTimeout:    30 * time.Second,
		HttpListenAddr: ":28567",
		ReloadKey:      reloadKey,
		MaxRetry:       3,
	}
}

type Fd struct {
	l   net.Listener
	err error
}

func (f *Fd) Listener() (net.Listener, error) {
	return f.l, f.err
}

type child struct {
	cmd      *exec.Cmd
	exitChan chan error
	spawn    bool
}

func (s *SuperServer) logf(format string, v ...interface{}) {
	if s.logger == nil {
		s.logger = log.New(os.Stderr, "[superserver] ", log.LstdFlags)
	}
	s.logger.Printf(format, v...)
}

func (s *SuperServer) Run() error {
	// 子进程
	if os.Getenv(sentinelEnvVar) != "" {
		var fds []*Fd
		for i := 0; i < len(s.ListenAddrs); i++ {
			fd := &Fd{}
			fileListener := os.NewFile(uintptr(3+i), "")
			if fileListener == nil {
				fd.err = fmt.Errorf("new file listner failed")
			} else {
				l, err := net.FileListener(fileListener)
				if err != nil {
					fd.err = fmt.Errorf("file listener: %s", err.Error())
				}
				err = fileListener.Close()
				if err != nil {
					fd.err = fmt.Errorf("file close: %s", err.Error())
				}
				fd.l = l
			}
			fds = append(fds, fd)
		}
		s.stat()
		s.Executor(fds)
	} else {
		s.pipeReader, s.pipeWriter, _ = os.Pipe()
		e := binary.Write(s.pipeWriter, binary.BigEndian, int32(0))
		if e != nil {
			s.logf("stat pipe init write err: %v", e)
		}

		// 父进程
		_ = s.initListener()

		err := s.startChild()
		if err != nil {
			s.logf("start child: %s", err.Error())
		}
		if s.Config.HttpListenAddr != "" {
			go func() {
				if err := s.handleHttp(); err != nil {
					s.logf("listen http: %s", err.Error())
				}
			}()
		}
		s.handleSignal()
	}

	return nil
}

func (s *SuperServer) initListener() error {
	var files []*os.File
	for _, addr := range s.ListenAddrs {
		var l net.Listener
		var f *os.File
		l, err := net.Listen("tcp", addr)
		if err != nil {
			s.logf("listen %s: %s", addr, err.Error())
			files = append(files, f)
			continue
		}
		f, err = l.(*net.TCPListener).File()
		if err != nil {
			s.logf("file listener %s: %s", addr, err.Error())
		}
		_ = l.Close()
		files = append(files, f)
	}

	s.files = files
	return nil
}

func (s *SuperServer) newChild() (*child, error) {
	cmdPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	env := os.Environ()
	env = append(env, fmt.Sprintf("%s=1", sentinelEnvVar))
	name := os.Args[0]
	s.files = append(s.files, s.pipeReader, s.pipeWriter)
	cmd := &exec.Cmd{
		Dir:        cmdPath,
		Path:       name,
		Args:       os.Args,
		ExtraFiles: s.files, // 传递文件描述符
		Env:        env,
		Stderr:     os.Stderr,
		Stdout:     os.Stdout,
	}
	if filepath.Base(name) == name {
		lp, err := exec.LookPath(name)
		if err != nil {
			return nil, err
		}
		cmd.Path = lp
	}

	return &child{
		cmd:      cmd,
		exitChan: make(chan error),
	}, nil
}

// Restart 平滑重启子进程
func (s *SuperServer) Restart() error {
	oldChi := s.child
	err := s.startChild()
	if err != nil {
		return err
	}
	if oldChi != nil {
		return s.stopChild(oldChi)
	}
	return nil
}

func (s *SuperServer) startChild() error {
	chi, err := s.newChild()
	if err != nil {
		return err
	}
	err = chi.cmd.Start()
	if err != nil {
		return err
	}
	go func() {
		err := chi.cmd.Wait()
		if err != nil {
			chi.exitChan <- err
		}
		close(chi.exitChan)
	}()

	select {
	case <-chi.exitChan:
		if !chi.spawn {
			err := s.shouldRetry()
			return err
		}

	case <-time.After(time.Second):
		s.retry = 0
		s.child = chi
		pid := chi.cmd.Process.Pid
		s.logf("start child %d", pid)
		go func() {
			err := <-chi.exitChan
			s.logf("child exit %d: %v", pid, err)
			if !chi.spawn {
				err := s.shouldRetry()
				s.logf("child %d retry: %v", pid, err)
			}
		}()
	}

	return nil
}

func (s *SuperServer) shouldRetry() error {
	if s.retry >= s.Config.MaxRetry {
		return errors.New("max retry exceed")
	}
	time.Sleep(time.Second * time.Duration(math.Pow(2, float64(s.retry))))
	s.retry++
	s.logf("retry %d", s.retry)
	return s.startChild()
}

func (s *SuperServer) stopChild(chi *child) error {
	pid := chi.cmd.Process.Pid
	chi.spawn = true
	err := chi.cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		s.logf("stop child %d: %s", pid, err.Error())
		return err
	}
	select {
	case err = <-chi.exitChan:
		s.logf("stop child %d", pid)
	case <-time.After(s.Config.StopTimeout):
		s.logf("stop child %d timeout, kill", pid)
		err = chi.cmd.Process.Kill()
	}

	return err
}

func (s *SuperServer) Stop() error {
	if s.child != nil {
		return s.stopChild(s.child)
	}
	return nil
}

func (s *SuperServer) handleSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-sigCh
		s.logf("got signal: %+v", sig)
		switch sig {
		case syscall.SIGHUP:
			err := s.Restart()
			if err != nil {
				s.logf("restart: %s", err.Error())
			}
		case syscall.SIGTERM, syscall.SIGINT:
			_ = s.Stop()
			os.Exit(0)
		}
	}
}

func (s *SuperServer) handleHttp() error {
	var lock sync.Mutex
	mu := http.NewServeMux()
	mu.HandleFunc("/-/reload", func(w http.ResponseWriter, r *http.Request) {
		key := r.FormValue("key")
		if key != s.Config.ReloadKey {
			_, _ = w.Write([]byte("failed"))
			return
		}
		lock.Lock()
		defer lock.Unlock()
		err := s.Restart()
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write([]byte("success"))
	})

	server := &http.Server{
		Handler: mu,
		Addr:    s.Config.HttpListenAddr,
	}

	return server.ListenAndServe()
}

func (s *SuperServer) stat() {
	pr := os.NewFile(uintptr(len(s.ListenAddrs)+3+0), "pipe_reader")
	pw := os.NewFile(uintptr(len(s.ListenAddrs)+3+1), "pipe_writer")
	defer func() {
		_ = pr.Close()
		_ = pw.Close()
	}()

	var size int32
	err := binary.Read(pr, binary.BigEndian, &size)
	if err != nil {
		s.logf("stat pipe read err: %v", err)
		return
	}
	superRestartCount.Set(float64(size))
	size++
	err = binary.Write(pw, binary.BigEndian, size)
	if err != nil {
		s.logf("stat pipe write err: %v", err)
	}
}
