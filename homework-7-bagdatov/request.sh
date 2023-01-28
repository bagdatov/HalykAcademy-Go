for i in {1..1000}
do
    # parallel requests   
    curl -X POST -F "user=test$i" http://localhost:8080/saveone & curl http://localhost:8080/getone?search="test$i" & curl http://localhost:8080/getall

    # serial requests
    # curl -X POST -F "user=test$i" http://localhost:8080/saveone
    # curl http://localhost:8080/getone?search="test$i"
    # curl http://localhost:8080/getall
done

# go tool pprof web http://localhost:8080/debug/pprof/profile?seconds=10
# go tool pprof -alloc_space web http://localhost:8080/debug/pprof/heap?seconds=10
# go tool pprof -inuse_space web http://localhost:8080/debug/pprof/heap?seconds=10