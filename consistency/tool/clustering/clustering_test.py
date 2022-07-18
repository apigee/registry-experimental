import unittest
import json
from mock import patch
from parameterized import parameterized
from clustering import ClusterWords
import os
class TestClusterWords(unittest.TestCase):
    ROOT_DIR = os.path.realpath(os.path.join(os.path.dirname(__file__), '..'))

    with open(os.path.join(ROOT_DIR, 'clustering', 'clustering_test.json'), 'r') as myfile:
            data=myfile.read()
    obj = json.loads(data)
    names = []
    expected_values = []

    for name in obj:
        names.append(name)
        expected_values.append(obj[name]["clusters"])

    @parameterized.expand(zip(names, expected_values))

    @patch.object(ClusterWords, 'get_words')
    def test_clustering(self, name, expected, mock_get_words):

        # PATCH
        # Construct mock_response
        ROOT_DIR = os.path.realpath(os.path.join(os.path.dirname(__file__), '..'))
        with open(os.path.join(ROOT_DIR, 'clustering', 'clustering_test.json'), 'r') as myfile:
                data=myfile.read()
        obj = json.loads(data)
    
        mock_get_words.return_value = obj[name]["words"]

        #CALL
        clustr = ClusterWords(stub = "stub")
        actual = clustr.create_word_groups()

        # ASSERT
        self.assertDictEqual(dict(sorted(actual.items())), dict(sorted(expected.items())))
 
 
if __name__ == '__main__':
    unittest.main()