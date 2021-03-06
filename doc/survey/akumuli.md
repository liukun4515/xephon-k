# Akumuli

- http://akumuli.org/akumuli/2016/09/14/next/
  - https://docs.google.com/document/d/1jFK8E3CZSqR5IPsMGojm2LknkNyUZA7tY51N6IgzW_g/pub
- https://github.com/at15/papers-i-read/issues/39

## TODO

- [ ] How many level would be in memory
- [ ] If want to be distributed, how to handle the address (pointer), file offset?
partition a large series? what's the physical limit of a single series?
- [ ] interleave different series would cause little performance penalty when read on SSD, but would
it cause trouble on spin drive? The block size is pretty small.
- [ ] In leaf node, how are the compressed data points stored? `t1, v1, t2, v2` or `t1, t2, v1, v2`

## Meta



## Vocabulary

Bear with my English

- fan-out: the number of inputs that can be connected to a specified output.

## NB+ tree

https://github.com/akumuli/Akumuli/blob/master/libakumuli/storage_engine/nbtree.h

- [ ] the topmost superblock is called rightmost https://github.com/akumuli/Akumuli/blob/master/libakumuli/storage_engine/nbtree.h#L62
which might explain the graph I didn't understand in the Google doc

Inner Node

- 4KB (2^12 bytes)
- 32 links, each contains aggregate (each link is 2^7 bytes)
  - first & last point timestamp and value (2^5 bytes = 64 bit * 2 * 2)
  - small & largest timestamp and value (2^5 bytes)
  - sum of all values in subtree (2^3 bytes)
    - [ ] what if it overflow
  - number of points in the subtree (2^3 bytes)
  - [ ] when to update the aggregate values, when the subtree is full?

Links

- inner node -> inner node
- inner node -> leaf node
- leaf node -> leaf node
- back references during crash recovery
- [ ] no link from leaf node to inner node?

## Stuff

@Lazin I came across the blog of older version of Akumuli that use delta and run length encoding for timestamp, and for the newer version it is dropped because the leaf node is small (4KB) I suppose?

Is it a good idea to use the delta and run length encoding for timestamp if I plan to store data in large chunk, like 64MB instead of 4KB of a single series and compress timestamp and value separately.

``````
series id | offset-t | offset-v | t1, t2, t3 .... | v1, v2, v3 .....
``````

And for current Akumuli implementation, since different series (tree) are store in same file interleaved, if I need to scan a long time range of a single series, I need to issue many concurrent reads because leaf nodes of a single series is not likely contiguous. For SSD you can have parallel read, but for HDD the performance might be poor. Also [this article](http://codecapsule.com/2014/02/12/coding-for-ssds-part-5-access-patterns-and-system-optimizations/) says
'a single large read is better than many small concurrent reads'.

@Lazin one file per series will also have the problem of running out of inode as mentioned by some guy from prometheus about their previous engine (https://fabxc.org/blog/2017-04-10-writing-a-tsdb/). But by increasing in memory cache size, you can have very large data block
