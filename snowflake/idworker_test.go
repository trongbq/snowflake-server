// go test -run none -bench . -benchtime 3s -benchmem

package snowflake

import (
	"testing"
)

var temp int64

// BenchmarkNextID provides performance testing for IDWorker#NextID function
func BenchmarkNextID(b *testing.B) {
	iw, _ := NewIDWorker(1)
	var id int64

	for i := 0; i < b.N; i++ {
		id = iw.NextID()
	}

	temp = id
}

func TestNextID(t *testing.T) {
	iw, _ := NewIDWorker(1)

	t.Log("Given now is large than w.lastTime")
	{
		t.Log("\tWhen generate next ID")
		{
			iw.lastTime = 999
			id := nextID(iw, 1000)
			var want int64 = (1000 << TimestampShift) | (1 << MachineIDShift) | 0
			if id != want {
				t.Fatalf("\tF\t Should generate an ID %d, instead of %d", want, id)
			}
			t.Logf("\tOK\t Should generate expected ID %d", id)
		}
	}

	t.Log("Given now is equal to w.lastTime")
	{
		t.Log("\tWhen generate next ID")
		{
			iw.lastTime = 999
			id := nextID(iw, 999)
			var want int64 = (999 << TimestampShift) | (1 << MachineIDShift) | 1
			if id != want {
				t.Fatalf("\tF\t Should generate an ID %d, instead of %d", want, id)
			}
			t.Logf("\tOK\t Should generate expected ID %d", id)
		}
	}
}
