#cassandra-stress write n=1000000 cl=QUORUM -schema "replication(strategy=NetworkTopologyStrategy,ap-northeast=3)" -node 172.31.30.22,172.31.28.32,172.31.30.100 -rate threads\>=4 threads\<=512
