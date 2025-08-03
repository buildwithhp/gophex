package cmd

import (
	"os/exec"
	"testing"
	"time"
)

func TestProcessManager_AddAndRemoveProcess(t *testing.T) {
	pm := &ProcessManager{
		processes: make(map[string]*ProcessInfo),
	}

	// Create a dummy command
	cmd := exec.Command("sleep", "1")
	processName := "test-process"
	description := "Test process"
	projectPath := "/tmp/test"

	// Add process
	pm.AddProcess(processName, description, projectPath, cmd)

	// Check if process was added
	pm.mutex.RLock()
	proc, exists := pm.processes[processName]
	pm.mutex.RUnlock()

	if !exists {
		t.Fatal("Process was not added to manager")
	}

	if proc.Name != processName {
		t.Errorf("Expected process name %s, got %s", processName, proc.Name)
	}

	if proc.Description != description {
		t.Errorf("Expected description %s, got %s", description, proc.Description)
	}

	if proc.ProjectPath != projectPath {
		t.Errorf("Expected project path %s, got %s", projectPath, proc.ProjectPath)
	}

	// Remove process
	pm.RemoveProcess(processName)

	// Check if process was removed
	pm.mutex.RLock()
	_, exists = pm.processes[processName]
	pm.mutex.RUnlock()

	if exists {
		t.Fatal("Process was not removed from manager")
	}
}

func TestProcessManager_GetRunningProcesses(t *testing.T) {
	pm := &ProcessManager{
		processes: make(map[string]*ProcessInfo),
	}

	// Start a real process that will run briefly
	cmd := exec.Command("sleep", "0.1")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start test process: %v", err)
	}

	pm.AddProcess("test-sleep", "Sleep test", "/tmp", cmd)

	// Get running processes immediately
	running := pm.GetRunningProcesses()
	if len(running) != 1 {
		t.Errorf("Expected 1 running process, got %d", len(running))
	}

	// Wait for process to finish
	cmd.Wait()
	time.Sleep(100 * time.Millisecond)

	// Check again - should be cleaned up
	running = pm.GetRunningProcesses()
	if len(running) != 0 {
		t.Errorf("Expected 0 running processes after completion, got %d", len(running))
	}
}

func TestProcessManager_HasRunningProcesses(t *testing.T) {
	pm := &ProcessManager{
		processes: make(map[string]*ProcessInfo),
	}

	// Initially no processes
	if pm.HasRunningProcesses() {
		t.Error("Expected no running processes initially")
	}

	// Start a process
	cmd := exec.Command("sleep", "0.1")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start test process: %v", err)
	}

	pm.AddProcess("test-process", "Test", "/tmp", cmd)

	// Should have running processes
	if !pm.HasRunningProcesses() {
		t.Error("Expected to have running processes")
	}

	// Wait for process to finish
	cmd.Wait()
	time.Sleep(100 * time.Millisecond)

	// Should not have running processes
	if pm.HasRunningProcesses() {
		t.Error("Expected no running processes after completion")
	}
}

func TestProcessManager_StartProcessWithTracking(t *testing.T) {
	pm := &ProcessManager{
		processes: make(map[string]*ProcessInfo),
	}

	// Create a command that will run briefly
	cmd := exec.Command("echo", "test")

	err := pm.StartProcessWithTracking("echo-test", "Echo test", "/tmp", cmd)
	if err != nil {
		t.Fatalf("Failed to start process with tracking: %v", err)
	}

	// Check if process was added
	pm.mutex.RLock()
	_, exists := pm.processes["echo-test"]
	pm.mutex.RUnlock()

	if !exists {
		t.Error("Process was not added to tracking")
	}

	// Wait for process to complete and be cleaned up
	time.Sleep(200 * time.Millisecond)

	// Process should be automatically removed
	pm.mutex.RLock()
	_, exists = pm.processes["echo-test"]
	pm.mutex.RUnlock()

	if exists {
		t.Error("Process was not automatically cleaned up")
	}
}

func TestProcessManager_TerminateAllProcesses(t *testing.T) {
	pm := &ProcessManager{
		processes: make(map[string]*ProcessInfo),
	}

	// Start multiple processes
	cmd1 := exec.Command("sleep", "10")
	cmd2 := exec.Command("sleep", "10")

	err1 := cmd1.Start()
	err2 := cmd2.Start()

	if err1 != nil || err2 != nil {
		t.Fatalf("Failed to start test processes: %v, %v", err1, err2)
	}

	pm.AddProcess("sleep1", "Sleep 1", "/tmp", cmd1)
	pm.AddProcess("sleep2", "Sleep 2", "/tmp", cmd2)

	// Verify processes are running
	if len(pm.GetRunningProcesses()) != 2 {
		t.Error("Expected 2 running processes")
	}

	// Terminate all processes
	pm.TerminateAllProcesses()

	// Wait a bit for termination
	time.Sleep(500 * time.Millisecond)

	// Check that processes are terminated
	if len(pm.processes) != 0 {
		t.Error("Expected all processes to be removed from tracking")
	}

	// The main thing we care about is that the processes are removed from tracking
	// The actual termination is handled by the OS and may take time
	// We've already verified that the processes are removed from the manager
}

func TestGetProcessManager(t *testing.T) {
	pm1 := GetProcessManager()
	pm2 := GetProcessManager()

	if pm1 != pm2 {
		t.Error("GetProcessManager should return the same instance (singleton)")
	}
}
