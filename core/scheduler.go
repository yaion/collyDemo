package core

import (
	"log"
	"sync"
	"time"
)

// ScheduledTask 定时任务结构
type ScheduledTask struct {
	ID          string
	Name        string
	Description string
	Schedule    string // cron表达式或时间间隔
	Handler     func() error
	LastRun     time.Time
	NextRun     time.Time
	Enabled     bool
	mu          sync.Mutex
}

// Scheduler 定时任务调度器
type Scheduler struct {
	tasks         map[string]*ScheduledTask
	dispatcher    *TaskDispatcher
	taskScheduler *TaskScheduler
	stop          chan struct{}
	wg            sync.WaitGroup
	mu            sync.RWMutex
}

// NewScheduler 创建定时任务调度器
func NewScheduler(dispatcher *TaskDispatcher, taskScheduler *TaskScheduler) *Scheduler {
	return &Scheduler{
		tasks:         make(map[string]*ScheduledTask),
		dispatcher:    dispatcher,
		taskScheduler: taskScheduler,
		stop:          make(chan struct{}),
	}
}

// AddTask 添加定时任务
func (s *Scheduler) AddTask(id, name, description, schedule string, handler func() error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	nextRun := s.calculateNextRun(schedule, now)

	task := &ScheduledTask{
		ID:          id,
		Name:        name,
		Description: description,
		Schedule:    schedule,
		Handler:     handler,
		LastRun:     time.Time{},
		NextRun:     nextRun,
		Enabled:     true,
	}

	s.tasks[id] = task
	log.Printf("添加定时任务: %s (%s), 下次执行时间: %s", name, schedule, nextRun.Format("2006-01-02 15:04:05"))
}

// RemoveTask 移除定时任务
func (s *Scheduler) RemoveTask(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, exists := s.tasks[id]; exists {
		task.Enabled = false
		delete(s.tasks, id)
		log.Printf("移除定时任务: %s", task.Name)
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	log.Println("启动定时任务调度器")

	s.wg.Add(1)
	go s.run()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	log.Println("停止定时任务调度器")
	close(s.stop)
	s.wg.Wait()
}

// run 运行调度器主循环
func (s *Scheduler) run() {
	defer s.wg.Done()

	ticker := time.NewTicker(1 * time.Minute) // 每分钟检查一次
	defer ticker.Stop()

	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			s.checkAndExecuteTasks()
		}
	}
}

// checkAndExecuteTasks 检查并执行到期的任务
func (s *Scheduler) checkAndExecuteTasks() {
	s.mu.RLock()
	tasks := make([]*ScheduledTask, 0, len(s.tasks))
	for _, task := range s.tasks {
		if task.Enabled {
			tasks = append(tasks, task)
		}
	}
	s.mu.RUnlock()

	now := time.Now()
	for _, task := range tasks {
		task.mu.Lock()
		if task.Enabled && now.After(task.NextRun) {
			// 执行任务
			go s.executeTask(task)
		}
		task.mu.Unlock()
	}
}

// executeTask 执行单个任务
func (s *Scheduler) executeTask(task *ScheduledTask) {
	task.mu.Lock()
	task.LastRun = time.Now()
	task.mu.Unlock()

	log.Printf("执行定时任务: %s", task.Name)

	if err := task.Handler(); err != nil {
		log.Printf("定时任务执行失败: %s, 错误: %v", task.Name, err)
	} else {
		log.Printf("定时任务执行成功: %s", task.Name)
	}

	// 计算下次执行时间
	task.mu.Lock()
	task.NextRun = s.calculateNextRun(task.Schedule, time.Now())
	task.mu.Unlock()

	log.Printf("定时任务下次执行时间: %s, %s", task.Name, task.NextRun.Format("2006-01-02 15:04:05"))
}

// calculateNextRun 计算下次执行时间
func (s *Scheduler) calculateNextRun(schedule string, from time.Time) time.Time {
	switch schedule {
	case "hourly":
		// 每小时执行一次
		return from.Add(1 * time.Hour)
	case "daily":
		// 每天执行一次
		return from.Add(24 * time.Hour)
	case "weekly":
		// 每周执行一次
		return from.Add(7 * 24 * time.Hour)
	case "monthly":
		// 每月执行一次
		return from.AddDate(0, 1, 0)
	default:
		// 尝试解析时间间隔 (如 "2h", "30m", "1d")
		if duration, err := time.ParseDuration(schedule); err == nil {
			return from.Add(duration)
		}
		// 默认每小时执行一次
		return from.Add(1 * time.Hour)
	}
}

// GetTaskStatus 获取任务状态
func (s *Scheduler) GetTaskStatus() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status := make(map[string]interface{})
	for id, task := range s.tasks {
		task.mu.Lock()
		status[id] = map[string]interface{}{
			"name":        task.Name,
			"description": task.Description,
			"schedule":    task.Schedule,
			"enabled":     task.Enabled,
			"last_run":    task.LastRun,
			"next_run":    task.NextRun,
		}
		task.mu.Unlock()
	}

	return status
}

// 预定义的定时任务处理器
func (s *Scheduler) createMainTasksHandler() func() error {
	return func() error {
		log.Println("执行主要数据采集任务")
		s.taskScheduler.AddMainTasks()
		return nil
	}
}

func (s *Scheduler) createRankTasksHandler() func() error {
	return func() error {
		log.Println("执行排名数据采集任务")
		s.taskScheduler.AddRankTasks()
		return nil
	}
}

func (s *Scheduler) createAuthorTasksHandler() func() error {
	return func() error {
		log.Println("执行达人数据采集任务")
		// 只添加达人相关任务
		tasks := GetMainTasks()
		if len(tasks) > 0 {
			task := &Task{
				URL:     tasks[0].URL,
				Method:  tasks[0].Method,
				Headers: s.taskScheduler.headers,
				Body:    tasks[0].Body,
				Handler: s.taskScheduler.handlers["author"],
				Meta:    tasks[0].Meta,
			}
			s.dispatcher.AddTask(task)
		}
		return nil
	}
}

func (s *Scheduler) createBrandTasksHandler() func() error {
	return func() error {
		log.Println("执行品牌数据采集任务")
		tasks := GetMainTasks()
		if len(tasks) > 1 {
			task := &Task{
				URL:     tasks[1].URL,
				Method:  tasks[1].Method,
				Headers: s.taskScheduler.headers,
				Body:    tasks[1].Body,
				Handler: s.taskScheduler.handlers["brand"],
				Meta:    tasks[1].Meta,
			}
			s.dispatcher.AddTask(task)
		}
		return nil
	}
}

func (s *Scheduler) createLiveTasksHandler() func() error {
	return func() error {
		log.Println("执行直播数据采集任务")
		tasks := GetMainTasks()
		if len(tasks) > 2 {
			task := &Task{
				URL:     tasks[2].URL,
				Method:  tasks[2].Method,
				Headers: s.taskScheduler.headers,
				Body:    tasks[2].Body,
				Handler: s.taskScheduler.handlers["live"],
				Meta:    tasks[2].Meta,
			}
			s.dispatcher.AddTask(task)
		}
		return nil
	}
}

func (s *Scheduler) createProductTasksHandler() func() error {
	return func() error {
		log.Println("执行商品数据采集任务")
		tasks := GetMainTasks()
		if len(tasks) > 3 {
			task := &Task{
				URL:     tasks[3].URL,
				Method:  tasks[3].Method,
				Headers: s.taskScheduler.headers,
				Body:    tasks[3].Body,
				Handler: s.taskScheduler.handlers["product"],
				Meta:    tasks[3].Meta,
			}
			s.dispatcher.AddTask(task)
		}
		return nil
	}
}

func (s *Scheduler) createStoreTasksHandler() func() error {
	return func() error {
		log.Println("执行店铺数据采集任务")
		tasks := GetMainTasks()
		if len(tasks) > 4 {
			task := &Task{
				URL:     tasks[4].URL,
				Method:  tasks[4].Method,
				Headers: s.taskScheduler.headers,
				Body:    tasks[4].Body,
				Handler: s.taskScheduler.handlers["store"],
				Meta:    tasks[4].Meta,
			}
			s.dispatcher.AddTask(task)
		}
		return nil
	}
}

func (s *Scheduler) createVideoTasksHandler() func() error {
	return func() error {
		log.Println("执行视频数据采集任务")
		tasks := GetMainTasks()
		if len(tasks) > 5 {
			task := &Task{
				URL:     tasks[5].URL,
				Method:  tasks[5].Method,
				Headers: s.taskScheduler.headers,
				Body:    tasks[5].Body,
				Handler: s.taskScheduler.handlers["video"],
				Meta:    tasks[5].Meta,
			}
			s.dispatcher.AddTask(task)
		}
		return nil
	}
}

// InitDefaultTasks 初始化默认定时任务
func (s *Scheduler) InitDefaultTasks() {
	// 主要数据采集任务 - 每天执行一次
	s.AddTask("main_tasks", "主要数据采集", "采集达人、品牌、直播、商品、店铺、视频数据", "daily", s.createMainTasksHandler())

	// 排名数据采集任务 - 每小时执行一次
	s.AddTask("rank_tasks", "排名数据采集", "采集各种排名数据", "hourly", s.createRankTasksHandler())

	// 单独的数据采集任务 - 可以设置不同的执行频率
	s.AddTask("author_tasks", "达人数据采集", "采集达人数据", "6h", s.createAuthorTasksHandler())
	s.AddTask("brand_tasks", "品牌数据采集", "采集品牌数据", "12h", s.createBrandTasksHandler())
	s.AddTask("live_tasks", "直播数据采集", "采集直播数据", "2h", s.createLiveTasksHandler())
	s.AddTask("product_tasks", "商品数据采集", "采集商品数据", "4h", s.createProductTasksHandler())
	s.AddTask("store_tasks", "店铺数据采集", "采集店铺数据", "8h", s.createStoreTasksHandler())
	s.AddTask("video_tasks", "视频数据采集", "采集视频数据", "3h", s.createVideoTasksHandler())
}

// InitTasksWithConfig 使用配置文件初始化定时任务
func (s *Scheduler) InitTasksWithConfig(mainSchedules, rankSchedules map[string]string) {
	// 主要数据采集任务
	for taskType, schedule := range mainSchedules {
		switch taskType {
		case "author":
			s.AddTask("author_tasks", "达人数据采集", "采集达人数据", schedule, s.createAuthorTasksHandler())
		case "brand":
			s.AddTask("brand_tasks", "品牌数据采集", "采集品牌数据", schedule, s.createBrandTasksHandler())
		case "live":
			s.AddTask("live_tasks", "直播数据采集", "采集直播数据", schedule, s.createLiveTasksHandler())
		case "product":
			s.AddTask("product_tasks", "商品数据采集", "采集商品数据", schedule, s.createProductTasksHandler())
		case "store":
			s.AddTask("store_tasks", "店铺数据采集", "采集店铺数据", schedule, s.createStoreTasksHandler())
		case "video":
			s.AddTask("video_tasks", "视频数据采集", "采集视频数据", schedule, s.createVideoTasksHandler())
		}
	}

	// 排名数据采集任务
	for taskType, schedule := range rankSchedules {
		switch taskType {
		case "author_fans_increase_rank":
			s.AddTask("author_fans_increase_rank_tasks", "达人涨粉榜采集", "采集达人涨粉榜数据", schedule, s.createRankTasksHandler())
		case "author_fans_decrease_rank":
			s.AddTask("author_fans_decrease_rank_tasks", "达人掉粉榜采集", "采集达人掉粉榜数据", schedule, s.createRankTasksHandler())
		case "author_potential_rank":
			s.AddTask("author_potential_rank_tasks", "达人带货潜力榜采集", "采集达人带货潜力榜数据", schedule, s.createRankTasksHandler())
		case "product_hot_sale_rank":
			s.AddTask("product_hot_sale_rank_tasks", "商品热销榜采集", "采集商品热销榜数据", schedule, s.createRankTasksHandler())
		case "product_real_time_sales_rank":
			s.AddTask("product_real_time_sales_rank_tasks", "商品实时销量榜采集", "采集商品实时销量榜数据", schedule, s.createRankTasksHandler())
		case "live_author_sales_rank":
			s.AddTask("live_author_sales_rank_tasks", "直播达人带货榜采集", "采集直播达人带货榜数据", schedule, s.createRankTasksHandler())
		case "live_hot_push_rank":
			s.AddTask("live_hot_push_rank_tasks", "直播热推榜采集", "采集直播热推榜数据", schedule, s.createRankTasksHandler())
		case "hot_video_rank":
			s.AddTask("hot_video_rank_tasks", "热门视频榜采集", "采集热门视频榜数据", schedule, s.createRankTasksHandler())
		case "ecommerce_video_rank":
			s.AddTask("ecommerce_video_rank_tasks", "电商视频榜采集", "采集电商视频榜数据", schedule, s.createRankTasksHandler())
		case "video_hot_push":
			s.AddTask("video_hot_push_tasks", "视频热推采集", "采集视频热推数据", schedule, s.createRankTasksHandler())
		case "hot_sale_shop":
			s.AddTask("hot_sale_shop_tasks", "热销小店采集", "采集热销小店数据", schedule, s.createRankTasksHandler())
		case "site_hourly_rank":
			s.AddTask("site_hourly_rank_tasks", "全站小时榜采集", "采集全站小时榜数据", schedule, s.createRankTasksHandler())
		case "sales_hourly_rank":
			s.AddTask("sales_hourly_rank_tasks", "带货小时榜采集", "采集带货小时榜数据", schedule, s.createRankTasksHandler())
		case "real_time_hot_spot":
			s.AddTask("real_time_hot_spot_tasks", "实时热点采集", "采集实时热点数据", schedule, s.createRankTasksHandler())
		case "soaring_hot_spot":
			s.AddTask("soaring_hot_spot_tasks", "飙升热点采集", "采集飙升热点数据", schedule, s.createRankTasksHandler())
		case "explore_hot_burst":
			s.AddTask("explore_hot_burst_tasks", "探测爆款采集", "采集探测爆款数据", schedule, s.createRankTasksHandler())
		}
	}
}
