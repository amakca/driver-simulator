package examples

import (
	"fmt"
	set "practice/internal/driver/simulator"
	m "practice/internal/models"
	storage "practice/internal/storage"
	"time"
)

func main() {
	gensettings := set.GeneralSettings{
		ProgramLiveTime: time.Minute * 2,
		GenOptimization: true,
	}

	sets := m.DriverSettings{
		General: &gensettings,
		Tags:    map[m.DataID]m.Formatter{},
	}
	sets.Tags[m.DataID(1)] = &set.TagSettings{
		PollTime:  25 * time.Millisecond,
		GenConfig: "rand:30ms:1.0:3.0",
	}
	sets.Tags[m.DataID(2)] = &set.TagSettings{
		PollTime:  25 * time.Millisecond,
		GenConfig: "rand:30ms:1.0:2.0",
	}
	sets.Tags[m.DataID(3)] = &set.TagSettings{
		PollTime:  30 * time.Millisecond,
		GenConfig: "rand:30ms:1.0:2.0",
	}
	sets.Tags[m.DataID(4)] = &set.TagSettings{
		PollTime:  35 * time.Millisecond,
		GenConfig: "rand:30ms:1.0:4.0",
	}
	sets.Tags[m.DataID(5)] = &set.TagSettings{
		PollTime:  40 * time.Millisecond,
		GenConfig: "rand:30ms:1.0:3.0",
	}

	storage, _ := storage.New()
	storage.Create(1)
	storage.Create(2)
	storage.Create(3)
	storage.Create(4)
	storage.Create(5)

	sim, _ := set.New(sets, storage)

	sim.Run()

	for i := 6; i < 2000; i++ {
		storage.Create(m.DataID(i))
		sim.TagCreate(m.DataID(i), sets.Tags[m.DataID(1)])
	}
	for i := 2001; i < 4000; i++ {
		storage.Create(m.DataID(i))
		sim.TagCreate(m.DataID(i), sets.Tags[m.DataID(2)])
	}
	for i := 4001; i < 6000; i++ {
		storage.Create(m.DataID(i))
		sim.TagCreate(m.DataID(i), sets.Tags[m.DataID(3)])
	}
	for i := 6001; i < 8000; i++ {
		storage.Create(m.DataID(i))
		sim.TagCreate(m.DataID(i), sets.Tags[m.DataID(4)])
	}
	for i := 8001; i < 10000; i++ {
		storage.Create(m.DataID(i))
		sim.TagCreate(m.DataID(i), sets.Tags[m.DataID(5)])
	}

	time.Sleep(time.Second)
	sim.Reset()

	sim.Run()
	sim.Stop()

	sim.Close()

	for k, v := range storage.List() {
		fmt.Println(k, v.Value)
	}

}
