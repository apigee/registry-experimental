from imghdr import tests
import unittest
import json
from parameterized import parameterized
from clustering import ClusterWords
import os
class TestClusterWords(unittest.TestCase):

    @parameterized.expand(["simple-test", "int-test", "short-test", "None-test", "dot-test"])
    
    def test_clustering(self, name):

        # PATCH
        # Construct mock_response
        ROOT_DIR = os.path.realpath(os.path.join(os.path.dirname(__file__), '..'))
        with open(os.path.join(ROOT_DIR, 'clustering', 'clustering_test.json'), 'r') as myfile:
                data=myfile.read()
        obj = json.loads(data)
        mock_words  = obj['clean-words'][name]["words"]
        expected = obj['clean-words'][name]["cleaned-words"]
        #CALL
        clustr = ClusterWords(stub = "stub", words=mock_words)
        actual = clustr.clean_words()

        # ASSERT
        self.assertEqual(actual,  expected)
 
 
if __name__ == '__main__':
    unittest.main()