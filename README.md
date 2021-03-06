# adaptive [![GoDoc](https://godoc.org/github.com/seiflotfy/adaptive?status.svg)](https://godoc.org/github.com/seiflotfy/adaptive)

A probabilistic datastructure to estimate **How many times did i see item "x" within the timerange "T"**, using a set Adaptive Count-Min Sketches.

The Adaptive Count-Min Sketch algorithm (Ada-CMS), is just CMS but with the update and query mechanisms adapted to use the pre-emphasis and de-emphasis mechanism.

For more information read the post by Adrian Colyer [Time Adaptive Sketches (Ada-Sketches) for Summarizing Data Streams](https://blog.acolyer.org/2016/07/21/time-adaptive-sketches-ada-sketches-for-summarizing-data-streams/) or [the official paper with the same title](https://www.cs.rice.edu/~as143/Papers/16-ada-sketches.pdf) by Anshumali Shrivastava, Arnd Christian König, Mikhail Bilenko) 

## Usage
```go
duration := time.Duration(720 * time.Hour) // 720 hours range
unit := time.Hour

// Create sketch queryable with
// duation = 720 hours range
// unit = 1 hour
// width per sketch = 2^9
// depth per sketch = 8
// alpha = 1.004 (used for emphasizing and de-emphasizing)
sks := NewSketches(duration, unit, 9, 7, 1.004)

item := []byte("foo")
t1 := time.Date(2017, 06, 03, 0, 0, 0, 0, time.UTC)
t2 := t1.Add(time.Hour)
t3 := t2.Add(time.Hour)
t4 := t3.Add(time.Hour)
count1 := uint64(1337)
count2 := uint64(100000)

// Update item for given timestamps
sks.Insert(item, t1, count1)
sks.Insert(item, t3, count2)

// Estimate count of item within time range [t1, t2]
got, _ := sks.Estimate(item, t1, t2)
fmt.Printf("Expected count for \"%s\" in timerange [%v, %v] to be %d, got %d\n",
    string(item), t1.Format(time.Kitchen), t2.Format(time.Kitchen), count1, got)

// Estimate count of item within time range [t1, t3]
got, _ = sks.Estimate(item, t1, t3)
fmt.Printf("Expected count for \"%s\" in timerange [%v, %v] to be %d, got %d\n",
    string(item), t1.Format(time.Kitchen), t3.Format(time.Kitchen), count1+count2, got)

// Estimate count of item within time range [t1, t3]
got, _ = sks.Estimate(item, t3, t4)
fmt.Printf("Expected count for \"%s\" in timerange [%v, %v] to be %d, got %d\n",
    string(item), t3.Format(time.Kitchen), t4.Format(time.Kitchen), count2, got)

// Output:
// Expected count for "foo" in timerange [12:00AM, 1:00AM] to be 1337, got 1337
// Expected count for "foo" in timerange [12:00AM, 2:00AM] to be 101337, got 101337
// Expected count for "foo" in timerange [2:00AM, 3:00AM] to be 100000, got 100000
```
