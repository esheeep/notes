# Examples of SSRF
Webhook: https://hackerone.com/reports/243277 
> 1. Create a webhook at https://app.mixmax.com/dashboard/settings/rules with url http://169.254.169.254/latest/meta-data/.
> 2. Trigger this webhook by sending/receiving an email. Wait a few hours.
> 3. Note that an email is not sent saying the webhook failed. I tried for other internal urls such as 'http://localhost', but they sent a failure email, indicating that http://169.254.169.254/latest/meta-data/ is open to the webhook.
> 4. n addition to verifying that this endpoint exists, an attacker could enumerate endpoints on this domain. For example, an attacker could enumerate MAC addresses at http://169.254.169.254/latest/meta-data/network/interfaces/macs/xx:xx:...