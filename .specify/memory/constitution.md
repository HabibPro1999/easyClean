# Unused Assets Detector Constitution

<!--
Sync Impact Report:
Version change: Initial → 1.0.0
Modified principles: N/A (initial version)
Added sections:
  - Core Principles (7 principles mapped from ultimate_principles.md)
  - Code Quality Standards
  - Development Workflow
  - Governance
Removed sections: N/A
Templates requiring updates:
  ✅ plan-template.md - Constitution Check section validated
  ✅ spec-template.md - Requirements alignment validated
  ✅ tasks-template.md - Task categorization validated
Follow-up TODOs: None
-->

## Core Principles

### I. Simplicity First (KISS)
**The simplest solution that works is the right solution.**

- Every feature starts with the most straightforward implementation possible
- Complex solutions require explicit justification in code reviews
- Before adding code, always attempt to delete or simplify existing code first
- Solutions with fewer moving parts are preferred over clever abstractions

**Rationale**: Simple code reduces cognitive load, minimizes bugs, and accelerates maintenance. Code that seems trivial is ideal—it means the solution is clear and understandable.

### II. Avoid Repetition (DRY with Rule of Three)
**Every piece of knowledge should exist once, but wait for patterns to emerge.**

- Copy once, paste twice, refactor on the third occurrence (Rule of Three)
- Abstractions MUST be justified by 3+ actual use cases
- No premature abstraction—let patterns emerge from working code
- Single source of truth for business logic, data schemas, and configuration

**Rationale**: Premature abstraction is worse than duplication. The wrong abstraction creates technical debt that's harder to fix than copied code. Three occurrences prove a pattern is real.

### III. Build Only What's Needed (YAGNI)
**Don't build features for imaginary futures.**

- Implement only features explicitly required by current specifications
- Speculative features, unused parameters, and "what if" code are prohibited
- Features are added when needed, not when anticipated
- Dead code paths must be removed immediately

**Rationale**: Speculative code bloats the codebase, increases maintenance burden, and often addresses problems that never materialize. Git preserves history if needed later.

### IV. Progressive Implementation (Make It Work → Right → Fast)
**Deliver working software, then improve iteratively.**

1. **First: Make It Work** - Achieve functional correctness
2. **Then: Make It Right** - Refactor for clarity and maintainability
3. **Only If Needed: Make It Fast** - Optimize based on measured performance

- Performance optimization requires profiling data showing actual bottlenecks
- Premature optimization is prohibited
- Code reviews enforce clarity before optimization

**Rationale**: Optimizing before understanding the problem wastes effort. Working code provides feedback; clear code ensures maintainability; fast code addresses proven needs.

### V. Explicit Over Magic
**Hidden complexity is still complexity. Make costs visible.**

- Dependencies must be explicitly passed as function parameters
- No hidden side effects in functions (especially getters)
- Framework "magic" (annotations, decorators, auto-wiring) requires justification
- Functions do exactly what their name suggests—nothing more

**Rationale**: Explicit dependencies make code testable, debuggable, and understandable. Magic behavior creates surprises and debugging nightmares at 3am.

### VI. Data Structures Over Algorithms
**Get the data model right, and the code writes itself.**

- Invest time designing appropriate data structures before writing algorithms
- Rich, well-structured data enables simple transformations
- Complex logic manipulating primitive types signals poor data modeling
- Data schemas are documented and versioned

**Rationale**: The right data structure naturally leads to simple, obvious code. Complex algorithms often compensate for inadequate data modeling.

### VII. Testability and Observability
**Code must be independently testable and debuggable.**

- Functions accept inputs as parameters and return outputs (no hidden state)
- External dependencies (databases, APIs, file systems) are passed explicitly
- Structured logging captures key decision points and state transitions
- Test isolation is mandatory—no shared mutable state between tests

**Rationale**: Testable code is maintainable code. Explicit dependencies enable mocking and testing. Observable systems are debuggable systems.

## Code Quality Standards

### Cognitive Load Limits
- Functions hold 7±2 concepts maximum
- Understanding a feature requires reading one file/module, not navigating 5+ files
- Indirection depth limited to 3 levels maximum
- No deep inheritance (maximum 1 level)

### Code Smells (Fix Immediately)
- **Abstraction Addiction**: Abstract classes with one implementation, interfaces for internal modules, factories creating one thing
- **State Soup**: Global variables, mutable shared state, objects that mutate themselves
- **Indirection Hell**: Wrappers around wrappers, functions that only call other functions
- **Magic Behavior**: Functions doing more than their name implies, side effects in getters

### Naming and Documentation
- Variable/function names should make code self-documenting
- Comments explain "why," not "what" (code should be clear enough to explain "what")
- If a comment is needed to explain what code does, refactor the code for clarity first
- Two Week Test: If you won't understand it in 2 weeks, simplify it now

## Development Workflow

### Before Writing Code
- [ ] Can I NOT write this code? (Apply YAGNI)
- [ ] Can I delete existing code instead? (Negative code is best)
- [ ] Is there a simpler solution? (Apply KISS)
- [ ] Have I designed the data structures first? (Data over algorithms)

### While Writing Code
- [ ] Am I repeating myself for the third time? (Apply Rule of Three, then DRY)
- [ ] Will I understand this in 2 weeks? (Two Week Test)
- [ ] Are dependencies explicit in function signatures? (No magic)
- [ ] Does this function do exactly what its name says? (No hidden side effects)

### Code Review Requirements
- [ ] Can any lines be deleted?
- [ ] Can any logic be simplified?
- [ ] Are variable/function names clear and accurate?
- [ ] Is each function independently testable?
- [ ] Are abstractions justified by 3+ use cases?
- [ ] Is complexity justified in writing? (See Complexity Tracking in plan.md)

### Architecture Evolution
- Start with functions over classes
- Start with copy/paste over wrong abstractions (wait for Rule of Three)
- Start with synchronous over async (until proven necessary)
- Start with monolith over microservices (until proven necessary)
- Let patterns emerge from repetition—extract abstractions from working code
- Design for deletion, not permanence

## Governance

### Amendment Process
This constitution supersedes all other development practices. Amendments require:
1. Documented justification for the change
2. Review and approval from project maintainers
3. Migration plan if changes affect existing code
4. Version increment following semantic versioning

### Compliance and Review
- All code reviews MUST verify compliance with these principles
- Violations of principles (complexity, abstraction, magic) MUST be justified in writing
- Use the Complexity Tracking table in `plan.md` when violations are necessary
- Unjustified complexity blocks merge approval

### Decision Framework
When stuck between solutions, choose based on:
1. Which has fewer moving parts?
2. Which is easier to delete?
3. Which is more boring?
4. Which would confuse you at 3am? (Choose the other)

### Version Control
**Version**: 1.0.0 | **Ratified**: 2025-10-23 | **Last Amended**: 2025-10-23
