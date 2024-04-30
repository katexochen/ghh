# ghh - GitHub Helper tool

## Authentication

Authentication against GitHub can be done by either setting `GHH_TOKEN` as environment
variable on command invocation or by using the `set-auth` subcommand, which will read
your token form the env or interactive input and write it to the GHH config file (unencrypted!).

## `create-project-issue`

Create project issue creates a new draft issue in a GitHub
project board. Metadata can be added in the form of fields and
users can be assigned.

The command takes two inputs: a path to an issue body (markdown or plain text)
and a path to a JSON document of metadata. The metadata file looks like this:

```json
{
    "organization": "<org>",
    "user": "<login>",
    "projectNumber": 42,
    "issueTitle": "<your issue title>",
    "assignees": [
        "<login>"
    ],
    "fields": {
        "Status": "<board column>",
        "field1": "value1",
        "field2": "value2"
    }
}
```

Notice that either `organization` or `user`, and `projectNumber` as well
as `issueTitle` are required fields.

The project number is part of the URL of your target project board.

![GitHub project board URL](assets/project-url.png)

For boards with columns, the column is also a field. You can find the
field name and the possible values in your project settings.

![GitHub project settings](assets/project-settings.png)


## `sync-forks`

Sync all forks of a user with their upstream repository. It will
fast-forward the fork if possible, otherwise it will merge the upstream branch.


**Set merge target branches** using the `--target-branches` flag. Per default, the target of the merge is the default branch of the fork.
You can pass multiple branch names, comma separated, or by using the flag multiple times. The existence of these
branch names will be checked in order, and the first matching branch will be used as
merge target. In case none of the branches matches, the default branch will be used.
If you don't want to fall back to the default branch after no target branch matched,
set the `--dont-target-default` flag.

Example:

```shell
# Only sync branches called 'upstream' or 'sync', no fallback to default branch
ghh sync-forks --target-branches upstream,sync --dont-target-default
```

**Ignore repos** you don't want to sync with the `--ignore-repos` flag.
You can pass multiple repository names, comma separated, or by using the flag multiple times.
Pass the repo name without your user prefix.

Example:

```shell
# This won't sync katexochen/ghh when executed by me.
ghh sync-forks --ignore-repos ghh
```
**Run as GitHub workflow** to keep all your fork automatically up to date.
You can easily copy [this example workflow](.github/workflows/sync.yml) and fit it to your needs.


## delete-all-runs

Delete all runs of a workflow and get your workflows sidebar clean again. To run this command,
within the repository, run the following command

```sh
ghh delete-all-runs
```

This will drop you into an interactive selection menu where you can select the workflow to delete.
Notice that every run must be deleted on its own, and the command can take quite a while to finish
(30-45 min if you have multiple thousand runs).
