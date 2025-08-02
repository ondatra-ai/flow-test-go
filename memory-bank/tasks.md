# GitHub Actions Workflow for PR/Issue Comment Monitoring

## Description
Add a GitHub Actions workflow that automatically runs when new comments are added to Pull Requests or Issues. The workflow should have a single step that prints information about the comment to the console.

## GitHub Issue
Issue #8: https://github.com/ondatra-ai/flow-test-go/issues/8

## Complexity
Level: 1
Type: Quick Feature Addition

## Key Results
1. ✅ **Workflow Creation**: Created comment-monitor.yml in .github/workflows/
2. ✅ **Trigger Configuration**: Set up issue_comment event trigger for both PRs and Issues
3. ✅ **Output Format**: Print comment author, body, PR/Issue number, repository, timestamp
4. ✅ **CI Pattern Compliance**: Follow existing permissions and structure patterns from ci.yml

## Requirements Analysis
- Core Requirements:
  - [x] Trigger on issue_comment events (both PRs and Issues)
  - [x] Single job with minimal setup
  - [x] Output all required comment information
  - [x] Follow existing CI patterns
- Technical Constraints:
  - [x] Use ubuntu-latest runner
  - [x] Minimal permissions (read-only)
  - [x] No external dependencies

## Implementation Strategy
1. [x] Create .github/workflows/comment-monitor.yml
2. [x] Configure issue_comment trigger
3. [x] Set minimal permissions
4. [x] Add single job with output step
5. [x] Use GitHub context variables for comment data
6. [x] Test with commit validation

## Implementation Details

### Workflow File: `.github/workflows/comment-monitor.yml`
- **Event Trigger**: `issue_comment` with `types: [created]`
- **Permissions**: `contents: read` (minimal)
- **Runner**: `ubuntu-latest`
- **Timeout**: 5 minutes
- **Dependabot Exclusion**: `if: ${{ github.actor != 'dependabot[bot]' }}`

### Output Information
- Comment author: `@${{ github.event.comment.user.login }}`
- Comment body: `${{ github.event.comment.body }}`
- PR/Issue number: `#${{ github.event.issue.number }}`
- Repository: `${{ github.repository }}`
- Timestamp: `${{ github.event.comment.created_at }}`
- Additional details: Comment ID, URL, issue title, state, action

### Validation
- [x] Workflow file created successfully
- [x] Pre-commit hooks passed (including workflow validation)
- [x] Git commit completed successfully
- [x] Follows existing CI patterns from ci.yml

## Branch
- Name: task-20250111-pr-issue-comment-monitoring
- Created: ✅
- Committed: ✅

## Status
- [x] Initialization complete
- [x] Implementation complete
- [x] Workflow file created and validated
- [x] Pre-commit validation passed
- [x] Git commit successful
- [x] Ready for testing with actual comments

## Next Steps
- Ready for PR creation and testing with real comments
- Workflow will automatically trigger when comments are added to PRs or Issues
