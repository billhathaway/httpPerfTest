httpPerfTest is a simple test server that can be used to serve arbitrary sized responses.

Options are:  
* -s <responseSize>  defaults to 128 bytes  
* -p <httpPort>      defaults to 8080  
* -r <pprofPort>     defaults to 8081  
* -m <GOMAXPROCS>    defaults to 1  

The max speed I have seen so far is ~ 450k req/second when running with -m 32 and driving the test with  
`wrk -c 500 -d 10 -t 10`  
 against localhost on an EC2 c3.8xl instance  

 pprof is enabled and can be reached at http://localhost:8081/debug/pprof/profiling  

