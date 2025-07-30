#!/usr/bin/env node

import { execSync } from 'child_process';

type Comment = {
  file: string | null;
  line: number | null;
  author: string;
  body: string;
  createdAt: string;
  outdated: boolean;
  resolved: boolean;
  diffHunk: string;
  url: string;
};

type Conversation = {
  id: string;
  isResolved: boolean;
  comments: Comment[];
};

type CommentNode = {
  path: string | null;
  line: number | null;
  body: string;
  createdAt: string;
  outdated: boolean;
  diffHunk: string;
  url: string;
  author: {
    login: string;
  };
};

type ThreadNode = {
  id: string;
  isResolved: boolean;
  comments: {
    nodes: CommentNode[];
  };
};

type GraphQLResponse = {
  data: {
    repository: {
      pullRequest: {
        reviewThreads: {
          nodes: ThreadNode[];
        };
      };
    };
  };
};

function getPRNumber(): string {
  const args = process.argv.slice(2);
  if (args.length === 0) {
    console.error('Usage: list-pr-comments.ts [PR_NUMBER]');
    process.exit(1);
  }
  return args[0];
}

function buildGraphQLQuery(): string {
  return `
    query($prNumber: Int!) {
      repository(owner: "ondatra-ai", name: "flow-test") {
        pullRequest(number: $prNumber) {
          reviewThreads(first: 100) {
            nodes {
              id
              isResolved
              comments(first: 10) {
                nodes {
                  path
                  line
                  body
                  createdAt
                  outdated
                  diffHunk
                  url
                  author {
                    login
                  }
                }
              }
            }
          }
        }
      }
    }
  `;
}

function parseConversations(data: GraphQLResponse): Conversation[] {
  const conversations: Conversation[] = [];

  data.data.repository.pullRequest.reviewThreads.nodes
    .filter((thread: ThreadNode) => !thread.isResolved)
    .forEach((thread: ThreadNode) => {
      const comments: Comment[] = thread.comments.nodes.map(
        (comment: CommentNode) => ({
          file: comment.path,
          line: comment.line,
          author: comment.author.login,
          body: comment.body,
          createdAt: comment.createdAt,
          outdated: comment.outdated,
          resolved: thread.isResolved,
          diffHunk: comment.diffHunk,
          url: comment.url,
        })
      );

      conversations.push({
        id: thread.id,
        isResolved: thread.isResolved,
        comments,
      });
    });

  return conversations;
}

function getPRComments(prNumber: string): void {
  const query = buildGraphQLQuery();

  try {
    const result = execSync(
      `gh api graphql -f query='${query}' -F prNumber=${prNumber}`,
      { encoding: 'utf8' }
    );

    const data = JSON.parse(result) as GraphQLResponse;
    const conversations = parseConversations(data);

    // Sort by creation date of first comment in each conversation
    conversations.sort((a: Conversation, b: Conversation) => {
      const aFirstComment = a.comments[0]?.createdAt || '';
      const bFirstComment = b.comments[0]?.createdAt || '';
      return (
        new Date(aFirstComment).getTime() - new Date(bFirstComment).getTime()
      );
    });

    // Output JSON by default for programmatic use
    console.log(JSON.stringify(conversations, null, 2));
  } catch (error) {
    console.error('Error fetching PR comments:', error);
    process.exit(1);
  }
}

// Main execution
const prNumber = getPRNumber();
getPRComments(prNumber);
