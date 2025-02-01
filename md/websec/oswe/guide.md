# OSWE Study

## Syllabus

- [https://manage.offsec.com/app/uploads/2023/01/WEB-300-Syllabus-Google-Docs.pdf](https://manage.offsec.com/app/uploads/2023/01/WEB-300-Syllabus-Google-Docs.pdf)

## Tools

- Burp Suite
- dnSpy:
  - [Codingo - Decompiling with dnSpy](https://codingo.io/reverse-engineering/ctf/2017/07/25/Decompiling-CSharp-By-Example-with-Cracknet.html)
  - [krypt0mux - Reverse Engineering .NET Applications](https://www.youtube.com/watch?v=_HvqI3Bsgfs)
- Reverse Shells
  - [Reverse Shell Cheat Sheet](https://highon.coffee/blog/reverse-shell-cheat-sheet/)
  - [Upload Insecure Files](https://github.com/swisskyrepo/PayloadsAllTheThings/tree/master/Upload%20Insecure%20Files)

## Programming concepts

| Concept                    | What You Should Know:                                                              |
| -------------------------- | ---------------------------------------------------------------------------------- |
| **Data Types**             | • How are they declared?                                                           |
|                            | • How can they be casted/converted to other data types?                            |
|                            | • Which data types have the ability to hold multiple sets of data?                 |
| **Variables & Constants**  | • Why do some data types need to be dynamic?                                       |
|                            | • Why do some data types need to remain constant?                                  |
| **Keywords**               | • Which words are reserved and why can they not be used as a variable or constant? |
| **Conditional Statements** | • How is data compared to create logic?                                            |
|                            | • Which operators are used to make these comparisons?                              |
|                            | • How does logic branch from an if/then/else statement?                            |
| **Loops**                  | • What are loops primarily used for?                                               |
|                            | • How is a loop exited?                                                            |
| **Functions**              | • How are functions called?                                                        |
|                            | • How are they called from a different file in the codebase?                       |
|                            | • How is data passed to a function?                                                |
|                            | • How is data returned from a function?                                            |
| **Comments**               | • Which characters denote the start of a comment?                                  |

## Topics

### ATutor Authentication Bypass and RCE

#### Blind SQL Injections

#### Data Exfiltration

#### Subverting the Atutor Authentication

#### Bypassing File Upload Restrictions

#### Remote Code Execution

### Atutor LMS Type Juggling Vulnerablity

#### PHP Loose and Strict Comparisons

#### PHP String Conversion to Numbers

### ManageEngine Applications Manager AMUserResourcesSyncServlet SQL Injection RCE

#### Houdine Escapes

#### Blind Bats

#### Accessing the File System

#### PostgreSQL Extensions

#### UDF Reverse shell

### Bassmaster NodeJS Arbitrary JavaScript Injection Vulnerability

#### The Bassmaster Plugin

### DotNetNuke Cookie Deserialization RCE

### ERPNext Authentication Bypass and Server Side Template Injection

### openCRX Authentication Bypass and Remote Code Execution

### openITCOCKPIT XSS and OS Command Injection - Blackbox

### Concord Authentication Bypass to RCE

### Server Side Request Forgery

### Guacamole Lite Prototype Pollution

## Resources

- [Exploit Writing for OSWE](https://github.com/rizemon/exploit-writing-for-oswe): This is an amazing resource that breaks down all of the important concepts for the python requests library.
- [Java Runtime Exec Command Generator](https://ares-x.com/tools/runtime-exec/): t can be painful to make your reverse shell payload work with Runtime exec, this website makes it a breeze.
