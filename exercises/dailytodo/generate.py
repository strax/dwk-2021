import requests

r = requests.head("https://en.wikipedia.org/wiki/Special:Random")
uri = r.headers['location']

print("URI: {}".format(uri))

text = "Read {}".format(uri)
r = requests.post("http://api.dwk-project.svc/todos", json={"text": text})
r.raise_for_status()
print(r.json())
