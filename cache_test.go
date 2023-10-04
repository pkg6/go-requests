package requests

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

var (
	_newCache = NewFileCache("./.testdata/cache")
)

// 测试过期的缓存被驱逐
func TestCacheEvict(t *testing.T) {
	err := _newCache.Set("test", "value", 1*time.Second)
	require.NoError(t, err)
	time.Sleep(2 * time.Second) // 等待过期
	_, err = _newCache.Get("test")
	require.Error(t, err) // 应该获取不到
}

// 测试设置和获取字符串缓存
func TestSetGet(t *testing.T) {
	key := "testkey"
	value := "test value"
	err := _newCache.Set(key, value, 0)
	require.NoError(t, err)
	v, err := _newCache.Get(key)
	require.NoError(t, err)
	require.Equal(t, v, value)
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		_ = _newCache.Set(key, strings.Repeat("a", 100), 0)
	}
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		_, _ = _newCache.Get(key)
	}
}
