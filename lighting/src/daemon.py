import time
import sys
import artnet.dmx as dmx
import artnet.rig
import artnet.fixtures
import yaml
import threading
import OSC
import Queue

import logging
logging.basicConfig(level=logging.DEBUG)
log = logging.getLogger(__name__)


def main():
    message_queue = Queue.Queue()

    config_file = './config/config.yaml'
    with open(config_file, 'r') as filehandler:
        try:
            config = yaml.load(filehandler)
        except Exception as exc:
            print("Error parsing config")
            sys.exit()

    control = CaveController(config, message_queue)
    output_thread = threading.Thread(target=control.run)
    output_thread.setDaemon(True)
    output_thread.start()

    input_listener = CaveListener(config, message_queue)
    input_thread = threading.Thread(target=input_listener.run)
    input_thread.setDaemon(True)
    input_thread.start()

    try:
        while True:
            time.sleep(150)
    except:
        print "Killed"



class CaveController(object):
    def __init__(self, config, message_queue):
        self.light_controller = dmx.Controller(config['controller_ip'], fps=30.0)
        self.rig = artnet.rig.load(config['lighting_configuration'])
        self.light = self.rig.groups['all']
        self.light.setColor('#FFFFFF')
        self.light.setIntensity(config['min_intensity'])
        self.add_generator(self.get_frames)
        self.message_queue = message_queue

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
        while True:
            time.sleep(1)

    def start(self):
        self.light_controller.start()

    def stop(self):
        self.light_controller.stop()


class CaveListener(object):
    def __init__(self, config, message_queue):
        self.ip = config['listen_ip']
        self.port = config['listen_port']
        self.server = OSC.ThreadingOSCServer(
            (self.ip, self.port))
        self.server.addMsgHandler("/eeg", self.handler)
        self.message_queue = message_queue
        self.color_gen = ColorGenerator(config)

    def run(self):
        self.server.serve_forever()

    def stop(self):
        self.server.close()

    def handler(self, *vals):
        value = vals[2][0]
        messages = self.color_gen.get_colors(value)
        for i in messages:
            print i
            self.message_queue.put(i)


class ColorMessage(object):
    def __init__(self, red, green, blue, intensity):
        self.red = red
        self.green = green
        self.blue = blue
        self.intensity = intensity

    def hex(self):
        return artnet.fixtures.rgb_to_hex((self.red, self.green, self.blue))

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
        self.last_color = ColorMessage(255, 255, 255, self.start_color.intensity)
        self.granularity = config["color_granularity"]

    def get_colors(self, target):
        print target
        print self.start_color
        print self.end_color
        print self.last_color

        target_color = ColorMessage(
            int(target * (self.end_color.red - self.start_color.red)),
            int(target * (self.end_color.green - self.start_color.green)),
            int(target * (self.end_color.blue - self.start_color.blue)),
            self.start_color.intensity,
        )

        red_increment = (target_color.red - self.start_color.red)
        green_increment = (target_color.green - self.start_color.green)
        blue_increment = (target_color.blue - self.start_color.blue)

        color_list = [ColorMessage(
            (int(i * red_increment + self.last_color.red)),
            (int(i * green_increment + self.last_color.green)),
            (int(i * blue_increment + self.last_color.blue)),
            self.start_color.intensity,
        ) for i in xrange(self.granularity)]

        self.last_color = target_color
        return color_list


if __name__ == "__main__":
    main()
