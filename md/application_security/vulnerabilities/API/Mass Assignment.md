# Mass Assignment

## What is Mass Assignment?

Mass Assignment occurs when an application automatically binds user-supplied input to internal object properties without proper filtering. This vulnerability allows attackers to modify object properties that were never intended to be modified from the client side, potentially escalating privileges, bypassing security controls, or manipulating application data.

## Key Questions to Answer

- Does the application properly restrict which fields can be updated through API requests?
- Can sensitive properties like role, permissions, or ownership be modified via API requests?
- Does the backend validate and sanitise all input properties before processing them?
- Are some fields protected in the UI but accessible directly through API calls?

## Detailed Testing Methodology

### 1. Reconnaissance & Endpoint Mapping

#### 1.1 Identify All Endpoints That Accept Object Data

- Manually browse the application with proxy tools (Burp Suite, OWASP ZAP)
- Automate API endpoint discovery with tools like Kiterunner or FFUF
- Review API documentation (Swagger, OpenAPI) if available
- Analyse JavaScript files for object structures and API interactions
- Look for data submission in:
    - POST requests for object creation
    - PUT/PATCH requests for object updates
    - Custom actions that modify object properties

#### 1.2 Document Object Structures and Properties

- Create a comprehensive map of all object types (users, accounts, products, etc.)
- Document normal properties visible in the UI for each object
- Identify sensitive properties that should be protected:
    - `role`, `isAdmin`, `permissions`, `accessLevel`
    - `owner`, `createdBy`, `ownerId`
    - `price`, `cost`, `discount`, `isFree`
    - `verified`, `approved`, `status`
    - `accountBalance`, `credits`, `points`

#### 1.3 Analyse Request/Response Patterns

- Document the normal request structures for each endpoint
- Identify which properties are typically included in requests
- Note any differences between what's shown in UI and what's sent in requests
- Document server response patterns when invalid properties are submitted

### 2. Basic Mass Assignment Testing

#### 2.1 Property Enumeration

- Compare visible properties in UI with properties sent in requests
- Review API documentation for hidden properties
- Analyse JavaScript for additional object properties
- Review application responses for leaked property names
- Check for developer comments in HTML/JS that mention properties
- Look for property names in error messages

#### 2.2 Basic Property Injection

- Add new properties to normal requests:
    - Add `isAdmin: true` or `role: "admin"` to user update requests
    - Add `price: 0` to product purchase requests
    - Add `verified: true` to verification-related requests
- Test with various data types (boolean, string, number, object, array)
- Try both JSON and form data submissions when applicable

#### 2.3 Test Property Overwriting

- Identify write-once properties (creation timestamp, user ID, etc.)
- Attempt to overwrite these properties in update requests
- Test if properties can be reset to default values
- Check if null or empty values bypass validation

### 3. Advanced Mass Assignment Testing Techniques

#### 3.1 Property Discovery Through Error Messages

- Submit invalid property values to trigger validation errors
- Analyze error responses for property name leaks
- Test with extreme values (very long strings, negative numbers)
- Use special characters that might cause parsing errors

#### 3.2 Test Nested Properties and Objects

- Try to modify nested properties: `{"user": {"role": "admin"}}`
- Test array manipulation: `{"permissions": ["admin", "user"]}`
- Try adding unexpected nested objects
- Test with deeply nested properties: `{"profile": {"settings": {"privacy": {"isPrivate": false}}}}`

#### 3.3 Parameter Pollution Techniques

- Submit the same property multiple times with different values
- Mix different formats (JSON + form parameters)
- Try different case variations of property names
- Use array notation for single properties: `property[]=value1&property[]=value2`

#### 3.4 Framework-Specific Tests

- Test for Spring Framework binding issues (ModelAttribute)
- Check for Node.js/Express body-parser vulnerabilities
- Test for Ruby on Rails strong parameters bypass
- Look for PHP object injection via mass assignment
- Test for ASP.NET model binding weaknesses

### 4. Testing Different HTTP Methods

#### 4.1 Method Variation Testing

- Test POST endpoints with unexpected properties
- Try PATCH with partial updates containing sensitive fields
- Test if PUT handles unknown properties differently than PATCH
- Check if property filtering differs between creation and update operations

#### 4.2 Content-Type Manipulation

- Try different content types:
    - `application/json`
    - `application/x-www-form-urlencoded`
    - `multipart/form-data`
    - `application/xml` (if supported)
    - `text/plain`
    - `application/javascript`
- Test if property filtering changes based on content type
- Try manipulating the Content-Type header while keeping the same payload format
- Test with malformed Content-Type values (e.g., `application/json;charset=UTF-8;param=value`)
- Use incorrect Content-Type headers that don't match the actual payload
- Try Content-Type downgrade (e.g., submit JSON data with `application/x-www-form-urlencoded` header)
- Test with multiple conflicting Content-Type headers
- Check if the application handles charset parameters differently (e.g., `application/json;charset=UTF-16`)
- Test with rare or custom MIME types to bypass filters (e.g., `application/vnd.api+json`)
- For APIs that support multiple formats, test if validation differs between formats

### 5. Testing for Business Logic Impact

#### 5.1 Price and Payment Manipulation

- Attempt to modify price, cost, or discount values
- Try setting `total` or `finalAmount` properties directly
- Test for tax or shipping cost manipulation
- Check if payment status can be directly modified

#### 5.2 Status and Workflow Bypass

- Try to directly change object status values
- Test if approval workflows can be bypassed
- Check if verification steps can be skipped
- Test if date-based restrictions can be manipulated

#### 5.3 Ownership Manipulation

- Attempt to change resource ownership properties
- Test if creator IDs can be modified
- Check if team or group assignments can be manipulated
- Try to add yourself to restricted access lists

### 6. Response Analysis & Validation

#### 6.1 Confirm Property Acceptance

- Check if the modified properties appear in responses
- Verify if the changes persist when re-fetching objects
- Test if functionality changes reflect the property modifications
- Check database directly if possible to confirm changes

#### 6.2 Document Response Patterns

- Note which properties are reflected in responses
- Document error patterns for rejected properties
- Identify inconsistencies in property handling
- Look for information leakage in error responses

## Exploitation & Reporting

### Exploitation Templates

#### Basic Property Injection (JSON)

```
POST /api/users HTTP/1.1
Host: example.com
Authorization: Bearer [your_token]
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "role": "user",
  "isAdmin": true  // Injected sensitive property
}
```

#### Nested Object Manipulation (JSON)

```
PUT /api/users/profile HTTP/1.1
Host: example.com
Authorization: Bearer [your_token]
Content-Type: application/json

{
  "displayName": "John",
  "bio": "Regular user",
  "settings": {
    "theme": "dark",
    "permissions": {
      "canEditUsers": true  // Nested permission property
    }
  }
}
```

#### Form Parameter Pollution

```
POST /api/orders HTTP/1.1
Host: example.com
Authorization: Bearer [your_token]
Content-Type: application/x-www-form-urlencoded

productId=123&quantity=1&price=10.00&status=COMPLETED&paymentVerified=true
```

#### Content-Type Manipulation Attacks

```
# JSON payload with form-encoded Content-Type
POST /api/users HTTP/1.1
Host: example.com
Authorization: Bearer [your_token]
Content-Type: application/x-www-form-urlencoded

{"name":"John Doe","email":"john@example.com","isAdmin":true}
```

```
# Mixed Content-Type with boundary
POST /api/users HTTP/1.1
Host: example.com
Authorization: Bearer [your_token]
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW

------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="user"
Content-Type: application/json

{"name":"John","email":"john@example.com","role":"admin"}
------WebKitFormBoundary7MA4YWxkTrZu0gW--
```

```
# XML payload with JSON content type
POST /api/users HTTP/1.1
Host: example.com
Authorization: Bearer [your_token]
Content-Type: application/json

<user>
  <name>John Doe</name>
  <email>john@example.com</email>
  <role>admin</role>
  <isAdmin>true</isAdmin>
</user>
```

### Impact Assessment

For each Mass Assignment vulnerability, assess the impact:

|Impact Level|Description|Example|
|---|---|---|
|Critical|Privilege escalation or complete security bypass|Setting admin role, bypassing payment|
|High|Significant data manipulation or business rule bypass|Changing prices, modifying ownership|
|Medium|Manipulation of object state or workflow|Changing approval status, dates|
|Low|Minor data manipulation with limited impact|Modifying non-critical fields|

### Reporting Template

For each Mass Assignment vulnerability found:

1. **Vulnerability Title**: Clear description of the mass assignment vulnerability
2. **Affected Endpoint**: The specific API endpoint or function affected
3. **Vulnerability Description**: Technical explanation of the property manipulation
4. **Reproduction Steps**:
    - Detailed step-by-step instructions
    - Original vs. modified request examples
    - Properties that were successfully manipulated
5. **Impact**: Description of what an attacker could achieve
6. **Remediation**:
    - Implement allowlist approach for property binding
    - Use explicit property mapping instead of automatic binding
    - Create separate DTOs (Data Transfer Objects) for client-side data
    - Implement property-level access control

## Remediation Guidance

### Implementation Patterns

#### Allowlist Approach (JavaScript/Node.js)

```javascript
// Only allow specific properties to be updated
function updateUser(userId, userData) {
    // Define allowed fields
    const allowedFields = ['name', 'email', 'displayName', 'bio'];
    
    // Create a clean object with only allowed properties
    const sanitizedData = {};
    allowedFields.forEach(field => {
        if (userData[field] !== undefined) {
            sanitizedData[field] = userData[field];
        }
    });
    
    // Update user with sanitized data only
    return db.users.update(userId, sanitizedData);
}
```

#### Explicit DTO Mapping (Java/Spring)

```java
@PostMapping("/users")
public ResponseEntity<User> createUser(@RequestBody UserCreateDTO userDTO) {
    // UserCreateDTO only contains safe fields
    // Explicitly map from DTO to entity
    User user = new User();
    user.setName(userDTO.getName());
    user.setEmail(userDTO.getEmail());
    
    // Set defaults for sensitive properties
    user.setRole("user");
    user.setActive(false);
    
    return ResponseEntity.ok(userService.save(user));
}
```

#### Property-Level Authorization (Python)

```python
def update_user(user_id, data, current_user):
    # Field-level permission checks
    allowed_fields = ['name', 'email', 'bio']
    
    # Add additional fields based on user permissions
    if current_user.is_admin:
        allowed_fields.extend(['role', 'status', 'permissions'])
    
    # Filter data to only include allowed fields
    filtered_data = {k: v for k, v in data.items() if k in allowed_fields}
    
    # Update with filtered data
    return db.users.update(user_id, filtered_data)
```

### Framework-Specific Protections

#### Ruby on Rails

```ruby
# Using Strong Parameters
class UsersController < ApplicationController
  def update
    @user = User.find(params[:id])
    @user.update(user_params)
    
    redirect_to @user
  end
  
  private
  
  # Explicitly define allowed parameters
  def user_params
    params.require(:user).permit(:name, :email, :bio)
  end
end
```

#### ASP.NET

```csharp
// Using [Bind] attribute to explicitly include only safe properties
public ActionResult UpdateUser([Bind(Include = "Name,Email,DisplayName")] UserViewModel model)
{
    if (ModelState.IsValid)
    {
        var user = db.Users.Find(User.Identity.GetUserId());
        user.Name = model.Name;
        user.Email = model.Email;
        user.DisplayName = model.DisplayName;
        
        db.SaveChanges();
        return RedirectToAction("Index");
    }
    return View(model);
}
```

#### Content-Type Validation (Node.js/Express)

```javascript
// Middleware to validate Content-Type
function validateContentType(req, res, next) {
    const contentType = req.headers['content-type'] || '';
    
    // For endpoints expecting JSON
    if (req.path.startsWith('/api/') && req.method !== 'GET') {
        if (!contentType.includes('application/json')) {
            return res.status(415).json({ 
                error: 'Unsupported Media Type. API only accepts application/json'
            });
        }
        
        // Check if content actually is valid JSON
        try {
            if (req.body && typeof req.body === 'string') {
                JSON.parse(req.body);
            }
        } catch (e) {
            return res.status(400).json({ 
                error: 'Invalid JSON format'
            });
        }
    }
    
    next();
}

// Usage
app.use(validateContentType);
```

#### Content-Type Security Headers (Web Server Config)

```
# Nginx configuration to enforce Content-Type
add_header X-Content-Type-Options "nosniff" always;

# Additional security headers
add_header Content-Security-Policy "default-src 'self'" always;
```

## Tools for Mass Assignment Testing

- **Burp Suite Pro**: JSON Beautifier, Content Type Converter
- **OWASP ZAP**: Active Scanner with Parameter Pollution rules
- **Postman**: For systematic API testing
- **Arjun**: For parameter discovery
- **ParamMiner**: Burp extension for parameter discovery
- **JSON Wizard**: For manipulating JSON payloads
- **MassAssignment Scanner**: Specialized testing tools

## Mass Assignment Testing Cheatsheet

1. **Map all endpoints that accept object data (POST/PUT/PATCH)**
2. **Document normal object properties from UI and API documentation**
3. **Identify sensitive properties through code analysis and testing**
4. **Test by adding unexpected properties to requests**
5. **Try different property types and nested objects**
6. **Test with different content types and HTTP methods**
7. **Verify if property changes persist in the database**
8. **Test all privileged operations like status changes and approvals**
9. **Check if security controls can be bypassed via property injection**
10. **Document findings with clear impact and exploitation steps**