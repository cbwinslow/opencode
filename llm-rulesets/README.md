# LLM Rulesets Repository

A comprehensive collection of rules and guidelines designed to optimize Large Language Model (LLM) interactions for safety, efficiency, and robust development practices.

## 🎯 Purpose

This repository provides structured rulesets that:
- Prevent LLM hallucinations and task drift
- Enforce safety protocols and prevent dangerous actions
- Ensure efficient and focused LLM responses
- Mandate proper backup and documentation practices
- Require comprehensive testing for task completion

## 📁 Repository Structure

```
llm-rulesets/
├── vscode/                  # VSCode-specific rules and configurations
│   └── llm-rules.json      # Main ruleset configuration
├── rules/                   # Detailed rule specifications
│   ├── task-management-with-testing.md
│   └── backup-documentation-requirements.md
├── docs/                    # Documentation and guides
├── examples/                # Implementation examples
└── README.md               # This file
```

## 🔧 Key Features

### Safety & Control
- **Hallucination Prevention**: Fact verification and uncertainty acknowledgment
- **Task Focus**: Prevent scope creep and maintain objective alignment
- **Danger Prevention**: Block potentially catastrophic actions
- **Irreversible Action Protection**: Mandatory backups and confirmations

### Efficiency & Intelligence
- **Smart Responses**: Concise, actionable communication
- **Context Optimization**: Relevant information maintenance
- **Quality Assurance**: Self-review and validation protocols

### Development Workflow
- **Task Management**: Structured breakdown with testing requirements
- **Backup Protocols**: Mandatory safeguards for destructive operations
- **Documentation Updates**: Synchronized documentation with code changes
- **Testing Requirements**: Comprehensive testing as completion criteria

## 🚀 Quick Start

### For VSCode Users

1. Copy the VSCode ruleset to your workspace:
```bash
cp llm-rulesets/vscode/llm-rules.json .vscode/
```

2. Install the LLM Rules extension (hypothetical):
```bash
code --install-extension llm-rules.vscode-extension
```

3. Reload VSCode to activate the rules

### For General Use

1. Review the core rules in `rules/` directory
2. Adapt the rules to your specific use case
3. Implement the task management template for your projects
4. Follow backup and documentation protocols

## 📋 Core Rules Summary

### Safety Rules
- ✅ Verify facts before stating them
- ✅ Distinguish between facts and speculation
- ✅ Admit uncertainty when information is unknown
- ✅ Maintain focus on primary objectives
- ✅ Block destructive operations without backup
- ✅ Require explicit confirmation for sensitive actions

### Efficiency Rules
- ✅ Provide direct, concise answers
- ✅ Use code examples over lengthy explanations
- ✅ Prioritize actionable advice
- ✅ Maintain relevant context without noise

### Development Rules
- ✅ Create tasks.md with testing requirements
- ✅ Tasks incomplete until tests pass
- ✅ Mandatory backups before risky operations
- ✅ Update documentation with all changes
- ✅ Follow testing-first development approach

## 🧪 Testing Requirements

All tasks must include:
- **Unit Tests**: For individual functions and components
- **Integration Tests**: For multi-component interactions
- **Edge Case Testing**: Boundary conditions and error scenarios
- **Performance Tests**: When applicable to the task
- **Security Tests**: For authentication and data handling

Tasks are **not complete** until all tests pass and are validated.

## 🔄 Task Management Template

Use this structure for all complex tasks:

```markdown
## Task: [Task Name]
**Priority**: [High/Medium/Low]
**Status**: [Pending/In Progress/Completed]

### Testing Requirements
- [ ] Unit tests for new functions
- [ ] Integration tests for component interactions
- [ ] Edge case testing
- [ ] Performance benchmarks (if applicable)

### Acceptance Criteria
- [ ] All tests pass
- [ ] Code review completed
- [ ] Documentation updated
- [ ] Backup created (if required)
```

## 🛡️ Backup Protocols

Mandatory backups before:
- Database schema changes
- Major code refactoring
- Production deployments
- Configuration changes
- Any destructive operations

## 📚 Documentation Standards

Required documentation updates:
- README.md for new features
- API documentation for interface changes
- Changelog for all modifications
- Inline comments for complex logic

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Follow the task management template
4. Include comprehensive tests
5. Update relevant documentation
6. Submit a pull request

## 📄 License

This repository is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Related Resources

- [LLM Safety Guidelines](https://example.com/llm-safety)
- [Testing Best Practices](https://example.com/testing)
- [Documentation Standards](https://example.com/docs)

## 📞 Support

For questions or support:
- Create an issue in this repository
- Review the documentation in `docs/`
- Check the examples in `examples/`

---

**Note**: These rules are designed to be adapted to your specific context and requirements. Review and modify them as needed for your particular use case and environment.