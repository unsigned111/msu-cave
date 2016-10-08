import random

MIN_VALUE = 0
MAX_VALUE = 255


def walk(start):
    state = start
    while True:
        if random.choice([True, False]):
            if state + 10 < MAX_VALUE:
                state += 10
        else:
            if state - 10 > MIN_VALUE:
                state += -10
        yield state