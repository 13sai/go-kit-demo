
# 高可用设计

系统出现故障的原因多种多样，主要有以下这些：

- 网络问题，网络连接故障、网络带宽出现超时拥塞等；
- 性能问题，数据库出现慢查询、Java Full GC 导致执行长时间等待、CPU 使用率过高、硬盘 IO 过载、内存分配失败等；
- 安全问题，被网络攻击，如 DDoS 等；异常客户端请求，如爬虫等；
- 运维问题，需求变更频繁不可控，架构也在不断地被调整，监控问题等；
- 管理问题，没有梳理出关键服务以及服务的依赖关系，运行信息没有和控制系统同步；
- 硬件问题，硬盘损坏导致数据读取失败、网卡出错导致网络 IO 处理失败、交换机出问题、机房断电导致服务器失联，甚至是人祸（比如挖掘机挖断机房光缆，导致一整片机房网络中断）等。


## 系统可用性指标

系统可用性指标是衡量分布式系统高可用性的重要因素，它通常是指系统可用时间与总运行时间之比，即Availability=MTTF/(MTTF+MTTR)。

其中，MTTF（Mean Time To Failure）是指平均故障前的时间，一般是指系统正常运行的时间。系统的可靠性越高，MTTF 越长，即系统正常运行的时间越长。

MTTR（Mean Time To Recovery）是指平均修复时间，即从故障出现到故障修复的这段时间，也就是系统不可用的时间。MTTR 越短说明系统的可用性越高。

## 冗余设计

> 如何降低分布式中出现单点故障的可能性

分布式系统中单点故障不可取的，而降低单点故障的不二法门就是冗余设计，通过多点部署的方式，并且最好部署在不同的物理位置上，避免单机房中多点同时失败。冗余设计不仅可以提高服务的吞吐量，还可以在出现灾难时快速恢复。目前常见的冗余设计有主从设计和对等治理设计，其中主从设计又可以细分为一主多从、多主多从。

冗余设计中一个不可避免的问题是考虑分布式系统中的数据一致性，多个节点中冗余的数据追求强一致性还是最终一致性。即使节点提供无状态服务，也需要借助外部服务，比如数据库、分布式缓存等维护数据状态。

CAP 是描述分布式系统下节点数据同步的基本原则，分别指：

- Consistency，数据强一致性，各个节点中对于同一份数据在任意时刻都是一致的；
- Availablity，可用性，系统在任何情况下接收到客户端请求后，都能够给出响应；
- Partition Tolerance，分区容忍性，系统允许节点网络通信失败。


## 熔断设计

> 如何防止服务雪崩，保护服务调用者的资源。

在分布式系统中，一次完整的请求可能需要经过多个服务模块的通力合作，请求在多个服务中传递，服务对服务的调用会产生新的请求，这些请求共同组成了这次请求的调用链。当调用链中的某个环节，特别是下游服务不可用时，将会导致上游服务调用方不可用，最终将这种不可用的影响扩大到整个系统，导致整个分布式系统的不可用，引发服务雪崩现象。

为了避免这种情况，在下游服务不可用时，保护上游服务的可用性显得极其重要。对此，我们可以参考电路系统的断路器机制，在必要的时候“壮士断腕”，当下游服务因为过载或者故障出现各种调用失败或者调用超时现象时，及时“熔断”服务调用方和服务提供方的调用链，保护服务调用方资源，防止服务雪崩现象的出现。

断路器的基本设计图如下所示，由关闭、打开、半开三种状态组成。
![断路器](https://s0.lgstatic.com/i/image/M00/4D/C7/CgqCHl9bJX6AYR6WAACD8asiP4k125.png)

## 限流设计

- 拒绝服务，把多出来的请求拒绝掉。一般来说，好的限流系统在经受流量暴增情况时，会暂时拒绝周期时间内请求数量最大的客户端，这样可以在一定程度上把一些不正常的或者是带有恶意的高并发访问挡在“门外”。
- 服务降级，关闭或是把后端做降级处理，释放资源给主流程服务以支持更多的请求。降级有很多方式，一种是把一些不重要的服务给停掉，把 CPU、内存或是数据的资源让给更重要的功能；一种是数据接口只返回部分关键数据，减少数据查询处理链路；还有更快的一种是直接返回预设的缓存或者静态数据，不需要经过复杂的业务查询处理获取数据，从而能够响应更多的用户请求。
- 优先级请求，是指将目前系统的资源分配给优先级更高的用户，优先处理权限更高的用户的请求。
- 延时处理，在这种情况下，一般来说会使用缓冲队列来缓冲大量的请求，系统根据自身负载能力异步消费队列中的请求。如果该队列也满了，那么就只能拒绝用户请求。使用缓冲队列只是为了减缓压力，一般用于应对瞬时大量的流量削峰。
- 弹性伸缩，采用自动化运维的方式对相应的服务做自动化的伸缩。这种方案需要应用性能监控系统，能够感知到目前最繁忙的服务，并自动伸缩它们；还需要一个快速响应的自动化发布、部署和服务注册的运维系统。如果系统的处理压力集中在数据库这类不易自动扩容的外部服务，服务弹性伸缩意义不大。


限流设计最主要的思想是保证系统处理自身承载能力内的请求访问，拒绝或者延缓处理过量的流量，而这种思想主要依赖于它的限流算法。那接下来我们介绍两种常用的限流算法：漏桶算法和令牌桶算法。


## 降级设计

在应对大流量冲击时，可以尝试对请求的处理流程进行裁剪，去除或者异步化非关键流程的次要功能，保证主流程功能正常运转。

一般来说，降级时可以暂时“牺牲”的有：
- 降低一致性。从数据强一致性变成最终一致性，比如说原本数据实时同步方式可以降级为异步同步，从而系统有更多的资源处理响应更多请求。
- 关闭非关键服务。关闭不重要功能的服务，从而释放出更多的资源。
- 简化功能。把一些功能简化掉，比如，简化业务流程，或是不再返回全量数据，只返回部分数据。也可以使用缓存的方式，返回预设的缓存数据或者静态数据，不执行具体的业务数据查询处理。

## 无状态设计
在分布式系统设计中，倡导使用无状态化的方式设计开发服务模块。这里“无状态”的意思是指对于功能相同的服务模块，在服务内部不维护任何的数据状态，只会根据请求中携带的业务数据从外部服务比如数据库、分布式缓存中查询相关数据进行处理，这样能够保证请求到任意服务实例中处理结果都是一致的。

无状态设计的服务模块可以简单通过多实例部署的方式进行横向扩展，各服务实例完全对等，可以有效提高服务集群的吞吐量和可用性。但是如此一来，服务处理的性能瓶颈就可能出现在提供数据状态一致性的外部服务中。

## 幂等性设计
幂等性设计是指系统对于相同的请求，一次和多次请求获取到的结果都是一样的。幂等性设计对分布式系统中的超时重试、系统恢复有重要的意义，它能够保证重复调用不会产生错误，保证系统的可用性。一般我们认为声明为幂等性的接口或者服务出现调用失败是常态，由于幂等性的原因，调用方可以在调用失败后放心进行重新请求。

举个简单的例子，在一笔订单的支付中，订单服务向支付服务请求支付接口，由于网络抖动或者其他未知的因素导致请求没能及时返回，那么此时订单服务并不了解此次支付是否成功。如果支付接口是幂等性的，那我们就可以放心使用同一笔订单号重新请求支付，如果上次支付请求已经成功，将会返回支付成功；如果上次支付请求未成功，将会重新进行金额扣费。这样就能保证请求的正确进行，避免重复扣费的错误。

## 超时设计
鉴于目前网络传播的不稳定性，在服务调用的过程中，很容易出现网络包丢失的现象。如果在服务调用者发起调用请求处理结果时出现网络丢包，在请求结果返回之前，服务调用者的调用线程会一直被操作系统挂起；或者服务提供者处理时间过长，迟迟没返回结果，服务调用者的调用线程也会被同样挂起。当服务调用者中出现大量的这样被挂起的服务调用时，服务调用者中的线程资源就可能被耗尽，导致服务调用者无法创建新的线程处理其他请求。这时就需要超时设计了。

超时设计是指给服务调用添加一个超时计时器，在超时计时器到达之后，调用结果还没返回，就由服务调用者主动结束调用，关闭连接，释放资源。通过超时设计能够有效减少系统等待时间过长的服务调用，使服务调用者有更多的资源处理其他请求，提高可用性。但是需注意的是，要根据下游服务的处理和响应能力合理设置超时时间的长短，过短将会导致服务调用者难以获取到处理结果，过长将会导致超时设计失去意义。

## 重试设计
在很多时候，由于网络不可靠或者服务提供者宕机，服务调用者的调用很可能会失败。如果此时服务调用者中存在一定的重试机制，就能够在一定程度上减少服务失败的概率，提高服务可用性。

比如业务系统在某次数据库请求中，由于临时的网络原因，数据请求超时了，如果业务系统中具备一定的超时重试机制，根据请求参数再次向数据库请求数据，就能正常获取到数据，完成业务处理流程，避免该次业务处理失败。

使用重试设计的时候需要注意以下问题：

待重试的服务接口是否为幂等性。对于某些超时请求，请求可能在服务提供者中执行成功了，但是返回结果却在网络传输中丢失了，此时若重复调用非幂等性服务接口就很可能会导致额外的系统错误。

服务提供者是否只是临时不可用。对于无法快速恢复的服务提供者或者网络无法立即恢复的情况下，盲目的重试只会使情况更加糟糕，无脑地消耗服务调用方的 CPU 、线程和网络 IO 资源，过多的重试请求甚至可能会把不稳定的服务提供者打垮。在这种情况下建议你结合熔断设计对服务调用方进行保护。

## 接口缓存
接口缓存是应对大并发量请求，降低接口响应时间，提高系统吞吐量的有效手段。基本原理是在系统内部，对于某部分请求参数和请求路径完成相同的请求结果进行缓存，在周期时间内，这部分相同的请求结果将会直接从缓存中读取，减少业务处理过程的负载。

最简单的例子是在一些在线大数据查询系统中，查询系统会将周期时间内系统查询条件相同的查询结果进行缓存，加快访问速度。

但接口缓存同样有着它不适用的场景。接口缓存牺牲了数据的强一致性，因为它返回的过去某个时间节点的数据缓存，并非实时数据，这对于实时性要求高的系统并不适用。另外，接口缓存加快的是相同请求的请求速率，这对于请求差异化较大的系统同样无能为力，过多的缓存反而会大量浪费系统内存等资源。

## 实时监控和度量
由于分布式中服务节点众多，问题的定位变得异常复杂，对此建议对每台服务器资源使用情况和服务实例的性能指标进行实时监控和度量。最常见的方式是健康检查，通过定时调用服务提供给健康检查接口判断服务是否可用。

目前业内也有开源的监控系统 Prometheus，它监控各个服务实例的运行指标，并根据预设的阈值自动报警，及时通知相关开发运维人员进行处理。

## 常规化维护
定期清理系统的无用代码，及时进行代码评审，处理代码中 bad smell，对于无状态服务可以定期重启服务器减少内存碎片和防止内存泄漏……这些都是非常有效的提高系统可用性的运维手段。
