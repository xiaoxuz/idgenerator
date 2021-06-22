package idgenerator

import (
	"fmt"
	//"github.com/imroc/biu"
	"math"
	"testing"
)

func TestGeneGenerator(t *testing.T) {

	var nodeID int64 = 1
	username := "18600000000"
	g, _ := NewGeneGenerator(nodeID, []byte(username), 1)
	testCnt := 10
	ch := make(chan int64, testCnt)
	for i := 0; i < testCnt; i++ {
		go func() {
			ch <- g.Generate()
		}()
	}

	idMap := map[int64]int64{}
	for i := 0; i < testCnt; i++ {
		id := <-ch
		fmt.Println(id)
		if _, ok := idMap[id]; ok == true {
			t.Errorf("id repeat: %d", id)
		}
		idMap[id] = 1
	}
	sampleGeneID := ExtractGene([]byte(username))
	// 水平分N张
	partNum := int64(math.Pow(2, 4))
	samplePart := sampleGeneID % partNum
	t.Logf("partNum:%d sample part:%d", partNum, samplePart)
	for k, _ := range idMap {
		idPart := k % partNum
		if samplePart != idPart {
			t.Errorf("hash table part discord. samplePart:%d idPart:%d", samplePart, idPart)
		}
	}

	t.Log("Done")
}
