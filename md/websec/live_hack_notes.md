# Notes

Attack vector
User get added to one organisation, then assign admin privilages, get added to another organisation.

# caido
filtering out requests
```bash
req.host.ncont:"f-log" and req.host.ncont:"gnar"
```
searching in the responses
```bash
resp.raw.cont:"redirect_location"
```
conditioning match and replace to validate theory
e.g. https://app.grammarly.com  replace javascript:alert(1)
condition: `req.path.cont:"/redirect/way`


interesting: domain for redirect

when you see base64 string starting with "ey"  thats mean you're dealing base65 encode json
normally its "eyj" if its using single it can be different which means its not a valid json format
- probably its python, someone taking python dict and running `str()`


Interesting to look at:
post request with a list, response with json with more information.
Look for toggle true, false.
the list in the request, search that in the js code.
try to turn on all the features - and see the changes and deep dive into the functionality

gadgets
see /redirect/ but no 302 means it can possibly mean client side redirection
possible to csp, xss

if someone have an xss in a subdomain of grammarly.com 
and we can set the relocation cookie
our cookie has more granular path and get priorties when the client side is reading 
redirect location that it should redirect to after login
we can convert xss from a random subdomain to the main app

fuzz with
all url encoded ascii
see what character works
e.g ; 

```javascript
invite_key=test;Path=/abc123" -> invite_key=test;Path=/abc123\"
invite_key=test;Path=/abc123\" -> invite_key=test;Path=/abc123\\"
```
Keep this in mind when you're trying to escape encoding
