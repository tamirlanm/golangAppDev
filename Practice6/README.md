Problem 2 Explanation
The final value is not always 1000. Because multiple goroutines increment the same shared counter concurrently without synchronization causing a data race and lost updates.