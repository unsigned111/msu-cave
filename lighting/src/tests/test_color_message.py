import unittest
import numpy as np
from lightdaemon import ColorMessage


class TestColorMessage(unittest.TestCase):
    def color_message_constructor_valid(self):
        x = ColorMessage(255, 255, 255, 255)
        self.assertEqual(x.red, 255)
        self.assertEqual(x.green, 255)
        self.assertEqual(x.blue, 255)
        self.assertEqual(x.intensity, 255)

    def color_message_constructor_out_of_bounds(self):
        with self.assertRaises(ValueError):
            ColorMessage(256, 255, 255, 255)
        with self.assertRaises(ValueError):
            ColorMessage(255, 256, 255, 255)
        with self.assertRaises(ValueError):
            ColorMessage(255, 255, 256, 255)
        with self.assertRaises(ValueError):
            ColorMessage(255, 255, 255, 256)
        with self.assertRaises(ValueError):
            ColorMessage(-1, 255, 255, 255)
        with self.assertRaises(ValueError):
            ColorMessage(255, -1, 255, 255)
        with self.assertRaises(ValueError):
            ColorMessage(255, 255, -1, 255)
        with self.assertRaises(ValueError):
            ColorMessage(255, 255, 255, -1)

    def color_message_constructor_wrong_types(self):
        with self.assertRaises(TypeError):
            TypeError, ColorMessage(255.0, 255, 255, 255)
        with self.assertRaises(TypeError):
            ColorMessage(255, 255.0, 255, 255)
        with self.assertRaises(TypeError):
            ColorMessage(255, 255, 255.0, 255)
        with self.assertRaises(TypeError):
            ColorMessage(255, 255, 255, 255.0)

    def interpolation_contains_start(self):
        x = ColorMessage(0, 0, 0, 0)
        y = ColorMessage(255, 255, 255, 255)
        z = x.interp(y, 2)
        self.assertEqual(x, z[0])

    def interpolation_contains_end(self):
        x = ColorMessage(0, 0, 0, 0)
        y = ColorMessage(255, 255, 255, 255)
        z = x.interp(y, 2)
        self.assertEqual(y, z[1])

    def interpolation_middle(self):
        x = ColorMessage(0, 0, 0, 0)
        y = ColorMessage(255, 255, 255, 255)
        z = x.interp(y, 3)
        self.assertEqual(z[1].red, 128)
        self.assertEqual(z[1].green, 128)
        self.assertEqual(z[1].blue, 128)
        self.assertEqual(z[1].intensity, 128)

    def interpolation_length(self):
        x = ColorMessage(0, 0, 0, 0)
        y = ColorMessage(255, 255, 255, 255)
        z = x.interp(y, 15)
        self.assertEqual(len(z), 15)

    def interpolation_bug(self):
        x = ColorMessage(255, 0, 0, 30)
        y = ColorMessage(0, 0, 255, 30)
        z = x.interp(y, 15)
        self.assertEqual(z[-1], y)

    def from_numpy_success(self):
        x = ColorMessage.from_np(np.array([255, 255, 255, 255]))
        self.assertEqual(x.red, 255)
        self.assertEqual(x.green, 255)
        self.assertEqual(x.blue, 255)
        self.assertEqual(x.intensity, 255)

    def from_numpy_wrong_type(self):
        with self.assertRaises(TypeError):
            ColorMessage.from_np([255, 255, 255, 255])

    def from_numpy_wrong_dtype(self):
        with self.assertRaises(TypeError):
            ColorMessage.from_np(np.array([255, 255, 255, 255.0]))

    def as_numpy_success(self):
        x = ColorMessage(255, 255, 255, 255)
        y = x.as_np()
        self.assertIsInstance(y, np.ndarray)
        self.assertTrue(y.dtype == np.dtype('int32') or y.dtype == np.dtype('int64'))
        self.assertEqual(y.all(), np.array([255, 255, 255, 255]).all())

    def hex_white(self):
        x = ColorMessage(255, 255, 255, 255)
        self.assertEqual(x.hex(), "#ffffff")

    def hex_black(self):
        x = ColorMessage(0, 0, 0, 0)
        self.assertEqual(x.hex(), "#000000")

    def blackout(self):
        x = ColorMessage(0, 0, 0, 0)
        y = ColorMessage.blackout()
        self.assertEqual(x, y)

    def hex_whiteout(self):
        x = ColorMessage(255, 255, 255, 255)
        y = ColorMessage.whiteout()
        self.assertEqual(x, y)

    def equality_true(self):
        x = ColorMessage(255, 255, 255, 255)
        y = ColorMessage(255, 255, 255, 255)
        self.assertEqual(x, y)

    def equality_false(self):
        x = ColorMessage(255, 255, 255, 255)
        y = ColorMessage(0, 255, 255, 255)
        self.assertNotEqual(x, y)
        y = ColorMessage(255, 0, 255, 255)
        self.assertNotEqual(x, y)
        y = ColorMessage(255, 255, 0, 255)
        self.assertNotEqual(x, y)
        y = ColorMessage(255, 255, 255, 0)
        self.assertNotEqual(x, y)

    def repr(self):
        x = ColorMessage(255, 255, 255, 255)
        self.assertEqual(x.__repr__(), "(255, 255, 255) ~ 255")

    def string_cast(self):
        x = ColorMessage(255, 255, 255, 255)
        self.assertEqual(str(x), "(255, 255, 255) ~ 255")

    def unicode(self):
        x = ColorMessage(255, 255, 255, 255)
        self.assertEqual(unicode(x), "(255, 255, 255) ~ 255")