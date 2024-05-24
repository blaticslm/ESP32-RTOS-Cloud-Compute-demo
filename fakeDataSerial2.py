import json
import re
import numpy as np
import itertools
import serial
import time
import random

def main():
    global serialPort2

    serialPort2 = serial.Serial(
        port = "COM4", baudrate = 115200, bytesize=8, timeout=0.001, stopbits=serial.STOPBITS_ONE # COM7
    )
    print ("Collecting Sim Data: ")

    start_time = time.time()
    end_time = start_time + 300

    number_pattern = r"\d+\.\d+|\d+"
    letter_pattern = r"[A-Z]+"
        
        # Wait until there is data waiting in the serial buffer
    while time.time() < end_time:
            

        # print(type(line))
        # print(len(line))
        X = random.randint(0,100)
        Y = random.randint(0,100)
        Z = random.randint(0,100)
        line = f"X:{X}, Y:{Y}"
        

        letters = re.findall(letter_pattern, line)
        numbers = re.findall(number_pattern, line)

        numbers = np.array(numbers, dtype=float)

        combined_dict = dict(itertools.zip_longest(letters, numbers))

        res = json.dumps(combined_dict) + "\n"
        print(res)
        # serialPort2.write(serialData)
        time.sleep(1)
        # new_data = line.encode()
        serialPort2.write(res.encode())



if __name__ == "__main__":
    main()



