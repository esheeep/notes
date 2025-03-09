# Escalating XSS via Cookie Manipulation for Cross-Domain Redirection

You found a cross-site scripting vulnerability in a subdomain of grammarly.com. 
XSS allows you to inject malicious scripts into a webpage viewed by others.

Many web applications use cookies to store information about the user’s session or navigation. 
In this case, the "relocation cookie" stores the URL to which the user should be redirected after logging in.

Cookies can have a path attribute that specifies which URL paths can access the cookie. 
A more granular (specific) path might mean a cookie is accessible only to certain sections of the website.

If you can set a cookie with a more granular path, it may take priority when the client-side code checks which cookie to use. 
This is important because web applications often prioritize cookies with more specific paths over those with broader paths.

By controlling the redirection URL through the relocation cookie, you could potentially redirect the user to a URL of your choice after they log in to the main application, 
enabling further attacks like session hijacking or credential theft.

Even though the XSS vulnerability is in a subdomain, by exploiting the relocation cookie, 
you could pivot the attack to affect the main application (e.g., grammarly.com), escalating the severity of the vulnerability.

## PoC
XSS in subdomain
```javascript
document.cookie = "relocation=https://malicious-site.com; path=/; domain=.grammarly.com;";
```

