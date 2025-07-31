# Ondatra Code Rules

This directory contains all the rules and specifications for the Ondatra Code project.

## Rule Categories

### Core Architecture
- **Configuration Structure**: How the `.flows` folder works
- **Flow System Architecture**: Step structure and conditional execution
- **MCP Integration**: Server management and tool discovery

### Development Standards
- **TypeScript Configuration**: Strictest settings and type safety
- **Code Quality Standards**: ESLint, Prettier, testing requirements
- **Error Handling**: Comprehensive error management

### Project Structure
- **Directory Organization**: How files should be organized
- **Naming Conventions**: File and code naming standards
- **Implementation Priorities**: What to build in what order

## Key Features
- **Chat Interface**: Interactive conversational UI like claude-code
- **MCP Servers**: Multiple server connections for extended capabilities
- **Flow Engine**: JSON-based workflow definitions with branching
- **Tool Discovery**: Automatic detection of available MCP tools
- **TypeScript**: Strict type safety throughout the application

## Rule Files

1. **project-overview.mdc** - High-level project description and core features
2. **architecture-rules.mdc** - System architecture and design patterns
3. **development-standards.mdc** - TypeScript configuration, testing, and code quality standards
4. **flow-format-rules.mdc** - JSON schema and validation rules for flow definitions
5. **mcp-server-rules.mdc** - MCP server configuration and integration guidelines
6. **project-structure.mdc** - Directory structure and implementation priorities

## Quick Reference

### Core Technologies
- **Language**: TypeScript (strict mode)
- **Testing**: Vitest
- **Code Quality**: ESLint + Prettier
- **CLI Framework**: Ink (React for CLI)
- **MCP Integration**: @modelcontextprotocol/sdk

### Key Concepts
- **Flows**: JSON-based workflow definitions with steps and conditions
- **MCP Servers**: External processes providing tools and capabilities
- **Chat Interface**: Interactive conversational UI like claude-code
- **Configuration**: Read from `.flows` folder in current directory

### Development Workflow
1. Follow strict TypeScript configuration
2. Write tests for all features (unit + integration)
3. Use meaningful names and add documentation
4. Handle errors gracefully with clear messages
5. Keep functions small and focused

For detailed information, refer to the individual rule files. 