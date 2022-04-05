## 第六章 并行

### Golang scheduler 原理

![scheduler 原理](./images/ch6/scheduler.png)

如上图所示：

- 每个`M` 对应着一个计算机核心线程
- 每个`P` 绑定一个 `M` ，基于自己的队列完成基于 `G` 的调度 
- 每个`G` 是一个协程，类似于线程，但是更轻量级，这里这要理解，是`Golang` 应用程序调度的最小单位即可。
- `Golang` 的调度器充分考虑了高并发场景
  - `G1` 调用网络IO 操作时会将其挂到`Net Puller`的`P`上，现有`P`可以调度其他`G`继续执行，避免CPU浪费，同时当`G1`网络IO结束后，自动会调度到活动的`P`上，完成指定工作；
  - `G2` 进行非`异步system call`时，会新建出一个新的`M2`，对应的P会绑定到`M2`上，`P`基于`M2`继续执行代码，当`G2`完成`system call`后，将其放入某个`P`的`LRQ`上；
  - `Golang`的Scheduler是用户态的调度，所以没有内核态和用户态之间的频繁切换；所以基于`go routine`的调度，context switching的消耗小很多。操作系统context switching 需要消耗1000~1500ns 大概消耗12K-18K的系统指令执行时间；golang的调度时间只有400ns，消耗只有2.4K左右的系统指令；

