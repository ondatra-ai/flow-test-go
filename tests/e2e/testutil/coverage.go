package testutil

import (
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

// Package-level errors
var (
	ErrNoCoverageData = errors.New("no coverage data found")
)

// CoverageCollector manages coverage data collection and aggregation
type CoverageCollector struct {
	baseDir      string
	manifestPath string
	manifest     *CoverageManifest
}

// CoverageManifest tracks all test coverage data
type CoverageManifest struct {
	Version   string              `json:"version"`
	Timestamp string              `json:"timestamp"`
	Tests     []CoverageTestEntry `json:"tests"`
	Summary   CoverageSummary     `json:"summary"`
}

// CoverageTestEntry represents a single test's coverage data
type CoverageTestEntry struct {
	Package  string `json:"package"`
	Test     string `json:"test"`
	CoverDir string `json:"coverDir"`
	Duration string `json:"duration"`
	Status   string `json:"status"`
}

// CoverageSummary provides overall coverage statistics
type CoverageSummary struct {
	TotalTests int    `json:"totalTests"`
	Passed     int    `json:"passed"`
	Failed     int    `json:"failed"`
	Coverage   string `json:"coverage"`
}

// NewCoverageCollector creates a new coverage collector
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

// LoadExistingManifest loads an existing coverage manifest if it exists
func (c *CoverageCollector) LoadExistingManifest() error {
	if _, err := os.Stat(c.manifestPath); os.IsNotExist(err) {
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

// AddTestResult adds a test result to the coverage manifest
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

// SaveManifest saves the coverage manifest to disk
func (c *CoverageCollector) SaveManifest() error {
	// Ensure directory exists
	err := os.MkdirAll(filepath.Dir(c.manifestPath), 0o750)
	if err != nil {
		return fmt.Errorf("failed to create manifest directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(c.manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	// Write to file
	err = os.WriteFile(c.manifestPath, data, 0o600)
	if err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// AggregateCoverage merges all coverage data and generates reports
func (c *CoverageCollector) AggregateCoverage() error {
	// Create merged coverage directory
	mergedDir := filepath.Join("coverage", "e2e-merged")
	err := os.MkdirAll(mergedDir, 0o750)
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

// mergeCoverageData merges coverage data from multiple directories
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
	cmd := exec.Command("go", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to merge coverage data: %w\nOutput: %s", err, output)
	}

	// Convert to text format
	return c.generateReports(mergedDir)
}

// generateReports creates text and HTML coverage reports
func (c *CoverageCollector) generateReports(mergedDir string) error {
	coverageOut := filepath.Join("coverage", "e2e.out")
	coverageHTML := filepath.Join("coverage", "e2e.html")
	coverageSummary := filepath.Join("coverage", "e2e-summary.txt")

	// Convert to text format
	cmd := exec.Command("go", "tool", "covdata", "textfmt",
		"-i="+mergedDir, "-o="+coverageOut)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to convert coverage to text format: %w\nOutput: %s", err, output)
	}

	// Generate HTML report
	cmd = exec.Command("go", "tool", "cover",
		"-html="+coverageOut, "-o="+coverageHTML)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate HTML report: %w\nOutput: %s", err, output)
	}

	// Generate summary report
	cmd = exec.Command("go", "tool", "cover", "-func="+coverageOut)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate summary report: %w\nOutput: %s", err, output)
	}

	// Save summary to file
	err = os.WriteFile(coverageSummary, output, 0o600)
	if err != nil {
		return fmt.Errorf("failed to write summary file: %w", err)
	}

	// Extract total coverage percentage and update manifest
	coverage := c.extractCoveragePercentage(string(output))
	c.manifest.Summary.Coverage = coverage

	return nil
}

// extractCoveragePercentage extracts the total coverage percentage from go tool cover output
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

// updateSummary recalculates the coverage summary
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

// RecordTestExecution is a helper function for tests to record their execution
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

// RecordTestResult is a helper function to record test execution with standardized logic
func RecordTestResult(t *testing.T, testType string, result *FlowTestResult, duration time.Duration) {
	t.Helper()
	RecordTestExecution(t, testType, "", duration, result.ExitCode == 0)
}
