package main

import (
	"os/exec"
	"os"
	"fmt"
	"syscall"
	"path/filepath"
	"io/ioutil"
	"strconv"
)

const rootFsLocation = "/home/martin/oci-con/alpine-minirootfs"

func main() {
	checkArgs()
	switch os.Args[1] {
	case "run":
		//run()
		runContainer()
	case "child":
		child()
	default:
		panic("Usage: ocscratch run <command>")
	}
}

// Run given command as child process
func run() {
	fmt.Printf("Running %v\n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	bindStd(cmd)
	setUpNameSpaces(cmd, syscall.CLONE_NEWUTS)

	orPanic(cmd.Run())
}

func bindStd(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

// Create process attributes with all needed name spaces.
// NEWUTS - Unix Time Sharing system (hostname) to isolate hostname of container
// NEWPID - ProcessID to limit seen processes
// NEWNS - Mounts to limit seen mounts
func setUpNameSpaces(cmd *exec.Cmd, cloneFlags uintptr) {
	cmd.SysProcAttr = &syscall.SysProcAttr {
		Cloneflags: cloneFlags,
	}
}

// Run self as child with applied namespaces (in isolation) as container
func runContainer() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	bindStd(cmd)
	setUpNameSpaces(cmd, syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS)

	orPanic(cmd.Run())
}

// Prepare container environment
// Run given command in container
func child() {
	fmt.Printf("Running %v\n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	bindStd(cmd)

	orPanic(setUpHostName())
	chroot()
	mountProc()
	defer unmountProc()
	mountTmp()
	defer unmountTmp()

	orPanic(cmd.Run())
}

func setUpHostName() error {
	return syscall.Sethostname([]byte("container"))
}

func chroot() {
	orPanic(syscall.Chroot(rootFsLocation))
	orPanic(os.Chdir("/"))
}

func mountProc() {
	orPanic(syscall.Mount("proc", "proc", "proc", 0 , ""))
}

func unmountProc() {
	orPanic(syscall.Unmount("proc", 0))
}

func mountTmp() {
	orPanic(syscall.Mount("container_tmp", "tmp", "tmpfs", 0 , ""))
}

func unmountTmp() {
	orPanic(syscall.Unmount("tmp", 0))
}

// Set up cgroups memory limit for container
func cg() {
	cgroups := "/sys/fs/cgroup/"

	mem := filepath.Join(cgroups, "memory")
	os.Mkdir(filepath.Join(mem, "me"), 0755)
	orPanic(ioutil.WriteFile(filepath.Join(mem, "me/memory.limit_in_bytes"), []byte("999424"), 0700))
	//orPanic(ioutil.WriteFile(filepath.Join(mem, "me/memory.memsw.limit_in_bytes"), []byte("999424"), 0700))
	orPanic(ioutil.WriteFile(filepath.Join(mem, "me/notify_on_release"), []byte("1"), 0700))

	pid := strconv.Itoa(os.Getpid())
	orPanic(ioutil.WriteFile(filepath.Join("mem", "me/cgroup.procs"), []byte(pid), 0700))
}

// Create process attributes with all needed name spaces as non root user.
func getNameSpacesUser() *syscall.SysProcAttr {
	return &syscall.SysProcAttr {
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Credential: &syscall.Credential{Uid: 0, Gid: 0},
		UidMappings: []syscall.SysProcIDMap {
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap {
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
	}
}

func checkArgs() {
	// TODO implement
}

// Panic in case of error
func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}
