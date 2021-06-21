package idgenerator

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	GENE_NODEIDID_BITS int64 = 10
	GENE_SEQ_BITS      int64 = 13
	GENE_BITS          int64 = 12

	GENE_NODEID_MAX int64 = -1 ^ (-1 << GENE_NODEIDID_BITS)
	GENE_SEQ_MAX    int64 = -1 ^ (-1 << GENE_SEQ_BITS)

	/**
	 * 1 符号位  |  28 timestemp                      |  10 nodeID  		|  13 自增ID 		| 12 基因
	 * 0        |   00000000 00000000 00000000 0000  | 00000000 00      |  00000000 00000 	| 00000000 0000
	 */
	// 这块放弃了时间，为了保留基因。timestemp 单位 s，28个 bit 大约7 8年
	GENE_TIMESTEMP_SHIFT       = GENE_NODEIDID_BITS + GENE_SEQ_BITS + GENE_BITS
	GENE_NODEIDID_SHIFT  int64 = GENE_SEQ_BITS + GENE_BITS
	GENE_SEQ_SHIFT       int64 = GENE_BITS

	// 拒绝浪费，珍惜时间
	GENE_EPOCH int64 = 1624258189

	// 默认步长
	GENE_DEFAULT_STEP_LONG = 1
)

type GeneID struct {
	m         sync.Mutex
	timestemp int64
	nodeID    int64
	seq       int64
	geneID    int64
	step      int64
}

func NewGeneGenerator(nodeID int64, geneSample []byte, step int64) *GeneID {
	if nodeID > GENE_NODEID_MAX {
		panic("params invalid")
	}
	g := &GeneID{
		m:         sync.Mutex{},
		timestemp: 0,
		nodeID:    nodeID,
		seq:       0,
		step:      step,
		geneID:    ExtractGene(geneSample),
	}
	if step <= 0 {
		g.step = GENE_DEFAULT_STEP_LONG
	}
	return g
}

// 基于基因样本提取基因ID
func ExtractGene(geneSample []byte) int64 {
	gene := md5.Sum(geneSample)
	hashGeneValue := fmt.Sprintf("%x", gene)[29:32]
	geneID, _ := strconv.ParseInt(hashGeneValue, 16, 64)
	return geneID
}

func (g *GeneID) Generate() int64 {
	g.m.Lock()
	defer g.m.Unlock()

	now := time.Now().UnixNano() / 1e9 // 纳秒转秒
	if now == g.timestemp {
		g.seq = g.seq + g.step
		if g.seq > GENE_SEQ_MAX {
			for now <= g.timestemp {
				now = time.Now().UnixNano() / 1e9
			}
			g.seq = 0
		}
	} else {
		g.seq = 0
	}
	g.timestemp = now
	return g.timeBlock() | g.nodeBlock() | g.seqBlock() | g.geneBlock()
}

func (g *GeneID) timeBlock() int64 {
	return (g.timestemp - GENE_EPOCH) << GENE_TIMESTEMP_SHIFT
}

func (g *GeneID) nodeBlock() int64 {
	return g.nodeID << GENE_NODEIDID_SHIFT
}

func (g *GeneID) seqBlock() int64 {
	return g.seq << GENE_SEQ_SHIFT
}

func (g *GeneID) geneBlock() int64 {
	return g.geneID
}
