# test_math.py
from nose.tools import assert_equal
from parameterized import parameterized, parameterized_class

import unittest
import math


def test_pow(base, exponent, expected):
   assert_equal(math.pow(base, exponent), expected)

class TestMathUnitTest(unittest.TestCase):
   @parameterized.expand([
       ("negative", -1.5, -2.0),
       ("integer", 1, 1.0),
       ("large fraction", 1.6, 1),
   ])
   def test_floor(self, name, input, expected):
       self.assertEqual(math.floor(input), expected)
