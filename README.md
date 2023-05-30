# binaryornot

Detects if an input from file or stdin is text or binary. Wrapper around Go's [`http.DetectContentType`](https://pkg.go.dev/net/http#DetectContentType).

## Usage

```bash
binaryornot [-s] [file]
```

options:

```text
-s    use stdin (instead file)
```

Exit codes:

- `0`: text detected
- `5`: binary detected
- `1`: error

## Examples

```bash
$ binaryornot file1.txt
text
```

```bash
$ binaryornot -s < file2.dump
binary
```

```bash
$ echo hello | binaryornot -s
text
```
