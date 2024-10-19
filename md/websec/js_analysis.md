# JavaScript Analysis

## Identifying JS files
When identifying JavaScript (JS) files, it's essential to have a solid understanding
of the application and its functionality before diving in.

**Script Tags**: Pay close attention to the `<script>` tags in the HTML as they often 
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

4. Analyzing Client-Side Paths
Client-side URLs (routes) are often not reflective of server-side routes. 
They are part of the front-end routing system (e.g., React Router, Angular Router)
and dictate what happens on the client without necessarily making a server request.
- **Redirecting and Triggering Code Paths**: Analyze how redirecting to certain client-side URLs triggers different code paths. 
- For example, visiting `/signin` might set certain cookies or initiate session storage depending on the client-side logic.
5. Understanding the Structure of Path Definitions
Look at the structure of the route where you're currently located. 
For example, if you're on `/signin`, find where that path is defined in the code. Take note of how that path is structured and used.
Use that pattern to identify other client-side paths. Many applications follow a consistent structure for defining their routes, components, and behavior.
6. Pattern Recognition
As you analyze the JavaScript code, you will notice common patterns. For example, paths with associated components may have:
- Class definitions
- Constructors that initialize data
- Methods that handle route parameters (queryParams, headers, router)

Search for keywords such as:
- `queryParams`
- `headers`
- `cookies`
- `router`
- `sessionStorage`

These will help you understand what kind of data is being processed when navigating between routes.
7. Look for Big Classes or Methods
In large JavaScript files, it can be helpful to search for big classes or methods related to routing and components. 
Once you find them, set breakpoints and try to correlate how different parts of the code work together.
For example, large classes might handle routing, UI updates, and state management, making them good targets for deeper analysis.

8. Triggering Cookie Setting
If you can trigger a specific endpoint (e.g., visiting `/signin`), and it sets a cookie, analyze how and where that cookie is being set.
Look for similar patterns in other parts of the application to see if there are other paths that trigger cookies, modify session data, or interact with browser storage.

### Identifying server paths
#### API endpoints
Look for API endpoints by tracking HTTP requests made to the server. Focus on identifying server-side paths (e.g., `/api/login`, `/auth/verify`) 
that are being hit when requests are sent from the client.
Once you find a server-side path, search for that path in the JavaScript code to locate where it is being set. 
This will help you understand how the client triggers requests to the server.
Identify how the API endpoints are structured to get a full picture of the app’s communication with the server.
#### HTTP verbs
Check which methods the application is using to make requests (e.g., `fetch`, `httpClient`).
Look for common HTTP methods like:
- `POST`: For creating or sending data.
- `GET`: For retrieving data.
- `PUT`: For updating data.
These methods help you understand what kind of data is being sent to and retrieved from the server.

**Tools for endpoint discovery**
- Use tools such as LinkFinder or JSluice to discover hidden or less obvious endpoints you might have missed in manual analysis.

**Monitoring for New Endpoints**
Develop a pattern based on keywords associated with API calls (e.g., fetch, url, axios), and use regular expressions (regex) to search for them in the codebase. 
This can help you spot new or undocumented API endpoints.

## Sources & sinks
### Sources
These are areas in the code where user-controlled data can enter, potentially leading to security risks if not properly validated or sanitized:
- `URLSearchParams`: User-controlled query parameters passed in the URL.
- `location.*`
  - `location.assign`
  - `location.replace`
  - User-controlled navigation through the location object, which can allow redirection or manipulation of paths.
  - Search for `location\.` in the code to find where the location object is used.
- `window.open`: User-controlled URL when opening a new window or tab via JavaScript.
- **Cookies**: Data stored in cookies can be manipulated by users if cookies are not secured properly.
- **LocalStorage/SessionStorage**: Users can interact with or manipulate stored data in the browser's localStorage or sessionStorage. 
These can act as both sources (retrieving data) and sinks (storing data).
  By identifying user-controlled data, you can trace how this data is being handled and whether there are any 
- potential risks (e.g., injection attacks, XSS) due to insufficient validation or sanitization.

### Sinks
These are places where untrusted data can be used in a potentially insecure way:
- `window.location.href`: Always check Content Security Policy (CSP) to ensure secure handling of redirects or URL manipulations.
- `innerHTML`: Directly inserting untrusted data into the DOM, which can lead to cross-site scripting (XSS) attacks.
- `.html`: Similar to innerHTML, risky if data is inserted dynamically.
- `Unsafe Templating`: Using JavaScript functions that allow for unsafe HTML or script insertion like `dangerouslySetInnerHTML`.
- `createElement`: Dynamically creating elements like `iframe`, `script`, or `a` can introduce security risks if not properly sanitized.

## JS Adjacents
### Feature Flags
Feature flags control the activation of certain features based on conditions. To find and analyze feature flags in JavaScript code, look for patterns like:
`function isFeatureFlagEnabled(){...}`: This function typically checks whether a feature flag is active or not.

### Monitoring and Response
- **Response Body**: Monitor the response body to see if any feature flag details are being sent or controlled by the server.
- **Feature Flag Patterns**: 
  - `isFeatureFlagEnabled(){`: This checks if a feature flag function is present and what conditions it evaluates.
  - `isFeatureFlagEnabled(){return true;`: This is an important pattern where the feature flag is hardcoded to true, indicating that the feature is permanently enabled in the environment you're analyzing.
#### Monitoring Feature Flags
Create an alert system to track when a specific feature flag is being activated or pushed. 
This will help you stay informed about changes or updates to feature flag statuses in the application

## Dynamic wordlist generation
Parse automation scripts to scan and extract relevant words, such as parameter names, function names, and keywords, from the codebase and documentation. 
This will help in identifying important terms for monitoring or further analysis.
## Monitoring
**Python Script + Regex**: Write a Python script that uses regular expressions (regex) to search for specific patterns in the codebase, such as feature flags, API calls, or user-controlled data.
**Crontab**: Set up the script with a crontab to run at regular intervals (e.g., every hour). This helps you monitor the website for any changes or new patterns, such as feature flags being toggled or new endpoints being introduced.

## Exploit examples
### Query Parameter Exploitation

Suppose a query parameter is embedded directly into the JavaScript source code dynamically. This can lead to exploitation as follows:
- Initial Exploit: You set a query parameter in the URL, such as ?param=value, and it gets embedded into the JavaScript code without proper sanitization.
- Breaking Out: By manipulating the query parameter value, you can "break out" of its intended context. For instance, by injecting special characters, you can alter the structure of the JavaScript code.
- Chaining with Additional Parameters: After breaking out, you can define additional query parameters, such as ?param=value&maliciousParam=<script>. These parameters can now be used to inject malicious JavaScript, such as inserting a marketing URL or dynamically creating and executing harmful scripts.

### Using XSS to Control Cookies and Hijack Sessions
- XSS Attack: Suppose you've identified a cross-site scripting (XSS) vulnerability on the site. You exploit this XSS to modify cookies directly.
- Cookie as Part of Authentication: Many sites associate login sessions with cookies. If the cookie is part of the authentication process and is also set via headers, you may be able to manipulate it.
- Bypassing 2FA: In some cases, after a login session is established via a cookie, the user is prompted for two-factor authentication (2FA). If you can control the cookie, you may bypass this 2FA step.
- Partial Session Hijacking: By manipulating only part of the authentication flow (e.g., hijacking the session via cookies), you could access specific endpoints without fully authenticating, allowing actions such as changing the user's email or performing unauthorized actions within the account.

## Useful links & extra notes
[Reversing and Tooling a Signed Request Hash in Obfuscated JavaScript](https://buer.haus/2024/01/16/reversing-and-tooling-a-signed-request-hash-in-obfuscated-javascript/)

[jsluice](https://github.com/BishopFox/jsluice): Parse js files and pulling out interesting strings.
No need to use this tool, it's better to assess the files manually. 

You'll need to understand the codebase better than the people who wrote to find as many bugs as you can. 
Break the app apart. Understand all the little pieces, understand all the functions, all the client side paths, all the ways it interact with the apis. 

To find as many vulnerabilities as possible, you need to break the application down and understand it better than the original developers. Focus on dissecting every part of the app, including:
- Understanding how each function works
- Mapping all client-side paths
- Identifying all interactions with APIs
By doing this, you can uncover hidden flaws and discover bugs that automated tools might miss.



