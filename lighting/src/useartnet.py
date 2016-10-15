import artnet.dmx as dmx
import artnet.rig
import artnet.fixtures
import yaml
import threading
from pythonosc import dispatcher, osc_server

import walker
import logging
logging.basicConfig(level=logging.DEBUG)
log = logging.getLogger(__name__)


def main():
    config_file = './config/config.yaml'
    control = CaveController(config_file)
    output_thread = threading.Thread(target=control.run)
    output_thread.start()

    input = CaveListener(config_file)
    input_thread = threading.Thread(target=input.run)
    input_thread.start()

    while True:
        pass


class CaveController:
    def __init__(self, config_file):
        with open(config_file, 'r') as file:
            try:
                config = yaml.load(file)
            except Exception as exc:
                print(exc)
        self.light_controller = dmx.Controller(config['controller_ip'], fps=15.0)
        self.rig = artnet.rig.load(config['lighting_configuration'])
        self.light = self.rig.groups['all']
        self.light.setColor('#000000')
        self.light.setIntensity(0)
        self.add_generator(self.get_frames)

    @staticmethod
    def get_frames(group, clock):
        light = group[0]
        clk = clock()
        red = walker.walk(127)
        green = walker.walk(127)
        blue = walker.walk(127)
        intensity = walker.walk(127)

        while clk['running']:
            next_colors = (red.next(), green.next(), blue.next())
            next_intensity = intensity.next()
            group.setColor(artnet.fixtures.rgb_to_hex(next_colors))
            group.setIntensity(next_intensity)
            yield light.getFrame()

    def add_generator(self, function):
        self.light_controller.add(function(
            self.light, self.light_controller.get_clock()))

    def run(self):
        self.start()
        while True:
            try:
                pass
            except Exception:
                self.stop()

    def start(self):
        self.light_controller.start()

    def stop(self):
        self.light_controller.stop()


class CaveListener:
    def __init__(self, config_file):
        with open(config_file, 'r') as file:
            try:
                config = yaml.load(file)
            except Exception as exc:
                print(exc)
        self.ip = config['listen_ip']
        self.port = config['listen_port']
        self.dispatcher = dispatcher.Dispatcher()
        self.dispatcher.map("/eeg", self.handler, "handler")

    def run(self):
        server = osc_server.ThreadingOSCUDPServer(
            (self.ip, self.port), self.dispatcher)
        server.serve_forever()

    def handler(self, args, *vals):
        print("{0} ~ {1}").format(args[0], vals)


if __name__ == "__main__":
    main()
