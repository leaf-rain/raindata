package clickhouse_sqlx

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/leaf-rain/raindata/common/ecode"
	"log"
	"math/rand"
	"sync"
	"time"
)

type ClickhouseCluster struct {
	lock        sync.Mutex
	clusterConn []*ShardConn
}

// Each shard has a pool.Conn which connects to one replica inside the shard.
// We need more control than replica single-point-failure.
func InitClusterConn(chCfg *ClickhouseConfig) (cc *ClickhouseCluster, err error) {
	cc = new(ClickhouseCluster)
	cc.lock.Lock()
	defer cc.lock.Unlock()
	cc.freeClusterConn()

	proto := clickhouse.Native
	if chCfg.Protocol == clickhouse.HTTP.String() {
		proto = clickhouse.HTTP
	}
	if chCfg.MaxOpenConns <= 0 {
		chCfg.MaxOpenConns = defaultMaxOpenConns
	}
	for _, replicas := range chCfg.Hosts {
		numReplicas := len(replicas)
		replicaAddrs := make([]string, numReplicas)
		for i, ip := range replicas {
			// Changing hostnames to IPs breaks TLS connections in many cases
			if !chCfg.Secure {
				if ips2, err := GetIP4Byname(ip); err == nil {
					ip = ips2[0]
				}
			}
			replicaAddrs[i] = fmt.Sprintf("%s", ip)
		}
		sc := &ShardConn{
			replicas: replicaAddrs,
			chCfg:    chCfg,
			opts: clickhouse.Options{
				Auth: clickhouse.Auth{
					Database: chCfg.DB,
					Username: chCfg.Username,
					Password: chCfg.Password,
				},
				Protocol:    proto,
				DialTimeout: time.Minute * 10,
			},
			writingPool: NewWorkerPool(chCfg.MaxOpenConns, 1),
		}
		if chCfg.Secure {
			tlsConfig := &tls.Config{}
			tlsConfig.InsecureSkipVerify = chCfg.InsecureSkipVerify
			sc.opts.TLS = tlsConfig
		}
		if proto == clickhouse.Native {
			sc.opts.MaxOpenConns = chCfg.MaxOpenConns
			sc.opts.MaxIdleConns = chCfg.MaxOpenConns
			sc.opts.ConnMaxLifetime = time.Minute * 10
		}
		sc.protocol = proto
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		idx := r.Intn(numReplicas)
		sc.nextRep = idx
		if _, _, err = sc.NextGoodReplica(idx); err != nil {
			return
		}
		cc.clusterConn = append(cc.clusterConn, sc)
	}
	return
}

func (cc *ClickhouseCluster) freeClusterConn() {
	for _, sc := range cc.clusterConn {
		sc.Close()
	}
	cc.clusterConn = []*ShardConn{}
}

func (cc *ClickhouseCluster) FreeClusterConn() {
	cc.lock.Lock()
	defer cc.lock.Unlock()
	cc.freeClusterConn()
}

func (cc *ClickhouseCluster) NumShard() (cnt int) {
	cc.lock.Lock()
	defer cc.lock.Unlock()
	return len(cc.clusterConn)
}

// GetShardConn select a clickhouse shard based on batchNum
func (cc *ClickhouseCluster) GetShardConn(batchNum int64) (sc *ShardConn) {
	cc.lock.Lock()
	defer cc.lock.Unlock()
	sc = cc.clusterConn[batchNum%int64(len(cc.clusterConn))]
	return
}

// CloseAll closed all connection and destroys the pool
func (cc *ClickhouseCluster) CloseAll() {
	cc.FreeClusterConn()
}

// ShardConn a datastructure for storing the clickhouse connection
type ShardConn struct {
	lock        sync.Mutex
	conn        *Conn
	dbVer       int
	opts        clickhouse.Options
	replicas    []string    //ip:port list of replicas
	nextRep     int         //index of next replica
	writingPool *WorkerPool //the all tasks' writing ClickHouse, cpu-net balance
	protocol    clickhouse.Protocol
	chCfg       *ClickhouseConfig
}

func (sc *ShardConn) SubmitTask(fn func()) (err error) {
	return sc.writingPool.Submit(fn)
}

// GetReplica returns the replica to which db connects
func (sc *ShardConn) GetReplica() (replica string) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	if sc.conn != nil {
		curRep := (len(sc.replicas) + sc.nextRep - 1) % len(sc.replicas)
		replica = sc.replicas[curRep]
	}
	return
}

// Close closes the current replica connection
func (sc *ShardConn) Close() {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	if sc.conn != nil {
		sc.conn.Close()
		sc.conn = nil
	}
	if sc.writingPool != nil {
		sc.writingPool.StopWait()
	}
}

func (sc *ShardConn) NextGoodReplica(failedVer int) (db *Conn, dbVer int, err error) {
	// todo:支持集群单节点
	sc.lock.Lock()
	defer sc.lock.Unlock()
	if sc.conn != nil {
		if sc.dbVer > failedVer {
			// Another goroutine has already done connection.
			// Notice: Why recording failure version instead timestamp?
			// Consider following scenario:
			// conn1 = NextGood(0); conn2 = NexGood(0); conn1.Exec failed at ts1;
			// conn3 = NextGood(ts1); conn2.Exec failed at ts2;
			// conn4 = NextGood(ts2) will close the good connection and break users.
			return sc.conn, sc.dbVer, nil
		}
		sc.conn.Close()
		sc.conn = nil
	}
	savedNextRep := sc.nextRep
	// try all replicas, including the current one
	conn := Conn{
		protocol: sc.protocol,
		ctx:      context.Background(),
	}
	for i := 0; i < len(sc.replicas); i++ {
		replica := sc.replicas[sc.nextRep]
		sc.opts.Addr = []string{replica}
		sc.nextRep = (sc.nextRep + 1) % len(sc.replicas)
		if sc.protocol == clickhouse.HTTP {
			// refers to https://github.com/ClickHouse/clickhouse-go/issues/1150
			// An obscure error in the HTTP protocol when using compression
			// disable compression in the HTTP protocol
			conn.db = clickhouse.OpenDB(&sc.opts)
			conn.db.SetMaxOpenConns(sc.chCfg.MaxOpenConns)
			conn.db.SetMaxIdleConns(sc.chCfg.MaxOpenConns)
			conn.db.SetConnMaxLifetime(time.Minute * 10)
		} else {
			sc.opts.Compression = &clickhouse.Compression{
				Method: clickhouse.CompressionLZ4,
			}
			conn.c, err = clickhouse.Open(&sc.opts)
		}
		if err != nil {
			log.Printf("clickhouse.Open failed, err:%v, replice:%s", err, replica)
			continue
		}
		sc.dbVer++
		log.Printf("clickhouse.Open succeeded, dvVer:%d, replica:%s", sc.dbVer, replica)
		sc.conn = &conn
		return sc.conn, sc.dbVer, nil
	}
	err = ecode.Newf("no good replica among replicas %v since %d", sc.replicas, savedNextRep)
	return nil, sc.dbVer, err
}
