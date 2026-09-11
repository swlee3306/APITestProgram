# Authentication request samples

These JSON files are request-shape examples, not working credentials.
`REPLACE_WITH_YOUR_PASSWORD` is a literal placeholder; these files do not
automatically expand environment variables or fetch values from a secret manager.

Use a separate, untracked local copy and supply credentials only for a service you
are authorized to access. Do not commit that copy, print request bodies containing
passwords, or use these examples against a production endpoint as a smoke test.
Review the requested scope before use: `token.json` requests system scope and
`projectTkn.json` requests project scope. Neither example implies permission to
request privileged access.

Removing a value from the current files does not remove prior Git history or prove
that a previously committed credential has expired. No live credential validation
or account changes were performed as part of this sample cleanup.
