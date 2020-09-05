package snowflake

import (
    "time"
    "errors"
    "sync"
    "encoding/json"
)

const (
    SequenceBits = 12
    MachineBits = 10

    MachineIDMax = -1 ^ (-1 << MachineBits)
    MachineIDShift int64 = SequenceBits
    TimestampShift int64 = MachineBits + SequenceBits
    SequenceMask = -1 ^ (-1 << SequenceBits)

    // Our schema ID has 41 bits, equivalent to 69 years support, using TIMESTAMP_EPOCH,
    // we can support generate ID til year 2079, you can change this value to recent time
    // in order to increase supported upper bound time
    // 1288834974657
    TwEpoch int64 = 1288834974657
)

var ErrInvalidMachineID = errors.New("invalid machine ID")

var TwEpochTime = time.Unix(TwEpoch/1000, (TwEpoch % 1000) * 1000000)

type IDWorker struct {
    machineID int
    lastTime int64
    sequence int
    mu sync.Mutex
}

func NewIDWorker(mID int) (*IDWorker, error) {
    if mID < 0 || mID > MachineIDMax {
        return nil, ErrInvalidMachineID
    }

    w := IDWorker {
        machineID: mID,
        lastTime: -1,
    }

    return &w, nil
}

func (w *IDWorker) NextID() int64 {
    w.mu.Lock()
    defer w.mu.Unlock()

    now := time.Since(TwEpochTime).Milliseconds()
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

func (w IDWorker) Stats() ([]byte, error) {
    stats := struct {
        MachineID int
        LastTime int64
        Sequence int
    } {
        MachineID: w.machineID,
        LastTime: w.lastTime,
        Sequence: w.sequence,
    }

    return json.Marshal(stats)
}

func tilNextMillis(lt int64) int64 {
    now := time.Since(TwEpochTime).Milliseconds()
    for now <= lt {
        now = time.Since(TwEpochTime).Milliseconds()
    }
    return now
}
