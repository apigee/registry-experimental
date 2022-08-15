import unittest
import json
from mock import patch
from word_extraction import ExtractWords
from metrics import vocabulary_pb2
from google.protobuf.json_format import ParseDict
from parameterized import parameterized


class TestExtractWords(unittest.TestCase):
    @parameterized.expand(
        [
            (
                "simple-multiple-artifacts",
                [
                    "ab",
                    "ab",
                    "bc",
                    "cd",
                    "ab",
                    "ab",
                    "bc",
                    "cd",
                    "ab",
                    "ab",
                    "bc",
                    "cd",
                    "ab",
                    "ab",
                    "bc",
                    "cd",
                ]
                * 2,
            ),
            (
                "simple-one-artifact",
                [
                    "ab",
                    "ab",
                    "bc",
                    "cd",
                    "ab",
                    "ab",
                    "bc",
                    "cd",
                    "ab",
                    "ab",
                    "bc",
                    "cd",
                    "ab",
                    "ab",
                    "bc",
                    "cd",
                ],
            ),
            (
                "empty-schema",
                [
                    "ab",
                    "ab",
                    "bc",
                    "cd",
                    "ab",
                    "ab",
                    "bc",
                    "cd",
                    "ab",
                    "ab",
                    "bc",
                    "cd",
                ],
            ),
            ("empty-artifact", []),
            ("null-artifact", []),
            ("no-return-values", []),
        ]
    )
    @patch.object(ExtractWords, "extract_vocabs")
    def test_vocab(self, name, expected, mock_extract_vocabs):
        # PATCH
        # Construct mock_response
        mock_response = []
        with open("word_extraction_test.json", "r") as myfile:
            data = myfile.read()
        obj = json.loads(data)
        for json_vocab in obj[name]:
            vocab = vocabulary_pb2.Vocabulary()
            mock_response.append(ParseDict(json_vocab, vocab))
        mock_extract_vocabs.return_value = mock_response

        # CALL
        extrct = ExtractWords(stub="stub", project_name="-")
        actual = extrct.get_vocabs()

        # ASSERT
        self.assertListEqual(actual, expected)

    @parameterized.expand([("simple-None-test", None)])
    @patch.object(ExtractWords, "extract_vocabs")
    def test_vocab_none(self, name, expected, mock_extract_vocabs):

        # PATCH
        # Construct mock_response
        mock_response = None
        mock_extract_vocabs.return_value = mock_response

        # CALL
        extrct = ExtractWords(stub="stub", project_name="-")
        actual = extrct.get_vocabs()

        # ASSERT
        self.assertEqual(actual, expected)


if __name__ == "__main__":
    unittest.main()
