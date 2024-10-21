# Recon

## ASN
Autonomous system number - 
code for a company that is grown large that they need a routing data on the internet.

Note: Don't automate ASN because tools only scope at the domain you give it or it 
over scopes and find different company.

### Tools
[Hurricane Electrice](https://bgp.he.net/)

Searching on will give you ASNs and IP ranges.
Copy the ip ranges.

Scanning the ip ranges
- get web servers
- full port scans for services, ssh, ftp, sftp, rdp

### Port Scan
```bash
echo AS46489 | asnmap -silent | naabu -silent
```

### Service Scan
```bash
echo AS46489 | asnmap -slient | naabu -silent -nmap-cli 'nmap -sV'
```

Note: use ASNmap to convert

[naabuu](https://github.com/projectdiscovery/naabu)

### Passive Port Scan (SMAP)
Get subdomain from subfinder
Then smap it - use shodan.io and see what ports are open

## Cloud recon
SSL certificate enumeration to find subdomains.

cert.sh
keeping track of certificates on the internet.

### Tool
[CloudRecon](https://github.com/g0ldencybersec/CloudRecon)

to scan any server a pull out ssl fields and parse and output to command line.
Run cloud recon on the cloud ranges, and it will go and visit each ip and grab its cert

```bash
CloudRecon scrape -i <file>
```
ip ranges for cloud vendor

Use grep search the file from CloudRecon
```bash
cat file.txt | grep "twitch.tv"
```

parse out subs:
```bash
grep -F 'DOMAIN.com' 12_11_2023_DB.txt | awk -F '[][]' '{print $2}' | sed 's#\n#g' | grep ".DOMAIN.com" | sort -fu | cut -d ';' -f1 | sort -u
```
parse out all domains
```bash
grep -F 'DOMAIN.com' 12_11_2023_DB.txt | awk -F '[][]' '{print $2}' | sed 's#\n#g' | sort -fu | cut -d ';' -f1 | sort -u
```

[prips](https://github.com/honzahommer/prips.sh)
parse ip ranges into individual ip addresses

[hakip2host]()
look at certificate data -> reverse dns lookup

```bash
prips 27.126.144.0/21 | hakip2host
```






