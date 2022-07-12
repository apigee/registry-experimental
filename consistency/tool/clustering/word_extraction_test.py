print("\033c")
from glob import glob
import unittest
import json
from mock import patch
from word_extraction import ExtractWords
from metrics import vocabulary_pb2
from google.protobuf.json_format import Parse, ParseDict


class TestExtractWords(unittest.TestCase):
    
    global tests
    def parse_json_to_vocab(self):

        
        with open('consistency/tool/clustering/word_extraction_test.json', 'r') as myfile:
            data=myfile.read()

        # parse file
        ids = []
        tests = []
        obj = json.loads(data)
        for key, value in obj.items():
            ids.append(key)
            test = []
            for i in range (len(obj[key])):
                vocab = vocabulary_pb2.Vocabulary()
                fake_vocab = ParseDict(obj[key][i], vocab)
                test.append(fake_vocab)
            tests.append(test)
        return ids, tests

    @patch.object(ExtractWords, 'extract_vocabs')
    def test_fake_vocabs(self, mock_extract_vocabs):


        longMessage = True 
        ids, tests= self.parse_json_to_vocab()
        for i in range(len(tests)):

            mock_extract_vocabs.return_value = tests[i]

            extrct = ExtractWords(stub = "stub", linearize=True)
            actual = extrct.get_vocabs()

            actual.sort()
            expected = ["ab", "abc", "bc", "cd", "ab", "ab", "bc", "cd", "ab", "ab", "bc", "cd", "ab", "ab", "bc", "cd"]*2
            expected.sort()

            self.assertEqual(actual, expected, "failed test: " + str(ids[i]) )







if __name__ == '__main__':
    unittest.main() 