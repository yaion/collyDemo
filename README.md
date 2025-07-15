# 考古家数据采集系统

这是一个基于Go语言和Colly框架的数据采集系统，用于定时采集考古家网站的各种数据。

## 功能特性

### 数据采集类型
- **主要数据**: 达人、品牌、直播、商品、店铺、视频
- **排名数据**: 各种排行榜数据（涨粉榜、热销榜、带货榜等）
- **热点数据**: 实时热点、飙升热点、探测爆款等

### 定时任务系统
- 支持灵活的时间间隔配置
- 自动重试机制
- 任务状态监控
- 优雅关闭

## 系统架构

```
collyDemo/
├── main.go              # 主程序入口
├── core/                # 核心模块
│   ├── scheduler.go     # 定时任务调度器
│   ├── task_config.go   # 任务配置管理
│   ├── task_dispatcher.go # 任务分发器
│   ├── account_pool.go  # 账号池管理
│   └── types.go         # 类型定义
├── handlers/            # 数据处理器
│   ├── author_handler.go    # 达人数据处理器
│   ├── brand_handler.go     # 品牌数据处理器
│   ├── live_handler.go      # 直播数据处理器
│   ├── product_handler.go   # 商品数据处理器
│   ├── store_handler.go     # 店铺数据处理器
│   ├── video_handler.go     # 视频数据处理器
│   └── rank_handler.go      # 排名数据处理器
├── mongodb/             # 数据库模型
│   ├── author.go        # 达人模型
│   ├── brand.go         # 品牌模型
│   ├── live.go          # 直播模型
│   ├── product.go       # 商品模型
│   ├── store.go         # 店铺模型
│   ├── video.go         # 视频模型
│   └── ...              # 其他排名模型
├── config/              # 配置文件
│   └── schedule_config.go # 定时任务配置
└── pkg/utils/           # 工具函数
    └── utils.go         # 通用工具
```

## 安装和运行

### 1. 环境要求
- Go 1.23.1+
- MongoDB 4.0+

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 配置MongoDB
修改 `mongodb/mongodb.go` 中的连接字符串：
```go
clientOptions := options.Client().ApplyURI("mongodb://your-mongodb-host:27017")
```

### 4. 配置账号信息
在 `main.go` 中配置你的账号信息：
```go
accounts := []*core.Account{
    {
        ID:       "1",
        UserName: "your-username",
        Password: "your-password",
        Token:    "your-token",
        // ... 其他配置
    },
}
```

### 5. 运行程序
```bash
go run main.go
```

## 定时任务配置

### 默认配置
系统提供了默认的定时任务配置：

```go
// 主要数据采集任务
config.MainTasks.Author = "6h"   // 每6小时采集一次达人数据
config.MainTasks.Brand = "12h"   // 每12小时采集一次品牌数据
config.MainTasks.Live = "2h"     // 每2小时采集一次直播数据
config.MainTasks.Product = "4h"  // 每4小时采集一次商品数据
config.MainTasks.Store = "8h"    // 每8小时采集一次店铺数据
config.MainTasks.Video = "3h"    // 每3小时采集一次视频数据

// 排名数据采集任务
config.RankTasks.AuthorFansIncreaseRank = "1h"    // 每小时采集一次
// ... 其他排名任务
```

### 自定义配置
你可以通过修改 `config/schedule_config.go` 来自定义任务执行频率：

```go
func GetDefaultConfig() *ScheduleConfig {
    config := &ScheduleConfig{}
    
    // 自定义主要任务频率
    config.MainTasks.Author = "4h"   // 改为每4小时
    config.MainTasks.Live = "1h"     // 改为每小时
    
    // 自定义排名任务频率
    config.RankTasks.AuthorFansIncreaseRank = "30m"  // 改为每30分钟
    
    return config
}
```

### 支持的时间格式
- `"1h"` - 1小时
- `"30m"` - 30分钟
- `"2h30m"` - 2小时30分钟
- `"daily"` - 每天
- `"weekly"` - 每周
- `"monthly"` - 每月

## 系统监控

系统提供了实时监控功能，每5分钟输出一次状态信息：

```
=== 系统状态监控 ===
任务队列长度: 15
活跃任务数: 3
定时任务状态:
  author_tasks: 达人数据采集 (下次执行: 2024-01-15 14:00:00)
  brand_tasks: 品牌数据采集 (下次执行: 2024-01-15 18:00:00)
  live_tasks: 直播数据采集 (下次执行: 2024-01-15 12:00:00)
==================
```

## 数据存储

所有采集的数据都存储在MongoDB中，数据库名为 `kaogujia`，包含以下集合：

- `authors` - 达人数据
- `brands` - 品牌数据
- `live` - 直播数据
- `products` - 商品数据
- `stores` - 店铺数据
- `videos` - 视频数据
- `author_fans_increase_ranks` - 达人涨粉榜
- `author_fans_decrease_ranks` - 达人掉粉榜
- `product_hot_sale_ranks` - 商品热销榜
- `live_author_sales_ranks` - 直播达人带货榜
- `hot_video_ranks` - 热门视频榜
- `real_time_hot_spots` - 实时热点
- `soaring_hot_spots` - 飙升热点
- `explore_hot_bursts` - 探测爆款
- 等等...

## 错误处理

系统具备完善的错误处理机制：

1. **自动重试**: 任务失败时自动重试，最多重试3次
2. **延迟重试**: 重试间隔递增，避免频繁请求
3. **错误日志**: 详细的错误日志记录
4. **优雅降级**: 单个任务失败不影响其他任务

## 性能优化

1. **并发控制**: 可配置的并发数量
2. **账号池**: 多账号轮换使用，避免频率限制
3. **请求限流**: 内置频率限制器
4. **内存管理**: 批量数据处理，减少内存占用

## 扩展开发

### 添加新的数据处理器

1. 在 `handlers/` 目录下创建新的处理器文件
2. 实现处理器函数：
```go
func NewDataHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
    // 处理逻辑
    return nil
}
```

3. 在 `main.go` 中注册处理器：
```go
taskScheduler.RegisterHandler("new_data", handlers.NewDataHandler)
```

### 添加新的数据模型

1. 在 `mongodb/` 目录下创建新的模型文件
2. 定义数据结构和DAO方法
3. 在处理器中使用新的模型

## 注意事项

1. **遵守网站规则**: 请遵守目标网站的robots.txt和使用条款
2. **合理频率**: 不要设置过于频繁的采集间隔
3. **数据备份**: 定期备份MongoDB数据
4. **监控资源**: 注意系统资源使用情况

## 许可证

本项目仅供学习和研究使用，请勿用于商业用途。 