# Path "/" Method : GET x4
curl localhost:2112/
curl localhost:2112/
curl localhost:2112/
curl localhost:2112/

# Path "/" Method : POST x2
curl -XPOST localhost:2112/
curl -XPOST localhost:2112/

# Path "/" Method : PUT x1
curl -XPUT localhost:2112/

# Path "/" Method : DELETE x1
curl -XDELETE localhost:2112/


# Path "/test1" Method : GET x2
curl localhost:2112/test1
curl localhost:2112/test1

# Path "/test2" Method : POST x2
curl -XPOST localhost:2112/test2
curl -XPOST localhost:2112/test2

# Path "/test3" Method : PUT x2
curl -XPUT localhost:2112/test3
curl -XPUT localhost:2112/test3

# Path "/test4" Method : DELETE x2
curl -XDELETE localhost:2112/test4
curl -XDELETE localhost:2112/test4