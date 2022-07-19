import unittest
import json
from parameterized import parameterized
from clustering import ClusterWords
import os
class TestClusterWords(unittest.TestCase):

    @parameterized.expand(["simple-test", "int-test", "short-test", "None-test", "dot-test"])
    def test_cleaning(self, name):

        # PATCH
        # Construct mock_response
        ROOT_DIR = os.path.realpath(os.path.join(os.path.dirname(__file__), '..'))
        with open('clustering_test.json', 'r') as myfile:
                data=myfile.read()
        obj = json.loads(data)
        mock_words  = obj['clean-words'][name]["words"]
        expected = obj['clean-words'][name]["cleaned-words"]
        #CALL
        clustr = ClusterWords(stub = "stub", words=mock_words)
        clustr.clean_words()

        # ASSERT
        self.assertEqual(clustr.words,  expected)

    #Simple Assertions test
    @parameterized.expand(["warning-test1", "warning-test2"])
    def test_cluster_assertions(self, name):

        # PATCH
        # Construct mock_response
        with open('clustering_test.json', 'r') as myfile:
                data=myfile.read()
        obj = json.loads(data)
        mock_words  = obj['cluster-words'][name]["words"]

        #CALL
        clustr = ClusterWords(stub= "stub", words=mock_words)

        #ASSERT
        with self.assertRaises(AssertionError):
            clustr.cluster()

    #Simple clustering labels test
    @parameterized.expand(["simple-test", "nonconvergence-test"])
    def test_cluster_simple(self, name):

        # PATCH
        # Construct mock_response
        with open('clustering_test.json', 'r') as myfile:
                data=myfile.read()
        obj = json.loads(data)
        mock_words  = obj['cluster-words'][name]["words"]


        #CALL
        clustr = ClusterWords(stub= "stub", words=mock_words)
        labels = clustr.cluster()
        expected_labels = obj['cluster-words'][name]["labels"]

        #ASSERT
        self.assertListEqual(list(labels), expected_labels)
if __name__ == '__main__':
    unittest.main() 