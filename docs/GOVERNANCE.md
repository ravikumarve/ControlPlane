# ControlPlane AI — Project Governance

This document describes the governance model for ControlPlane AI open-source and community projects. It defines roles, decision-making processes, and contribution pathways.

---

## Maintainer

**Ravi Kumar** is the sole maintainer and project owner. As the founder of ControlPlane AI, Ravi holds final authority over:

- Project direction and roadmap
- Release decisions and versioning
- Security vulnerability response
- Licensing and intellectual property
- Maintainer succession

Ravi may delegate maintainer responsibilities to trusted contributors over time as the project grows.

---

## Roles

### Maintainer

The maintainer is responsible for the overall health and direction of the project.

**Responsibilities:**
- Make final decisions on project direction and architecture
- Manage the release process and versioning
- Respond to security vulnerabilities
- Review and merge significant changes
- Moderate community discussions and enforce the Code of Conduct
- Appoint new maintainers (future)

### Contributor

A contributor is anyone who has had a pull request merged, submitted a meaningful issue report, or provided significant review feedback.

**Privileges:**
- Label and manage issues
- Review pull requests
- Participate in RFC discussions
- Vote on non-technical decisions (see Voting section)

**Becoming a contributor:**
- Have 3+ pull requests merged, or
- Demonstrate sustained, high-quality issue reporting and community participation

### Community Member

A community member is anyone who participates in discussions, files issues, or uses the project.

**Privileges:**
- Open issues and participate in discussions
- Submit pull requests
- Provide feedback on RFCs
- Vote on non-technical decisions (see Voting section)

**Responsibilities:**
- Follow the Code of Conduct
- Respect project decisions and governance

---

## Decision Process

Decisions are made at the appropriate level of formality depending on the impact of the change.

### RFC Process (Significant Changes)

Significant changes require a written RFC (Request for Comments) before implementation.

**What requires an RFC:**
- New features or capabilities
- Breaking API changes
- Policy model changes
- Architecture changes
- License or governance changes
- Addition or removal of core protocols

**RFC process:**
1. **Proposal** — Author writes a structured RFC document and opens a pull request in the project repository.
2. **Discussion** — The RFC is discussed for a minimum of 5 business days. All community members may comment.
3. **Resolution** — The maintainer accepts, rejects, or requests revisions. Accepted RFCs are merged and tracked as design documents.
4. **Implementation** — Implementation follows the approved RFC. Significant deviations require an amended RFC.

### Lazy Consensus (Routine Changes)

Routine changes use a lazy consensus model.

**What uses lazy consensus:**
- Bug fixes
- Documentation improvements
- Refactoring without behavior changes
- Test additions
- Dependency updates (non-security)

**Process:**
- Submit a pull request with the change
- Request review from relevant contributors
- If no objections are raised within 48 hours, the change is approved
- Objections require discussion and resolution before merge

### Breaking Changes

Any change that breaks backward compatibility requires:
1. A written RFC (see above)
2. A minimum 72-hour comment period from the time the RFC is announced
3. Explicit approval from the maintainer

Breaking changes are tracked in the changelog and released in major version bumps.

---

## Voting

Voting is used for non-technical decisions that affect the community.

**Applicable to:**
- Code of Conduct enforcement (see separate process)
- Community event planning
- Project branding and messaging
- Recognition of contributors

**Voting rules:**
- All community members may vote
- Simple majority (50% + 1) carries the decision
- Voting period is 7 days unless otherwise specified
- The maintainer may veto any vote with a written explanation

Technical decisions are not subject to voting. They are resolved through the RFC process described above.

---

## Release Process

ControlPlane AI projects follow semantic versioning (SemVer 2.0).

### Version Scheme

```
MAJOR.MINOR.PATCH
```

- **MAJOR** — Breaking changes
- **MINOR** — New features, backward-compatible
- **PATCH** — Bug fixes, backward-compatible

### Release Cadence

- **Major releases** — As needed, with at least 2 weeks notice
- **Minor releases** — Approximately every 4-6 weeks
- **Patch releases** — As needed, especially for security fixes

### Release Steps

1. **Changelog** — All changes are documented in CHANGELOG.md, organized by version.
2. **Release candidate** — A release candidate (RC) is tagged and announced at least 48 hours before the final release.
3. **Testing** — Release candidates undergo basic integration tests. Contributors are encouraged to test RC builds.
4. **Tag and release** — The final release is tagged, binaries are built, and the release is published on GitHub.
5. **Announcement** — The release is announced via GitHub and community channels.

### Security Releases

Security fixes follow an expedited process:
- Fixes are developed in a private fork
- A patch release is prepared and tested
- The fix is released simultaneously with public disclosure
- CVEs are filed as appropriate

---

## Code of Conduct

All participants in ControlPlane AI projects and communities are expected to follow the [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md). 

The maintainer is responsible for enforcing the Code of Conduct. Reports should be sent to the maintainer at the contact address specified in the Code of Conduct document.

---

## Contribution Guidelines

Contributors should read [CONTRIBUTING.md](./CONTRIBUTING.md) before submitting pull requests. Key expectations:

- All changes must be associated with an issue or RFC
- Pull requests should include tests and documentation
- Commit messages should follow conventional commits format
- Significant changes require an approved RFC first

---

## Dispute Resolution

In the event of a dispute that cannot be resolved through normal discussion:

1. **Escalation** — The matter is escalated to the maintainer.
2. **Mediation** — The maintainer may appoint a mediator from the contributor pool.
3. **Final decision** — The maintainer makes a final written decision.

All disputes are handled with respect for the community and the project's goals.

---

## Amendments

Changes to this governance document follow the RFC process for significant changes, with a minimum 10-day comment period. The maintainer must explicitly approve all amendments.

---

*This governance model is inspired by the [Rust project governance](https://github.com/rust-lang/rfcs) and [CNCF project governance](https://contribute.cncf.io/maintainers/governance/) patterns, adapted for a solo-founded open-core company.*
