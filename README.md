# sireus
## Sireus - SRE Utility System - Utility AI for executing SRE and DevOps commands & API calls

### Sireus goals

- Execute commands or API calls based on a scoring system, which is deterministic and inspectable
- Scalable to large amounts of tests and commands, giving priority to commands that have a higher score
- Locking commands per Agent or Agent class, to stop conflicting commands from running at once, or within a window to verify results of previous commands
- Automatic creation and detection of Agents from input data (ex: Prometheus or other monitoring software)
