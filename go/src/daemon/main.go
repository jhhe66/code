// daemon project main.go
package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func Log(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func Daemonize(nochdir, noclose int) error {

	var ret, ret2 uintptr
	var err error

	Log(fmt.Sprintf("common.Daemonize][current ppid:%d", syscall.Getppid()))

	//already a daemon
	if syscall.Getppid() == 1 {
		return nil
	}

	Log("common.Daemonize][will daemonize")

	//fork off the parent process
	ret, ret2, err = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
	if err != nil {
		var s string
		s = fmt.Sprintf("%T", err)
		if s != "syscall.Errno" {
			Log("common.Daemonize][fork err:" + s)
			return err
		} else {
			Log("common.Daemonize][fork no err:" + err.Error())
			//no problem see http://www.ibm.com/developerworks/aix/library/au-errnovariable/
		}
	}

	//failure
	if ret2 < 0 {
		os.Exit(-1)
	}

	//if we got a good PID, then we call exit the parent process.
	if ret > 0 {
		os.Exit(0)
	}
	Log("common.Daemonize][forked,we in forked process")

	/* Change the file mode mask */
	_ = syscall.Umask(0)
	Log("common.Daemonize][umask zero")

	//create a new SID for the child process
	s_ret, s_errno := syscall.Setsid()
	if s_ret < 0 || s_errno != nil {
		log.Printf("common.Daemonize][Error: syscall.Setsid errno: %d", s_errno)
	}
	Log("common.Daemonize][sid seted")

	if nochdir == 0 {
		os.Chdir("/")
		Log("common.Daemonize][chdir to root")
	}

	if noclose == 0 {
		f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
		if e == nil {
			fd := f.Fd()
			syscall.Dup2(int(fd), int(os.Stdin.Fd()))
			syscall.Dup2(int(fd), int(os.Stdout.Fd()))
			syscall.Dup2(int(fd), int(os.Stderr.Fd()))
			Log("common.Daemonize][fs closed")
		}
	}

	Log("common.Daemonize][daemonize done")
	return nil
}

func main() {
	if e := Daemonize(1, 0); e == nil {

	} else {
		Log("daemon error\n")
	}
}
