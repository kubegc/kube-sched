### 

#### Problem 1
![](http://latex.codecogs.com/gif.latex?Gandiva_{fair})中提到，只针对加速比最大和加速比最小的任务，进行GPU交易，交易的代价（一块好的GPU换几块差的GPU）是加速比第二大的任务的加速比。因为只考虑两个任务，交易代价也是确定的，相对好实现。

~~如果想要最大化吞吐量，那么则要求加速比最大的任务的好的GPU尽可能多，同时能够保证公平性，如何实现还没想太好。~~

目前使用简单的贪心算法，加速比最大的任务为victor，其他任务都是victims，victims按照加速比排序，从加速比最小的任务开始，在满足公平的前提下，尽可能分配差的GPU。剩下的GPU都分配给victor

#### TODO
GPU数据监控
感觉 DCGM + Prometheus 可以
