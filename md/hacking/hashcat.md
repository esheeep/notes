# hashcat

Cracking JWT signature secret

```bash
hashcat -a 0 -m 16500  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiJ1c2VyIiwiaWF0IjoxNzM5NTI3MzgyfQ.viHWm4mWio03aKiFGRDNZ_81HbrRBLmDVIE6JNBnteo /wordlist/rockyou.txt --show
```
