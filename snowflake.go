package idgenerator

import (
	"sync"
	"time"
)

const (
	CLUSTERID_BITS int64 = 5
	NODEIDID_BITS  int64 = 5
	SEQ_BITS       int64 = 12

	CLUSTERID_MAX int64 = -1 ^ (-1 << CLUSTERID_BITS)
	NODEID_MAX    int64 = -1 ^ (-1 << NODEIDID_BITS)
	SEQ_MAX       int64 = -1 ^ (-1 << SEQ_BITS)

	// 41个字节存储时间，大约69年
	TIMESTEMP_SHIFT       = CLUSTERID_BITS + NODEIDID_BITS + SEQ_BITS
	CLUSTERID_SHIFT int64 = NODEIDID_BITS + SEQ_BITS
	NODEID_SHIFT    int64 = SEQ_BITS

	// 拒绝浪费，珍惜时间
	EPOCH int64 = 1624258189000

	// 默认步长
	DEFAULT_STEP_LONG = 1
)

type SnowFlaker struct {
	m         sync.Mutex
	timestemp int64
	clusterID int64
	nodeID    int64
	seq       int64
	step      int64
}

func NewGenerator(clusterID int64, nodeID int64, step int64) *SnowFlaker {
	if clusterID > CLUSTERID_MAX || nodeID > NODEID_MAX {
		panic("params invalid")
	}
	s := &SnowFlaker{
		m:         sync.Mutex{},
		timestemp: 0,
		clusterID: clusterID,
		nodeID:    nodeID,
		seq:       0,
		step:      step,
	}
	if step <= 0 {
		s.step = DEFAULT_STEP_LONG
	}

	return s
}

func (s *SnowFlaker) Generate() int64 {
	s.m.Lock()
	defer s.m.Unlock()

	now := time.Now().UnixNano() / 1e6 // 纳秒转毫秒
	if now == s.timestemp {
		s.seq = s.seq + s.step
		if s.seq > SEQ_MAX {
			for now <= s.timestemp {
				now = time.Now().UnixNano() / 1e6
			}
			s.seq = 0
		}
	} else {
		s.seq = 0
	}
	s.timestemp = now
	return s.timeBlock() | s.clusterBlock() | s.nodeBlock() | s.seqBlock()
}

func (s *SnowFlaker) timeBlock() int64 {
	return (s.timestemp - EPOCH) << TIMESTEMP_SHIFT
}

func (s *SnowFlaker) clusterBlock() int64 {
	return s.clusterID << CLUSTERID_SHIFT
}

func (s *SnowFlaker) nodeBlock() int64 {
	return s.nodeID << NODEID_SHIFT
}

func (s *SnowFlaker) seqBlock() int64 {
	return s.seq
}
