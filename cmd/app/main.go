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
		UseGenManager: true,
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

	sim, err := set.New(sets, str)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = sim.Run(); err != nil {
		fmt.Println(err)
		return
	}
	if err = sim.Close(); err != nil {
		fmt.Println(err)
		return
	}
}
