import json

f = open("results.txt", "w")
r = {"score": 0.31}
results = json.dumps(r)
f.write(results)

