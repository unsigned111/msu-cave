import unittest
from lightdaemon import ColorGenerator, ColorMessage


class TestColorGenerator(unittest.TestCase):

    def setUp(self):
        self.config = {
            "start_color": {"red": 0, "green": 0, "blue": 0},
            "end_color": {"red": 255, "green": 255, "blue": 255},
            "min_intensity": 30,
            "max_intensity": 255,
            "color_granularity": 3,
        }
        self.colorgen = ColorGenerator(self.config)

    def test_init(self):
        self.assertEqual(self.colorgen.start_color.red, 0)
        self.assertEqual(self.colorgen.start_color.green, 0)
        self.assertEqual(self.colorgen.start_color.blue, 0)
        self.assertEqual(self.colorgen.start_color.intensity, 30)
        self.assertEqual(self.colorgen.end_color.red, 255)
        self.assertEqual(self.colorgen.end_color.green, 255)
        self.assertEqual(self.colorgen.end_color.blue, 255)
        self.assertEqual(self.colorgen.end_color.intensity, 30)
        self.assertEqual(self.colorgen.granularity, 3)

    def test_get_colors(self):
        x = self.colorgen.get_colors(1.0)
        self.assertEqual(len(x), 3)
        y = [
            ColorMessage(0, 0, 0, 30),
            ColorMessage(127, 127, 127, 30),
            ColorMessage(255, 255, 255, 30),
        ]
        self.assertEqual(x, y)
