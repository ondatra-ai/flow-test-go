#!/usr/bin/env ts-node

import { execSync } from 'child_process';

type GitHubThread = {
  id: string;
  isResolved: boolean;
};

type GitHubError = {
  message: string;
};

type GitHubResponse = {
  data?: {
    resolveReviewThread?: {
      thread: GitHubThread;
    };
    addPullRequestReviewThreadReply?: {
      comment: {
        id: string;
      };
    };
  };
  errors?: GitHubError[];
};

/**
 * Adds a comment to a specific review thread
 * @param threadId - The thread ID to add comment to
 * @param comment - The comment text to add
 */
function addCommentToThread(threadId: string, comment: string): void {
  try {
    console.log(`Adding comment to thread: ${threadId}`);

    const query = `mutation {
      addPullRequestReviewThreadReply(input: {
        pullRequestReviewThreadId: "${threadId}",
        body: "${comment.replace(/\\/g, '\\\\').replace(/"/g, '\\"').replace(/\n/g, '\\n')}"
      }) {
        comment {
          id
        }
      }
    }`;

    const result = execSync(`gh api graphql -f query='${query}'`, {
      encoding: 'utf8',
      stdio: 'pipe',
    });

    const response = JSON.parse(result) as GitHubResponse;

    if (response.data?.addPullRequestReviewThreadReply?.comment) {
      console.log(`✅ Successfully added comment to thread`);
    } else if (response.errors) {
      console.error('❌ GitHub API errors:');
      response.errors.forEach((error: GitHubError) => {
        console.error(`  - ${error.message}`);
      });
      throw new Error('Failed to add comment to thread');
    }
  } catch (error) {
    console.error('❌ Failed to add comment to thread:', error);
    throw error;
  }
}

/**
 * Resolves a GitHub review thread (conversation) using the provided thread ID
 * @param conversationId - The GitHub review thread ID to resolve
 * @param comment - Optional comment to add before resolving
 */
function resolveConversation(conversationId: string, comment?: string): void {
  if (!conversationId) {
    console.error('Error: conversationId is required');
    console.log(
      'Usage: npm run resolve-conversation <conversationId> [comment]'
    );
    console.log('');
    console.log(
      'If a comment is provided, it will be added to the conversation thread before resolving.'
    );
    process.exit(1);
  }

  try {
    // Add comment to thread if provided
    if (comment) {
      addCommentToThread(conversationId, comment);
    }

    // Resolve the thread
    const query = `mutation { 
      resolveReviewThread(input: {threadId: "${conversationId}"}) { 
        thread { id isResolved } 
      } 
    }`;

    console.log(`Resolving conversation: ${conversationId}`);

    const result = execSync(`gh api graphql -f query='${query}'`, {
      encoding: 'utf8',
      stdio: 'pipe',
    });

    const response = JSON.parse(result) as GitHubResponse;

    if (response.data?.resolveReviewThread?.thread) {
      const thread = response.data.resolveReviewThread.thread;
      console.log(`✅ Successfully resolved conversation ${thread.id}`);
      console.log(`Status: ${thread.isResolved ? 'Resolved' : 'Not resolved'}`);
    } else if (response.errors) {
      console.error('❌ GitHub API errors:');
      response.errors.forEach((error: GitHubError) => {
        console.error(`  - ${error.message}`);
      });
      process.exit(1);
    } else {
      console.error('❌ Unexpected response format');
      process.exit(1);
    }
  } catch (error) {
    console.error('❌ Failed to resolve conversation:', error);
    process.exit(1);
  }
}

// Get conversationId and optional comment from command line arguments
const conversationId = process.argv[2];
const comment = process.argv[3];

// Call the function and handle it properly
try {
  resolveConversation(conversationId, comment);
} catch (error) {
  console.error('Failed to resolve conversation:', error);
  process.exit(1);
}
