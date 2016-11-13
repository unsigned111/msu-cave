from lightdaemon import CaveController
import unittest
import Queue


class TestCaveController(unittest.TestCase):
    def test_init(self):
        pass

    def test_init_bad_lighting_configuration_location(self):
        with self.assertRaises(IOError):
            x = CaveController(
                {
                    'controller_ip': '127.0.0.1',
                    'lighting_configuration': 'nofile.yaml',
                },
                Queue.Queue(),
            )

    def test_get_frames(self):
        pass

    def test_add_generator(self):
        pass

    def test_run(self):
        pass

    def test_start(self):
        pass

    def test_stop(self):
        pass

    def test_blackout(self):
        pass

    @staticmethod
    def __build_default_controller():
        return CaveController(
            {
                "listen_ip": "192.168.1.109",
                "listen_port": 57121,
                "controller_ip": "10.7.153.129",
                "lighting_configuration": "./config/rigs/pod-rig.yaml",
                "min_intensity": 30,
                "max_intensity": 255,
                "start_color": {
                    "red": 255,
                    "green": 0,
                    "blue": 0,
                },
                "end_color": {
                    "red": 0,
                    "green": 0,
                    "blue": 255,
                },
                "color_granularity": 5,
            },
            Queue.Queue(),
        )
