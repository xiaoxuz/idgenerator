package idgenerator

import (
	"testing"
)

func TestGenerator(t *testing.T) {

	var clusterID int64 = 1
	var nodeID int64 = 1
	g, _ := NewGenerator(clusterID, nodeID, 2)

	testCnt := 10000
	ch := make(chan int64, testCnt)
	for i := 0; i < testCnt; i++ {
		go func() {
			ch <- g.Generate()
		}()
	}

	idMap := map[int64]int64{}
	for i := 0; i < testCnt; i++ {
		id := <-ch
		//fmt.Println(id)
		if _, ok := idMap[id]; ok == true {
			t.Errorf("id repeat: %d", id)
		}
		idMap[id] = 1
	}
	t.Log("Done")
}
