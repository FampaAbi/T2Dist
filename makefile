datanode:
	go run DataNode.go
namenode:
	go run NameNode.go
cliente:
	go run Cliente.go
reset:
	 rm -r ./JoinBooks/*
	 rm -r ./Partes/*
	 rm -r ./SplitBooks/*
	 rm -r ./ChunksDownload/*
	 rm -fv LOG.txt
	 > LOG.txt
