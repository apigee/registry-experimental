# zero

This directory contains an experimental command-line tool that explores the
[Google Service Infrastructure](https://cloud.google.com/service-infrastructure/docs/overview).

This tool is temporarily named `zero` to indicate that this is a minimalist
effort to explore the Service Infrastructure APIs out of the context of an API
gateway.

Areas to explore include:

- Automatically creating managed services and manage service configurations for
  APIs in an API registry.
- Automatically importing information into an API registry from the
  [Service Management API](https://cloud.google.com/service-infrastructure/docs/service-management/getting-started).
- Directly calling the
  [Service Control API](https://cloud.google.com/service-infrastructure/docs/service-control/getting-started)
  as an alternative to using an API proxy for the most basic API management
  needs.

Custom components used by the tool are in the `pkg` directory. Implementations
of CLI subcommands are in the `cmd` directory.

Code is published for transparency and community review with no promises of
support or continued existence.
