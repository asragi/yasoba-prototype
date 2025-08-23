package debug

import (
	"fmt"
	"runtime"
	"time"
)

type MemoryLeakDetector struct {
	lastGCAlloc      uint64
	lastNumGC        uint32
	leakThreshold    uint64 // MB単位
	lastAlertTime    time.Time
	displayTime      int
	detectedIncrease uint64
}

func NewMemoryLeakDetector(thresholdMB uint64, displayTime int) *MemoryLeakDetector {
	return &MemoryLeakDetector{
		leakThreshold: thresholdMB * 1024 * 1024, // MBをバイトに変換
		lastNumGC:     0,
		displayTime:   displayTime,
	}
}

func (d *MemoryLeakDetector) Check(memory runtime.MemStats) (bool, string) {
	record := func() {
		d.lastGCAlloc = memory.Alloc
		d.lastNumGC = memory.NumGC
	}

	displayNotification := func() string {
		return fmt.Sprintf("MEMORY LEAK DETECTED! +%d MB after GC", d.detectedIncrease/1024/1024)
	}
	returnText := func(detected bool) (bool, string) {
		if d.displayTime > 0 {
			d.displayTime--
		}
		if !detected {
			if d.displayTime > 0 {
				return false, displayNotification()
			}
			return false, "Memory Leak: OK"
		}
		d.displayTime = 300
		return true, displayNotification()

	}

	// GCが発生したかチェック（NumGCの変化を監視）
	if memory.NumGC <= d.lastNumGC {
		return returnText(false)
	}

	// 前回のGC後のAllocと比較
	fmt.Printf("Alloc: %d, lastAlloc: %d\n", memory.Alloc, d.lastGCAlloc)
	if memory.Alloc < d.lastGCAlloc {
		record()
		return returnText(false)
	}

	increase := memory.Alloc - d.lastGCAlloc
	if increase < d.leakThreshold {
		record()
		return returnText(false)
	}

	fmt.Printf("Increase: %d, alloc: %d, lastAlloc: %d\n", increase, memory.Alloc, d.lastGCAlloc)
	d.lastAlertTime = time.Now()
	record()
	d.detectedIncrease = increase
	return returnText(true)
}
