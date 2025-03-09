# FFUF

Basic

```
ffuf -u http://ffuf.me/cd/basic/FUZZ -w common.txt
```

Recursion

```
ffuf -u http://ffuf.me/cd/basic/FUZZ -w common.txt -recursion
```

Looking for certain files

```
ffuf -w common.txt -e .log -u http://ffuf.me/cd/ext/logs/FUZZ
```

Using saved request file

```
ffuf -request req.txt -request-proto http -w passwords.txt -mc 200
```

Filter by size

```
ffuf -request req.txt -request-proto http -w passwords.txt -fs 30
```
