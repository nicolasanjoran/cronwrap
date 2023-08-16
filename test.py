import time

for i in range(3):
    print(i)
    time.sleep(1)

raise Exception('test exception')