package lru

import "testing"

func TestNew(t *testing.T) {
    l := New(4)
    l.OnEvicted = func(key string, value interface{}, fullType RemoveReason) {
        t.Logf("remove:%s %d %v", key, value.(int), fullType)
    }
    l.Add("1", 1) // [1]
    l.Add("2", 2) // [2, 1]
    l.Add("3", 3) // [3, 2, 1]
    l.Range(func(key string, value interface{}) {
        t.Logf("[range 1] key:%s val:%v", key, value)
    })

    l.Get("2") // [2, 3, 1]
    l.Range(func(key string, value interface{}) {
        t.Logf("[range 2] key:%s val:%v", key, value)
    })
    l.Get("3") // [3, 2, 1]

    l.Add("4", 11) // [11, 3, 2, 1]
    l.Add("5", 22) // [22, 11, 3, 2]
    l.Add("6", 33) // [33, 22, 11, 3]

    l.Range(func(key string, value interface{}) {
        t.Logf("[range 3]key:%s val:%v", key, value)
    })

    t.Log("cache length:", l.Len())
    e, _ := l.Get("3")
    t.Log("get cache", e)
}
