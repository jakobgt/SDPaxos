# SDPaxos
The prototype implementation and the extended version paper of SDPaxos,
a new state machine replication protocol for efficient geo-replication.

## What's in the repo?
- The extended version of our paper. This version provides a correctness proof of our protocol (in the Appendix).
- The source code of our prototype implementation. This implementation is based on the [codebase of EPaxos](https://github.com/efficient/epaxos).

## How to build and run
```
export GOPATH=[...]/SDPaxos
go install master
go install server
go install client

bin/master &
bin/server -port 7070 &
bin/server -port 7071 &
bin/server -port 7072 &
bin/client
```
The above commands (`bin/server`) by default execute Multi-Paxos. You can add an argument `-n` to run SDPaxos. For more argument options, you can see `src/server/server.go`.

## Related paper
[SDPaxos: Building Efficient Semi-Decentralized Geo-replicated State Machines](https://dl.acm.org/citation.cfm?id=3267837) (ACM Symposium on Cloud Computing 2018, SoCC '18)


## Performance numbers:

Setup with AWS, one node in each of

- us-east-1
- us-east-2 (leader is here)
- us-west-2

The node is a t2.large instance.

Running `bin/clientread -maddr <master_ip> -check -q 1000 -c 0 -l <local_replica_id> -T 10` gives

in us-east-1:

|  PARAMETERS   |          VALUES          |
|---------------|:------------------------:|
| Concurrency   | 10                       |
| Duration      | 13.8574489s              |
| Iterations    | 5000                     |
| Successes     | 5000 (100.00%)           |
| Errors        | 0 (0.00%)                |
| IPS           | 360.82                   |
| Latency (max) | 36.485518ms              |
| Latency (p99) | 25ms                     |
| Latency (p95) | 20ms                     |
| Latency (p50) | 15ms                     |
| Latency (min) | 11.817146ms              |


in us-east-2 (where the leader is)

|  PARAMETERS   |     VALUES     |
|---------------|:--------------:|
| Concurrency   | 10             |
| Duration      | 12.134331.02s  |
| Iterations    | 5000           |
| Successes     | 5000 (100.00%) |
| Errors        | 0 (0.00%)      |
| IPS           | 412.05         |
| Latency (max) | 21.430887ms    |
| Latency (p99) | 15ms           |
| Latency (p95) | 15ms           |
| Latency (p50) | 15ms           |
| Latency (min) | 11.77583ms     |

in us-west-2:

|  PARAMETERS   |     VALUES     |
|---------------|:--------------:|
| Concurrency   | 10             |
| Duration      | 76.778272.02s  |
| Iterations    | 5000           |
| Successes     | 5000 (100.00%) |
| Errors        | 0 (0.00%)      |
| IPS           | 65.12          |
| Latency (max) | 127.492946ms   |
| Latency (p99) | 90ms           |
| Latency (p95) | 85ms           |
| Latency (p50) | 75ms           |
| Latency (min) | 28.134µs       |
