### Future Plans

1. Implement a minimum and maximum version mechanism for project builds, this will prevent failures in the bootstrapping process.

2. ~~Implement logic for Agnos to generate libs that are installable within Agnos itself, thus allowing the creation of highly complex projects with Agnos.~~ Done: `agnos add-dep <module>@<version> --as <name>` copies another agnos repo's `sandbox/api/` into `sandbox/deps/<name>/` and generates the adapter that converts it. See [Adapters](../Adapters/doc.md).

3. Replace dependencies (deps) with an API at the function input, allowing the API to be visible to all functions.

4. Improve the code generation mechanisms to be more organized in large projects.
 