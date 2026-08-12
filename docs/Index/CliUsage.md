# CLI Usage

## Description
Index of the documentation for people who drive Agnos-Cli from a terminal: installing the binary, recording a budget with it, and looking up what each command does. Consuming the same behavior from Go code is indexed by [LibUsage.md](/docs/Index/LibUsage.md).

The interface itself is `api.Lib.Sandboxmain` — one field of the library like any other — so everything documented here is produced by the same closed sandbox the library exposes.

---

## Tutorials

- [InstallCli.md](/docs/Tutorials/InstallCli.md)
  - **description:** Install the CLI globally, or build and run it from a checkout
  - [macOS / Linux](/docs/Tutorials/InstallCli.md#macos--linux)
  - [Windows (PowerShell)](/docs/Tutorials/InstallCli.md#windows-powershell)
  - [Verify after reboot](/docs/Tutorials/InstallCli.md#verify-after-reboot)
  - [Troubleshooting](/docs/Tutorials/InstallCli.md#troubleshooting)
  - [Install from a Clone](/docs/Tutorials/InstallCli.md#install-from-a-clone)
- [UseCli.md](/docs/Tutorials/UseCli.md)
  - **description:** Create categories, record transactions, and read balances from the terminal
- [RunCliSample.md](/docs/Tutorials/RunCliSample.md)
  - **description:** Run one of the shipped CLI examples from the source tree
  - [Run CLI Examples](/docs/Tutorials/RunCliSample.md#run-cli-examples)

---

## References

- [Commands.md](/docs/References/Commands.md)
  - **description:** Every command, flag, amount format, and exit code of the interface
  - [Commands](/docs/References/Commands.md#commands)
  - [Flags](/docs/References/Commands.md#flags)
  - [Amounts](/docs/References/Commands.md#amounts)
  - [Exit Codes](/docs/References/Commands.md#exit-codes)
  - [Environment](/docs/References/Commands.md#environment)
- [SamplesList.md](/docs/References/SamplesList.md)
  - **description:** Every shell example shipped in `examples/cliExamples/`
  - [Examples](/docs/References/SamplesList.md#examples)
