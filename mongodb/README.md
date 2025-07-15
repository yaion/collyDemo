# MongoDB 模型文件说明

本目录包含了根据中文名称JSON文件生成的MongoDB模型文件，每个文件都遵循了`author.go`的模板结构。

## 已生成的文件列表

### 1. 地区小店和热销小店
- **文件**: `area_shop.go`
- **结构体**: `AreaShop`
- **DAO**: `AreaShopDAO`
- **集合名**: `area_shops`
- **对应JSON**: `地区小店.json`, `热销小店.json`

### 2. 实时热点和飙升热点
- **文件**: `hot_spot.go`
- **结构体**: `HotSpot`, `TrendPoint`
- **DAO**: `HotSpotDAO`
- **集合名**: `hot_spots`
- **对应JSON**: `实时热点.json`, `飙升热点.json`

### 3. 商品热销榜和商品实时销量榜
- **文件**: `product_rank.go`
- **结构体**: `ProductRank`
- **DAO**: `ProductRankDAO`
- **集合名**: `product_ranks`
- **对应JSON**: `商品热销榜.json`, `商品实时销量榜.json`

### 4. 达人涨粉榜、达人掉粉榜、达人带货潜力榜
- **文件**: `author_rank.go`
- **结构体**: `AuthorRank`
- **DAO**: `AuthorRankDAO`
- **集合名**: `author_ranks`
- **对应JSON**: `达人涨粉榜.json`, `达人掉粉榜.json`, `达人带货潜力榜.json`

### 5. 直播达人带货榜（日、周、月）
- **文件**: `live_author_rank.go`
- **结构体**: `LiveAuthorRank`
- **DAO**: `LiveAuthorRankDAO`
- **集合名**: `live_author_ranks`
- **对应JSON**: `直播达人带货榜日.json`, `直播达人带货榜周.json`, `直播达人带货榜月.json`

### 6. 图文带货榜、电商视频榜、热门视频榜、视频热推
- **文件**: `media_rank.go`
- **结构体**: `MediaRank`
- **DAO**: `MediaRankDAO`
- **集合名**: `media_ranks`
- **对应JSON**: `图文带货榜.json`, `电商视频榜.json`, `热门视频榜.json`, `视频热推.json`

### 7. 直播热推榜
- **文件**: `live_hot_rank.go`
- **结构体**: `LiveHotRank`
- **DAO**: `LiveHotRankDAO`
- **集合名**: `live_hot_ranks`
- **对应JSON**: `直播热推榜.json`

### 8. 全站小时榜、带货小时榜
- **文件**: `hourly_rank.go`
- **结构体**: `HourlyRank`
- **DAO**: `HourlyRankDAO`
- **集合名**: `hourly_ranks`
- **对应JSON**: `全站小时榜.json`, `带货小时榜.json`

### 9. 探测爆款
- **文件**: `explore_hot.go`
- **结构体**: `ExploreHot`
- **DAO**: `ExploreHotDAO`
- **集合名**: `explore_hots`
- **对应JSON**: `探测爆款.json`

## 每个DAO文件包含的方法

所有DAO文件都包含以下标准方法：

1. **NewXXXDAO(db *mongo.Database)** - 创建DAO实例
2. **Create(data *XXX)** - 创建单条记录
3. **BatchCreate(ctx context.Context, data []interface{})** - 批量创建记录
4. **GetByID(ctx context.Context, id string)** - 根据ID获取记录
5. **Update(id string, updateData *XXX)** - 更新记录
6. **Delete(id string)** - 删除记录
7. **ListAll(ctx context.Context, filter bson.M, page, limit int64)** - 分页查询所有记录

## 使用示例

```go
// 创建DAO实例
areaShopDAO := NewAreaShopDAO(db)

// 创建记录
areaShop := &AreaShop{
    ShopID: "test123",
    Name: "测试店铺",
    // ... 其他字段
}
err := areaShopDAO.Create(areaShop)

// 查询记录
result, err := areaShopDAO.GetByID(context.Background(), "test123")

// 分页查询
filter := bson.M{"rank": bson.M{"$lte": 10}}
result, err := areaShopDAO.ListAll(context.Background(), filter, 1, 20)
```

## 注意事项

1. 所有结构体都包含了`json`和`bson`标签
2. 批量创建方法使用`ReplaceOneModel`实现upsert操作
3. 分页查询默认按排名升序排序
4. 所有方法都包含适当的错误处理和日志记录
5. 集合名称使用英文命名，便于数据库管理 