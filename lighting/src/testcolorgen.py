import unittest
from daemon import ColorGenerator


class TestColorGenerator(unittest.TestCase):
    def testInit(self):
        config = {
            "start_color": {"red": 0, "green": 0, "blue": 0},
            "end_color": {"red": 255, "green": 255, "blue": 255},
            "min_intensity": 30,
            "max_intensity": 255,
            "color_granularity": 2,
        }
        colorgen = ColorGenerator(config)
        self.assertEqual(colorgen.start_color.red, 0)
        self.assertEqual(colorgen.start_color.green, 0)
        self.assertEqual(colorgen.start_color.blue, 0)
        self.assertEqual(colorgen.end_color.red, 255)
        self.assertEqual(colorgen.end_color.green, 255)
        self.assertEqual(colorgen.end_color.blue, 255)
        self.assertEqual(colorgen.start_color.intensity, 30)
        self.assertEqual(colorgen.end_color.intensity, 30)
        self.assertEqual(colorgen.granularity, 10)

    def testInterpolation(self):
        x = ColorGenerator.interpolate(0, 255, 0.5)
        self.assertEqual(x, 127)