# sireus
## Sireus - SRE Utility System - Dynamic Bot manager for executing SRE and DevOps commands conditionally

Replaces cron jobs, Jenkins, or other less sophisticated execution methods for SRE/DevOps automation.

### Sireus goals

- Bots execute commands or API calls, designed for SRE and DevOps environments
- Dynamically create Bots for any Platform, Service, Process, Host, etc from monitoring software (ex: Prometheus)
- Bots have rulesets for prioritizing conditional commands to respond to detected issues
- Scalable to large amounts of tests and commands, with deterministic execution, and inspectable with historical or test data
- Locking commands per Bot or Bot Group, to stop conflicting commands from running at once, or within a window to verify results of previous commands.
- Uses the "Utility AI" behavior system, which provides a method scoring for N conditions per command, to prioritize execution based on collected Bot information 

