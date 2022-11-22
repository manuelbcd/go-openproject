package openproject

import (
	"fmt"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
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

	box := make(chan saveBox)
	// use more goroutine for speed up
	for i := offset; i <= totalPage; i++ {
		go request(box, filter, pageSize, i, fetch)
	}

	// Sort by red-black tree
	t := addResToTree(box, res.TotalPage())
	var tmpRes T
	for _, key := range t.Keys() {
		b, found := t.Get(key)
		if found {
			bx := b.(T)
			if isInterfaceZero(tmpRes) {
				tmpRes = bx
			} else {
				tmpRes.ConcatEmbed(bx)
			}
		}
	}

	return tmpRes, nil
}

func addResToTree(collection <-chan saveBox, size int) *rbt.Tree {
	t := rbt.NewWithIntComparator()
	for i := 0; i < size; i++ {
		res := <-collection
		fmt.Printf("get res %+v", res)
		t.Put(res.Index, res.Res)
	}
	return t
}

func request[T IPaginationResponse](ch chan<- saveBox, filter *FilterOptions, pageSize int, idx int, fetch func(*FilterOptions, int, int) (T, *Response, error)) {
	pageRes, _, err := fetch(filter, idx, pageSize)
	if err != nil {
		return
	}
	ch <- saveBox{Index: idx, Res: pageRes}
}

func isInterfaceZero[T any](val T) bool {
	switch v := any(val).(type) {
	case *SearchResultProject:
		return v == nil
	case *SearchResultQuery:
		return v == nil
	case *SearchResultStatus:
		return v == nil
	case *SearchResultUser:
		return v == nil
	case *SearchResultWP:
		return v == nil
	}
	return false
}
