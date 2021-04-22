# just a quick script to test parallel requests
# need to be parametrized (hardcoded for now)
import gevent
from gevent.monkey import patch_all

patch_all()

import requests

def transfer():
    resp = requests.post("http://localhost:8000/api/v1/transfer", json={"from": "06d8b60d-5af6-4d56-a0fc-e779fe88647d","to": "459f4752-5163-48b3-afff-24b9511eccc2","amount": "1"})
    print(resp.content)


gs = []

for i in range(1000):
    gs.append(gevent.spawn(transfer))


gevent.joinall(gs)
