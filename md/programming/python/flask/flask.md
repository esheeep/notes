# Flask

## Setup
Create virtual environment
```bash
$ mkdir flaskproject
$ cd flaskproject
$ python3 -m venv .venv
```
Activate environment
```bash
$ source .venv/bin/activate   
```
Install flask
```bash
pip install flask
```
Page setup `main.py`
```python
from flask import Flask

app = Flask(__name__)

@app.route("/")
def hello_world():
    return "<p>Hello, World!</p>"
```
Run flask application
```bash
$ flask --app main run
```
