package zutils

import (
	"fmt"
	"sync"
	"testing"
)

func TestMapShardInt(t *testing.T) {
	m := NewMap[uint64]()

	m.Set(1000, 1100)

	wg := sync.WaitGroup{}
	for i := 0; i < 5000000; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			m.Set(uint64(x), 111111)
		}(i)
	}

	for i := 0; i < 5000000; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			_, _ = m.Get(uint64(x))
		}(i)
	}

	for i := 0; i < 5000000-5; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			m.Delete(uint64(x))
		}(i)
	}

	wg.Wait()

	cnt := m.Len()
	fmt.Println("cnt", cnt)

	v, ok := m.Get(2222)
	fmt.Println("v", v, "ok", ok)

	cnt1 := m.Len()
	fmt.Println("cnt1", cnt1)

	m.Range(func(key uint64, value uint64) bool {
		fmt.Println("---key", key, "value", value)
		return true
	})

	all := m.All()
	fmt.Printf("---all:%+v\n", all)
}

func TestMapShardStr(t *testing.T) {
	m := NewMapStr[int]()

	m.Set("1", 1)

	wg := sync.WaitGroup{}
	for i := 0; i < 5000000; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			m.Set(fmt.Sprintf(`%v`, x), 666666)
		}(i)
	}

	for i := 0; i < 5000000; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			_, _ = m.Get(fmt.Sprintf(`%v`, x))
		}(i)
	}

	for i := 0; i < 5000000-5; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			m.Delete(fmt.Sprintf(`%v`, x))
		}(i)
	}

	wg.Wait()

	cnt := m.Len()
	fmt.Println("cnt", cnt)

	m.Range(func(key string, value int) bool {
		fmt.Println("---key", key, "value", value)
		return true
	})

	all := m.All()
	fmt.Printf("---all:%+v\n", all)
}
