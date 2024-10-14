# DNS
DNS is like a phone book it associate name with a number, in this case a domain name and an ip address.

## Email
DNS is not meant to handle email verification but to do so there are 3 checks that the big email provider uses:

1. SPF (Sender policy framework): list of domains and ips allowed to send on your behalf.
2. DKIM (DomainKeys Identified Mail): Generate key or hash, mail server will look up signature dns.
3. DMARC (Domain-based Message Authentication): A policy that states what happens when dkim and spf pass or fail. Action include quarantine, reject or none.

Not all email provider will honor these checks but the big 3 do, Google, Microsoft and Yahoo.

**Domain alignment**
- from
- reply-to
- x-sender-from

ip address matches with the email domain.

**Some issues to note with email marketing**
- Email marketing goal is to sender bulk emails on behalf of customers
- Currently doing SPF, not every customer is using DMARC that means not everybody is receiving the emails.

Steps:
- Correct spf record
- Supply dkim hash
- Customer need to have a policy