package lru_bytes

import "testing"

func TestNew(t *testing.T) {
    l := New(4, 6)
    l.OnEvicted = func(key string, value []byte, fullType RemoveReason) {
        t.Logf("remove:%s %s %s", key, value, fullType)
    }
    l.Add("1", []byte("1")) // [1]
    l.Add("2", []byte("2")) // [2, 1]
    l.Add("3", []byte("3")) // [3, 2, 1]
    l.Range(func(key string, value interface{}) {
        t.Logf("key:%s val:%v", key, value)
    })

    l.Add("4", []byte("11")) // [11, 3, 2, 1]
    l.Add("5", []byte("22")) // [22, 11, 3, 2]
    l.Add("6", []byte("33")) // [33, 22, 11, 3]

    l.Range(func(key string, value interface{}) {
        t.Logf("key:%s val:%v", key, value)
    })

    t.Log("cache length:", l.Len())
    e, _ := l.Get("3")
    t.Log("get cache", e)
}
