package cache

import (
	"context"
	"strings"
	"time"
	"user/internal/model"

	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"github.com/zhufuyi/sponge/pkg/utils"
)

const (
	// cache prefix key, must end with a colon
	tbVirtualHostInfoCachePrefixKey = "tbVirtualHostInfo:"
	// TbVirtualHostInfoExpireTime expire time
	TbVirtualHostInfoExpireTime = 10 * time.Minute
)

var _ TbVirtualHostInfoCache = (*tbVirtualHostInfoCache)(nil)

// TbVirtualHostInfoCache cache interface
type TbVirtualHostInfoCache interface {
	Set(ctx context.Context, id uint64, data *model.TbVirtualHostInfo, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.TbVirtualHostInfo, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.TbVirtualHostInfo, error)
	MultiSet(ctx context.Context, data []*model.TbVirtualHostInfo, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// tbVirtualHostInfoCache define a cache struct
type tbVirtualHostInfoCache struct {
	cache cache.Cache
}

// NewTbVirtualHostInfoCache new a cache
func NewTbVirtualHostInfoCache(cacheType *model.CacheType) TbVirtualHostInfoCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	var c cache.Cache
	if strings.ToLower(cacheType.CType) == "redis" {
		c = cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.TbVirtualHostInfo{}
		})
	} else {
		c = cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.TbVirtualHostInfo{}
		})
	}

	return &tbVirtualHostInfoCache{
		cache: c,
	}
}

// GetTbVirtualHostInfoCacheKey cache key
func (c *tbVirtualHostInfoCache) GetTbVirtualHostInfoCacheKey(id uint64) string {
	return tbVirtualHostInfoCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *tbVirtualHostInfoCache) Set(ctx context.Context, id uint64, data *model.TbVirtualHostInfo, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetTbVirtualHostInfoCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *tbVirtualHostInfoCache) Get(ctx context.Context, id uint64) (*model.TbVirtualHostInfo, error) {
	var data *model.TbVirtualHostInfo
	cacheKey := c.GetTbVirtualHostInfoCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *tbVirtualHostInfoCache) MultiSet(ctx context.Context, data []*model.TbVirtualHostInfo, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetTbVirtualHostInfoCacheKey(v.FuniqueID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *tbVirtualHostInfoCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.TbVirtualHostInfo, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetTbVirtualHostInfoCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.TbVirtualHostInfo)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.TbVirtualHostInfo)
	for _, id := range ids {
		val, ok := itemMap[c.GetTbVirtualHostInfoCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *tbVirtualHostInfoCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetTbVirtualHostInfoCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *tbVirtualHostInfoCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetTbVirtualHostInfoCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
