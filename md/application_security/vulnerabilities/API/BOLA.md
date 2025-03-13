# BOLA

## What is BOLA?

Broken Object Level Authorization (BOLA) occurs when an application fails to properly verify that a user has appropriate permissions to access, modify, or delete specific objects or resources. This vulnerability allows attackers to manipulate identifiers to access unauthorized data or perform unauthorized actions on resources belonging to other users.
## Key Questions to Answer

- Does the API properly verify ownership before allowing access to objects?
- Can a user access another user's data through API endpoints?
- Can a user modify, delete, or perform actions on another user's resources?
- Are authorization checks consistent across all endpoints and HTTP methods?

## Detailed Testing Methodology

### 1. Reconnaissance & Endpoint Mapping

#### 1.1 Identify All Endpoints with Object Identifiers

- Manually browse the application with proxy tools (Burp Suite, OWASP ZAP)
- Automate API endpoint discovery with tools like Kiterunner or FFUF
- Review API documentation (Swagger, OpenAPI) if available
- Analyze JavaScript files for hidden API endpoints and object structures
- Look for identifiers in:
    - URL paths: `/api/users/123/profile`
    - Query parameters: `?user_id=123`
    - Request bodies: `{"order_id": "abc-123"}`
    - Headers: `X-User-ID: 123`
    - Cookies: `user_session=user_123`

#### 1.2 Classify Object Types and Their Identifiers

- Create a comprehensive list of all object types (users, orders, files, comments, etc.)
- Document identifier patterns for each object type:
    - Sequential integers (1, 2, 3...)
    - UUIDs (8f9d6b69-9586-482f-b7e5-2ae8056571db)
    - Base64 encoded values
    - Custom formats (order-20250310-123456)

#### 1.3 Document Access Control Patterns

- Identify how the application manages sessions (JWT, cookies, custom tokens)
- Document all authentication headers or parameters
- Note which user roles exist and their expected permissions
- Map which endpoints should require ownership verification

### 2. Testing for Horizontal Privilege Escalation

#### 2.1 Create Multiple Test Accounts

- Create at least two accounts with the same privilege level
- Generate and document objects owned by each account (e.g., orders, posts, messages)
- Record all object identifiers for cross-account testing

#### 2.2 Test GET Requests (Read Access)

- Replace object IDs in requests with IDs belonging to another user
- Systematically test all endpoints that return user-specific data
- Try variations of the original ID:
    - For sequential IDs: increment/decrement by 1, 10, 100
    - For UUIDs: replace with UUIDs from other user accounts
    - For encoded IDs: decode, modify, re-encode
- Check for information leakage in both successful and error responses
- Observe HTTP status codes and response sizes for anomalies

#### 2.3 Test Modification Requests (PUT/PATCH)

- Attempt to update another user's objects
- Try partial updates to modify specific fields
- Test if validation differs for full vs. partial updates
- Check if responses reveal information about other users' objects
- Verify if successful modifications persist by re-fetching the object

#### 2.4 Test DELETE Requests

- Attempt to delete objects owned by other users
- Verify deletion by trying to access the object afterward
- Test if "soft delete" functionality behaves differently than permanent deletion
- Check for race conditions where objects can be accessed briefly after deletion

#### 2.5 Test POST Requests (Creation/Actions)

- Attempt to create objects in another user's context
- Try to perform actions on behalf of another user
- Test functions like "share," "transfer," or "assign" for BOLA
- Check if object ownership can be specified during creation

### 3. Advanced BOLA Testing Techniques

#### 3.1 Parameter Manipulation

- Modify multiple parameters simultaneously
- Test with empty values, null, or special characters
- Try removing identifiers completely
- Use arrays where single values are expected: `{"id": [123, 456]}`
- Mix identifiers from different objects

#### 3.2 HTTP Method Manipulation

- Change the HTTP method while keeping the same endpoint
    - GET → POST, PUT → GET, etc.
- Add or remove parameters when changing methods
- Test non-standard methods (HEAD, OPTIONS, PATCH)
- Check if the application enforces different authorization rules for different HTTP methods

#### 3.3 UUID/Non-Sequential ID Testing

- Analyze UUID generation patterns for predictability
- Check if UUIDs contain encoded information (timestamps, user info)
- Use tools to generate valid UUIDs in the same pattern
- Look for leaked UUIDs in:
    - HTML source code
    - JavaScript files
    - Public APIs or endpoints
    - Error messages
    - Application logs or debug information
    - Public repositories or documentation

#### 3.4 Mass Assignment Testing

- Test if you can modify restricted fields during update operations
- Try adding unexpected properties to request bodies
- Check if you can manipulate ownership fields directly
- Test if bulk operations bypass authorization checks

#### 3.5 Test for Indirect Object References

- Look for reference maps or indirect access mechanisms
- Test if you can manipulate these indirect references
- Check if the server properly validates references against user context

### 4. Testing for Vertical Privilege Escalation

#### 4.1 Test Admin/Privileged Endpoints

- Identify endpoints that should be restricted to administrators
- Attempt to access these endpoints with regular user credentials
- Test if adding admin-specific parameters to requests grants access
- Check if changing identifiers can access administrative functions

#### 4.2 Role-Based Access Testing

- Test endpoints with different user role accounts
- Attempt to access higher-privileged role functionality
- Check if authorization is based solely on client-side information
- Test if API endpoints enforce the same restrictions as the UI

### 5. Bypassing Protection Mechanisms

#### 5.1 Test for Inconsistent Authorization

- Check if authorization is implemented at the controller vs. model level
- Test if different API versions have different authorization rules
- Look for newer or legacy endpoints that might bypass checks

#### 5.2 Header and Parameter Manipulation

- Test if adding/modifying headers affects authorization:
    - Add `X-Original-User: [admin_user]`
    - Set `X-Forwarded-For` to internal IP addresses
    - Try `X-Authorization-Skip: true` or similar
- Check for debug parameters that might bypass authorization:
    - `debug=true`
    - `test_mode=1`
    - `internal=yes`

#### 5.3 Session Manipulation

- Test if you can manipulate session tokens to change identity
- Check if the application properly validates all claims in JWTs
- Try to reuse expired tokens or sessions

#### 5.4 Race Condition Testing

- Send multiple simultaneous requests to test race conditions
- Check if timing attacks can bypass authorization
- Test if temporary authorization states can be exploited

### 6. Response Analysis & Validation

#### 6.1 Document Response Patterns

- Compare successful vs. failed authorization responses
- Document error messages, status codes, and response formats
- Look for information leakage in error responses
- Check response times for timing side-channels

#### 6.2 Confirm Vulnerabilities

- Verify that unauthorized access was actually achieved
- Document the exact steps to reproduce
- Test multiple times to ensure consistency
- Identify the authorization flaw location

## Exploitation & Reporting

### Exploitation Templates

#### Simple ID Manipulation (GET)

```
GET /api/users/123/profile HTTP/1.1
Host: example.com
Authorization: Bearer [your_token]

# Change to:
GET /api/users/456/profile HTTP/1.1
Host: example.com
Authorization: Bearer [your_token]
```

#### Cross-User Action (POST)

```
POST /api/transfer HTTP/1.1
Host: example.com
Authorization: Bearer [your_token]
Content-Type: application/json

{
  "from_account_id": "123",
  "to_account_id": "789",
  "amount": 100
}

# Change to use another user's account:
{
  "from_account_id": "456", # Another user's account
  "to_account_id": "789",
  "amount": 100
}
```

#### UUID Manipulation (PUT)

```
PUT /api/documents/8f9d6b69-9586-482f-b7e5-2ae8056571db HTTP/1.1
Host: example.com
Authorization: Bearer [your_token]
Content-Type: application/json

{
  "title": "Updated Document",
  "content": "New content..."
}

# Change to another user's document UUID:
PUT /api/documents/7a8b9c0d-1e2f-3a4b-5c6d-7e8f9a0b1c2d HTTP/1.1
```

### Impact Assessment

For each BOLA vulnerability, assess the impact:

|Impact Level|Description|Example|
|---|---|---|
|Critical|Unauthorized access to highly sensitive data or functionality|Access to financial information, ability to transfer funds, PII exposure|
|High|Access to sensitive user data or significant functionality|Reading private messages, accessing medical records|
|Medium|Limited data exposure or functionality impact|Viewing other users' non-sensitive settings|
|Low|Minimal data exposure with limited business impact|Seeing another user's public profile data|

### Reporting Template

For each BOLA vulnerability found:

1. **Vulnerability Title**: Clear description of the vulnerability
2. **Affected Endpoint**: The specific API endpoint or function affected
3. **Vulnerability Description**: Technical explanation of the authorization bypass
4. **Reproduction Steps**:
    - Detailed step-by-step instructions
    - Sample requests and responses
    - Parameters that were manipulated
5. **Impact**: Description of what an attacker could access or modify
6. **Remediation**:
    - Implement proper object-level authorization checks
    - Verify user ownership of resources before any operation
    - Use indirect reference maps instead of direct object references
    - Implement a centralized authorization mechanism

## Remediation Guidance

### Implementation Patterns

#### Centralized Authorization

```python
# Before any object operation, verify ownership
def get_resource(resource_id):
    resource = db.find_resource(resource_id)
    if resource is None:
        return {"error": "Resource not found"}, 404
        
    # Critical authorization check
    if resource.owner_id != current_user.id and not current_user.is_admin:
        # Log potential BOLA attempt
        security_log.warning(f"BOLA attempt: User {current_user.id} attempted to access resource {resource_id} owned by {resource.owner_id}")
        return {"error": "Access denied"}, 403
        
    return resource, 200
```

#### Indirect Reference Maps

```python
# Use an indirection layer instead of exposing real IDs
user_resource_map = {
    "user123": {
        "doc_a": "8f9d6b69-9586-482f-b7e5-2ae8056571db",
        "doc_b": "7a8b9c0d-1e2f-3a4b-5c6d-7e8f9a0b1c2d"
    }
}

def get_document(user_id, doc_reference):
    # Get the actual ID from the user's reference map
    if user_id not in user_resource_map or doc_reference not in user_resource_map[user_id]:
        return {"error": "Document not found"}, 404
        
    real_doc_id = user_resource_map[user_id][doc_reference]
    return db.find_document(real_doc_id), 200
```

#### Database-Level Checks

```sql
-- Always include owner check in queries
SELECT * FROM orders 
WHERE order_id = ? AND user_id = ?;
```

## Tools for BOLA Testing

- **Burp Suite Pro**: Autorize extension, AuthMatrix
- **OWASP ZAP**: Access Control Testing
- **Postman**: For systematic API testing
- **Custom Scripts**: For automated ID fuzzing
- **IAmAuthority**: Specialized BOLA testing tool
- **Autorize**: Burp extension for automated authorization testing

## BOLA Testing Cheatsheet

1. **Map all endpoints and object identifiers**
2. **Create multiple test accounts at same privilege level**
3. **Systematically test all CRUD operations with other users' IDs**
4. **Check all ID locations: URL, body, headers, parameters**
5. **Test different HTTP methods against the same endpoint**
6. **Try parameter pollution and mixed mode attacks**
7. **Analyze UUID patterns and test with generated values**
8. **Check for authorization bypass in error conditions**
9. **Test bulk operations and mass assignments**
10. **Document everything systematically for proper reporting**