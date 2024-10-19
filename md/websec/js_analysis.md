# JavaScript Analysis

## Identifying JS files
When identifying JavaScript (JS) files, it's essential to have a solid understanding
of the application and its functionality before diving in.

**Script Tags**: Pay close attention to the <script> tags in the HTML as they often 
contain important references to JS file
**Hidden Functionalities**: Be mindful that JS files may reference hidden functionalities, 
such as API endpoints, that aren't immediately visible.

Example of stupid code
```javascript
if (window.location.href.includes('verified.capitalone'))
```
Whenever you see the includes method checking part of the URL, especially if you can control or influence the URL,
investigate it for potential vulnerabilities.

### App relevant
#### Lazy loaded js
Lazy loading is a method where JS files are loaded only when needed, instead of all at once. This technique is often used to optimize performance, 
but it can hide key JS files that might be critical to your analysis.

These files are loaded dynamically and often contain important functionality. They're typically identified by a unique entity name followed by a hash.
Example:
```javascript
<script src="runtime.d8ba6c6599cb3a9a.js" type="module"></script>
<script src="polyfills.244c7c2108cacf1c.js" type="module">
</script><script src="main.0b4c369979ae0e88.js" type="module"></script>
```
You might also see file paths like: 
```bash
/auth/assets/js/smartBanner.js
/auth/runtime.d8ba6c6599cb3a9a.js
```
These paths indicate multiple segments of JS files within the application.

Lazy-loaded files are often found in runtime.js and are commonly identifiable by the letter "u" in their structure.

```javascript
u = e => e + '.' + {
      76: 'd82157288b70d9fa',
      183: 'dd70866c3e5cd13c',
      285: 'e1f996b97a77964e',
      407: '88557614ddfa3f0d',
      413: 'a14edcc382816519',
      471: '84b1656beee57893',
      476: '3f2e9ad8a00364ea',
      477: '0ec2192254914875',
      585: 'c9ad8769bbe643ff',
      626: '498bb4356a7007cf',
      716: 'f0a0f2ac5ba8d374',
      804: 'd375052183a003a2',
      848: 'a76dd26a40d63e91',
      914: '26f9cea2ff180428'
    }
    [e] + '.js'
```
Lazy-loaded JS files may contain hidden functionalities that aren't immediately visible in the main JS files. Make sure to examine these to avoid missing key vulnerabilities.

Script to generate the js files
```javascript
prefix = "https://verified.capitalone.com/auth/";
functionName = "u";
data = `
u = e => e + '.' + {
      76: 'd82157288b70d9fa',
      183: 'dd70866c3e5cd13c',
      285: 'e1f996b97a77964e',
      407: '88557614ddfa3f0d',
      413: 'a14edcc382816519',
      471: '84b1656beee57893',
      476: '3f2e9ad8a00364ea',
      477: '0ec2192254914875',
      585: 'c9ad8769bbe643ff',
      626: '498bb4356a7007cf',
      716: 'f0a0f2ac5ba8d374',
      804: 'd375052183a003a2',
      848: 'a76dd26a40d63e91',
      914: '26f9cea2ff180428'
    }[e] + '.js'
`;

// Use a more specific regex to match ONLY the object keys (1-4 digit numbers before a colon).
numbers = [...data.matchAll(/(\d{1,4})(?=:)/g)].map(e => parseInt(e[1]));

// Remove duplicates
numbers = numbers.filter((item, index) => numbers.indexOf(item) === index);

// Define the function u
eval(data);
f = eval(functionName);
output = "";

// Generate URLs
for (var i = 0; i < numbers.length; i++) {
    const result = f(numbers[i]);
    if (result) { // Ensure only valid mappings are included
        output += prefix + result + "\n";
    }
}

// Print the output
console.log(output);

```
Replace the prefix and the data.

There can be a lot of hidden functionalities that you don't often have.
Particularly references to api endpoints. 

Note: this is more applicable to single page loaded apps. 

#### Vendor libraries
**Don't Overlook Vendor Libraries**: Even third-party vendor libraries, like New Relic, often have multiple JS files. 
It's important to examine them as they may contain exploitable vulnerabilities. Some vendors, like New Relic, even run their own bug bounty programs.

### Third party
**Can you pivot it to xss?**
When you encounter third-party JS, especially if it’s served within an iFrame, check if you can exploit Cross-Site Scripting (XSS) vulnerabilities. 
An example attack is using the trust relationship between an iFrame and the main page to inject an XSS on both.
Like “open-faced iFrame sandwich” attack.

#### Steal relevant info
When analyzing JS files, always consider whether they expose sensitive information. 
Tracking scripts can sometimes leak information, such as `window.location.href`, 
especially if the URL is included in logs or passed to third-party services. 
Investigate any includes statements in JS code to see if they can be controlled and lead to information leakage.

## Analysis
### Beautification
1. Download JS files
Run the following command to get the list of JS files:
```bash
node lazyAssFiles.js > output.txt
```
Download all the JS files:
```bash
wget -i output.txt
```
Don’t forget to include the main files.

2. Beautify the JS files
Use [pprettier](https://github.com/microsoft/parallel-prettier)

```bash
pprettier --write *.js*
```
After beautifying the files, open them in VS Code:
```bash
code .
```
### Identifying client-side path
1. Search for known urls
   Start by searching for paths in the JavaScript files that correspond to known URLs. For example, if you have a URL like:
```bash
https://verified.capitalone.com/auth/signin
```
You would search for the path definition in the JavaScript files, using the pattern:
```bash
path: "signin"
```
This helps you locate where the client-side route is defined.

2. Component Loading and Constructor
After finding the path definition, the next step is to identify the component associated with that path. For example:
```bash
{ path: "signin", component: SomeComponent }
```
The component (`SomeComponent` in this case) often has a class definition with a constructor function. 
The constructor is the first function that runs when the component is instantiated.

3. Setting Breakpoints in the Constructor
Set a breakpoint at the constructor of the component class and visit the URL (e.g., `/signin`). 
This allows you to observe the behavior of the component and see which other functions get 
triggered as part of the page load.
By tracing the constructor, you can:
- Understand what initial setup happens when the route is accessed.
- See if there are any functions related to cookies, query parameters, or other important actions triggered by visiting the route.
```javascript
{ path: "success", component: v.DN },
{ path: "not-now", component: v.DN },
{ path: "no-mobile", component: v.DN },
{ path: "snag", component: v.xT },
```


### Identifying server paths
Find apis and http verbs.
Note a server side paths in a request that being sent to the server side.
Then find that text in the actual code. 
Set the point and figure out what js is setting the server url.

Try to understand the structure of API endpoints.

#### HTTP verbs
Are they using `fetch` or `httpClient`
`post`, `get`, `put`

Also use linkfinder, jsluice to find endpoints that you might miss.

#### Monitoring
Find a pattern with the keywords and write a regex to search for that pattern.
Try to find new http endpoints.  

## Sources & sinks
### Sources
`URLSearchParams`
`location.*`
- `location.assign`
- `location.replace`
search on vscode `location\.`
`window.open`
`cookies`
`localstorage`/`sessionstorage` - can be source and sink

### Sinks
`windows.location.href` always check CSP - or some location relate sink
use js uri to direct to that 
`innerhtml`
`.html`
`unsafe templating`
`dangerouslysethtml`
`createElement (iframe, a, script, etc)`

Example
Set a query parameter in the url, query get embedded into the script source dynamically. 
Then break out of the query parameter value context, define another query parameter, 
that query parameter with a specific marketing url allow insert to js to a dynamically created js. 
    
Examples
Got XSS on a site use that to set cookie 
The cookie is set in a header - which can be used as part of the auth session
Login session is associated with a cookie, then prompt to do 2FA, if you can control the cookie, 
you can potentially hijack half auth session and sometimes there are endpoint you can hit with half off
like changing your email and do something funky. 

## JS Adjacents
Feature Flags
`function isFeatureFlagEnabled(){...}`
M&R rule: 
- `Response Body`
- `isFeatureFlagEnabled(){` 
- `isFeatureFlagEnabled(){return true;`

Create an alert for the Feature Flag that is being pushed.

## Dynamic wordlist generation
Parse doc, parse js files use words from those files.
## Monitoring
Use python script + regex with crontab
Request page every hour.


## Useful links
[Reversing and Tooling a Signed Request Hash in Obfuscated JavaScript](https://buer.haus/2024/01/16/reversing-and-tooling-a-signed-request-hash-in-obfuscated-javascript/)

[jsluice](https://github.com/BishopFox/jsluice)

Parse js files and pulling out interesting strings.
No need to use this tool, it's better to assess the files manually. 


You'll need to understand the codebase better than the people who wrote to find as many bugs as you can. 
Break the app apart. Understand all the little pieces, understand all the functions, all the client side paths, all the ways it interact with the apis. 





