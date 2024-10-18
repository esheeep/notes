# JavaScript Analysis

## Identifying JS files
Only after you understand the application and know what you're dealing with. 

### Script tag
Stupid code
```javascript
if (window.location.href.includes('verified.capitalone'))
```
Each time you see includes and something you can control

Don't ignore the js in the script tag.

### Lazy loaded js
Format of a webpack, has name of a specific entity and a hash after it, which is very common. 
```javascript
<script src="runtime.d8ba6c6599cb3a9a.js" type="module"></script><script src="polyfills.244c7c2108cacf1c.js" type="module"></script><script src="main.0b4c369979ae0e88.js" type="module"></script>
```
Usually looks at main and app js but also look at runtime js as well. 

Look at the path. 
for example
`/auth/assets/js/smartBanner.js`
`/auth/runtime.d8ba6c6599cb3a9a.js`

Means there are multiple segments of js files in this application that I need to be aware of.

In the js you can find lazy loaded js files, the format is always similar to the below code, identified with a `u`. 
Most of the time the lazy loaded files are in runtime.js. 
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

### Vendor libraries
Don't ignore vendor libraries
New relic has alot of js files and they have their own bug bounty program

### Third party
Can you pivot it into xss? 
If you see an iframe from a third party, can you pop an XSS on that iframe and use the trust relationship between the iframe and the main page to pop an xss on the main page. 
Open faced iframe sandwich

### Tracking
Most you can get out of tracking files is `window.location.href` leak potentially.

## Analysis
1. Download all the js files
Save js urls
```bash
node lazyAssFiles.js > output.txt
```
Download the js from the files
```bash
wget -i output.txt
```
Remember to download the main files as well

2. Beautify the js
Use [pprettier](https://github.com/microsoft/parallel-prettier)

```bash
pprettier --write *.js*
```
Then `code .` to open VS code and do the analysis there.

### Identifying client-side path
Search for `path: "`
```javascript
{ path: "success", component: v.DN },
{ path: "not-now", component: v.DN },
{ path: "no-mobile", component: v.DN },
{ path: "snag", component: v.xT },
```
Search for known paths 
e.g `https://verified.capitalone.com/auth/signin`
search for `path: "signin"`
takes to this code
```javascript
{
    path: "signin",
        data: { title: "Sign In" },
    component: Es,
        resolve: {
        addMetaData: (Ee, Y) => {
            const K = (0, t.f3M)(Ln.VU);
            return (0, t.f3M)(g.Aq)
                .isPasskeySupported()
                .pipe(
                    (0, ss.b)((Ue) => {
                        K.pageViewMetaData = {
                            uiFeatures: [
                                {
                                    ui_features_element_names: `Fido Capability:${Ue}`,
                                    ui_features_product_id: `Fido Capability:${Ue}`,
                                },
                            ],
                        };
                    })
                );
        },
    },
},
```
How to find the structure for each individual application? Reverse engineer the client-side paths using the js files
You need to read the code until you find a unique string e.g. `other-products` to help find where the path is defined in the files.
From there you can get the idea how client-side path is structured.
This is mainly to single side page.

Client-side urls are routes that aren't reflective of serverside routes. 

Need to understand what code path you can trigger on the client side by redirecting the user to that url.

Path with a component that is a class,then theres a constructor. 
Put a break point at the constructor and visit the page. 
The constructor is the first thing going to call when the component class is instantiated. 
Then you can see what other functions get called. 

If you can trigger visiting a pacific endpoint and  making a cookie get set.
Look for if the structure somewhere else, understand the patterns that is used within the application.

keywords: `queryParams`, `headers`, `cookie`, `router`, `sessionStorage`

Look at the path you're already at.
Find it in the code. Where that path is being defined.
Notice that structure. Use that structure to identify all the client-side paths.

In the js looks for something big like classes and throw it in the search. 
Set a breakpoint and try to correlate between the codes. 

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
- `search on vscode `location\.``
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





