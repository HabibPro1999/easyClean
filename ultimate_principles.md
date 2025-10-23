# Ultimate Engineering Principles
*Development constitution for writing clean, simple, scalable code*

## Core Rules (Always Apply First)

### The Big Three
**#KISS**: Keep It Simple, Stupid. The simplest solution that works is the right solution.

**#DRY**: Don't Repeat Yourself. Every piece of knowledge should exist once.

**#YAGNI**: You Aren't Gonna Need It. Don't build features for imaginary futures.

### The Process
**#MakeItWork** → **#MakeItRight** → **#MakeItFast**
1. First: Get it working (correctness)
2. Then: Clean it up (clarity)
3. Only if needed: Optimize (performance)

### The Thresholds
**#RuleOfThree**: Copy once, paste twice, refactor on the third.

**#TwoWeekTest**: If you won't understand it in 2 weeks, simplify it now.

**#DeleteFirst**: Before adding code, try to delete code.

## Deep Principles (Why The Rules Work)

### Cognitive Load
**Your brain can hold 7±2 concepts at once**. Every abstraction, every file jump, every piece of state consumes this limited resource.

```
BAD:  Navigate through 5 files to understand one feature
GOOD: Read one file/function and understand completely
```

### Data Over Logic
**Data structures matter more than algorithms**. Get the data right, and the code writes itself.

```
BAD:  Complex logic manipulating simple arrays
GOOD: Rich data structures with simple transformations
```

### Explicit Over Magic
**Hidden complexity is still complexity**. Make costs visible.

```
BAD:  @AutoMagic @DoEverything function process() { }
GOOD: function process(data, db, cache) { } // Dependencies visible
```

## Code Smells (Fix Immediately)

### 🚨 Abstraction Addiction
- Abstract classes with one implementation
- Interfaces for internal modules
- Factories that create one thing
- **Fix**: Delete the abstraction, use the concrete implementation

### 🚨 State Soup
- Global variables
- Mutable shared state
- Objects that change themselves
- **Fix**: Pass values, return new values

### 🚨 Indirection Hell
- Code that just calls other code
- Wrappers around wrappers
- 5+ levels of function calls
- **Fix**: Inline the indirection

### 🚨 Speculative Generality
- "What if we need..." code
- Unused parameters
- Dead code paths
- **Fix**: Delete it. Git remembers everything.

### 🚨 Magic Behavior
- Functions that do more than their name
- Side effects in getters
- Surprising mutations
- **Fix**: Make functions do one thing, name them honestly

## Decision Framework

When stuck between solutions:

1. **Which has fewer moving parts?** → Choose that
2. **Which is easier to delete?** → Choose that
3. **Which is more boring?** → Choose that
4. **Which would confuse you at 3am?** → Choose the other

## Writing New Code Checklist

Before writing:
- [ ] Can I NOT write this code? (YAGNI)
- [ ] Can I delete code instead? (Negative code is best code)
- [ ] Is there a simpler solution? (KISS)

While writing:
- [ ] Am I repeating myself? (DRY - but wait for 3 instances)
- [ ] Will I understand this in 2 weeks? (Two Week Test)
- [ ] Are dependencies explicit? (No magic)

After writing:
- [ ] Can I delete any lines?
- [ ] Can I simplify any logic?
- [ ] Can I improve any names?

## Architecture Guidelines

### Start Simple
- Functions over classes
- Copy/paste over wrong abstraction
- Synchronous over async (until proven necessary)
- Monolith over microservices (until proven necessary)

### Grow Naturally
- Let patterns emerge from repetition
- Extract abstractions from working code
- Design for deletion, not permanence

### Stay Focused
- Each module does ONE thing
- Dependencies flow ONE direction
- Changes affect ONE place

## The Anti-Patterns Never to Use

❌ **Premature Optimization**: Optimizing without measuring
❌ **Premature Abstraction**: Abstracting without repetition
❌ **Framework First**: Reaching for frameworks before trying simpler solutions
❌ **Clever Code**: Code that shows off rather than solves
❌ **God Objects**: Classes that do everything
❌ **Deep Inheritance**: More than 1 level of inheritance
❌ **Shared Mutable State**: The root of all bugs

## The Golden Rule

**When in doubt, choose the solution with fewer parts.**

Less code = Less bugs = Less maintenance = More happiness

## Quick Reference Card

```
DECISION TIME:
1. Can I avoid writing this?     → Don't write it
2. Can I delete code instead?     → Delete it
3. Is this the simplest approach? → Simplify it
4. Will I understand this later?  → Clarify it
5. Am I guessing at the future?   → Stop guessing

CODE REVIEW:
- Every file: Can it be shorter?
- Every function: Does it do one thing?
- Every line: Is it necessary?
- Every abstraction: Is it used 3+ times?
- Every comment: Could better naming replace it?
```

Remember: **The best code is code that doesn't exist. The second best is code that's so simple it seems trivial.**