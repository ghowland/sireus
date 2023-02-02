# sireus
## Sireus - SRE Utility System - Utility AI for prioritizing and conditionally executing SRE and DevOps commands & API calls

### Sireus goals

- Dynamically create Bots for any Service, Process, Host, etc from monitoring software (Prometheus, etc)
- Bots have rulesets for prioritizing conditional commands to respond to detected issues
- Scalable to large amounts of tests and commands, with deterministic execution, and inspectable with past or test situations
- Locking commands per Bot or Bot Group, to stop conflicting commands from running at once, or within a window to verify results of previous commands.

