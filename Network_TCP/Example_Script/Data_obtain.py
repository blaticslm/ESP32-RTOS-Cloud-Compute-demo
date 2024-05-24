# This program is for continues getting the data

import requests
import json
import numpy as np


def isPrinting(dictionary):
    return dictionary["IsPrint"] != False

IP = "3.21.128.133"     # changes everytime I turn on the EC2 instance
PORT = 8080              # Specific IP port for accessing the server
API = "getGroupRange"    # API for accessing the database
MACHINE_ID = 30          # Specific Machine database
JOB_ID = 7               # Specific printing job according to MACHINE_ID
start = 0                # The query will first locating at the position at {start + 1}. Ex:range [1, limit] when start = 0
limit = 4000                # How many data will come each time starts from {Start + 1}
x_pos = np.array([])
y_pos = np.array([])

while(1):
    URL = f"http://{IP}:{PORT}/{API}/{MACHINE_ID}/{JOB_ID}/?start={start}&limit={limit}"
    print(URL)

    response = requests.get(URL)
    # print(response)

    if(response.json() is None):
        break

    list_of_dict = json.loads(json.dumps(response.json(), indent=1))
    x_input_pos = np.array([d['X_input_pos'] for d in list_of_dict])
    y_input_pos = np.array([d['Y_input_pos'] for d in list_of_dict])
    
    # x_pos = np.concatenate(x_pos, x_input_pos)
    print(type(x_input_pos))
    print(x_input_pos)
    print(y_input_pos)

     

    start += limit







