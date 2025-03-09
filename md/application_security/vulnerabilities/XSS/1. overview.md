# XSS 

1. Where does the payload end up?
2. What validation happens between sending and response.

## XSS context
### inside HTML tags
Add html tags to trigger javascript 
1. create a list of payloads with a list of tag types https://portswigger.net/web-security/cross-site-scripting/cheat-sheet 
2. create a list of payloads with a list of attributes 
3. construct payload to deliver to victim

### inside HTML tag attributes
terminate the attribute and close tag

`href="javascript: alert(1)"`

### inside JavaScript
terminating the existing script or breaking out of the current string.

### DOM
`location.search`
`URLSearchParams`