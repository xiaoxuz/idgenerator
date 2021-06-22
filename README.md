### idgenerator
A simple ID generator.

### 常见模型
  - 数据库自增
  - UUID
  - 号段模式
  - 雪花算法(snowflake)
  - 基因算法(gene)
  
### How to use? 
## 雪花算法
```go
    var clusterID int64 = 1
	var nodeID int64 = 1
    var step int64 = 1
    // 创建实例
	g := NewGenerator(clusterID, nodeID, step)
    // 生成 ID
    id := g.Generate()
```
`clusterID: 集群 ID`

`nodeID: 集群节点 ID`

`step： 序列号自增步长`

## 基因算法
提取样本基因，融入到分布式全局唯一 ID 中。
> tips: 此 demo 沿用snowflake 算法实现全局唯一 ID
```go
    var nodeID int64 = 1
    var step int64 = 1
	sample := "18600000000"
    // 创建实例
	g := NewGeneGenerator(nodeID, []byte(sample), step)
    // 生成基因 ID
    id := g.Generate()
```
`nodeID: 集群节点 ID`

`sample: 样本`

`step： 序列号自增步长`

Testing:
```go
    sampleGeneID := ExtractGene([]byte(username))
	// 水平分N张
	partNum := int64(math.Pow(2, 4))
	samplePart := sampleGeneID % partNum
	t.Logf("partNum:%d sample part:%d", partNum, samplePart)
	for k,_ := range idMap{
		idPart := k % partNum
		if samplePart != idPart {
			t.Errorf("hash table part discord. samplePart:%d idPart:%d", samplePart, idPart)
		}
	}
```

### 更多信息
![avatar](https://github.com/xiaoxuz/limiter/blob/main/wechat.png)

  

