// +build windows

package core

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procCreateMutex  = kernel32.NewProc("CreateMutexW")
	procReleaseMutex = kernel32.NewProc("ReleaseMutex")
	procCloseHandle  = kernel32.NewProc("CloseHandle")
)

const (
	ERROR_ALREADY_EXISTS = 183
)

var mutexHandle syscall.Handle

// acquireLock uses a Windows named mutex to ensure single instance
func (si *SingleInstance) acquireLock() (bool, error) {
	mutexName, err := syscall.UTF16PtrFromString("Global\\TyrDesktopMutex")
	if err != nil {
		return false, fmt.Errorf("failed to create mutex name: %w", err)
	}

	// CreateMutexW(SECURITY_ATTRIBUTES, BOOL, LPCWSTR)
	ret, _, err := procCreateMutex.Call(
		0,                       // lpMutexAttributes (NULL)
		0,                       // bInitialOwner (FALSE)
		uintptr(unsafe.Pointer(mutexName)), // lpName
	)

	if ret == 0 {
		return false, fmt.Errorf("CreateMutex failed: %w", err)
	}

	mutexHandle = syscall.Handle(ret)

	// Check if mutex already existed
	if err.(syscall.Errno) == ERROR_ALREADY_EXISTS {
		// Another instance is already running
		log.Println("Another instance of Tyr is already running")
		procCloseHandle.Call(uintptr(mutexHandle))
		return false, nil
	}

	log.Println("Single instance lock acquired (Windows mutex)")
	return true, nil
}

// releaseLock releases the Windows mutex
func (si *SingleInstance) releaseLock() error {
	if mutexHandle == 0 {
		return nil
	}

	// Release and close the mutex
	procReleaseMutex.Call(uintptr(mutexHandle))
	procCloseHandle.Call(uintptr(mutexHandle))
	mutexHandle = 0

	return nil
}
