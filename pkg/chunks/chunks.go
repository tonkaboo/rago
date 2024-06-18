package chunks

import (
	"fmt"
	"gems/pkg/chunk"
	"gems/pkg/model"
	"gems/pkg/repo"
	"sync"
)

// type chunks struct {
// 	chunks []Chunk
// }

type Chunks []*chunk.Chunk

// 如果已经抽象出chunks这个结构，那么可以把embed和insert分开处理了
// 不过都还是需要一个db实例来真正调用对吧，因为我既要支持mo又要支持faiss
// 所以要定义一个接口来兼容mo和faiss
func (cs Chunks) EmbedAndInsert(model model.Model, chunkRepo repo.Repo) error {
	var wg sync.WaitGroup
	errors := make(chan error, len(cs))

	for _, c := range cs {
		wg.Add(1)
		go func(c *chunk.Chunk) {
			defer wg.Done()

			embedding, err := model.Embed(c.Content)
			if err != nil {
				errors <- err
				return
			}
			c.Embedding = embedding

			if err := chunkRepo.InsertChunk(c); err != nil {
				errors <- err
			}
		}(c)
	}

	wg.Wait()
	close(errors)

	var resultError error
	for err := range errors {
		if err != nil {
			if resultError == nil {
				resultError = err
			} else {
				resultError = fmt.Errorf("%v; %v", resultError, err)
			}
		}
	}

	return resultError
}

// 我刚才一直在想的是需不需要在主函数里面去直接调用vecStore.Retrieve(chunks)，但看来更好的做法是chunks.Retrieve吧！
// 在哪里去初始化store！！！
// 一个xx对应的是一张表呢？chunks吗？vecStore吗？
// 设想这样的场景：mo里面已经存储了，那就不需要新加这个chunks结构，
// 只需要去加载数据库，直接Retrieve就行了
// 综上，Retrieve应该定义在vecStore实例中
