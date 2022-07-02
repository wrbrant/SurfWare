import random, logging
from more_itertools import divide
from mpl_toolkits.mplot3d import axes3d
import matplotlib.pyplot as plt
from matplotlib import style
import numpy as np

class tris3:

    #def setup():
    def __init__(self, fig, ax1, x, y, z):
        fig = plt.figure()                               # create a new figure for plotting
        self.ax1 = fig.add_subplot(111, projection='3d') # create a new subplot on our figure and set projection as 3d
        self.x = [0, 10, 5, 5]
        self.y = [0, 0, 10, 5]
        self.z = [0, 0, 0, 10]
        #return ax1, x, y, z
    def midpoint(self):
        a = random.randint(0,3) #select first starting point
        for i in range(5000):
            b = random.randint(0, 3) #select other starting point, which must be one of the orginal vertexes #TODO: make sure they are not the same point
            logging.debug(f"a is {a} and b is {b} ")
            # [-1] will grab the last element in the list
            self.x.append((self.x[-1] + self.x[b]) / 2)
            self.y.append((self.y[-1] + self.y[b]) / 2)
            self.z.append((self.z[-1] + self.z[b]) / 2)

    def display(self):
        self.ax1.scatter(self.x, self.y, self.z, c = 'm', marker = 'o')
        self.ax1.set_xlabel('x-axis')
        self.ax1.set_ylabel('y-axis')
        self.ax1.set_zlabel('z-axis')
        logging.info(self.x)
        logging.info(self.y)
        logging.info(self.z)
        plt.show()
if __name__ == "__main__":
    logging.basicConfig(filename='example.log', encoding='utf-8', level=logging.INFO)
    j = tris3(1, 2, 3, 4, 5)
    j.midpoint()
    j.display()