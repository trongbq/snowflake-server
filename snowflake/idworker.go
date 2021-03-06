package snowflake

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

const (
	SequenceBits = 12
	MachineBits  = 10

	MachineIDMax         = -1 ^ (-1 << MachineBits)
	MachineIDShift int64 = SequenceBits
	TimestampShift int64 = MachineBits + SequenceBits
	SequenceMask         = -1 ^ (-1 << SequenceBits)

	// Our schema ID has 41 bits, equivalent to 69 years support, using TIMESTAMP_EPOCH,
	// we can support generate ID til year 2079, you can change this value to recent time
	// in order to increase supported upper bound time
	// 1288834974657
	TwEpoch int64 = 1288834974657
)

var ErrInvalidMachineID = errors.New("invalid machine ID")

var TwEpochTime time.Time

func init() {
	// https://go.googlesource.com/proposal/+/master/design/12914-monotonic.md
	n := time.Now()
	TwEpochTime = n.Add(time.Unix(TwEpoch/1000, (TwEpoch%1000)*1000000).Sub(n))
}

type IDWorker struct {
	machineID int
	lastTime  int64
	sequence  int
	mu        *sync.Mutex
}

func NewIDWorker(mID int) (*IDWorker, error) {
	if mID < 0 || mID > MachineIDMax {
		return nil, ErrInvalidMachineID
	}

	w := IDWorker{
		machineID: mID,
		lastTime:  -1,
		mu:        &sync.Mutex{},
	}

	return &w, nil
}

func (w *IDWorker) NextID() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Since(TwEpochTime).Milliseconds()
	return nextID(w, now)
}

func (w IDWorker) Stats() ([]byte, error) {
	stats := struct {
		MachineID int
		LastTime  int64
		Sequence  int
	}{
		MachineID: w.machineID,
		LastTime:  w.lastTime,
		Sequence:  w.sequence,
	}

	return json.Marshal(stats)
}

// Create a private function to handle real work to make it easier for testing
func nextID(w *IDWorker, now int64) int64 {
	if now == w.lastTime {
		w.sequence = (w.sequence + 1) & SequenceMask
		if w.sequence == 0 {
			now = tilNextMillis(w.lastTime)
		}
	} else {
		w.sequence = 0
	}

	w.lastTime = now

	return (now << TimestampShift) | (int64(w.machineID) << MachineIDShift) | int64(w.sequence)
}

func tilNextMillis(lt int64) int64 {
	now := time.Since(TwEpochTime).Milliseconds()
	for now <= lt {
		now = time.Since(TwEpochTime).Milliseconds()
	}
	return now
}
