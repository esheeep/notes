# BFLA - Broken Function Level Authorization

## What is BFLA?

Broken Function Level Authorization (BFLA) occurs when an application fails to properly restrict access to functionality based on user roles and permissions. This vulnerability allows attackers to access features and API endpoints that should be restricted, enabling them to perform unauthorized operations. While BOLA focuses on unauthorized access to objects, BFLA focuses on unauthorized access to functions or features.

## Key Questions to Answer

- Are API endpoints properly enforcing role-based access controls?
- Can users access functionality intended for higher privilege roles?
- Is authorization enforced consistently across all endpoints and HTTP methods?
- Do multi-step processes validate authorization at each step?
- Can authorization checks be bypassed through alternative access paths?

## Detailed Testing Methodology

### 1. Reconnaissance & Role Mapping

#### 1.1 Identify User Roles and Permissions

- Document all user roles in the application (e.g., guest, user, moderator, admin)
- Create test accounts for each user role
- Map expected permissions for each role:
    - What functions should each role access?
    - What actions should each role perform?
    - What data should each role see?
- Review documentation for role-based permissions, if available
- Interview developers/stakeholders to understand intended authorization model

#### 1.2 Map API Endpoints and Functions

- Manually browse the application with proxy tools (Burp Suite, OWASP ZAP)
- Automate API endpoint discovery with tools like Kiterunner or FFUF
- Review API documentation (Swagger, OpenAPI) if available
- Analyze JavaScript files for hidden API endpoints and functions
- Look for function indicators in:
    - URL paths: `/api/admin/users/create`
    - Endpoint names: `/createUser`, `/deleteItem`
    - HTTP methods: GET (read), POST (create), PUT/PATCH (update), DELETE
    - Query parameters: `?action=approve`
    - Request bodies: `{"operation": "approve_user"}`

#### 1.3 Create Authorization Matrix

- Create a comprehensive table mapping:
    - API endpoints/functions
    - Required user roles for each function
    - HTTP methods supported for each endpoint
    - Expected authorization behavior
    - What happens on unauthorized access (block, error, etc.)
- Use tools like AuthMatrix (Burp extension) to organize your testing

Example matrix:

|Endpoint|Function|Required Role|GET|POST|PUT|DELETE|Expected Behavior|
|---|---|---|---|---|---|---|---|
|/api/users|List users|Admin|✓|✗|✗|✗|Normal users should not access|
|/api/items/{id}|View item|Any|✓|✗|✗|✗|All authenticated users can view|
|/api/items|Create item|User|✗|✓|✗|✗|Guests should not create items|
|/api/users/{id}/admin|Set admin|Admin|✗|✓|✗|✗|Only admins can promote users|

### 2. Testing for Vertical Privilege Escalation

#### 2.1 Baseline Testing

- Access each endpoint with appropriate role credentials
- Document normal responses, status codes, and behaviors
- Note any unique identifiers or tokens returned

#### 2.2 Unauthorized Access Testing

- Systematically test each endpoint with lower-privileged accounts
- Try to access admin-only features with regular user accounts
- Test application functionality that should be restricted:
    - User management
    - Configuration settings
    - Financial operations
    - Moderation actions
    - Bulk operations
    - Export/import functionality
    - Reporting features
    - System administration

#### 2.3 Test All HTTP Methods

- For each endpoint, test all HTTP methods (GET, POST, PUT, PATCH, DELETE, OPTIONS)
- Try alternative methods than those documented or observed
- Check if different HTTP methods bypass authorization checks:
    
    ```
    # Original admin requestPOST /api/users/create HTTP/1.1Host: example.comAuthorization: Bearer [admin_token]# Test with user token and alternative methodGET /api/users/create HTTP/1.1Host: example.comAuthorization: Bearer [user_token]
    ```
    

#### 2.4 Test for Role Parameter Manipulation

- Check if role information is passed in the request
- Attempt to manipulate role parameters:
    
    ```
    # OriginalPOST /api/action HTTP/1.1Content-Type: application/json{"role": "user", "action": "view_report"}# ModifiedPOST /api/action HTTP/1.1Content-Type: application/json{"role": "admin", "action": "view_report"}
    ```
    
- Look for hidden parameters that might bypass authorization:
    - `admin=true`
    - `role=administrator`
    - `access_level=9`
    - `is_admin=1`

### 3. Testing Multi-Step Operations

#### 3.1 Identify Process Flows

- Map out multi-step operations:
    - Registration and onboarding
    - Purchase flows
    - Approval processes
    - Configuration wizards
    - Import/export processes

#### 3.2 Test Authorization at Each Step

- Begin process flow with authorized role
- Complete initial steps of the process
- Switch to unauthorized role for later steps
- Check if authorization is verified at each step or only at process initiation
- Test for authorization "carry-over" between steps

#### 3.3 Test Process Manipulation

- Skip steps in multi-step processes
- Change the order of steps
- Return to previous steps after advancing
- Test if authorization state persists correctly throughout the process
- Check if directly accessing the final endpoint bypasses authorization

### 4. Advanced BFLA Testing Techniques

#### 4.1 Session and Token Analysis

- Analyze how authorization is implemented:
    - JWT tokens and claims
    - Session cookies
    - Custom authorization headers
- Test token manipulation:
    - Modify JWT payload claims to elevate privileges
    - Change role/permission attributes in the token
    - Test if expired tokens are properly rejected
    - Check if tokens are properly validated

#### 4.2 Test Hidden or Undocumented Functionality

- Look for "hidden" admin features:
    - Commented-out UI elements in HTML/JS
    - Debug or development endpoints
    - Legacy endpoints that may still be active
- Check for feature flags or parameters:
    
    ```
    # Try adding debug or admin parametersGET /api/users?debug=true HTTP/1.1GET /api/users?admin_view=1 HTTP/1.1GET /api/users?show_all=true HTTP/1.1
    ```
    

#### 4.3 Test for Insecure Direct Function Invocation

- Analyze if backend functions can be called directly
- Test for server-side parameter manipulation:
    
    ```
    # OriginalPOST /api/endpoint HTTP/1.1Content-Type: application/json{"action": "allowed_function"}# Modified to call restricted functionPOST /api/endpoint HTTP/1.1Content-Type: application/json{"action": "restricted_admin_function"}
    ```
    

#### 4.4 Test URL Path Traversal

- Check if paths can be manipulated to access admin functions:
    
    ```
    # Try traversing to admin areasGET /api/users/profile HTTP/1.1# ModifiedGET /api/admin/users/profile HTTP/1.1
    ```
    
- Test for alternative paths to the same functionality

#### 4.5 Test API Versioning Bypass

- Check if older API versions have weaker authorization
    
    ```
    # Current versionGET /api/v2/admin/users HTTP/1.1# Try older versionGET /api/v1/admin/users HTTP/1.1
    ```
    
- Test if authorization mechanisms differ between versions
- Look for deprecated but still functional endpoints

### 5. Testing Authorization Bypass Techniques

#### 5.1 Test HTTP Header Manipulation

- Add or modify headers that might affect authorization:
    
    ```
    X-Original-URL: /admin/usersX-Rewrite-URL: /admin/usersX-Override-URL: /admin/usersX-HTTP-Method-Override: PUTContent-Type: application/jsonX-Forwarded-For: 127.0.0.1X-Remote-Addr: 127.0.0.1X-Originating-IP: 127.0.0.1
    ```
    
- Test custom headers that might enable admin access:
    
    ```
    X-Admin: trueX-Role: adminX-Access-Level: 9X-Debug: trueX-Internal: trueX-Auth-Override: true
    ```
    

#### 5.2 Test for Client-Side Authorization

- Check if authorization is enforced only on the client side
- Look for hidden UI elements that could be re-enabled
- Test if API endpoints have proper server-side checks
- Compare UI-available functionality vs. direct API access

#### 5.3 Test for Race Conditions

- Send multiple simultaneous requests to check for race conditions
- Test if authorization can be bypassed during high load
- Check if temporary authorization states can be exploited

#### 5.4 Test for Forced Browsing

- Directly access restricted pages or endpoints
- Build URL/path lists based on patterns observed in allowed areas
- Use tools like DirBuster with authenticated sessions

### 6. Response Analysis & Validation

#### 6.1 Document Response Patterns

- Compare authorized vs. unauthorized responses:
    - HTTP status codes (200 vs. 401/403)
    - Response bodies
    - Response times
    - Error messages
- Look for information leakage in error responses
- Check for inconsistent error handling

#### 6.2 Confirm Vulnerabilities

- Verify that unauthorized access was actually achieved
- Document exact steps to reproduce
- Test multiple times to ensure consistency
- Identify the root cause of the authorization flaw

## Exploitation & Reporting

### Exploitation Templates

#### Direct Function Access

```
# Original (authorized for admin)
POST /api/admin/createUser HTTP/1.1
Host: example.com
Authorization: Bearer [admin_token]
Content-Type: application/json

{
  "username": "newuser",
  "email": "user@example.com",
  "role": "user"
}

# Unauthorized attempt (with user token)
POST /api/admin/createUser HTTP/1.1
Host: example.com
Authorization: Bearer [user_token]
Content-Type: application/json

{
  "username": "newuser",
  "email": "user@example.com",
  "role": "user"
}
```

#### Method Manipulation

```
# Standard method (blocked)
POST /api/admin/users HTTP/1.1
Host: example.com
Authorization: Bearer [user_token]

# Try alternative method (might bypass checks)
GET /api/admin/users HTTP/1.1
Host: example.com
Authorization: Bearer [user_token]
```

#### Process Manipulation

```
# Step 1 (authorized)
POST /api/order/create HTTP/1.1
Host: example.com
Authorization: Bearer [user_token]
Content-Type: application/json

{
  "product_id": "123",
  "quantity": 1
}

# Response includes order_id
{"order_id": "ORD12345", "status": "created"}

# Step 3 (directly accessing approval, skipping verification)
POST /api/order/approve HTTP/1.1
Host: example.com
Authorization: Bearer [user_token]
Content-Type: application/json

{
  "order_id": "ORD12345"
}
```

### Impact Assessment

For each BFLA vulnerability, assess the impact:

|Impact Level|Description|Example|
|---|---|---|
|Critical|Unauthorized access to administrative functions|Creating admin accounts, accessing all user data|
|High|Access to sensitive functionality|Approving transactions, modifying system settings|
|Medium|Access to elevated user functions|Moderator actions, special user privileges|
|Low|Minor function access with limited impact|Accessing non-sensitive reports or views|

### Reporting Template

For each BFLA vulnerability found:

1. **Vulnerability Title**: Clear description of the vulnerability
2. **Affected Endpoint/Function**: The specific API endpoint or function affected
3. **Vulnerability Description**: Technical explanation of the authorization bypass
4. **Reproduction Steps**:
    - Detailed step-by-step instructions
    - Sample requests and responses
    - Parameters or methods that were manipulated
5. **Impact**: Description of what unauthorized actions can be performed
6. **Remediation**:
    - Implement proper function-level authorization checks
    - Use role-based access control consistently
    - Centralize authorization logic
    - Apply the principle of least privilege

## Remediation Guidance

### Implementation Patterns

#### Centralized Authorization Middleware

```javascript
// Express.js example of authorization middleware
function requireRole(role) {
  return (req, res, next) => {
    // Get user from previous authentication middleware
    const user = req.user;
    
    if (!user) {
      return res.status(401).json({ error: "Authentication required" });
    }
    
    // Check if user has required role
    if (user.role !== role && user.role !== 'admin') {
      // Log potential BFLA attempt
      console.warn(`BFLA attempt: User ${user.id} (${user.role}) attempted to access ${req.method} ${req.path} which requires ${role} role`);
      return res.status(403).json({ error: "Insufficient permissions" });
    }
    
    // User is authorized, proceed
    next();
  };
}

// Usage in routes
app.get('/api/users', requireRole('admin'), (req, res) => {
  // Only admins reach this point
  res.json(getAllUsers());
});
```

#### Role-Based Decorators

```python
# Python/Flask example using decorators
def require_role(role):
    def decorator(f):
        @wraps(f)
        def decorated_function(*args, **kwargs):
            if not current_user.is_authenticated:
                return jsonify({"error": "Authentication required"}), 401
                
            if current_user.role != role and current_user.role != 'admin':
                # Log attempt
                app.logger.warning(f"BFLA attempt: User {current_user.id} ({current_user.role}) attempted to access function requiring {role} role")
                return jsonify({"error": "Insufficient permissions"}), 403
                
            return f(*args, **kwargs)
        return decorated_function
    return decorator

# Usage
@app.route('/api/admin/createUser', methods=['POST'])
@require_role('admin')
def create_user():
    # Function body here
    pass
```

#### Authorization Matrix Implementation

```java
// Java example using Spring Security
@Configuration
@EnableWebSecurity
public class SecurityConfig extends WebSecurityConfigurerAdapter {

    @Override
    protected void configure(HttpSecurity http) throws Exception {
        http
            .authorizeRequests()
                // Define URL patterns and required roles
                .antMatchers(GET, "/api/users").hasRole("ADMIN")
                .antMatchers(POST, "/api/users").hasRole("ADMIN")
                .antMatchers(GET, "/api/items/**").authenticated()
                .antMatchers(POST, "/api/items").hasAnyRole("USER", "ADMIN")
                .antMatchers(PUT, "/api/items/**").hasAnyRole("USER", "ADMIN")
                .antMatchers(DELETE, "/api/items/**").hasRole("ADMIN")
                .antMatchers(GET, "/api/public/**").permitAll()
                .anyRequest().authenticated()
            .and()
                .formLogin()
            .and()
                .csrf().disable();
    }
}
```

## Tools for BFLA Testing

- **Burp Suite Pro**: Autorize extension, AuthMatrix
- **OWASP ZAP**: Access Control Testing
- **Postman**: For systematic API testing with different tokens
- **JWT Tool**: For analyzing and manipulating JWT tokens
- **Autorize**: Burp extension for automated authorization testing
- **403 Bypasser**: Burp extension to try common auth bypass techniques
- **Authz**: CLI tool for testing authorization in web applications

## BFLA Testing Cheatsheet

1. **Map all user roles and permissions**
2. **Create test accounts for each role**
3. **Document all API endpoints and their intended access levels**
4. **Test each endpoint with each user role**
5. **Try alternative HTTP methods for restricted endpoints**
6. **Test multi-step processes for authorization at each step**
7. **Look for parameter and header manipulation opportunities**
8. **Check if older API versions have weaker authorization**
9. **Test direct function invocation through parameter manipulation**
10. **Document findings with clear reproduction steps and impact**