import json
import math
import time
import matplotlib.pyplot as plt
from matplotlib import style
import numpy as np
from sklearn.cluster import KMeans
from os import listdir
from os.path import isfile, join
from numba import jit, float32, cuda

# Indicate that these scripts are just for practice
print("This is just a practice for plotting graphs")
print("This has nothing to do with the project")

x = [1, 2, 3, 4, 5, 6, 7, 8]
y = [99.303, 119.007, 127.119, 176.093, 236.828, 260.154, 328.448, 341.47]

plt.plot(x, y)
plt.show()
