import time

import artnet.dmx as dmx
import artnet.rig
import artnet.fixtures

import walker
import logging
logging.basicConfig(level=logging.DEBUG)


log = logging.getLogger(__name__)
c = dmx.Controller('10.7.153.129')
rig = artnet.rig.load("./rigs/pod-rig.yaml")


def get_frames(group, clock):
    t = time.time()
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
        yield group[0].getFrame()


if __name__ == "__main__":
    group = rig.groups['all']
    group.setColor('#000000')
    group.setIntensity(0)

    c.add(get_frames(group, c.get_clock()))
    c.start()
    while True:
        try:
            pass
        except:
            c.stop()