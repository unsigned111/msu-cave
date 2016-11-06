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

#import logging
#logging.basicConfig(level=logging.DEBUG)
#log = logging.getLogger(__name__)


def main():
    message_queue = Queue.Queue()

    config = MasterControl.load_config()
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
            time.sleep(10000)
    except KeyboardInterrupt:
        control.stop()
        print("Server Killed")


class MasterControl(object):
    @staticmethod
    def load_config(path='./config/config.yaml'):
        with open(path, 'r') as filehandler:
            try:
                config = yaml.load(filehandler)
            except IOError as exc:
                print("Error parsing config")
                print exc
                sys.exit()
        return config


class CaveController(object):
    def __init__(self, config, message_queue):
        self.light_controller = dmx.Controller(config['controller_ip'], fps=45.0)
        self.rig = artnet.rig.load(config['lighting_configuration'])
        self.light = self.rig.groups['all']
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
                frame = light.getFrame()
                yield frame
            except Queue.Empty:
                pass

    def add_generator(self, function):
        self.light_controller.add(function(
            self.light, self.light_controller.get_clock()))

    def run(self):
        self.start()

    def start(self):
        self.light_controller.start()

    def stop(self):
        self.blackout()
        self.light_controller.stop()

    def blackout(self):
        self.message_queue.empty()
        self.message_queue.put(ColorMessage.blackout())
        self.light.setIntensity(0)


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
        for i in self.color_gen.get_headset_trans():
            self.controller.message_queue.put(i)

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
            messages = self.color_gen.get_headset_pulse()
            for i in messages:
                self.controller.message_queue.put(i)

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
            messages = self.color_gen.get_headset_trans()
            for i in messages:
                self.controller.message_queue.put(i)


class ColorMessage(object):
    def __init__(self, red, green, blue, intensity):
        self.__boundary_test(red, green, blue, intensity)
        self.red = red
        self.green = green
        self.blue = blue
        self.intensity = intensity

    @staticmethod
    def __boundary_test(*args):
        for i in args:
            if not isinstance(i, int):
                raise TypeError("Expected a type of int, but got {0}".format(type(i)))
            if i > 255 or i < 0:
                raise ValueError("Expected a value between 0 and 255, but got {0}".format(i))

    @staticmethod
    def from_np(np_array):
        if not isinstance(np_array, np.ndarray):
            raise TypeError('Expected an object of type numpy.ndarray, but got ' + type(np_array))
        if not (np_array.dtype == np.dtype('int32') or np_array.dtype == np.dtype('int64')):
            raise TypeError('Expected an numpy.ndarray object of dtype<int32/int64>, but got dtype<{0}>'.format(
                np_array.dtype))
        return ColorMessage(*np_array.flatten())

    def as_np(self):
        return np.array([self.red, self.green, self.blue, self.intensity])

    def hex(self):
        hex_string = artnet.fixtures.rgb_to_hex((self.red, self.green, self.blue))
        return hex_string

    def interp(self, next_color, count):
        red_interp = self._interp(self.red, next_color.red, count)
        green_interp = self._interp(self.green, next_color.green, count)
        blue_interp = self._interp(self.blue, next_color.blue, count)
        intensity_interp = self._interp(self.intensity, next_color.intensity, count)
        combined = np.vstack((red_interp, green_interp, blue_interp, intensity_interp))
        return [ColorMessage.from_np(i) for i in combined.transpose().round().astype(int)]

    @staticmethod
    def _interp(start, end, num):
        return np.interp(
            np.array(np.linspace(0, 1, num)),
            np.array([0, 1]),
            np.array([start, end]))

    @staticmethod
    def blackout():
        new_color = ColorMessage(0, 0, 0, 0)
        return new_color

    @staticmethod
    def whiteout():
        new_color = ColorMessage(255, 255, 255, 255)
        return new_color

    def __eq__(self, other):
        return self.as_np().all() == other.as_np().all()

    def __unicode__(self):
        string = "({0}, {1}, {2}) ~ {3}".format(self.red, self.green, self.blue, self.intensity)
        return string

    def __str__(self):
        return self.__unicode__()

    def __repr__(self):
        return self.__unicode__()


class ColorGenerator(object):
    def __init__(self, config):
        self.start_color = ColorMessage(
            config["start_color"]["red"],
            config["start_color"]["green"],
            config["start_color"]["blue"],
            config["min_intensity"])
        self.end_color = ColorMessage(
            config["end_color"]["red"],
            config["end_color"]["green"],
            config["end_color"]["blue"],
            config["min_intensity"])
        self.off_color = ColorMessage(
            config["headsetoff_color"]["red"],
            config["headsetoff_color"]["green"],
            config["headsetoff_color"]["blue"],
            config["headsetoff_color"]["intensity"])
        self.last_color = self.off_color
        self.granularity = config["color_granularity"]

    def get_colors(self, target):
        next_color = ColorMessage(
            int(target * (self.end_color.red - self.start_color.red) + self.start_color.red),
            int(target * (self.end_color.green - self.start_color.green) + self.start_color.green),
            int(target * (self.end_color.blue - self.start_color.blue) + self.start_color.blue),
            self.start_color.intensity,
        )
        to_return = self.last_color.interp(next_color, self.granularity)
        print "{0} => {1}".format(self.last_color, next_color)
        self.last_color = next_color
        return to_return

    def get_headset_trans(self):
        to_return = self.last_color.interp(self.off_color, self.granularity)
        self.last_color = self.off_color
        return to_return

    def get_headset_pulse(self):
        off_color_low = self.off_color
        off_color_high = ColorMessage(
            self.off_color.red, self.off_color.green, self.off_color.blue,
            self.off_color.intensity * 1.5)
        pulse_up = off_color_low.interp(off_color_high, self.granularity/2)
        pulse_down = off_color_high.interp(off_color_low, self.granularity/2)
        return pulse_up + pulse_down

if __name__ == "__main__":
    main()
