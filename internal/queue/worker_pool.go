package queue

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Task 任务接口
type Task interface {
	Execute() error
	GetID() string
}

// WorkerPool 工作池
type WorkerPool struct {
	workers    int
	taskQueue  chan Task
	wg         sync.WaitGroup
	logger     *zap.Logger
	ctx        context.Context
	cancel     context.CancelFunc
	retryCount int
	retryDelay time.Duration
}

// NewWorkerPool 创建新的工作池
func NewWorkerPool(workers int, queueSize int, logger *zap.Logger) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:    workers,
		taskQueue:  make(chan Task, queueSize),
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		retryCount: 3,
		retryDelay: time.Second * 2,
	}
}

// Start 启动工作池
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	wp.logger.Info("工作池已启动",
		zap.Int("workers", wp.workers),
		zap.Int("queue_size", cap(wp.taskQueue)),
	)
}

// Stop 停止工作池
func (wp *WorkerPool) Stop() {
	wp.cancel()
	close(wp.taskQueue)
	wp.wg.Wait()
	wp.logger.Info("工作池已停止")
}

// Submit 提交任务
func (wp *WorkerPool) Submit(task Task) error {
	select {
	case wp.taskQueue <- task:
		return nil
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	default:
		wp.logger.Warn("任务队列已满，任务被拒绝",
			zap.String("task_id", task.GetID()),
		)
		return ErrQueueFull
	}
}

// worker 工作协程
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case task, ok := <-wp.taskQueue:
			if !ok {
				return
			}

			// 执行任务，带重试
			if err := wp.executeWithRetry(task); err != nil {
				wp.logger.Error("任务执行失败",
					zap.Int("worker_id", id),
					zap.String("task_id", task.GetID()),
					zap.Error(err),
				)
			} else {
				wp.logger.Debug("任务执行成功",
					zap.Int("worker_id", id),
					zap.String("task_id", task.GetID()),
				)
			}
		}
	}
}

// executeWithRetry 带重试执行任务
func (wp *WorkerPool) executeWithRetry(task Task) error {
	var lastErr error
	for i := 0; i < wp.retryCount; i++ {
		if err := task.Execute(); err != nil {
			lastErr = err
			if i < wp.retryCount-1 {
				time.Sleep(wp.retryDelay)
				wp.logger.Debug("任务重试",
					zap.String("task_id", task.GetID()),
					zap.Int("attempt", i+1),
					zap.Int("max_retries", wp.retryCount),
				)
			}
		} else {
			return nil
		}
	}
	return lastErr
}

// ErrQueueFull 队列已满错误
var ErrQueueFull = &QueueFullError{}

// QueueFullError 队列已满错误类型
type QueueFullError struct{}

func (e *QueueFullError) Error() string {
	return "任务队列已满"
}

