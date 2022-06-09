


### the followign code gets 300 manually tagged words from Google's API, uses dbscan clustering to generate labels,
# and reveals the words corresponding to each point as you hover. Words in the same cluster colored the same. 
# If needed, remove the noises ( points with -1) just to plot words with defined clusters. 
# Intended for visualization of the algorithm we used for clustering.  
import numpy as np
import pandas as pd
from sklearn.cluster import dbscan
from strsimpy.sorensen_dice import SorensenDice
from sklearn.metrics.cluster import homogeneity_score
from sklearn.metrics.cluster import completeness_score 
from sklearn.metrics.cluster import v_measure_score
from sklearn.cluster import dbscan

tagged_df = pd.read_csv (r'/home/gelaw/work-stuff/gocode/src/registry-experimental/consistency/rpc/google/cloud/apigeeregistry/v1/similarity/analysis/vocab1000.csv')
tagged_df = tagged_df.drop(tagged_df.index[300:])
word_labels = tagged_df.iloc[:, 0]
word_labels = word_labels.to_numpy()
tagged_words = tagged_df.iloc[:, 1]
tagged_words = tagged_words.to_numpy()
data = tagged_words
X = np.arange(len(data)).reshape(-1, 1)
dice = SorensenDice(2)
def compute_predicted_lables(data, algorithm, dbscan_eps, dbscan_min_samples):
    db = dbscan(data, metric=algorithm, eps=dbscan_eps, min_samples=dbscan_min_samples, algorithm='brute')
    return db
def extract_indices_dice(x, y):
    i, j = int(x[0]), int(y[0])     # extract indices
    return dice.distance(data[i], data[j])

db = compute_predicted_lables(data = X, algorithm = extract_indices_dice, dbscan_eps = .3, dbscan_min_samples = 2)

def calculate_dissimilarity_matrix(api_strings, pairwise_dissimilarity_measure):
    size = len(api_strings)
    inconsistency_matrix = np.zeros((size, size))
    for i in range(size):
        for j in range(size):
            if i < j:
                string1  = api_strings[i]
                string2 = api_strings[j]
                if len(string1) == 0:
                    return len(string2)
                if(len(string2) == 0):
                    return len(string1)

                dissimilarity = pairwise_dissimilarity_measure(string1, string2)
                inconsistency_matrix[i][j] = dissimilarity
    inconsistency_matrix = inconsistency_matrix + inconsistency_matrix.T - np.diag(np.diag(inconsistency_matrix))
    return inconsistency_matrix

inconsistency_matrix = calculate_dissimilarity_matrix(data, dice.distance)
from sklearn.manifold import MDS

embedding = MDS(n_components=2, dissimilarity  = "precomputed")
fitted_strings = embedding.fit_transform(inconsistency_matrix)
x = fitted_strings[:,0]
y = fitted_strings[:,1]


import matplotlib.pyplot as plt 
names = tagged_words
c=db[1]

norm = plt.Normalize(1,4)
cmap = plt.cm.RdYlGn

fig,ax = plt.subplots()
sc = plt.scatter(x,y,c=c)

annot = ax.annotate("", xy=(0,0), xytext=(20,20),textcoords="offset points",
                    bbox=dict(boxstyle="round", fc="w"),
                    arrowprops=dict(arrowstyle="->"))
annot.set_visible(False)

def update_annot(ind):

    pos = sc.get_offsets()[ind["ind"][0]]
    annot.xy = pos
    text = "{}, {}".format(" ".join(list(map(str,ind["ind"]))), 
                           " ".join([str(names[n]) for n in ind["ind"]]))
    annot.set_text(text)
    annot.get_bbox_patch().set_facecolor(cmap(norm(c[ind["ind"][0]])))
    annot.get_bbox_patch().set_alpha(0.4)


def hover(event):
    vis = annot.get_visible()
    if event.inaxes == ax:
        cont, ind = sc.contains(event)
        if cont:
            update_annot(ind)
            annot.set_visible(True)
            fig.canvas.draw_idle()
        else:
            if vis:
                annot.set_visible(False)
                fig.canvas.draw_idle()

fig.canvas.mpl_connect("motion_notify_event", hover)

plt.show()