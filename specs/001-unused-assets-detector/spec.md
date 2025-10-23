# Feature Specification: Unused Assets Detector

**Feature Branch**: `001-unused-assets-detector`
**Created**: 2025-10-23
**Status**: Draft
**Input**: User description: "Find and remove unused assets from your codebase - automatically and safely."

## Clarifications

### Session 2025-10-23

- Q: How should assets referenced only in comments or dead code branches be handled? → A: Detect and warn separately ("potentially unused - found only in comments")
- Q: What safety mechanisms are employed when deleting unused assets? → A: Delete from filesystem only, preserve git history
- Q: How should assets with dynamic path construction patterns be handled? → A: Flag as "needs manual review", exclude from auto-deletion
- Q: What configuration format and location should be used for ignore patterns? → A: Standard config file (.unusedassets or similar) in project root
- Q: How should common build/dependency directories be handled? → A: Auto-exclude common build/dependency directories by default (node_modules, dist, build, etc.)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Scan Project and Identify Unused Assets (Priority: P1)

A developer wants to clean up their project by identifying which asset files (images, fonts, videos, etc.) are not referenced anywhere in their codebase. They run a scan command and receive a clear report of unused files.

**Why this priority**: This is the core value proposition and minimum viable product. Without the ability to identify unused assets, nothing else matters. This delivers immediate value by answering the question "what can I safely delete?"

**Independent Test**: Can be fully tested by running the scan command on a test project with known unused files and verifying the report correctly identifies them. Delivers standalone value even without deletion capability.

**Acceptance Scenarios**:

1. **Given** a project with 100 asset files where 20 are not referenced in code, **When** user runs the scan command, **Then** system generates a report listing exactly those 20 unused files
2. **Given** a project with all assets actively used, **When** user runs the scan command, **Then** system reports "No unused assets found"
3. **Given** a large project with 10,000+ files, **When** user runs the scan command, **Then** system completes the scan and provides progress updates during analysis
4. **Given** assets referenced using different patterns (direct imports, dynamic paths, CSS references), **When** user runs the scan, **Then** system correctly identifies all reference types and marks those assets as used

---

### User Story 2 - Preview and Confirm Before Deletion (Priority: P2)

A developer reviews the scan results and wants to safely remove unused assets. They can preview what will be deleted, see file sizes/savings, and confirm the operation with confidence that they won't break anything.

**Why this priority**: Safety is critical for adoption. Users need confidence that deleting assets won't break their app. Preview functionality provides that safety net and shows the value (disk space savings).

**Independent Test**: Can be tested independently by scanning a project, requesting a preview of deletions, and verifying the preview shows accurate file information without actually deleting anything. Delivers value by showing potential savings and building user confidence.

**Acceptance Scenarios**:

1. **Given** a scan has identified 50 unused assets totaling 15MB, **When** user requests deletion preview, **Then** system displays the list with file sizes and total savings before confirming
2. **Given** a deletion preview showing 20 files, **When** user confirms deletion, **Then** system removes only those files and provides a summary of what was deleted
3. **Given** a deletion preview showing critical-looking files, **When** user cancels the operation, **Then** no files are modified and user can re-examine the results
4. **Given** unused assets in version control, **When** user confirms deletion, **Then** system removes files from filesystem (preserving git history) and provides instructions for staging and committing the deletions

---

### User Story 3 - Exclude/Ignore Specific Assets (Priority: P3)

A developer knows certain assets are used dynamically (e.g., loaded by filename from a database) or wants to keep specific files for future use. They can mark these assets to be excluded from scan results.

**Why this priority**: Handles real-world edge cases where static analysis can't detect usage. Important for production use but not required for initial validation of the core concept.

**Independent Test**: Can be tested by creating an ignore configuration, running a scan, and verifying that ignored assets don't appear in the unused report even if they're not statically referenced. Delivers value by making the tool usable for complex projects.

**Acceptance Scenarios**:

1. **Given** a .unusedassets config file in project root specifying certain paths to ignore, **When** user runs a scan, **Then** system excludes those paths from the unused assets report
2. **Given** a scan result showing a false positive, **When** user adds that file pattern to .unusedassets config file, **Then** subsequent scans don't flag that file as unused
3. **Given** multiple patterns to ignore (e.g., "*.backup.png", "legacy/*"), **When** user runs a scan, **Then** system correctly excludes all matching files
4. **Given** an ignore configuration, **When** user wants to verify what's being excluded, **Then** system can show which files match ignore rules

---

### User Story 4 - Generate Scan Reports for Team Review (Priority: P4)

A developer runs a scan and wants to share the results with their team before making deletion decisions. They can export the results in formats suitable for code review or team discussion.

**Why this priority**: Enables team workflows and decision-making but not essential for individual developer use. Nice-to-have for organizational adoption.

**Independent Test**: Can be tested by running a scan and exporting results to different formats (JSON, CSV, HTML) and verifying each format contains complete and accurate information. Delivers value by enabling team collaboration.

**Acceptance Scenarios**:

1. **Given** a completed scan, **When** user exports results to JSON, **Then** system creates a structured file containing all unused assets with metadata
2. **Given** scan results, **When** user generates an HTML report, **Then** system creates a human-readable report with sortable/filterable results
3. **Given** multiple scan runs over time, **When** user compares reports, **Then** they can see trends in unused assets (growing/shrinking)

---

### Edge Cases

- Asset files referenced only in comments or dead code branches are flagged separately as "potentially unused" with a warning, distinct from confirmed unused assets
- How does the system handle assets with the same name in different directories?
- What happens when symbolic links or hard links point to assets?
- How does the system handle assets referenced in minified or obfuscated code?
- Common build and dependency directories (node_modules, dist, build, .next, target, vendor, etc.) are automatically excluded from scanning to avoid analyzing generated or third-party assets
- How does the system handle assets referenced in configuration files (JSON, YAML, XML)?
- Assets referenced via template strings or concatenated paths that can't be statically analyzed are flagged as "needs manual review" and excluded from automatic deletion
- How does the system distinguish between production assets and development/test assets?
- What happens when the user lacks file permissions to read certain directories?
- How does the system handle very large binary files (videos, archives) during scanning?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST scan a specified project directory recursively to identify all asset files
- **FR-001a**: System MUST automatically exclude common build and dependency directories by default (e.g., node_modules, dist, build, .next, target, vendor) unless explicitly overridden in configuration
- **FR-002**: System MUST detect references to assets across multiple file types (source code, stylesheets, markup, configuration files)
- **FR-003**: System MUST distinguish between used and unused assets based on detected references
- **FR-003a**: System MUST categorize assets into three tiers: used (active references), potentially unused (references only in comments/dead code), and unused (no references found)
- **FR-004**: System MUST generate a report listing all unused assets with their file paths
- **FR-005**: System MUST display file sizes and total disk space savings for unused assets
- **FR-006**: System MUST support configurable asset file extensions (images, fonts, videos, audio, etc.)
- **FR-007**: System MUST provide a preview mode that shows what would be deleted without actually deleting files
- **FR-008**: System MUST support safe deletion of unused assets after user confirmation by removing files from filesystem while preserving version control history (git does not delete from history)
- **FR-009**: System MUST allow users to specify ignore patterns for files/directories that should never be flagged as unused via a standard configuration file (e.g., .unusedassets) in the project root directory
- **FR-010**: System MUST show progress indication during long-running scan operations
- **FR-011**: System MUST handle common reference patterns including direct imports, require statements, URL references, and path strings
- **FR-012**: System MUST support exporting scan results in multiple formats (JSON, CSV, plain text)
- **FR-013**: System MUST operate without requiring network connectivity
- **FR-014**: System MUST work on projects of varying sizes (from small apps to large codebases)
- **FR-015**: System MUST preserve version control history when deleting files (not delete from git history)
- **FR-016**: System MUST handle projects with multiple source directories and asset locations
- **FR-017**: System MUST detect dynamic path construction patterns where feasible, flag affected assets as "needs manual review", and exclude them from automatic deletion operations
- **FR-018**: System MUST provide clear error messages when it cannot access files or directories due to permissions

### Assumptions

- **A-001**: Users have read access to all project files they want to scan
- **A-002**: Asset references follow standard conventions (import statements, file paths, URL patterns)
- **A-003**: Users understand that dynamic path construction (runtime string concatenation) may cause false positives
- **A-004**: Projects follow common directory structures (assets in dedicated folders like /images, /fonts, /static, etc.)
- **A-005**: Users run the tool from the project root directory or specify the correct path
- **A-006**: The tool operates on the current working copy of files, not analyzing git history
- **A-007**: Users have write permissions when attempting to delete files
- **A-008**: Standard file extensions are sufficient for most use cases (jpg, png, svg, ttf, woff, mp4, etc.)
- **A-009**: Users do not want to analyze assets in build/dependency directories (node_modules, dist, etc.) as these contain generated or third-party files

### Key Entities

- **Asset File**: A file in the project that is not source code but is used by the application (images, fonts, videos, audio, icons, etc.). Attributes include: file path, file size, file type/extension, last modified date.

- **Reference**: An occurrence in source code or configuration where an asset file is referenced. Attributes include: source file path, line number, reference type (import, URL, path string), matched pattern.

- **Scan Result**: The output of analyzing a project. Contains lists of used assets, potentially unused assets (references only in comments/dead code), assets needing manual review (dynamic path construction detected), unused assets, scan metadata (timestamp, project path, configuration used), and statistics (total assets, unused count, potentially unused count, manual review count, total size savings).

- **Ignore Rule**: A pattern or path specification that excludes certain files/directories from being flagged as unused. Attributes include: pattern (glob or regex), reason for exclusion, whether it's enabled.

- **Project Configuration**: User-specified settings for how to scan their project, stored in a configuration file (e.g., .unusedassets) in the project root directory. Includes: asset directories to scan, file extensions to consider, ignore patterns, reference detection rules, and output preferences. Can be version-controlled for team consistency.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can scan a project with 1,000 files and receive results in under 10 seconds on a standard developer machine
- **SC-002**: System correctly identifies 95%+ of unused assets in typical web/mobile projects (measured against hand-verified test projects)
- **SC-003**: False positive rate (flagging used assets as unused) is below 5% for common reference patterns
- **SC-004**: Users can safely delete unused assets and confirm their application still functions correctly
- **SC-005**: Users can complete a full workflow (scan → review → delete) in under 5 minutes for a typical project
- **SC-006**: System handles projects with up to 100,000 files without crashing or hanging
- **SC-007**: 90%+ of users can run the tool successfully without consulting documentation (measure via user testing)
- **SC-008**: Deleted files are recoverable via version control commands (e.g., git restore, git checkout) since deletion only removes from filesystem, not from git history
- **SC-009**: Tool reduces project size by an average of 10-30% for projects with asset bloat (measure across test projects)
- **SC-010**: Progress indication updates at least every 2 seconds during long-running scans to maintain user confidence
