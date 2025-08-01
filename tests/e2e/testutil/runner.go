package testutil

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

const (
	defaultRunnerTimeout = 30 * time.Second
	coverageDirMode      = 0o750
)

// FlowRunner handles subprocess execution of flow-test-go binary.
type FlowRunner struct {
	t           *testing.T
	flowFile    string
	configDir   string
	workDir     string
	timeout     time.Duration
	binaryPath  string
	coverageDir string
	stdout      bytes.Buffer
	stderr      bytes.Buffer
}

// NewFlowRunner creates a new flow runner.
func NewFlowRunner(t *testing.T) *FlowRunner {
	t.Helper()

	// Find project root for binary path
	wd, _ := os.Getwd()

	projectRoot := wd
	for {
		_, err := os.Stat(filepath.Join(projectRoot, "go.mod"))
		if err == nil {
			break
		}

		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			break
		}

		projectRoot = parent
	}

	binaryPath := filepath.Join(projectRoot, "bin", "flow-test-go-e2e")

	return &FlowRunner{
		t:           t,
		flowFile:    "",
		configDir:   "",
		workDir:     "",
		timeout:     defaultRunnerTimeout,
		binaryPath:  binaryPath, // Use absolute path to coverage-instrumented binary
		coverageDir: "",
		stdout:      bytes.Buffer{},
		stderr:      bytes.Buffer{},
	}
}

// SetFlowFile sets the flow file to execute.
func (r *FlowRunner) SetFlowFile(flowFile string) {
	r.flowFile = flowFile
}

// SetConfigDir sets the config directory.
func (r *FlowRunner) SetConfigDir(configDir string) {
	r.configDir = configDir
}

// SetWorkDir sets the working directory.
func (r *FlowRunner) SetWorkDir(workDir string) {
	r.workDir = workDir
}

// SetTimeout sets the execution timeout.
func (r *FlowRunner) SetTimeout(timeout time.Duration) {
	r.timeout = timeout
}

// SetBinaryPath sets the path to the binary to execute.
func (r *FlowRunner) SetBinaryPath(binaryPath string) {
	r.binaryPath = binaryPath
}

// Execute runs the flow and returns the result.
func (r *FlowRunner) Execute() *FlowTestResult {
	start := time.Now()

	// Reset buffers for this execution
	r.stdout.Reset()
	r.stderr.Reset()

	// Setup coverage collection
	r.setupCoverage()

	// Build and configure command
	cmd := r.buildAndConfigureCommand()

	// Execute command with timeout
	err := cmd.Run()

	// Calculate duration
	duration := time.Since(start)

	// Determine exit code
	exitCode := r.determineExitCode(err)

	return &FlowTestResult{
		ExitCode: exitCode,
		Stdout:   r.stdout.String(),
		Stderr:   r.stderr.String(),
		Error:    err,
		Duration: duration,
	}
}

// CleanupCoverage removes coverage files for a test (optional cleanup).
func (r *FlowRunner) CleanupCoverage() {
	if r.coverageDir != "" {
		err := os.RemoveAll(r.coverageDir)
		if err != nil {
			r.t.Logf("Warning: Failed to cleanup coverage directory %s: %v", r.coverageDir, err)
		}
	}
}

// buildAndConfigureCommand creates and configures the command to execute.
func (r *FlowRunner) buildAndConfigureCommand() *exec.Cmd {
	// Build command
	cmd := r.buildCommand()

	// Setup output capture
	cmd.Stdout = &r.stdout
	cmd.Stderr = &r.stderr

	// Setup working directory if specified
	if r.workDir != "" {
		cmd.Dir = r.workDir
	}

	// Setup coverage environment
	if r.coverageDir != "" {
		if cmd.Env == nil {
			cmd.Env = os.Environ()
		}

		cmd.Env = append(cmd.Env, "GOCOVERDIR="+r.coverageDir)
	}

	// Create context with timeout and update command
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	if ctx != context.Background() {
		cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)
		cmd.Stdout = &r.stdout
		cmd.Stderr = &r.stderr

		cmd.Dir = r.workDir
		if r.coverageDir != "" {
			if cmd.Env == nil {
				cmd.Env = os.Environ()
			}

			cmd.Env = append(cmd.Env, "GOCOVERDIR="+r.coverageDir)
		}
	}

	return cmd
}

// determineExitCode extracts exit code from command error.
func (r *FlowRunner) determineExitCode(err error) int {
	exitCode := 0

	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode = exitError.ExitCode()
		} else {
			// Other error (e.g., binary not found, timeout)
			exitCode = -1
		}
	}

	return exitCode
}

// buildCommand constructs the command to execute.
func (r *FlowRunner) buildCommand() *exec.Cmd {
	// Start with the binary path
	args := []string{r.binaryPath}

	// For now, we can only test the existing 'list' command
	// since the 'run' command doesn't exist yet
	args = append(args, "list")

	// Note: flow file and config dir parameters are not supported by list command
	// This is a limitation of the current implementation

	// Create command
	cmd := exec.CommandContext(context.Background(), args[0], args[1:]...)

	return cmd
}

// setupCoverage creates a unique coverage directory for this test.
func (r *FlowRunner) setupCoverage() {
	if r.t == nil {
		return
	}

	// Create unique coverage directory for this test
	testName := r.t.Name()
	coverageBase := filepath.Join("coverage", "e2e")
	r.coverageDir = filepath.Join(coverageBase, testName)

	// Create coverage directory with parent directories
	err := os.MkdirAll(r.coverageDir, coverageDirMode)
	if err != nil {
		r.t.Fatalf("Failed to create coverage directory %s: %v", r.coverageDir, err)
	}

	r.t.Logf("Created coverage directory: %s", r.coverageDir)
}

// EnsureBinaryExists checks if the test binary exists and builds it if needed.
func EnsureBinaryExists(t *testing.T) {
	t.Helper()

	// Find project root first
	wd, _ := os.Getwd()

	projectRoot := wd
	for {
		_, err := os.Stat(filepath.Join(projectRoot, "go.mod"))
		if err == nil {
			break
		}

		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			break
		}

		projectRoot = parent
	}

	binaryPath := filepath.Join(projectRoot, "bin", "flow-test-go-e2e")

	// Check if binary exists
	_, err := os.Stat(binaryPath)
	if os.IsNotExist(err) {
		t.Logf("Binary %s not found, building...", binaryPath)

		// Build the binary with coverage - use project root we already found
		cmd := exec.CommandContext(context.Background(), "make", "build-e2e-coverage")
		cmd.Dir = projectRoot

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
		}

		t.Logf("Binary built successfully")
	}
}
