import unittest
import json
from mock import patch
from comparison import Comparison
from parameterized import parameterized


class TestComparison(unittest.TestCase):
    @parameterized.expand(
        ["simple", "none-ids", "none-words", "unique-terms", "both-null"]
    )
    @patch.object(Comparison, "get_word_group_ids")
    def test_find_word_groups(self, name, mock_get_word_group_ids):

        # PATCH
        # Construct mock_response
        with open("comparison_test.json", "r") as myfile:
            data = myfile.read()
        obj = json.loads(data)
        test_suite = obj[name]
        ids = test_suite["ids"]
        new_words = test_suite["words"]
        expected = test_suite["closest_word_groups"]

        mock_get_word_group_ids.return_value = ids

        # CALL
        cmprsn = Comparison(stub="stub", new_words=new_words)
        actual = cmprsn.find_word_groups()

        # ASSERT
        self.assertEqual(actual, expected)


if __name__ == "__main__":
    unittest.main()
