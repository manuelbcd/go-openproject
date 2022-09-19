package openproject

import (
	"sync"
)

type saveBox struct {
	Index int
	Res   interface{}
}

// AutoPageTurn auto page turn
// @notice Use careful when dealing large amounts of data because it will set all objects in memory.
// Usage case:
//
//	users, err := AutoPageTurn(nil, 10, testClient.User.GetList)
func AutoPageTurn[T IPaginationResponse](filter *FilterOptions, pageSize int,
	fetch func(*FilterOptions, int, int) (T, *Response, error)) (T, error) {
	var res T
	offset := 1
	if pageSize == 0 {
		pageSize = 10
	}
	// First request get total count
	var err error
	res, _, err = fetch(filter, offset, pageSize)
	if err != nil {
		return res, err
	}
	if res.TotalPage() < 2 {
		// less 2 page, return directly
		return res, nil
	}
	totalPage := res.TotalPage()
	wg := &sync.WaitGroup{}
	offset += 1
	wg.Add(totalPage - 1)

	var box []saveBox
	box = append(box, saveBox{Index: 1, Res: res})

	// use more goroutine for speed up
	for i := offset; i <= totalPage; i++ {
		go func(offset int) {
			defer wg.Done()
			pageRes, _, err := fetch(filter, offset, pageSize)
			if err != nil {
				return
			}
			box = append(box, saveBox{Index: offset, Res: pageRes})
			offset += 1
		}(i)
	}

	// wait all goroutine done
	wg.Wait()
	// sort, use simple bubble sort -_-
	for i := 0; i < len(box); i++ {
		for j := i; j < len(box); j++ {
			if box[i].Index > box[j].Index {
				box[i], box[j] = box[j], box[i]
			}
		}
	}
	var tmpRes = box[1].Res.(T)
	for i := 1; i < len(box); i++ {
		tmpRes.ConcatEmbed(box[i].Res.(T))
	}
	return tmpRes, nil
}
