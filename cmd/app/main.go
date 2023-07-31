package main

import (
	"fmt"
	set "practice/internal/driver/simulator"
	m "practice/internal/models"
	str "practice/internal/storage"
	"time"
)

// Просто для реквеста
// Простой мейн, примеры запуска на разных стартовых условиях в simulator_test.go
func main() {
	gensettings := set.GeneralSettings{
		MaxLiveTime:   time.Minute * 2,
		UseGenManager: false,
	}

	sets := m.DriverSettings{
		General: &gensettings,
		Tags:    map[m.DataID]m.Formatter{},
	}
	sets.Tags[m.DataID(1)] = &set.TagSettings{
		PollTime: 25 * time.Millisecond,
		Settings: "rand:25ms:1.0:3.0",
	}
	sets.Tags[m.DataID(2)] = &set.TagSettings{
		PollTime: 25 * time.Millisecond,
		Settings: "rand:25ms:1.0:2.0",
	}

	str, _ := str.New()
	str.Create(1)
	str.Create(2)

	sim, _ := set.New(sets, str)

	sim.Run()

	for i := 3; i < 5000; i++ {
		str.Create(m.DataID(i))
		sim.TagCreate(m.DataID(i), sets.Tags[m.DataID(1)])
	}

	for i := 5000; i < 10000; i++ {
		str.Create(m.DataID(i))
		sim.TagCreate(m.DataID(i), sets.Tags[m.DataID(2)])
	}

	time.Sleep(time.Second)
	sim.Reset()
	time.Sleep(time.Second * 2)
	sim.Run()
	//time.Sleep(time.Second)
	sim.Close()

	//time.Sleep(time.Second)

	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

}
