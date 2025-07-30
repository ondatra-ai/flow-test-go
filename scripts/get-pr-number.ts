#!/usr/bin/env node

/**
 * Get PR Number Script
 *
 * This script reveals the Pull Request number associated with
 * the current git branch. It uses the GitHub CLI (gh) to search
 * for open PRs linked to the current branch.
 *
 * Usage:
 *   Direct execution: ./scripts/get-pr-number.ts
 *   Via npm: npm run pr
 *
 * Requirements:
 *   - GitHub CLI (gh) must be installed: https://cli.github.com/
 *   - Must be run from within a git repository
 *   - Repository must be hosted on GitHub
 *
 * Output:
 *   - Shows PR number, title, state, and URL
 *   - If no open PR is found, shows closed/merged PRs
 *   - Displays just the PR number at the end for easy scripting
 */

import { execSync } from 'child_process';

type PullRequest = {
  number: number;
  title: string;
  state: string;
  url: string;
};

function getCurrentBranch(): string {
  try {
    const branch = execSync('git rev-parse --abbrev-ref HEAD', {
      encoding: 'utf8',
    }).trim();
    return branch;
  } catch (error) {
    console.error('Error getting current branch:', error);
    process.exit(1);
  }
}

function handleError(error: unknown): void {
  if (error instanceof Error) {
    if (error.message.includes('gh: command not found')) {
      console.error(
        'GitHub CLI (gh) is not installed. Please install it first:'
      );
      console.error('https://cli.github.com/manual/installation');
    } else if (error.message.includes('Could not resolve to a Repository')) {
      console.error(
        'This directory does not appear to be a GitHub repository.'
      );
    } else {
      console.error('Error fetching PR information:', error.message);
    }
  } else {
    console.error('Unknown error:', error);
  }
  process.exit(1);
}

function displayPR(pr: PullRequest, branch: string): void {
  console.log(`\nCurrent branch: ${branch}`);
  console.log(`PR #${pr.number}: ${pr.title}`);
  console.log(`State: ${pr.state}`);
  console.log(`URL: ${pr.url}`);

  // Output just the PR number for easy scripting
  console.log(`\nPR Number: ${pr.number}`);
}

function checkForClosedPRs(branch: string): void {
  const allPRsResult = execSync(
    `gh pr list --head "${branch}" --state all ` +
      `--json number,title,state,url --limit 5`,
    { encoding: 'utf8' }
  );

  const allPRs = JSON.parse(allPRsResult) as PullRequest[];
  if (allPRs.length > 0) {
    console.log('\nFound closed/merged PRs:');
    allPRs.forEach((pr: PullRequest) => {
      console.log(`  PR #${pr.number}: ${pr.title} (${pr.state})`);
      console.log(`  URL: ${pr.url}`);
    });
  }
}

function getPRForBranch(branch: string): void {
  if (branch === 'main' || branch === 'master') {
    console.log('Current branch is the default branch. No PR associated.');
    return;
  }

  try {
    // First, try to get PR using GitHub CLI search
    const result = execSync(
      `gh pr list --head "${branch}" --json number,title,state,url --limit 1`,
      { encoding: 'utf8' }
    );

    const prs = JSON.parse(result) as PullRequest[];

    if (prs.length === 0) {
      console.log(`No PR found for branch: ${branch}`);

      // Try to find PRs that might have been merged or closed
      checkForClosedPRs(branch);
      return;
    }

    displayPR(prs[0], branch);
  } catch (error) {
    handleError(error);
  }
}

// Main execution
const currentBranch = getCurrentBranch();
console.log(`Checking for PR associated with branch: ${currentBranch}`);
getPRForBranch(currentBranch);
