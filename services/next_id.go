package services

import (
	"fmt"
	"github.com/sony/sonyflake"
	"sync"
	"time"
)

var (
	NextIDGenerator *NextID
	once            sync.Once
)

// GetNextIDGenerator 返回唯一的 NextID 实例（单例）
func GetNextIDGenerator() *NextID {
	once.Do(func() {
		NextIDGenerator = newNextID()
	})
	return NextIDGenerator
}

// NextID 封装唯一 ID 生成器
type NextID struct {
	sf *sonyflake.Sonyflake
}

// newNextID 初始化 NextID 实例
func newNextID() *NextID {
	st := sonyflake.Settings{
		StartTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), // 自定义纪元时间
		MachineID: func() (uint16, error) {
			return 1, nil // 返回机器 ID，确保分布式唯一性
		},
	}

	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		fmt.Println("Failed to create Sonyflake instance")
		return nil
	}

	return &NextID{sf: sf}
}

// Generate 生成唯一 ID
func (n *NextID) Generate() (uint64, error) {
	id, err := n.sf.NextID()
	if err != nil {
		return 0, err
	}
	return id, nil
}
