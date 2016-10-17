import time
import sys
import artnet.dmx as dmx
import artnet.rig
import artnet.fixtures
import yaml
import threading
import OSC
import Queue
import numpy as np

import logging
#logging.basicConfig(level=logging.DEBUG)
#log = logging.getLogger(__name__)


def main():
    message_queue = Queue.Queue()

    config_file = './config/config.yaml'
    with open(config_file, 'r') as filehandler:
        try:
            config = yaml.load(filehandler)
        except IOError as exc:
            print("Error parsing config")
            print exc
            sys.exit()

    control = CaveController(config, message_queue)
    output_thread = threading.Thread(target=control.run)
    output_thread.setDaemon(True)
    output_thread.start()

    input_listener = CaveListener(config, control)
    input_thread = threading.Thread(target=input_listener.run)
    input_thread.setDaemon(True)
    input_thread.start()

    try:
        while True:
            time.sleep(150)
    except:
        print "Killed"
        control.stop()
        time.sleep(0.5)


class CaveController(object):
    def __init__(self, config, message_queue):
        self.light_controller = dmx.Controller(config['controller_ip'], fps=10.0)
        self.rig = artnet.rig.load(config['lighting_configuration'])
        self.light = self.rig.groups['all']
        self.light.setColor('#FF0000')
        self.light.setIntensity(config['min_intensity'])
        self.add_generator(self.get_frames)
        self.message_queue = message_queue
        self.headset_on = False

    def get_frames(self, group, clock):
        light = group[0]
        clk = clock()
        while clk['running']:
            try:
                message = self.message_queue.get()
                light.setColor(message.hex())
                light.setIntensity(message.intensity)
                yield light.getFrame()
            except Queue.Empty:
                pass

    def add_generator(self, function):
        self.light_controller.add(function(
            self.light, self.light_controller.get_clock()))

    def run(self):
        self.start()

    def start(self):
        self.light.setIntensity(0)
        self.light_controller.start()

    def stop(self):
        self.message_queue.empty()
        self.message_queue.put(ColorMessage.blackout())
        self.light.setIntensity(0)
        self.light_controller.stop()


class CaveListener(object):
    def __init__(self, config, controller):
        self.ip = config['listen_ip']
        self.port = config['listen_port']
        self.server = OSC.ThreadingOSCServer(
            (self.ip, self.port))
        self.server.addMsgHandler("/eeg", self.handler)
        self.server.addMsgHandler("/occupied", self.occupied_handler)
        self.controller = controller
        self.color_gen = ColorGenerator(config)

    def run(self):
        self.server.serve_forever()

    def stop(self):
        self.server.close()

    def handler(self, *vals):
        value = vals[2][0]
        print vals
        if self.controller.headset_on:
            messages = self.color_gen.get_colors(value)
            for i in messages:
                self.controller.message_queue.put(i)
        else:
            pass

    def occupied_handler(self, *vals):
        value = vals[2][0]
        print vals
        if value:
            self.controller.headset_on = True
            messages = self.color_gen.get_colors(value)
            for i in messages:
                self.controller.message_queue.put(i)
        else:
            self.controller.headset_on = False
            self.controller.message_queue.empty()
            self.controller.message_queue.put(ColorMessage.blackout())


class ColorMessage(object):
    def __init__(self, red, green, blue, intensity):
        self.red = red
        self.green = green
        self.blue = blue
        self.intensity = intensity

    def hex(self):
        return artnet.fixtures.rgb_to_hex((self.red, self.green, self.blue))

    @staticmethod
    def blackout():
        return ColorMessage(0, 0, 0, 0)

    @staticmethod
    def whiteout():
        return ColorMessage(255, 255, 255, 255)

    def __unicode__(self):
        return "({0}, {1}, {2}) ~ {3}".format(self.red, self.green, self.blue, self.intensity)

    def __str__(self):
        return self.__unicode__()


class ColorGenerator(object):
    def __init__(self, config):
        self.start_color = ColorMessage(
            config["start_color"]["red"],
            config["start_color"]["blue"],
            config["start_color"]["green"],
            config["min_intensity"])
        self.end_color = ColorMessage(
            config["end_color"]["red"],
            config["end_color"]["blue"],
            config["end_color"]["green"],
            config["min_intensity"])
        self.last_color = ColorMessage(0, 0, 0, self.start_color.intensity)
        self.granularity = config["color_granularity"]
        self.last_target = 0.0

    def get_colors(self, target):
        target_color = ColorMessage(
            int(target * (self.end_color.red - self.start_color.red) + self.start_color.red),
            int(target * (self.end_color.green - self.start_color.green) + self.start_color.green),
            int(target * (self.end_color.blue - self.start_color.blue) + self.start_color.blue),
            self.start_color.intensity,
        )

        increment = (target - self.last_target) / self.granularity

        to_return = list()
        for i in xrange(self.granularity):
            new_red = self.interpolate(
                self.start_color.red, self.end_color.red, i * increment + self.last_target)
            new_green = self.interpolate(
                self.start_color.green, self.end_color.green, i * increment + self.last_target)
            new_blue = self.interpolate(
                self.start_color.blue, self.end_color.blue, i * increment + self.last_target)
            interm_color = ColorMessage(
                new_red, new_green, new_blue, self.start_color.intensity,)
            to_return.append(interm_color)

        self.last_color = target_color
        self.last_target = target
        return to_return

    @staticmethod
    def interpolate(y0, y1, x, x0=0.0, x1=1.0):
        return int(y0 + (x * (x1 - x0)) * ((y1 - y0)/(x1 - x0)))


if __name__ == "__main__":
    main()
