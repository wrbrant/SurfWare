import random
from more_itertools import divide
from mpl_toolkits.mplot3d import axes3d
import matplotlib.pyplot as plt
from matplotlib import style
import numpy as np
 
def setup():
    fig = plt.figure()                          # create a new figure for plotting
    ax1 = fig.add_subplot(111, projection='3d') # create a new subplot on our figure and set projection as 3d
    x = [0, 10, 5, 5]
    y = [0, 0, 10, 5]
    z = [0, 0, 0, 10]
    return ax1, x, y, z

def midpoint(x, y, z):
    
    for i in range(100):
        a = random.randint(0,len(x)-1) #select first starting point
        b = random.randint(0, len(x)-1) #select other starting point #TODO: make sure they are not the same point
        x.append((x[a] + x[b]) / 2)
        y.append((y[a] + y[b]) / 2)
        z.append((z[a] + z[b]) / 2)

    return x, y, z

def display(ax1, x, y, z):
    ax1.scatter(x, y, z, c = 'm', marker = 'o')
    ax1.set_xlabel('x-axis')
    ax1.set_ylabel('y-axis')
    ax1.set_zlabel('z-axis')
    # f = plt.figure()
    # f.set_figwidth(20)
    # f.set_figheight(20)
    plt.show()

 
# function to show the plot


if __name__ == "__main__":
    ax1, x, y, z = setup()
    midpoint(x,y,z)
    display(ax1, x, y, z)
    

