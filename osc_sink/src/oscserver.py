import argparse

from pythonosc import dispatcher
from pythonosc import osc_server


def make_print_handler(path):
    def print_handler(unused_addr, args, *vals):
        print(path + " [{0}] ~ {1}".format(args[0], vals))
    return print_handler


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--ip", default="127.0.0.1",
                        help="The ip to listen on")
    parser.add_argument("--port", type=int, default=57121,
                        help="The port to listen on")
    args = parser.parse_args()

    dispatcher = dispatcher.Dispatcher()
    dispatcher.map("/eeg", make_print_handler("eeg: "), "debug")
    dispatcher.map("/onoff", make_print_handler("onoff: "), "debug")

    server = osc_server.ThreadingOSCUDPServer((args.ip, args.port), dispatcher)
    print("Serving on {}".format(server.server_address))
    server.serve_forever()
