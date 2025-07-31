# 🎨🎨🎨 ARCHIVED: Custom Flow Engine Architecture 🎨🎨🎨

> **⚠️ IMPORTANT**: This document describes a custom flow engine design that was **REPLACED** by LangGraph integration.
> 
> **Current Decision**: Use LangGraph (github.com/tmc/langgraphgo) for all flow orchestration
> 
> **See**: `creative-langgraph-integration.md` for the current architecture

---

## Original Design (For Historical Reference)

This document originally explored three options for a custom flow engine:
1. Traditional State Machine Pattern
2. Event-Driven Architecture
3. Hybrid Orchestrator Pattern (was selected)

The Hybrid Orchestrator Pattern was chosen for its flexibility and production readiness.

## Why This Was Replaced

After the initial design phase, we discovered LangGraph, which provides:
- All the features we designed in the Hybrid Orchestrator
- Battle-tested production implementation
- Built-in checkpointing and state management
- Human-in-the-loop capabilities
- Better debugging and visualization tools

## Current Architecture

Instead of building a custom flow engine, we now:
1. Use LangGraph's graph structure for flow definition
2. Convert our JSON flow configs to LangGraph nodes and edges
3. Leverage LangGraph's built-in state management
4. Use LangGraph's checkpointing for fault tolerance

For the current implementation approach, see:
- `creative-langgraph-integration.md`
- `examples/langgraph-integration/`

---

🎨🎨🎨 ARCHIVED DOCUMENT - SEE LANGGRAPH INTEGRATION 🎨🎨🎨
