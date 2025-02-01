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

## Before OSWE

### Programming concepts

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

### Web App concepts

| Concept                  | What You Should Know:                                                                                                                                                              |
| ------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Input Validation**     | • How do web apps ensure user-provided data is valid? <br> • Which types of data can be dangerous to a web app?                                                                    |
| **Database Interaction** | • What kinds of databases can be used by a web app? <br> • How do database management systems differ? <br> • How does a web app create, retrieve, update, or delete database data? |
| **Authentication**       | • How does a web app authenticate users? <br> • What are hashes? Why is data often stored as hashes? <br> • How does an app                                                        |

| Language         | Sample Project for Code Review                                                                                                                                                              |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **PHP**          | • Beginner: [simple-php-website](https://github.com/banago/simple-php-website) <br> • Advanced: [Fuel CMS](https://www.getfuelcms.com/)                                                     |
| **ASP.NET & C#** | • Beginner: [SimpleWebAppMVC](https://github.com/adamajammary/simple-web-app-mvc-dotnet) <br> • Moderate: [Reddnet](https://github.com/moritz-mm/Reddnet)                                   |
| **NodeJS**       | • Beginner: [Employee Database](https://github.com/ijason/NodeJS-Sample-App) <br> • Moderate: [JS RealWorld Example App](https://github.com/gothinkster/node-express-realworld-example-app) |
| **Java**         | • Beginner: [Java Web App – Step by Step](https://github.com/in28minutes/JavaWebApplicationStepByStep) <br> • [Advanced: GeoStore](https://github.com/geosolutions-it/geostore)             |

### Vulnerabilities

| Vulnerability                      | Vulnerability Write-up                                                                                                                                                                                                                                               |
| ---------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Cross-Site Scripting (XSS)**     | • [From Reflected XSS to Account Takeover](https://medium.com/a-bugz-life/from-reflected-xss-to-account-takeover-showing-xss-impact-9bc6dd35d4e6) <br> • [XSS to Account Takeover](https://noobe.io/articles/2019-10/xss-to-account-takeover)                        |
| **Mass Assignment**                | • [KBID 147 – Parameter Binding](https://github.com/blabla1337/skf-labs/blob/master/kbid-147-parameter-binding.md) <br> • [Mass Assignment Cheat Sheet](https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Mass_Assignment_Cheat_Sheet.md)            |
| **Blind SQL Injection**            | • [KBID 156 – SQLI (Blind)](https://github.com/blabla1337/skf-labs/blob/master/kbid-156-sqli-blind.md)                                                                                                                                                               |
| **PHP Type Juggling**              | • [PHP Magic Tricks: Type Juggling](https://owasp.org/www-pdf-archive/PHPMagicTricks-TypeJuggling.pdf) <br> • [PHP Type Juggling Vulnerabilities](https://medium.com/swlh/php-type-juggling-vulnerabilities-3e28c4ed5c09)                                            |
| **Insecure Deserialization**       | • [KBID 271 – Deserialization YAML](https://github.com/blabla1337/skf-labs/blob/master/kbid-xxx-deserialisation-yaml.md) <br> • [Breaking .NET Through Serialization](https://media.blackhat.com/bh-us-12/Briefings/Forshaw/BH_US_12_Forshaw_Are_You_My_Type_WP.pdf) |
| **Business Logic Vulnerabilities** | • [KBID – Auth Bypass 2](https://github.com/blabla1337/skf-labs/blob/master/kbid-XXX-Auth-bypass-2.md)                                                                                                                                                               |
| **File Upload Vulnerabilities**    | • [Zorz VulnHub Writeup](https://berzerk0.github.io/GitPage/CTF-Writeups/ZorZ-Vulnhub.html)                                                                                                                                                                          |

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
