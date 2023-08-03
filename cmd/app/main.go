package main

import (
	"fmt"
	set "practice/internal/driver/simulator"
	m "practice/internal/models"
	str "practice/internal/storage"
	"time"
)

func main() {
	gensettings := set.GeneralSettings{
		ProgramLiveTime: time.Minute * 2,
		GenOptimization: false,
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

	str, _ := str.New()
	str.Create(1)
	str.Create(2)
	str.Create(3)
	str.Create(4)
	str.Create(5)

	sim, _ := set.New(sets, str)

	sim.Run()

	for i := 6; i < 2000; i++ {
		str.Create(m.DataID(i))
		sim.TagCreate(m.DataID(i), sets.Tags[m.DataID(1)])
	}
	for i := 2001; i < 4000; i++ {
		str.Create(m.DataID(i))
		sim.TagCreate(m.DataID(i), sets.Tags[m.DataID(2)])
	}
	for i := 4001; i < 6000; i++ {
		str.Create(m.DataID(i))
		sim.TagCreate(m.DataID(i), sets.Tags[m.DataID(3)])
	}
	for i := 6001; i < 8000; i++ {
		str.Create(m.DataID(i))
		sim.TagCreate(m.DataID(i), sets.Tags[m.DataID(4)])
	}
	for i := 8001; i < 10000; i++ {
		str.Create(m.DataID(i))
		sim.TagCreate(m.DataID(i), sets.Tags[m.DataID(5)])
	}

	time.Sleep(time.Second)
	sim.Reset()

	sim.Run()
	sim.Stop()

	sim.Close()

	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

}
