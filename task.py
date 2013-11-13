import sys
import json

args = sys.argv
output_file = args[-1]

f = open(output_file, "w")
r = {"score": 0.35}
results = json.dumps(r)
f.write(results)

