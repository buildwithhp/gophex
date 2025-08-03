package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/AlecAivazis/survey/v2"
)

// ProcessManager tracks and manages child processes started by Gophex
type ProcessManager struct {
	processes map[string]*ProcessInfo
	mutex     sync.RWMutex
}

// ProcessInfo contains information about a running process
type ProcessInfo struct {
	Cmd         *exec.Cmd
	Name        string
	Description string
	ProjectPath string
}

var globalProcessManager = &ProcessManager{
	processes: make(map[string]*ProcessInfo),
}

// GetProcessManager returns the global process manager instance
func GetProcessManager() *ProcessManager {
	return globalProcessManager
}

// AddProcess adds a process to the manager
func (pm *ProcessManager) AddProcess(name, description, projectPath string, cmd *exec.Cmd) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.processes[name] = &ProcessInfo{
		Cmd:         cmd,
		Name:        name,
		Description: description,
		ProjectPath: projectPath,
	}
}

// RemoveProcess removes a process from the manager
func (pm *ProcessManager) RemoveProcess(name string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	delete(pm.processes, name)
}

// GetRunningProcesses returns a list of currently running processes
func (pm *ProcessManager) GetRunningProcesses() []*ProcessInfo {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	var running []*ProcessInfo
	for _, proc := range pm.processes {
		if proc.Cmd != nil && proc.Cmd.Process != nil {
			// Check if process is still running
			if err := proc.Cmd.Process.Signal(syscall.Signal(0)); err == nil {
				running = append(running, proc)
			} else {
				// Process is dead, remove it
				delete(pm.processes, proc.Name)
			}
		}
	}

	return running
}

// HasRunningProcesses checks if there are any running child processes
func (pm *ProcessManager) HasRunningProcesses() bool {
	return len(pm.GetRunningProcesses()) > 0
}

// HandleGracefulShutdown handles shutdown when there are running processes
func (pm *ProcessManager) HandleGracefulShutdown() error {
	running := pm.GetRunningProcesses()
	if len(running) == 0 {
		fmt.Println("👋 Thank you for using Gophex!")
		return nil
	}

	fmt.Printf("\n⚠️  Found %d running process(es) started by Gophex:\n", len(running))
	for i, proc := range running {
		fmt.Printf("  %d. %s - %s (PID: %d)\n", i+1, proc.Name, proc.Description, proc.Cmd.Process.Pid)
	}
	fmt.Println()

	var action string
	shutdownPrompt := &survey.Select{
		Message: "What would you like to do with the running processes?",
		Options: []string{
			"🔄 Keep running in background and exit Gophex",
			"⏹️  Terminate all processes and exit",
			"❌ Cancel exit (return to menu)",
		},
	}

	err := survey.AskOne(shutdownPrompt, &action)
	if err != nil {
		if isUserInterrupt(err) {
			// Force terminate on interrupt
			pm.TerminateAllProcesses()
			fmt.Println("\n👋 All processes terminated. Goodbye!")
			return nil
		}
		return fmt.Errorf("shutdown prompt failed: %w", err)
	}

	switch {
	case action[:2] == "🔄":
		fmt.Println("📱 Processes will continue running in the background.")
		fmt.Println("💡 You can monitor them using your system's process manager.")
		for _, proc := range running {
			fmt.Printf("   • %s (PID: %d) in %s\n", proc.Name, proc.Cmd.Process.Pid, proc.ProjectPath)
		}
		fmt.Println("👋 Thank you for using Gophex!")
		return nil

	case action[:3] == "⏹️":
		fmt.Println("⏹️  Terminating all processes...")
		pm.TerminateAllProcesses()
		fmt.Println("👋 All processes terminated. Thank you for using Gophex!")
		return nil

	case strings.HasPrefix(action, "❌"):
		fmt.Println("↩️  Returning to menu...")
		return fmt.Errorf("exit cancelled")

	default:
		return fmt.Errorf("unknown shutdown action: %s", action)
	}
}

// TerminateAllProcesses terminates all tracked processes
func (pm *ProcessManager) TerminateAllProcesses() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	for name, proc := range pm.processes {
		if proc.Cmd != nil && proc.Cmd.Process != nil {
			fmt.Printf("   Terminating %s (PID: %d)...\n", proc.Name, proc.Cmd.Process.Pid)

			// Try graceful termination first
			if err := proc.Cmd.Process.Signal(os.Interrupt); err != nil {
				// Force kill if graceful termination fails
				proc.Cmd.Process.Kill()
			}
		}
		delete(pm.processes, name)
	}
}

// StartProcessWithTracking starts a process and adds it to the manager
func (pm *ProcessManager) StartProcessWithTracking(name, description, projectPath string, cmd *exec.Cmd) error {
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start %s: %w", name, err)
	}

	pm.AddProcess(name, description, projectPath, cmd)

	// Start a goroutine to clean up when process exits
	go func() {
		cmd.Wait()
		pm.RemoveProcess(name)
	}()

	return nil
}
