# Task Management Rules with Testing Requirements

## Core Principles

### 1. Task Creation Protocol
- **Mandatory Breakdown**: All complex tasks must be broken into smaller, manageable subtasks
- **Testing Integration**: Each task must include specific test requirements from the start
- **Dependency Mapping**: Clear identification of task dependencies and prerequisites
- **Time Estimation**: Realistic time estimates for each task component

### 2. Task Structure Template
```markdown
## Task: [Task Name]
**Priority**: [High/Medium/Low]
**Status**: [Pending/In Progress/Completed/Blocked]
**Estimated Time**: [Hours]

### Description
[Clear, concise description of what needs to be accomplished]

### Prerequisites
- [ ] Dependency 1
- [ ] Dependency 2

### Subtasks
1. [ ] Subtask 1
2. [ ] Subtask 2
3. [ ] Subtask 3

### Testing Requirements
- [ ] Unit tests for new functions
- [ ] Integration tests for component interactions
- [ ] Edge case testing
- [ ] Performance benchmarks (if applicable)
- [ ] Security validation (if applicable)

### Acceptance Criteria
- [ ] All tests pass
- [ ] Code review completed
- [ ] Documentation updated
- [ ] Backup created (if required)
```

### 3. Testing-First Development

#### Mandatory Test Requirements
- **Unit Tests**: Must be written before or alongside implementation
- **Integration Tests**: Required for multi-component features
- **Regression Tests**: Must pass for all existing functionality
- **Performance Tests**: Required for any performance-related changes
- **Security Tests**: Mandatory for authentication, data handling, or API changes

#### Test Coverage Standards
- **Minimum Coverage**: 80% line coverage for new code
- **Critical Path Coverage**: 100% coverage for business logic
- **Error Handling**: All error paths must be tested
- **Edge Cases**: Boundary conditions and null/undefined inputs

### 4. Task Completion Validation

#### Pre-Completion Checklist
- [ ] All implemented features work as specified
- [ ] All tests pass (unit, integration, regression)
- [ ] Code follows project conventions and style guides
- [ ] Documentation is updated (README, API docs, comments)
- [ ] No console errors or warnings
- [ ] Performance meets requirements
- [ ] Security review passed (if applicable)
- [ ] Backup created for any risky changes

#### Final Validation Steps
1. **Test Suite Execution**: Run full test suite and verify 100% pass rate
2. **Manual Testing**: Verify user interface and user experience
3. **Code Review**: Peer review or self-review using checklist
4. **Documentation Sync**: Ensure all documentation reflects changes
5. **Backup Verification**: Confirm backups are created and accessible

### 5. Task Status Management

#### Status Definitions
- **Pending**: Task is queued and not yet started
- **In Progress**: Active work is being performed
- **Testing**: Implementation complete, undergoing testing
- **Review**: Testing complete, awaiting review
- **Completed**: All requirements met, tests passing, documentation updated
- **Blocked**: Cannot proceed due to dependencies or issues

#### Status Transition Rules
- **Pending → In Progress**: Only when all prerequisites are met
- **In Progress → Testing**: Only when implementation is complete
- **Testing → Review**: Only when all tests pass
- **Review → Completed**: Only when acceptance criteria are met
- **Any Status → Blocked**: When blockers are identified

### 6. Backup and Safety Requirements

#### Mandatory Backup Scenarios
- Before any database schema changes
- Before major refactoring of core components
- Before deployment to production
- Before any destructive operations (data deletion, migrations)
- Before configuration changes to critical systems

#### Backup Verification
- Verify backup integrity before proceeding
- Test restore procedures periodically
- Document backup locations and restoration steps
- Maintain backup retention policies

### 7. Documentation Requirements

#### Automatic Documentation Updates
- README.md must reflect new features or changes
- API documentation must be updated for interface changes
- Changelog must record all modifications
- Inline comments for complex algorithms or business logic

#### Documentation Validation
- All code examples in documentation must be tested
- Installation instructions must be verified
- Configuration examples must be current
- Troubleshooting guides must be relevant

### 8. Quality Gates

#### Must-Pass Criteria
- All automated tests must pass
- Code coverage thresholds must be met
- Security scans must be clean
- Performance benchmarks must be met
- Documentation must be complete and accurate

#### Blocking Issues
- Any failing test blocks completion
- Security vulnerabilities block deployment
- Performance regressions block release
- Missing documentation blocks completion

### 9. Continuous Improvement

#### Task Retrospective
- Document lessons learned from each task
- Identify process improvements
- Update templates based on experience
- Share best practices with team

#### Metrics and Tracking
- Track task completion time vs estimates
- Monitor defect rates and test coverage
- Measure documentation quality
- Analyze backup and recovery success rates