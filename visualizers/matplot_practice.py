import matplotlib.pyplot as plt

x = [1,2,3]
y = [5,7,4]

x2 = [1,2,3]
y2 = [10,14,12]

plt.plot(x,y,label="First Line")
plt.plot(x2,y2,label="Second |")
plt.xlabel("Plot Number")
plt.ylabel("Plot points")
plt.title("Interesting Graph\nCheck it out")
plt.legend()
plt.plot()
plt.show()