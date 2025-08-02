# Q&A: System Design Clarifications

## Q1: Flow Execution Understanding
**Q**: Do I understand correctly that the flow execution works like this: A user defines a JSON flow file with steps and conditions, then runs the CLI command which reads this JSON, converts it to a LangGraph state machine, and each step in the flow can invoke tools from MCP servers that are running as separate processes?

**A**: I think that MCP servers should be included into current process

## Q2: MCP Servers as Embedded
**Q**: Do I understand correctly that MCP servers are embedded within the main flow-test-go process rather than running as external subprocesses? So when the flow engine needs to execute a tool, it directly calls the MCP server code that's compiled into the same binary, rather than communicating with a separate process via stdio or HTTP?

**A**: yes

## Q3: MCP Protocol Implementation
**Q**: Do I understand correctly that the MCP protocol implementation in this project is different from the standard MCP approach - instead of having separate server processes that communicate via stdio/HTTP, you've built MCP-compatible tools directly into the Go binary, so they're just function calls within the same process? This would mean the "MCP servers" are more like embedded tool providers that follow the MCP protocol interface but don't require process management or IPC?

**A**: yes

## Q4: Tool Discovery
**Q**: Do I understand correctly that since MCP servers are embedded in-process, the "tool discovery" mechanism is actually just registering available tools at startup from compiled-in providers, rather than querying external servers? And the flow configuration files that reference "mcpServer" names are actually referring to these internal tool providers rather than external server configurations?

**A**: yes, but every flow and step should define in json config to set of tools it has access to (subset of all registered tools)

## Q5: Tool Access Control
**Q**: Do I understand correctly that each step in the flow JSON must explicitly list which tools it's allowed to use via a "tools" array, like `"tools": ["github_create_issue", "github_list_prs"]`, and during execution, the flow engine enforces that a step can only access the tools it has declared? This provides security and clarity by making tool access explicit rather than giving every step access to all registered tools?

**A**: Correct

## Q6: LLM Integration in Steps
**Q**: Do I understand correctly that when a flow step with type "prompt" is executed, the system sends the prompt text to an LLM via OpenRouter, and the LLM can then call the tools that are specified in that step's "tools" array? So the LLM acts as the decision-maker that determines which tools to invoke based on the prompt instructions?

**A**: Correct.

In additional we need to discuss context and memory approach to every step, add this topic for furher qurestions

## Q7: Context Between Steps
**Q**: Do I understand correctly that each step in a flow has access to some form of context or state that gets passed between steps? For example, if step1 creates a GitHub issue, does step2 automatically have access to that issue number, or does the system maintain a shared memory/context object that accumulates data as the flow executes through its steps?

**A**: Correct, but every next step should have access to only certain pieces of context. this also should be configured in json, but i have no clear vision how

## Q8: Context Access Declaration
**Q**: Do I understand correctly that you want each step to declare which context variables it needs access to, similar to how it declares tools? For example, a step might have something like `"context": ["issue_number", "user_email"]` to explicitly state it only needs those two pieces of data from previous steps, rather than having access to the entire accumulated state? This would help with security, performance, and clarity of data flow?

**A**: Correct, this isn're related to security but to clarity of data. E.g. I want to have e system when each step has access to ALL information required for execution, but not more.

Also, context can be filled not only by LLM but also by other operation. Say it's possibe that there's can be a step that just reads body of issue from github and saves it to context. No LLM needed here at all

## Q9: Non-LLM Step Types
**Q**: Do I understand correctly that you want to support different step types beyond just "prompt" and "condition" - perhaps a "tool" or "action" type that directly executes a tool without LLM involvement? For example, a step like:
```json
{
  "type": "tool",
  "tool": "github_get_issue",
  "parameters": {
    "issue_number": "${context.issue_number}"
  },
  "output": "issue_body",
  "context": ["issue_number"]
}
```
This would read the issue and save its body to context without any LLM prompt?

**A**: Correct. Format should be specified later, but in general - yes.

## Q10: Context Isolation
**Q**: Do I understand correctly that the flow engine maintains a growing context object as it executes, where each step can read from specific context variables (as declared in its configuration) and write new values to context? So if step1 writes `issue_number: 123` to context, and step2 declares `"context": ["issue_number"]`, then step2 can read that value, but if step3 doesn't declare needing `issue_number`, it won't have access to it even though it's still in the global context?

**A**: yes

## Q11: Checkpointing with Context
**Q**: Do I understand correctly that when a flow fails mid-execution, the LangGraph checkpointing system saves both the current step position AND the accumulated context at that point? So when resuming from a checkpoint, the flow can continue from where it left off with all the context values that were built up to that point, rather than having to re-execute from the beginning?

**A**: correct.

## Q12: Documentation Update
**Q**: Do I understand correctly that I should now update the systemPatterns.md documentation to reflect all these clarifications we've discussed? Specifically:
- Clarify that MCP servers are embedded in-process, not external
- Document the context access pattern where each step declares which context variables it needs
- Add information about different step types (prompt, condition, tool/action)
- Update the flow JSON examples to show the "context" field for controlled data access
- Explain how checkpointing preserves both execution position and accumulated context

**A**: Yes, please update systemPatterns.md accordingly

Also, add: ./QnA.md file that exactly saves this dialog without any interpretations

## Q13: Initial Context Population
**Q**: Do I understand correctly that when a user executes a flow from the command line, they can pass initial values that populate the starting context? For example, something like `flow-test-go execute my-flow.json --context "repository_name=myrepo" --context "user_email=user@example.com"`, and these values would be available to the first step if it declares them in its context array?

**A**: correct. sources can be different: env variables, parameters, or path to yaml file with prepopulated context. Let's define later

## Q14: Tool Implementation
**Q**: Do I understand correctly that when implementing a new tool provider (like a Slack integration or AWS service), I would need to:
1. Create a Go struct that implements a specific Tool interface
2. Register this tool provider at startup with the Tool Registry
3. Define the tool's input parameters and output schema
4. Implement the execution logic that can read from context and write results back to context?

**A**: yes. In addition, you will need to create e2e tests, yes, it's possible to test what too sent message to slack and you catch this message via slack API.

## Q15: E2E Testing Requirements
**Q**: Do I understand correctly that e2e tests for tools should verify the actual side effects rather than just mocking? So for a Slack tool test, the e2e test would:
1. Execute a flow that uses the Slack tool to send a message
2. Use the Slack API to verify the message was actually received in the channel
3. Check that the message content matches what was expected
4. Clean up the test message afterwards

This ensures the tool integration actually works end-to-end in a real environment?

**A**: Absolutely correct. this is mandatory!

## Q16: Checkpointing Storage
**Q**: Do I understand correctly that checkpoints are stored on the local filesystem, perhaps in a `.checkpoints/` directory or similar, and they contain:
1. The serialized LangGraph state (current step position)
2. The complete context object as JSON
3. Metadata like timestamp, flow ID, and checkpoint version
4. Perhaps a hash to verify integrity

And when resuming, the system looks for the most recent valid checkpoint for a given flow ID?

**A**: In general - correct. But it heavily depends on how good flow that implements task based on expected result and plan, will do the work. maybe we don't need checkpoints at all.

## Q17: Checkpointing Decision
**Q**: Do I understand correctly that checkpointing is more of a safety net for long-running or failure-prone flows, but well-designed flows that are:
- Idempotent (can be safely re-run)
- Fast enough to complete reliably
- Built with proper error handling in each step

might not need checkpointing at all? So the checkpointing feature would be optional - perhaps enabled via a flag like `--enable-checkpoints` or configured per flow in the JSON definition?

**A**: nope, we either need checkpoints for all flows or don't need them. no need to add flag, ut always true

## Q18: System-wide Checkpointing
**Q**: Do I understand correctly that checkpointing will be a system-wide architectural decision made once during development - either:
1. The system ALWAYS saves checkpoints for every flow execution (providing resume capability universally), or
2. The system NEVER uses checkpoints (relying on flows being re-runnable from the start)

And this decision will be hardcoded into the implementation rather than being configurable?

**A**: Correct.

## Q19: Error Handling in Tools
**Q**: Do I understand correctly that when a flow execution encounters an error in a tool, the error handling strategy would be:
1. The tool returns an error result to the flow engine
2. The error gets written to context (perhaps as `context.last_error` or similar)
3. The flow can have conditional steps that check for errors and decide how to proceed
4. Some errors might be retryable (rate limits, network issues) while others are fatal (invalid credentials, not found)

And each tool would need to classify its errors into these categories?

**A**: correct

## Q20: Logging System
**Q**: Do I understand correctly that the logging system would have different levels of verbosity where:
1. User-facing output goes to stdout (progress messages, results)
2. Debug logs go to stderr or a log file (tool invocations, context changes, LLM prompts/responses)
3. Each step execution gets logged with timestamp, step ID, and status
4. Tool calls get logged with input parameters and output results (with sensitive data redacted)

And developers could enable debug logging with something like `--log-level debug` to troubleshoot failing flows?

**A**: yes

## Q21: User Input During Execution
**Q**: Do I understand correctly that there's no built-in mechanism for handling user input during flow execution? So if a flow needs user confirmation or input mid-execution, it would need to:
1. Save its state and exit with a special status
2. Prompt the user to provide input via command line args when resuming
3. Or alternatively, all user inputs must be provided upfront in the initial context

Meaning flows are designed to run non-interactively from start to finish?

**A**: all user inputs must be provided upfront, there shouldn't be a situation when it's needed in the middle

## Q22: Credential Management
**Q**: Do I understand correctly that for managing API keys and credentials, the system would:
1. Read them from environment variables (like `GITHUB_TOKEN`, `OPENROUTER_API_KEY`)
2. Pass them to the relevant tool providers during initialization
3. Never store them in flow JSON files or context
4. Tools access their credentials internally without exposing them to the flow context

So a GitHub tool would automatically use `GITHUB_TOKEN` from the environment without the flow needing to know about it?

**A**: Not really, I expect that flow.json should have something like "github_token": "{env.GITHUB_TOKEN}", so it's possible that user populate secret in json, but tool doesn't care about that, it's on user.

## Q23: Template Syntax
**Q**: Do I understand correctly that the flow engine supports a template syntax for variable substitution where:
1. Flow JSON can contain references like `"github_token": "{env.GITHUB_TOKEN}"` or `"issue_id": "{context.issue_number}"`
2. At runtime, the engine resolves `{env.*}` from environment variables and `{context.*}` from the current context
3. These resolved values are then passed to tools as actual values
4. If someone puts a plain credential directly in JSON (not recommended), it would still work but that's their security risk?

This gives flexibility for both environment variables and context values to be referenced in flow configurations?

**A**: yes. Systax isself should be defined later. note that anthopic doesn't allo dots, so we need to align templating syntax with staff provided by openrouter

## Q24: CLI Commands
**Q**: Do I understand correctly that the CLI would provide these main commands:
1. `flow-test-go execute <flow.json>` - Run a flow with optional context parameters
2. `flow-test-go list` - List available flows in the .flows directory
3. `flow-test-go validate <flow.json>` - Check if a flow JSON is valid without executing
4. `flow-test-go tools` - List all registered tools and their parameters
5. Maybe `flow-test-go resume <checkpoint-id>` - Resume from a checkpoint (if checkpointing is enabled)

Are there other commands planned for the CLI?

**A**: correct. I don't see other commands

## Q25: Concurrent Execution
**Q**: Do I understand correctly that the system is designed to run one flow at a time (no concurrent flow executions)? So if a user tries to execute a second flow while one is already running, it would either:
1. Queue the second flow to run after the first completes
2. Reject the execution with an error message
3. Or each flow execution is a separate process, so multiple can run simultaneously but they don't share any state or resources?

**A**: yes, one flow at time. system doesn't care if user run two flows and something breaks. internally, system SHOULD NOT use any global state! E.g. it must be possble to run 2 flows even though it's no a case.

## Q26: Stateless Architecture
**Q**: Do I understand correctly that the system architecture must be completely stateless and thread-safe, meaning:
1. No global variables or shared state between flow executions
2. Each flow execution gets its own instances of tool providers, context, and engine
3. While the CLI is intended for single flow execution, the underlying architecture doesn't prevent concurrent runs
4. If a user runs multiple instances of the CLI simultaneously, they might face external conflicts (like both trying to create the same GitHub issue) but the application itself won't crash due to internal state conflicts?

**A**: Yes. 100% completely stateless, not sure we need to care about thread-safe. e.g. it shouldn't be a cases when two flows has shared state at all.

## Q27: Resource Isolation
**Q**: Do I understand correctly that each flow execution creates completely isolated instances of everything - separate tool registry, separate context manager, separate LangGraph engine - with zero sharing between executions? So it's not about making shared resources thread-safe, but rather ensuring there are NO shared resources at all between flow executions?

**A**: not really. LangGraph engline is shared. but it already thread-safe, so we shouldn't care. There might be shared resource, e.g. 2 flows can work simultaniosly with the same Repo. but no need to care

## Q28: Rate Limiting
**Q**: Do I understand correctly that for rate limiting and resource management, the system would rely on:
1. OpenRouter's built-in rate limiting for LLM calls (they handle queuing/throttling)
2. Individual tools implementing their own rate limit handling (e.g., GitHub tool respecting GitHub API limits)
3. No global rate limiting or resource management at the flow engine level
4. If a tool hits a rate limit, it returns a retryable error and the flow can handle it with conditional logic?

**A**: yes

## Q29: Scheduling and Orchestration
**Q**: Do I understand correctly that the system has no built-in scheduling or workflow orchestration features? So if someone wants to:
1. Run a flow every hour - they'd use external tools like cron
2. Chain multiple flows together - they'd create one larger flow that encompasses all steps
3. Run flows in response to webhooks - they'd need to build a separate webhook receiver that calls the CLI

The tool is focused solely on executing individual flows when invoked, not on scheduling or event-driven orchestration?

**A**: correct

## Q30: Debugging Failed Flows
**Q**: Do I understand correctly that for debugging failed flows, developers would:
1. Run with `--log-level debug` to see detailed execution logs
2. Examine the context at each step to understand data flow
3. Look at the specific error messages from failed tools
4. Potentially add temporary "debug" steps that just log context values
5. Use the checkpoint data (if enabled) to inspect the state at failure

But there's no interactive debugger or step-by-step execution mode?

**A**: Correct, but if you see different appoeach, let me know. I eager to discuss

## Q31: Dry-Run Mode
**Q**: Do I understand correctly that you might be interested in a "dry-run" mode where the flow engine would:
1. Parse and validate the flow
2. Show each step that would be executed with its configuration
3. Display what context variables each step would access
4. Show which tools would be called (but not actually call them)
5. Potentially even send prompts to LLM with a prefix like "[DRY-RUN] " to get responses without executing tools

This could help developers understand flow behavior before running it for real, especially useful for flows that modify external systems?

**A**: It would be good, yes

## Q32: Flow Composition
**Q**: Do I understand correctly that there's no built-in mechanism for flow composition or reusability? So if I have common patterns like "authenticate and get user info" that multiple flows need, I would need to:
1. Copy-paste those steps into each flow JSON
2. Or potentially use a templating system outside the tool to generate flows
3. But there's no native support for importing sub-flows or defining reusable step groups?

Or would this be something worth considering for reducing duplication?

**A**: Correct, no need to reducre deduplication now

## Q33: Rollback Mechanisms
**Q**: Do I understand correctly that the system has no built-in rollback or compensation mechanisms? So if a flow fails after creating a GitHub issue and updating a database:
1. Those changes remain in place (issue stays created)
2. The flow author would need to manually design compensating steps
3. Or create a separate "cleanup" flow to undo changes
4. There's no automatic transaction-like behavior across tools?

**A**: Correct, no need to care about faults

## Q34: Local Development Testing
**Q**: Do I understand correctly that for local development and testing of flows, developers would:
1. Create a test flow JSON file
2. Set up mock/test environments (test GitHub repo, test Slack channel)
3. Run the flow against these test resources
4. Verify the results manually or with e2e tests
5. There's no built-in mock mode where tools pretend to execute but don't actually make external calls?

**A**: Almost correct, we won't verify manually, only using e2e tests. 95% of fuctionality must be covered by e2e tests

## Q35: Test Infrastructure
**Q**: Do I understand correctly that the e2e test infrastructure would need to:
1. Automatically provision test resources (create test repos, Slack workspaces, etc.) before tests
2. Execute flows against these test environments
3. Use APIs to verify actual side effects occurred
4. Clean up test resources after execution
5. Run in CI/CD with that 95% coverage target enforced

And for the remaining 5% not covered by e2e tests - would these be edge cases that are too complex/expensive to test, or would they be covered by unit tests instead?

**A**: Correct. Some resouces are just exists and provided to tests via .env file (e.g. slack, repository), some needs to be created (issue, comment etc..)

## Q36: Monitoring and Observability
**Q**: Do I understand correctly that for monitoring and observability in production, the system would:
1. Log all flow executions with start/end times and success/failure status
2. Track metrics like LLM token usage (via OpenRouter's tracking) and API calls made
3. But there's no built-in metrics collection or dashboards
4. Users would need to parse logs or integrate with external monitoring tools if they want dashboards?

Or is observability beyond the scope of this CLI tool?

**A**: Nope, This is CLI,  don't care about observability.

## Q37: Flow Validation
**Q**: Do I understand correctly that the flow JSON validation would check:
1. Schema compliance (required fields, correct types)
2. Step reference integrity (all nextStep/yes/no references point to valid steps)
3. Tool existence (all tools in "tools" arrays are registered)
4. Context variable usage consistency (if step2 needs "issue_id", some previous step should write it)
5. No circular references that would cause infinite loops?

And this validation happens both in the `validate` command and before execution?

**A**: Correct

## Q38: Documentation Update
**Q**: Do I understand correctly that I should now update both systemPatterns.md and QnA.md with all these additional clarifications we've discussed? This would include:
- Initial context sources (CLI args, env vars, YAML files)
- Tool implementation requirements including e2e test verification
- Stateless architecture with no shared resources
- Non-interactive execution model
- Template syntax for variable substitution
- CLI commands and debugging approaches
- Dry-run mode concept
- Test infrastructure with pre-existing and created resources
- Validation requirements

Or are there other architectural areas you'd like to explore first?

**A**: Correct, update files
