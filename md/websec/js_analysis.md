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
Save js urls
```bash
node lazyAssFiles.js > output.txt
```
Download the js from the files
```bash
wget -i output.txt
```


