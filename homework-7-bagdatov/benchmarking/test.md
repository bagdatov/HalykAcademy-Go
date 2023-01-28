## pprof cpu

```bash
Fetching profile over HTTP from http://localhost:8080/debug/pprof/profile?seconds=10
Saved profile in /home/alibi/pprof/pprof.web.samples.cpu.001.pb.gz
File: web
Type: cpu
Time: Oct 25, 2021 at 4:39pm (+06)
Duration: 10.13s, Total samples = 1.13s (11.16%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 630ms, 55.75% of 1130ms total
Showing top 10 nodes out of 225
      flat  flat%   sum%        cum   cum%
     210ms 18.58% 18.58%      210ms 18.58%  runtime.epollwait
     120ms 10.62% 29.20%      120ms 10.62%  syscall.Syscall
      70ms  6.19% 35.40%       70ms  6.19%  runtime.futex
      40ms  3.54% 38.94%      400ms 35.40%  runtime.findrunnable
      40ms  3.54% 42.48%       40ms  3.54%  syscall.Syscall6
      30ms  2.65% 45.13%       40ms  3.54%  fmt.(*pp).handleMethods
      30ms  2.65% 47.79%      120ms 10.62%  fmt.(*pp).printValue
      30ms  2.65% 50.44%       30ms  2.65%  runtime.epollctl
      30ms  2.65% 53.10%       30ms  2.65%  runtime.memclrNoHeapPointers
      30ms  2.65% 55.75%      240ms 21.24%  runtime.netpoll
```
## pprof heap

```bash

Fetching profile over HTTP from http://localhost:8080/debug/pprof/heap?seconds=10
Saved profile in /home/alibi/pprof/pprof.web.alloc_objects.alloc_space.inuse_objects.inuse_space.002.pb.gz
File: web
Type: alloc_space
Time: Oct 25, 2021 at 4:45pm (+06)
Duration: 10s, Total samples = 30.77MB 
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 28.07MB, 91.21% of 30.77MB total
Showing top 10 nodes out of 71
      flat  flat%   sum%        cum   cum%
       9MB 29.25% 29.25%        9MB 29.25%  reflect.packEface
    5.05MB 16.41% 45.66%     5.05MB 16.41%  web/api.(*userBase).showAll
    4.02MB 13.05% 58.71%     4.02MB 13.05%  bufio.NewReaderSize
       2MB  6.50% 65.21%        2MB  6.50%  net/textproto.(*Reader).ReadMIMEHeader
       2MB  6.50% 71.71%        2MB  6.50%  net/url.parseQuery
    1.50MB  4.88% 76.59%     3.50MB 11.38%  net/http.readRequest
    1.50MB  4.88% 81.46%     1.50MB  4.88%  mime.ParseMediaType
       1MB  3.25% 84.71%     4.50MB 14.63%  net/http.(*conn).readRequest
       1MB  3.25% 87.96%        1MB  3.25%  net.(*conn).Read
       1MB  3.25% 91.21%        1MB  3.25%  syscall.anyToSockaddr
```