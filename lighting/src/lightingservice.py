import argparse
import time

parser = argparse.ArgumentParser(description='Starts the lighting service')
parser.add_argument('start', help='Starts the lighting service')
parser.add_argument('stop', help='Ends the lighting service')

if __name__ == "__main__":

    while True:
        time.sleep(0.01)




# 10.7.153.129 Enttec
# 10.7.153.130 Laptop
# 10.7.153.63  Pi