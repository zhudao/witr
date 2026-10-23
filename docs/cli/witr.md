## witr

Why is this running?

### Synopsis

witr explains why a process is running, found by name, PID, port, container or open file, by tracing what started it and what keeps it running.

```
witr [process name...] [flags]
```

### Examples

```

  # Inspect a running process by name
  witr redis-server

  # Look up a process by PID
  witr --pid 1234

  # Find the process listening on a specific port
  witr --port 5432

  # Find the process holding a file open
  witr --file /var/lib/dpkg/lock

  # Inspect a container by name
  witr --container redis

  # Inspect a process by name with exact matching (no fuzzy search)
  witr bun --exact

  # Show the full process ancestry (who started whom)
  witr mysqld --tree

  # Show only warnings (running as root, public bind, restarts, ...)
  witr dockerd --warnings

  # Display only environment variables of the process
  witr --pid 1234 --env

  # Short, single-line output (useful for scripts)
  witr --port 22 --short

  # Disable colorized output (CI or piping)
  witr --port 5432 --no-color

  # Output machine-readable JSON
  witr --port 3000 --json

  # Show extended process information (memory, I/O, file descriptors)
  witr mysqld --verbose

  # Combine flags: inspect port, show environment variables, output JSON
  witr --port 8080 --env --json

  # Multiple inputs
  witr mysqld redis-server
  witr --port 8080 --port 3000
  witr --pid 1234 --pid 5678

  # Mixed inputs
  witr redis-server --pid 1234 --port 8080

```

### Options

```
  -c, --container strings   container(s) to look up (repeatable)
      --env                 show environment variables for the process
  -x, --exact               use exact name matching (no substring search)
  -f, --file strings        file(s) held open by a process (repeatable)
  -h, --help                help for witr
  -i, --interactive         interactive mode (TUI)
      --json                show result as JSON
      --no-color            disable colorized output
  -p, --pid strings         pid(s) to look up (repeatable)
  -o, --port strings        port(s) to look up (repeatable)
  -s, --short               show only ancestry
  -t, --tree                show only ancestry as a tree
      --verbose             show extended process information
  -v, --version             version for witr
      --warnings            show only warnings
```


### Exit codes

| Code | Meaning |
|------|---------|
| 0 | Clean: process found, no warnings |
| 1 | Warnings: process found but has one or more warnings |
| 2 | Not found: no matching process or service |
| 3 | Permission denied: insufficient privileges (retry with sudo) |
| 4 | Invalid input: bad arguments or ambiguous match |
| 5 | Internal error: an unexpected failure occurred |
| 6 | Cause unknown: process found, but what started it can't be traced (not used on Windows) |

With several targets, the most severe result wins. Cause unknown ranks above warnings but below every failure (2 to 5).
