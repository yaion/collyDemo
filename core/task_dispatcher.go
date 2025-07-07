package core

import (
	"log"
	"sync"
	"time"
)

type TaskDispatcher struct {
	accountPool *AccountPool
	taskChan    chan *Task
	wg          sync.WaitGroup
	mu          sync.Mutex
	stop        chan struct{}

	// 新增字段
	activeTasks  int
	activeMu     sync.Mutex
	currentTasks sync.Map
}

func NewTaskDispatcher(pool *AccountPool) *TaskDispatcher {
	return &TaskDispatcher{
		accountPool: pool,
		taskChan:    make(chan *Task, 1000000),
		stop:        make(chan struct{}),
	}
}

var addTaskLimiter = time.NewTicker(10 * time.Millisecond)

func (d *TaskDispatcher) AddTask(task *Task) {
	select {
	case <-d.stop:
		log.Printf("调度器已停止，拒绝新任务: %s", task.URL)
		return
	default:
	}

	<-addTaskLimiter.C

	// 非阻塞式添加任务，避免死锁
	select {
	case d.taskChan <- task:
		// 成功添加任务
	default:
		// 队列满时的处理
		log.Printf("任务队列已满，等待空间: %s", task.URL)
		for len(d.taskChan) > 900 { // 留一些缓冲
			time.Sleep(100 * time.Millisecond)
		}
		d.taskChan <- task
	}
}

func (d *TaskDispatcher) TaskStatus() (queueLen, active int) {
	d.mu.Lock()
	queueLen = len(d.taskChan)
	d.mu.Unlock()

	d.activeMu.Lock()
	active = d.activeTasks
	d.activeMu.Unlock()

	return queueLen, active
}

func (d *TaskDispatcher) Run(concurrency int) {
	log.Printf("启动任务调度器，并发数: %d", concurrency)

	// 启动监控协程
	go d.monitorTaskQueue()

	// 启动工作池
	for i := 0; i < concurrency; i++ {
		d.wg.Add(1)
		go d.worker(i)
	}

	d.wg.Wait()
	close(d.stop)
}

func (d *TaskDispatcher) worker(id int) {
	log.Printf("Worker %d 启动", id)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Worker %d panic: %v", id, r)
		}
		d.wg.Done()
		log.Printf("Worker %d 退出", id)
	}()

	for {
		select {
		case <-d.stop:
			return
		case task, ok := <-d.taskChan:
			if !ok {
				return
			}

			if task == nil {
				continue
			}

			log.Printf("Worker %d 接收到任务: %s", id, task.URL)

			// 增加活跃任务计数
			d.activeMu.Lock()
			d.currentTasks.Store(task.URL, time.Now())
			d.activeTasks++
			d.activeMu.Unlock()

			// 获取账号
			acc := d.accountPool.GetAccount()
			log.Printf("Worker %d 获取账号: %s, 执行任务: %s", id, acc.UserName, task.URL)

			// 带重试的执行
			retry := 0
			maxRetries := 3
			var lastErr error

			for retry < maxRetries {
				if err := ExecuteRequest(task, acc, d); err != nil {
					lastErr = err
					log.Printf("Worker %d 请求失败 (尝试 %d/%d): %v", id, retry+1, maxRetries, err)
					retry++
					if retry < maxRetries {
						time.Sleep(time.Duration(retry) * time.Second)
					}
				} else {
					lastErr = nil
					break
				}
			}

			if lastErr != nil {
				log.Printf("Worker %d 任务最终失败: %s, 错误: %v", id, task.URL, lastErr)
			}

			// 减少活跃任务计数
			d.activeMu.Lock()
			d.currentTasks.Delete(task.URL)
			d.activeTasks--
			d.activeMu.Unlock()

			log.Printf("Worker %d 完成任务: %s", id, task.URL)
		}
	}
}

func (d *TaskDispatcher) monitorTaskQueue() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	zeroCount := 0

	for {
		select {
		case <-d.stop:
			return
		case <-ticker.C:
			queueLen, active := d.TaskStatus()
			log.Printf("任务监控: 队列=%d, 执行中=%d, 连续空闲=%ds", queueLen, active, zeroCount)

			// 检查长时间执行的任务
			d.activeMu.Lock()
			d.currentTasks.Range(func(key, value interface{}) bool {
				startTime := value.(time.Time)
				if time.Since(startTime) > 2*time.Minute {
					log.Printf("警告: 任务执行时间过长: %s, 已执行: %v", key, time.Since(startTime))
				}
				return true
			})
			d.activeMu.Unlock()

			// 检查是否连续空闲
			if queueLen == 0 && active == 0 {
				zeroCount++
				if zeroCount >= 2 { // 10秒空闲就报告
					log.Println("所有任务已完成，连续10秒无新任务")
					zeroCount = 0
				}
			} else {
				zeroCount = 0
			}
		}
	}
}

func (d *TaskDispatcher) Stop() {
	close(d.stop)
	close(d.taskChan)
}
