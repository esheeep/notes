# Notes

Attack vector
User get added to one organisation, then assign admin privilages, get added to another organisation.

caido
filtering out requests
```bash
req.host.ncont:"f-log" and req.host.ncont:"gnar"
```

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



