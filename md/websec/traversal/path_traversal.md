# Path traversal

## Null-byte injection
Null-byte injection is a technique to bypass file inclusion when the application expects a file extension such as `.png`. Null byte characters (`$00`, `\x00`) are injected to terminate the string.