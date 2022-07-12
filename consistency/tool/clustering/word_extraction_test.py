print("\033c")
import unittest
import json
from mock import patch
from word_extraction import ExtractWords
from metrics import vocabulary_pb2
from google.protobuf.json_format import Parse, ParseDict


class TestExtractWords(unittest.TestCase):

    def parse_json_to_vocab(self):

        vocabs = []
        with open('consistency/tool/clustering/word_extraction_test.json', 'r') as myfile:
            data=myfile.read()

        # parse file
        ids = []
        obj = json.loads(data)
        for key, value in obj.items():
            ids.append(key)
            vocab = vocabulary_pb2.Vocabulary()
            fake_vocab = ParseDict(value[0], vocab)
            vocabs.append(fake_vocab)
        return vocabs

    @patch.object(ExtractWords, 'extract_vocabs')
    def test_fake_vocabs(self, mock_extract_vocabs):

        longMessage = True 
        fake_vocabs = self.parse_json_to_vocab()
        mock_extract_vocabs.return_value = fake_vocabs

        extrct = ExtractWords(stub = "stub", linearize=True)
        actual = extrct.get_vocabs()

        actual.sort()
        expected = ["ab", "ab", "bc", "cd", "ab", "ab", "bc", "cd", "ab", "ab", "bc", "cd", "ab", "ab", "bc", "cd"]
        expected.sort()

        self.assertEqual(actual, expected, "failed test" )







if __name__ == '__main__':
    unittest.main() 