import unittest
import json
from parameterized import parameterized
from clustering import ClusterWords
from mock import patch
class TestClusterWords(unittest.TestCase):

    # @parameterized.expand(["simple-test", "int-test", "short-test", "None-test", "dot-test"])
    # def test_cleaning(self, name):

    #     # PATCH
    #     # Construct mock_response
    #     ROOT_DIR = os.path.realpath(os.path.join(os.path.dirname(__file__), '..'))
    #     with open('clustering_test.json', 'r') as myfile:
    #             data=myfile.read()
    #     obj = json.loads(data)
    #     mock_words  = obj['clean-words'][name]["words"]
    #     expected = obj['clean-words'][name]["cleaned-words"]
    #     #CALL
    #     clustr = ClusterWords(stub = "stub", words=mock_words)
    #     clustr.clean_words()

    #     # ASSERT
    #     self.assertEqual(clustr.words,  expected)



    # #Simple Assertions test
    # @parameterized.expand(["less-than-3-words",  "null-words"])
    # def test_cluster_assertions(self, name):

    #     # PATCH
    #     # Construct mock_response
    #     with open('clustering_test.json', 'r') as myfile:
    #             data=myfile.read()
    #     obj = json.loads(data)
    #     mock_words  = obj['cluster-words'][name]["words"]

    #     #CALL
    #     clustr = ClusterWords(stub= "stub", words=mock_words)

    #     #ASSERT
    #     with self.assertRaises(AssertionError):
    #         clustr.cluster()



    # #Simple clustering labels test
    # @parameterized.expand(["simple-test", "nonconvergence-test"])
    # def test_cluster_simple(self, name):

    #     # PATCH
    #     # Construct mock_response
    #     with open('clustering_test.json', 'r') as myfile:
    #             data=myfile.read()
    #     obj = json.loads(data)
    #     mock_words  = obj['cluster-words'][name]["words"]


    #     #CALL
    #     clustr = ClusterWords(stub= "stub", words=mock_words)
    #     labels = clustr.cluster()
    #     print(labels)
    #     expected_labels = obj['cluster-words'][name]["labels"]

    #     #ASSERT
    #     self.assertListEqual(list(labels), expected_labels)



        #ID forming labels test
    @parameterized.expand(["noise-words", "duplicate-words"])
    @patch.object(ClusterWords, 'cluster')
    def test_cluster_simple(self, name, mock_cluster):

        # PATCH
        # Construct mock_response
        with open('clustering_test.json', 'r') as myfile:
            data=myfile.read()
        obj = json.loads(data)

        mock_words  = obj['form-ids'][name]["words"]
        mock_cluster.return_value = obj['form-ids'][name]["labels"]

        #CALL
        clustr = ClusterWords(stub= "stub", words=mock_words)

        actual = clustr.create_word_groups()
        print(actual)
        print("*****************************************************")
        expected_clusters = obj['form-ids'][name]["clustered-words"]
        print(expected_clusters)

        #ASSERT
        self.assertDictEqual(actual, expected_clusters) 

if __name__ == '__main__':
    unittest.main() 