---
name: Propose a new check
about: Suggest a macOS security control fort should check
title: "Check: "
labels: new-check
---

**The control**
What macOS security setting should fort check? (e.g. "Firmware password enabled")

**Why it matters**
What risk does it reduce, and who cares about it (individuals, SOC 2 / ISO teams, consultants)?

**How to detect it**
The documented command or API that reveals the current state, if you know it. fort only uses stable, documented macOS interfaces, no private frameworks.

```
$ (command that reveals the setting)
(example output)
```

**Can it be auto-fixed?**
- [ ] Yes, and here's the command that would safely change it: ...
- [ ] No, it needs the user to do something manually
- [ ] Not sure

**Framework mapping (optional)**
Relevant SOC 2 / ISO 27001 / NIST CSF / CIS controls, if known.
