# Channels and Goroutines
Just a self-reminder program to reference when working with goroutines and implementing context driven gracefull shutdown

## What it does
    1. Takes an n number of threads from args and instantiates n goroutines
    2. Creates a worker that notifies a broadcast channel every 1 second to trigger all of the current goroutines
    3. Goroutines will act when receiveing data from the channel and when context is canceled to gracefully shutdown
