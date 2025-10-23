# Implementation Tasks: Unused Assets Detector

**Feature**: Unused Assets Detector
**Branch**: `001-unused-assets-detector`
**Generated**: 2025-10-23

## Overview

This document contains actionable implementation tasks for the Unused Assets Detector CLI tool. Tasks are organized by user story to enable independent, incremental development and testing.

**MVP Scope**: User Story 1 (Scan Project and Identify Unused Assets)

---

## Task Statistics

- **Total Tasks**: 87
- **Setup Tasks**: 6
- **Foundational Tasks**: 8
- **User Story 1 (P1)**: 29 tasks
- **User Story 2 (P2)**: 15 tasks
- **User Story 3 (P3)**: 12 tasks
- **User Story 4 (P4)**: 10 tasks
- **Polish Tasks**: 7 tasks

---

## Phase 1: Setup (Project Initialization)

**Goal**: Initialize Go project structure, install dependencies, and set up development environment.

**Tasks**:

- [X] T001 Initialize Go module with `go mod init github.com/yourusername/asset-cleaner` in repository root
- [X] T002 [P] Install Cobra CLI framework v1.9.1+ via `go get github.com/spf13/cobra@v1.9.1`
- [X] T003 [P] Install Bubble Tea v1.3.10 via `go get github.com/charmbracelet/bubbletea@v1.3.10`
- [X] T004 [P] Install Lip Gloss v2.x via `go get github.com/charmbracelet/lipgloss@latest`
- [X] T005 [P] Install Viper v1.21.0 via `go get github.com/spf13/viper@v1.21.0`
- [X] T006 Create project directory structure as defined in plan.md (cmd/, internal/, tests/, etc.)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Goal**: Build core data structures and utilities that all user stories depend on.

**Dependencies**: Must complete Phase 1 first.

**Tasks**:

- [X] T007 Create AssetFile struct in internal/models/asset.go with all attributes from data-model.md
- [X] T008 Create Reference struct in internal/models/reference.go with all attributes from data-model.md
- [X] T009 Create ScanResult struct in internal/models/scan_result.go with statistics and methods
- [X] T010 Create ProjectConfig struct in internal/models/project_config.go with configuration fields
- [X] T011 Create ProjectType enum and detection helpers in internal/detector/project_type.go
- [X] T012 Implement DefaultConfig() function in internal/config/defaults.go returning default project configuration
- [X] T013 [P] Implement file system utilities (exists, isDir, etc.) in internal/utils/fileutil.go
- [X] T014 [P] Implement path manipulation utilities (normalize, relative, etc.) in internal/utils/pathutil.go

---

## Phase 3: User Story 1 - Scan Project and Identify Unused Assets (P1)

**User Story**: A developer wants to clean up their project by identifying which asset files are not referenced anywhere in their codebase. They run a scan command and receive a clear report of unused files.

**Why P1**: Core value proposition and MVP. Delivers immediate value by answering "what can I safely delete?"

**Independent Test Criteria**:
- ✅ Can scan a test project with known unused files
- ✅ Report correctly identifies unused assets (95%+ accuracy)
- ✅ Completes scan of 1,000 files in <10 seconds
- ✅ Shows progress updates during long scans
- ✅ Delivers standalone value without deletion capability

**Dependencies**: Phase 2 complete.

**Tasks**:

### 3.1: CLI Command Setup (US1)

- [X] T015 [US1] Create root command with global flags in cmd/asset-cleaner/commands/root.go
- [X] T016 [US1] Implement scan command skeleton in cmd/asset-cleaner/commands/scan.go
- [X] T017 [US1] Add scan command flags (--extensions, --exclude, --output, --format, --no-progress)
- [X] T018 [US1] Implement global flag handling (--config, --verbose, --quiet, --no-color, --version, --help)

### 3.2: Project Type Detection (US1)

- [X] T019 [P] [US1] Implement DetectProjectType() for package.json analysis in internal/detector/project_type.go
- [X] T020 [P] [US1] Add detection for Flutter (pubspec.yaml) in internal/detector/project_type.go
- [X] T021 [P] [US1] Add detection for iOS (.xcodeproj) in internal/detector/project_type.go
- [X] T022 [P] [US1] Add detection for Android (build.gradle) in internal/detector/project_type.go
- [X] T023 [US1] Implement DefaultAssetPaths() method for each ProjectType in internal/config/defaults.go
- [X] T024 [US1] Implement DefaultExtensions() method for each ProjectType in internal/config/defaults.go

### 3.3: Asset Discovery (US1)

- [X] T025 [US1] Implement single-threaded asset scanner using filepath.WalkDir in internal/scanner/asset_finder.go
- [X] T026 [US1] Add extension filtering logic in asset_finder.go
- [X] T027 [US1] Add directory exclusion logic (node_modules, dist, build, etc.) in asset_finder.go
- [X] T028 [US1] Implement AssetFile creation with metadata (size, modTime, hash) in asset_finder.go

### 3.4: Reference Detection (US1)

- [X] T029 [P] [US1] Implement regex-based parser patterns in internal/parser/patterns.go (simplified, no AST yet)
- [X] T030 [P] [US1] Extract import statements using regex patterns
- [X] T031 [P] [US1] Extract string literals containing asset paths using regex patterns
- [X] T032 [P] [US1] Extract template literals with asset paths using regex patterns
- [X] T033 [US1] Implement regex-based fallback parser patterns
- [X] T034 [US1] Add common reference patterns (CSS url(), HTML src/href, require()) in internal/parser/patterns.go
- [X] T035 [US1] Implement reference finder that scans source files in internal/scanner/reference_finder.go
- [X] T036 [US1] Create Reference objects with source location and context in reference_finder.go

### 3.5: Asset Classification (US1)

- [X] T037 [US1] Implement ClassifyAsset() function using classification logic from data-model.md in internal/classifier/classifier.go
- [X] T038 [US1] Add StatusUsed detection (has active references) in classifier.go
- [X] T039 [US1] Add StatusUnused detection (zero references) in classifier.go
- [X] T040 [US1] Add StatusPotentiallyUnused detection (only in comments) in classifier.go
- [X] T041 [US1] Add StatusNeedsManualReview detection (dynamic paths) in classifier.go

### 3.6: Progress Indicators & Output (US1)

- [X] T042 [US1] Create basic text formatter in internal/ui/formatter.go (simplified, no Bubble Tea progress yet)
- [X] T043 [US1] Implement basic progress output during scan
- [X] T044 [US1] Create basic styled output with header and separators in internal/ui/formatter.go
- [X] T045 [US1] Implement text format output for scan results in commands/scan.go
- [X] T046 [US1] Implement JSON format output for scan results in commands/scan.go
- [X] T047 [US1] Implement CSV format output for scan results in commands/scan.go
- [X] T048 [US1] Add ScanStatistics computation in internal/models/scan_result.go

---

## Phase 4: User Story 2 - Preview and Confirm Before Deletion (P2)

**User Story**: A developer reviews scan results and wants to safely remove unused assets with preview, file sizes/savings, and confirmation.

**Why P2**: Safety is critical for adoption. Preview builds user confidence and shows value.

**Independent Test Criteria**:
- ✅ Can preview deletions without modifying files
- ✅ Shows accurate file sizes and total savings
- ✅ Deletion only removes specified files
- ✅ Preserves git history when deleting
- ✅ Provides recovery instructions

**Dependencies**: User Story 1 complete.

**Tasks**:

### 4.1: Delete Command (US2)

- [X] T049 [US2] Create delete command skeleton in cmd/asset-cleaner/commands/delete.go
- [X] T050 [US2] Add delete command flags (--dry-run, --interactive, --force, --scan-file)
- [X] T051 [US2] Implement dry-run mode showing preview without deletion in delete.go
- [X] T052 [US2] Implement confirmation prompt for deletion in delete.go

### 4.2: File Deletion Logic (US2)

- [X] T053 [US2] Implement safe file deletion using os.Remove in delete.go
- [X] T054 [US2] Add git repository detection (check for .git/) in delete.go
- [X] T055 [US2] Add warning if not in git repository in delete.go
- [X] T056 [US2] Implement interactive mode (prompt per file) in delete.go
- [X] T057 [US2] Add deletion statistics (count, total bytes freed) in delete.go

### 4.3: Review Web UI (US2)

- [X] T058 [P] [US2] Create review command skeleton in cmd/asset-cleaner/commands/review.go
- [X] T059 [P] [US2] Implement embedded web UI files using go:embed in internal/ui/web/index.html
- [X] T060 [P] [US2] Create web server using net/http in internal/ui/server.go
- [X] T061 [US2] Add GET /api/results endpoint serving scan results as JSON in server.go
- [X] T062 [US2] Add POST /api/delete endpoint for batch deletion in server.go
- [X] T063 [US2] Implement browser auto-launch in review.go

---

## Phase 5: User Story 3 - Exclude/Ignore Specific Assets (P3)

**User Story**: A developer marks certain assets (used dynamically or kept for future) to be excluded from scan results.

**Why P3**: Handles edge cases where static analysis fails. Important for production use.

**Independent Test Criteria**:
- ✅ Can configure ignore patterns via .unusedassets.yaml
- ✅ Ignored assets excluded from unused report
- ✅ Multiple patterns supported (glob, regex)
- ✅ Can verify what's being excluded

**Dependencies**: User Story 1 complete.

**Tasks**:

### 5.1: Configuration File Support (US3)

- [X] T064 [P] [US3] Implement config file loader using Viper in internal/config/loader.go
- [X] T065 [P] [US3] Add support for .unusedassets.yaml parsing in loader.go
- [X] T066 [US3] Create init command to generate config file in cmd/asset-cleaner/commands/init.go
- [X] T067 [US3] Add config template generation (default, minimal, comprehensive) in init.go

### 5.2: Ignore Pattern Matching (US3)

- [ ] T068 [US3] Create IgnoreRule struct in internal/models/ignore_rule.go
- [ ] T069 [US3] Implement glob pattern matching in ignore_rule.go
- [ ] T070 [US3] Implement regex pattern matching in ignore_rule.go
- [ ] T071 [US3] Create ignore command to add/remove patterns in cmd/asset-cleaner/commands/ignore.go

### 5.3: Dynamic Path Detection (US3)

- [ ] T072 [US3] Implement constant declaration extraction in internal/parser/javascript.go
- [ ] T073 [US3] Implement constant usage tracking (template literals, concatenation) in javascript.go
- [ ] T074 [US3] Flag assets with dynamic references as StatusNeedsManualReview in internal/classifier/classifier.go
- [ ] T075 [US3] Add configuration option for base_path_vars in internal/config/config.go

---

## Phase 6: User Story 4 - Generate Scan Reports for Team Review (P4)

**User Story**: A developer exports scan results in formats suitable for code review or team discussion.

**Why P4**: Enables team workflows but not essential for individual use.

**Independent Test Criteria**:
- ✅ Can export to JSON with complete data
- ✅ Can export to CSV for spreadsheets
- ✅ Can generate HTML report
- ✅ Reports contain accurate information

**Dependencies**: User Story 1 complete.

**Tasks**:

### 6.1: Report Generation (US4)

- [X] T076 [P] [US4] Implement ToJSON() method in internal/models/scan_result.go (already implemented)
- [X] T077 [P] [US4] Implement ToCSV() method in internal/models/scan_result.go (already implemented)
- [ ] T078 [US4] Create HTML report template in cmd/asset-cleaner/web/report.html
- [ ] T079 [US4] Implement HTML report generation in internal/ui/server.go

### 6.2: Info Command (US4)

- [X] T080 [P] [US4] Create info command skeleton in cmd/asset-cleaner/commands/info.go
- [X] T081 [P] [US4] Implement project info display (type, paths, config) in info.go
- [X] T082 [US4] Add --show-config flag displaying current configuration in info.go
- [X] T083 [US4] Add --show-paths flag listing detected asset paths in info.go

### 6.3: Advanced Reporting (US4)

- [ ] T084 [US4] Add report comparison functionality (track trends over time) in scan_result.go
- [ ] T085 [US4] Implement sortable/filterable HTML table in web/report.html

---

## Phase 7: Polish & Cross-Cutting Concerns

**Goal**: Error handling, logging, documentation, performance optimization, and CI/CD setup.

**Dependencies**: All user stories complete.

**Tasks**:

- [ ] T086 Add structured logging at key decision points (project detection, reference matching) in all scanner files
- [ ] T087 [P] Create comprehensive error messages with context and suggestions per CLI contract in all commands
- [ ] T088 Implement concurrent scanning for projects >10K files in internal/scanner/concurrent.go
- [ ] T089 Add worker pool pattern with runtime.NumCPU() workers in concurrent.go
- [X] T090 Create example .unusedassets.yaml config file in repository root
- [X] T091 Write README.md with installation, quick start, and usage examples
- [X] T092 Set up GitHub Actions CI with matrix testing (macOS, Linux, Windows)

---

## Dependencies & Execution Order

### User Story Completion Order

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundational)
    ↓
Phase 3 (US1: Scan - P1) ← MVP
    ↓
    ├─→ Phase 4 (US2: Delete - P2)
    ├─→ Phase 5 (US3: Ignore - P3)
    └─→ Phase 6 (US4: Reports - P4)
    ↓
Phase 7 (Polish)
```

**Critical Path**: Phase 1 → Phase 2 → Phase 3 (US1)

**Parallelizable**: US2, US3, US4 can be implemented independently after US1 completes.

---

## Parallel Execution Opportunities

### Within User Story 1 (US1)

**Parallel Group 1** (after T018):
- T019-T022: Project type detection (different frameworks)
- T029-T032: Parser components (can work independently)

**Parallel Group 2** (after T024):
- T025-T028: Asset discovery
- T033-T034: Regex parser and patterns

**Parallel Group 3** (after T041):
- T042-T044: UI components
- T045-T047: Output formats

### Across User Stories

After US1 complete (T048), these can run in parallel:
- **Track A**: T049-T063 (US2: Deletion)
- **Track B**: T064-T075 (US3: Ignore patterns)
- **Track C**: T076-T085 (US4: Reports)

### Polish Phase

All polish tasks (T086-T092) marked with [P] can run concurrently.

---

## Implementation Strategy

### MVP-First Approach

**Iteration 1** (MVP - Deliver Value Fast):
- Phase 1: Setup (T001-T006)
- Phase 2: Foundational (T007-T014)
- Phase 3: User Story 1 only (T015-T048)
- **Outcome**: Working CLI that scans and identifies unused assets

**Iteration 2** (Add Safety):
- Phase 4: User Story 2 (T049-T063)
- **Outcome**: Users can safely delete assets with preview

**Iteration 3** (Handle Edge Cases):
- Phase 5: User Story 3 (T064-T075)
- **Outcome**: Support for dynamic assets and ignore patterns

**Iteration 4** (Team Workflows):
- Phase 6: User Story 4 (T076-T085)
- **Outcome**: Export and sharing capabilities

**Iteration 5** (Production Ready):
- Phase 7: Polish (T086-T092)
- **Outcome**: Robust, well-documented, performant tool

---

## Testing Strategy

### Unit Tests (Per Task)

Each task should include unit tests covering:
- Happy path (normal inputs)
- Edge cases (empty inputs, large inputs)
- Error cases (invalid inputs, file not found)

**Example**:
- T025 (Asset scanner): Test with 0 files, 1 file, 1000 files, excluded directories
- T037 (Classifier): Test all status transitions (Used, Unused, PotentiallyUnused, NeedsManualReview)

### Integration Tests (Per User Story)

Each user story phase should include integration tests:

**US1**:
- Scan fixture project with known unused files
- Verify report accuracy (95%+ detection rate)
- Measure performance (<10s for 1K files)

**US2**:
- Test dry-run doesn't modify files
- Test deletion removes correct files
- Test git history preservation

**US3**:
- Test ignore patterns exclude files correctly
- Test dynamic path detection flags assets

**US4**:
- Test JSON/CSV export completeness
- Test HTML report generation

### Performance Benchmarks

- T088-T089: Benchmark concurrent vs. single-threaded scanning
- Target: 1K files in <10s, handle 100K files without crashing

---

## Success Metrics Validation

Map tasks to success criteria from spec.md:

| Metric | Task(s) | Validation Method |
|--------|---------|-------------------|
| **SC-001** (1K files <10s) | T025, T088-T089 | Benchmark test with timer |
| **SC-002** (95%+ accuracy) | T029-T041 | Integration tests with fixtures |
| **SC-003** (<5% false positives) | T037-T041 | Accuracy measurement |
| **SC-004** (Safe deletion) | T053-T055 | Manual testing + user feedback |
| **SC-005** (Workflow <5min) | All US1+US2 | User testing with timer |
| **SC-006** (Handle 100K files) | T088-T089 | Stress test with synthetic project |
| **SC-007** (90%+ no docs) | T015-T018, T045 | Usability testing with 10 users |
| **SC-008** (Git recovery) | T054-T055 | Manual test: delete, restore via git |
| **SC-009** (10-30% reduction) | All US1+US2 | Measure on 5 real projects |
| **SC-010** (Progress every 2s) | T042-T043 | Timer validation in long scans |

---

## Risk Mitigation per Task

### High Risk: False Positives

**Tasks**: T029-T041 (Reference detection & classification)

**Mitigation**:
- Conservative detection (prefer false negatives)
- Three-tier classification (Used, PotentiallyUnused, Unused)
- Comprehensive test fixtures with edge cases
- StatusNeedsManualReview for ambiguous cases

### Medium Risk: Performance

**Tasks**: T025, T088-T089 (Scanning)

**Mitigation**:
- Early directory skipping (T027)
- Concurrent scanning for large projects (T088-T089)
- Progress indicators (T042-T043)
- Memory profiling and optimization

### Medium Risk: Language Support

**Tasks**: T029-T034 (Parsing)

**Mitigation**:
- Start with JS/TS (80% of use cases)
- Regex fallback (T033-T034)
- Document limitations clearly
- Make parser system extensible

---

## Completion Checklist

Before marking feature complete, verify:

- [ ] All tasks T001-T092 completed
- [ ] All unit tests passing
- [ ] Integration tests passing for each user story
- [ ] Performance benchmarks meet targets (SC-001, SC-006, SC-010)
- [ ] Accuracy metrics validated (SC-002, SC-003)
- [ ] Documentation complete (README, quickstart, config examples)
- [ ] CI/CD pipeline running on all platforms
- [ ] Manual testing on real projects (React, Vue, Flutter)
- [ ] User testing validates <5min workflow (SC-005, SC-007)
- [ ] Git integration tested (SC-008)

---

## Next Steps

1. **Review tasks with team** - Ensure all requirements covered
2. **Set up tracking** - Move tasks to project management tool
3. **Start MVP** - Begin with Phase 1 (T001-T006)
4. **Iterate rapidly** - Complete US1 first, then add features
5. **Test continuously** - Write tests alongside implementation
6. **Gather feedback** - Test with real users after US1+US2 complete

---

**Tasks Ready for Implementation** ✅

**Total Estimated Effort**: ~6-8 weeks for full feature (1-2 weeks for MVP)
