import argparse
import random
import time
import sys

from pythonosc import osc_message_builder
from pythonosc import udp_client


def send_message():
    msg = osc_message_builder.OscMessageBuilder(address="/eeg")
    val = random.randint(0, 1023)
    msg.add_arg(val)
    msg = msg.build()
    client.send(msg)
    print("Sending a message! " + str(val))


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--ip", default="127.0.0.1",
                        help="The ip of the OSC server")
    parser.add_argument("--port", type=int, default=57121,
                        help="The port the OSC server is listening on")
    parser.add_argument("--push", default=False,
                        help="Push to send")
    parser.add_argument("--count", default=10,
                        help="Number of messages to send")
    args = parser.parse_args()

    client = udp_client.UDPClient(args.ip, args.port)

    for x in range(args.count):
        if args.push:
            input("<Press a key to send a message>")
            send_message()
        else:
            send_message()
            time.sleep(1)
