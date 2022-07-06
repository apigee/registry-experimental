import unittest
import word_extracting
class TestWordExtracting(unittest.TestCase):
 
   extractor = word_extracting.ExtractWords()
   def test_upper(self):
       self.assertEqual('foo'.upper(), 'FOO')
 
   def test_isupper(self):
       self.assertTrue('FOO'.isupper())
       self.assertFalse('Foo'.isupper())
 
if __name__ == '__main__':
   unittest.main()