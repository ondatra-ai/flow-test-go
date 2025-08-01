package testutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// File permissions constants.
const (
	dirMode  = 0o750
	fileMode = 0o600
)

// Package-level errors.
var (
	ErrNoCoverageData = errors.New("no coverage data found")
)

// Static errors for better error handling.
var (
	ErrInvalidFilePath        = errors.New("invalid file path")
	ErrNullBytesInPath        = errors.New("file path contains null bytes")
	ErrPathTraversal          = errors.New("file path contains directory traversal")
	ErrInvalidCommandName     = errors.New("invalid command name")
	ErrInvalidMergedDir       = errors.New("invalid merged directory")
	ErrInvalidCoverageOut     = errors.New("invalid coverage output path")
	ErrInvalidCoverageHTML    = errors.New("invalid coverage HTML path")
	ErrInvalidCoverageSummary = errors.New("invalid coverage summary path")
	ErrCreateCovdataCommand   = errors.New("failed to create covdata command")
	ErrCreateHTMLCommand      = errors.New("failed to create HTML cover command")
	ErrCreateSummaryCommand   = errors.New("failed to create summary cover command")
)

// CoverageCollector manages coverage data collection and aggregation.
type CoverageCollector struct {
	baseDir      string
	manifestPath string
	manifest     *CoverageManifest
}

// CoverageManifest tracks all test coverage data.
type CoverageManifest struct {
	Version   string              `json:"version"`
	Timestamp string              `json:"timestamp"`
	Tests     []CoverageTestEntry `json:"tests"`
	Summary   CoverageSummary     `json:"summary"`
}

// CoverageTestEntry represents a single test's coverage data.
type CoverageTestEntry struct {
	Package  string `json:"package"`
	Test     string `json:"test"`
	CoverDir string `json:"coverDir"`
	Duration string `json:"duration"`
	Status   string `json:"status"`
}

// CoverageSummary provides overall coverage statistics.
type CoverageSummary struct {
	TotalTests int    `json:"totalTests"`
	Passed     int    `json:"passed"`
	Failed     int    `json:"failed"`
	Coverage   string `json:"coverage"`
}

// NewCoverageCollector creates a new coverage collector.
func NewCoverageCollector() *CoverageCollector {
	baseDir := filepath.Join("coverage", "e2e")
	manifestPath := filepath.Join(baseDir, "manifest.json")

	return &CoverageCollector{
		baseDir:      baseDir,
		manifestPath: manifestPath,
		manifest: &CoverageManifest{
			Version:   "1.0",
			Timestamp: time.Now().Format(time.RFC3339),
			Tests:     make([]CoverageTestEntry, 0),
			Summary: CoverageSummary{
				TotalTests: 0,
				Passed:     0,
				Failed:     0,
				Coverage:   "0.0%",
			},
		},
	}
}

// LoadExistingManifest loads an existing coverage manifest if it exists.
func (c *CoverageCollector) LoadExistingManifest() error {
	_, err := os.Stat(c.manifestPath)
	if os.IsNotExist(err) {
		return nil // No existing manifest, start fresh
	}

	data, err := os.ReadFile(c.manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	err = json.Unmarshal(data, c.manifest)
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	return nil
}

// AddTestResult adds a test result to the coverage manifest.
func (c *CoverageCollector) AddTestResult(packageName, testName, coverDir string, duration time.Duration, passed bool) {
	status := "passed"
	if !passed {
		status = "failed"
	}

	entry := CoverageTestEntry{
		Package:  packageName,
		Test:     testName,
		CoverDir: coverDir,
		Duration: duration.String(),
		Status:   status,
	}

	c.manifest.Tests = append(c.manifest.Tests, entry)
	c.updateSummary()
}

// SaveManifest saves the coverage manifest to disk.
func (c *CoverageCollector) SaveManifest() error {
	// Ensure directory exists
	err := os.MkdirAll(filepath.Dir(c.manifestPath), dirMode)
	if err != nil {
		return fmt.Errorf("failed to create manifest directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(c.manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	// Write to file
	err = os.WriteFile(c.manifestPath, data, fileMode)
	if err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// AggregateCoverage merges all coverage data and generates reports.
func (c *CoverageCollector) AggregateCoverage() error {
	// Create merged coverage directory
	mergedDir := filepath.Join("coverage", "e2e-merged")

	err := os.MkdirAll(mergedDir, dirMode)
	if err != nil {
		return fmt.Errorf("failed to create merged directory: %w", err)
	}

	// Find all coverage directories
	coverDirs := make([]string, 0)

	err = filepath.Walk(c.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Look for coverage data directories (contain covcounters files)
		if info.IsDir() && path != c.baseDir && path != mergedDir {
			// Check if this directory contains coverage files
			entries, err := os.ReadDir(path)
			if err == nil {
				for _, entry := range entries {
					if !entry.IsDir() && (filepath.Ext(entry.Name()) == "" ||
						strings.Contains(entry.Name(), "covcounters") ||
						strings.Contains(entry.Name(), "covmeta")) {
						coverDirs = append(coverDirs, path)

						break
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to find coverage directories: %w", err)
	}

	if len(coverDirs) == 0 {
		return ErrNoCoverageData
	}

	// Merge coverage data
	return c.mergeCoverageData(coverDirs, mergedDir)
}

// mergeCoverageData merges coverage data from multiple directories.
func (c *CoverageCollector) mergeCoverageData(sourceDirs []string, mergedDir string) error {
	// Build merge command
	args := []string{"tool", "covdata", "merge"}

	// Add input directories as single comma-separated value
	if len(sourceDirs) > 0 {
		inputDirs := "-i=" + strings.Join(sourceDirs, ",")
		args = append(args, inputDirs)
	}

	// Add output directory
	args = append(args, "-o="+mergedDir)

	// Execute merge command
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "go", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to merge coverage data: %w\nOutput: %s", err, output)
	}

	// Convert to text format
	return c.generateReports(mergedDir)
}

// generateReports creates text and HTML coverage reports.
func (c *CoverageCollector) generateReports(mergedDir string) error {
	// Validate the merged directory path
	err := validateFilePath(mergedDir)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidMergedDir, err)
	}

	coverageOut := filepath.Join("coverage", "e2e.out")
	coverageHTML := filepath.Join("coverage", "e2e.html")
	coverageSummary := filepath.Join("coverage", "e2e-summary.txt")

	// Validate all file paths
	err = validateCoveragePaths(coverageOut, coverageHTML, coverageSummary)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Convert to text format - use secure command construction
	err = executeConversionCommand(ctx, mergedDir, coverageOut)
	if err != nil {
		return err
	}

	// Generate HTML report - use secure command construction
	err = executeHTMLCommand(ctx, coverageOut, coverageHTML)
	if err != nil {
		return err
	}

	// Generate summary report - use secure command construction
	output, err := executeSummaryCommand(ctx, coverageOut)
	if err != nil {
		return err
	}

	// Save summary to file
	err = os.WriteFile(coverageSummary, output, fileMode)
	if err != nil {
		return fmt.Errorf("failed to save summary report: %w", err)
	}

	// Extract total coverage percentage and update manifest
	coverage := c.extractCoveragePercentage(string(output))
	c.manifest.Summary.Coverage = coverage

	return nil
}

// extractCoveragePercentage extracts the total coverage percentage from go tool cover output.
func (c *CoverageCollector) extractCoveragePercentage(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "total:") {
			// Extract percentage (should be at the end of the line)
			fields := strings.Fields(line)
			if len(fields) > 0 {
				lastField := fields[len(fields)-1]
				if strings.Contains(lastField, "%") {
					return lastField
				}
			}
		}
	}

	return "0.0%"
}

// updateSummary recalculates the coverage summary.
func (c *CoverageCollector) updateSummary() {
	c.manifest.Summary.TotalTests = len(c.manifest.Tests)
	c.manifest.Summary.Passed = 0
	c.manifest.Summary.Failed = 0

	for _, test := range c.manifest.Tests {
		if test.Status == "passed" {
			c.manifest.Summary.Passed++
		} else {
			c.manifest.Summary.Failed++
		}
	}

	// Update timestamp
	c.manifest.Timestamp = time.Now().Format(time.RFC3339)
}

// RecordTestExecution is a helper function for tests to record their execution.
func RecordTestExecution(t *testing.T, packageName string, coverDir string, duration time.Duration, passed bool) {
	t.Helper()

	collector := NewCoverageCollector()

	err := collector.LoadExistingManifest()
	if err != nil {
		t.Logf("Warning: Could not load existing manifest: %v", err)
	}

	collector.AddTestResult(packageName, t.Name(), coverDir, duration, passed)

	err = collector.SaveManifest()
	if err != nil {
		t.Logf("Warning: Could not save coverage manifest: %v", err)
	}
}

// RecordTestResult is a helper function to record test execution with standardized logic.
func RecordTestResult(t *testing.T, testType string, result *FlowTestResult, duration time.Duration) {
	t.Helper()
	RecordTestExecution(t, testType, "", duration, result.ExitCode == 0)
}

// validateFilePath ensures the file path is safe to use in commands.
func validateFilePath(path string) error {
	// Convert to absolute path to prevent path traversal
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidFilePath, err)
	}

	// Check for null bytes and other dangerous characters
	if strings.Contains(absPath, "\x00") {
		return fmt.Errorf("%w: %s", ErrNullBytesInPath, absPath)
	}

	// Check for path traversal attempts
	if strings.Contains(absPath, "..") {
		return fmt.Errorf("%w: %s", ErrPathTraversal, absPath)
	}

	return nil
}

// secureCommand creates a command with validated arguments.
func secureCommand(ctx context.Context, name string, args ...string) (*exec.Cmd, error) {
	// Validate the command name
	if strings.Contains(name, "\x00") || strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommandName, name)
	}

	// Validate all arguments
	var validArgs []string

	for _, arg := range args {
		// Remove null bytes and trim whitespace
		cleaned := strings.ReplaceAll(arg, "\x00", "")
		cleaned = strings.TrimSpace(cleaned)

		if cleaned != "" {
			validArgs = append(validArgs, cleaned)
		}
	}

	return exec.CommandContext(ctx, name, validArgs...), nil
}

// validateCoveragePaths validates all coverage-related file paths.
func validateCoveragePaths(coverageOut, coverageHTML, coverageSummary string) error {
	err := validateFilePath(coverageOut)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidCoverageOut, err)
	}

	err = validateFilePath(coverageHTML)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidCoverageHTML, err)
	}

	err = validateFilePath(coverageSummary)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidCoverageSummary, err)
	}

	return nil
}

// executeConversionCommand executes the coverage conversion command.
func executeConversionCommand(ctx context.Context, mergedDir, coverageOut string) error {
	cmd, err := secureCommand(ctx, "go", "tool", "covdata", "textfmt",
		"-i="+mergedDir, "-o="+coverageOut)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCreateCovdataCommand, err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to convert coverage to text format: %w\nOutput: %s", err, output)
	}

	return nil
}

// executeHTMLCommand executes the HTML coverage report generation command.
func executeHTMLCommand(ctx context.Context, coverageOut, coverageHTML string) error {
	cmd, err := secureCommand(ctx, "go", "tool", "cover",
		"-html="+coverageOut, "-o="+coverageHTML)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCreateHTMLCommand, err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate HTML report: %w\nOutput: %s", err, output)
	}

	return nil
}

// executeSummaryCommand executes the summary coverage report generation command.
func executeSummaryCommand(ctx context.Context, coverageOut string) ([]byte, error) {
	cmd, err := secureCommand(ctx, "go", "tool", "cover", "-func="+coverageOut)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateSummaryCommand, err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary report: %w\nOutput: %s", err, output)
	}

	return output, nil
}
