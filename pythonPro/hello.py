

def hello(x):
    if x == 5:
        print("end")
        return
    x = x+1
    print(x)
    hello(x)


if __name__ == "__main__":
    hello(0)
