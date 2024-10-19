# Testing Email domain associated with group id

## Email domain associated with group id
When testing for vulnerabilities in systems where email domains are associated with group IDs or similar logic, 
your goal is to identify weaknesses in how the application verifies and handles email domains. Here’s a structured approach to testing and what to look for:
1. Understand the Email Validation Process
- How does the application handle domain verification? Investigate how the system checks if an email domain (e.g., @company.com) is valid for a particular group or organization.
- Where is the logic implemented? Is the domain checking done on the client-side (JavaScript) or server-side? Client-side logic is usually more vulnerable to manipulation.

2. Test Domain Validation Rules
- Check for Subdomain Misuse: See if the system treats subdomains as valid domains. For example, if company.com is a valid domain, test with attacker.company.com or fake.company.com to see if the system improperly assigns the same group or organization ID.
Example test email: `attacker@fake.company.com`.
- Email with Additional TLDs: Test variations of the email domain with different top-level domains (TLDs), like .net, .org, or .edu.
Example test email: `attacker@company.net` instead of `attacker@company.com`.
- Special Characters and Domain Manipulation: Test if the system properly handles characters like dots, hyphens, or mixed cases.
Example test email: `attacker@company-com.com` or `attacker@Company.com`.

3. Fuzz Test the Email Domain Input
- Try Edge Cases: Input various malformed or unusual domain formats to see if the system incorrectly handles them. Tools like Burp Suite can help automate fuzz testing with crafted payloads.
Examples:
```bash
   attacker@company.com.fake.com
   attacker@company.com?extra=bad
   attacker@company..com
   attacker@company.-com
```

4. Inspect the Group/Organization ID Assignment
- Check if you can control the group/organization ID: Once you've manipulated the email domain, check whether the group/organization ID assigned to you changes as expected. If the system associates you with the wrong group or gives access to an unrelated account, it indicates a vulnerability.
- Monitor responses: Use tools like Burp Suite or browser developer tools to inspect server responses and see how the group/organization ID is assigned. Pay attention to headers, cookies, or hidden fields in the login or registration process.

5. Tamper with Email Input
- Bypass Client-Side Validation: If domain validation is done on the client side (JavaScript or HTML forms), try bypassing it by modifying requests. You can intercept and modify requests using a proxy (e.g., Burp Suite or OWASP ZAP).
- Input Email Without the Domain: Some systems might only check the format of the email. You could try bypassing the check altogether by submitting malformed emails like attacker@ or attacker@localhost.
- Test for Parameter Manipulation: If the system sends the email domain to the server in a query parameter or body, try modifying the parameter directly. For example, if you see something like:

```bash
POST /login
email=attacker%40company.com
```

You can try to modify the email domain directly in the request to `attacker@fake.company.com` and see how the server responds.

6. Test for Role and Privilege Escalation
- Different Groups, Different Permissions: Some groups or organizations might have different permissions or roles. After manipulating the email domain, test if you are assigned the wrong role or permissions for another group.
- Check for Unauthorized Data Access: If the group/organization ID changes, check if you can view or modify data belonging to that group.

7. Cross-Test with Known Group Emails

- Register with Known Email Formats: If you know the valid domain for a certain group (e.g., company.com for an organization), test if using a variant of the email allows you to bypass restrictions.
- Spoof Other Groups: If you can successfully manipulate the domain to match another group, you might be able to register as a user for that group. For example:
   If the application is meant to handle emails like @university.edu, try registering as attacker@university.edu.com and see if it assigns you to the university's organization ID.

8. Understand Group-ID Logic
- How are IDs generated or assigned? Test to see if there’s any pattern to how group/organization IDs are assigned when an email domain is verified. This could involve inspecting API requests or the structure of the response.
- Does the system use regular expressions or exact string matching? If it's using regex for domain matching, try variations to see if the regex is too permissive.

