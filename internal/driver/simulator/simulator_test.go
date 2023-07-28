package driver

import (
	"fmt"
	m "practice/internal/models"
	str "practice/internal/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimulator_NormalWorking(t *testing.T) {
	gensettings := GeneralSettings{
		MaxLiveTime:   time.Minute * 2,
		UseGenManager: false,
	}

	set := m.DriverSettings{
		General: &gensettings,
		Tags:    map[m.DataID]m.Formatter{},
	}
	set.Tags[m.DataID(1)] = &TagSettings{
		PollTime: 25 * time.Millisecond,
		Settings: "rand:25ms:1.0:2.0",
	}
	set.Tags[m.DataID(2)] = &TagSettings{
		PollTime: 25 * time.Millisecond,
		Settings: "rand:25ms:1.0:2.0",
	}

	str, _ := str.New()
	str.Create(1)
	str.Create(2)

	sim, err := New(set, str)
	assert.NoError(t, err)
	assert.NotNil(t, sim)

	err = sim.Run()
	assert.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	settings3 := &TagSettings{
		PollTime: 25 * time.Millisecond,
		Settings: "rand:25ms:1.0:3.0",
	}
	str.Create(3)
	_, err = sim.TagCreate(3, settings3)
	assert.NoError(t, err)

	fmt.Println("-----Running-----")
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

	undo, err := sim.TagDelete(3)
	fmt.Println("----Delete id3----")
	assert.NoError(t, err)
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

	fmt.Println("-----Undo id3-----")
	err = undo()
	assert.NoError(t, err)
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

	err = sim.Stop()
	assert.NoError(t, err)
	fmt.Println("-----Stopped-----")
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

	err = sim.Run()
	assert.NoError(t, err)
	fmt.Println("-----Running-----")
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

	err = sim.Reset()
	assert.NoError(t, err)
	fmt.Println("------Reset------")
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

	err = sim.Run()
	assert.NoError(t, err)
	fmt.Println("-----Running-----")
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

	err = sim.Close()
	assert.NoError(t, err)
	fmt.Println("-----Closing-----")
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}
}

func TestSimulator_NormalWorkingWithManager(t *testing.T) {
	gensettings := GeneralSettings{
		MaxLiveTime:   time.Minute * 2,
		UseGenManager: true,
	}

	set := m.DriverSettings{
		General: &gensettings,
		Tags:    map[m.DataID]m.Formatter{},
	}
	set.Tags[m.DataID(1)] = &TagSettings{
		PollTime: 25 * time.Millisecond,
		Settings: "rand:25ms:1.0:3.0",
	}
	set.Tags[m.DataID(2)] = &TagSettings{
		PollTime: 25 * time.Millisecond,
		Settings: "rand:25ms:1.0:2.0",
	}

	str, _ := str.New()
	str.Create(1)
	str.Create(2)

	sim, err := New(set, str)
	assert.NoError(t, err)
	assert.NotNil(t, sim)

	err = sim.Run()
	assert.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	settings3 := &TagSettings{
		PollTime: 25 * time.Millisecond,
		Settings: "rand:25ms:1.0:3.0",
	}
	str.Create(3)
	_, err = sim.TagCreate(3, settings3)
	assert.NoError(t, err)

	settings4 := &TagSettings{
		PollTime: 30 * time.Millisecond,
		Settings: "rand:25ms:1.0:4.0",
	}
	str.Create(4)
	_, err = sim.TagCreate(4, settings4)
	assert.NoError(t, err)

	fmt.Println("-----Running-----")
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

	undo, err := sim.TagDelete(3)
	fmt.Println("----Delete id3----")
	assert.NoError(t, err)
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

	fmt.Println("-----Undo id3-----")
	err = undo()
	assert.NoError(t, err)
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

	err = sim.Stop()
	assert.NoError(t, err)
	fmt.Println("-----Stopped-----")
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

	err = sim.Run()
	assert.NoError(t, err)
	fmt.Println("-----Running-----")
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

	err = sim.Reset()
	assert.NoError(t, err)
	fmt.Println("------Reset------")
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

	err = sim.Run()
	assert.NoError(t, err)
	fmt.Println("-----Running-----")
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}

	err = sim.Close()
	assert.NoError(t, err)
	fmt.Println("-----Closing-----")
	time.Sleep(50 * time.Millisecond)
	for k, v := range str.List() {
		fmt.Println(k, v.Value)
	}
}

func TestSimulator_ErrorWorking(t *testing.T) {
	gensettings := GeneralSettings{
		MaxLiveTime:   time.Hour * 5,
		UseGenManager: true,
	}

	set := m.DriverSettings{
		General: &gensettings,
		Tags:    map[m.DataID]m.Formatter{},
	}
	set.Tags[m.DataID(1)] = &TagSettings{
		PollTime: 25 * time.Millisecond,
		Settings: "rand:25ms:1.0:3.0",
	}
	set.Tags[m.DataID(2)] = &TagSettings{
		PollTime: 25 * time.Millisecond,
		Settings: "rand:25ms:1.0:2.0",
	}

	str, _ := str.New()
	str.Create(1)
	str.Create(2)

	sim, err := New(set, str)
	assert.ErrorIs(t, err, errLiveTimeLong)

	gensettings.MaxLiveTime = time.Minute * 5
	sim, err = New(set, str)
	assert.NoError(t, err)

	err = sim.Run()
	assert.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	_, err = sim.TagCreate(2, set.Tags[2])
	assert.ErrorIs(t, err, errDataExist)

	_, err = sim.TagDelete(3)
	assert.ErrorIs(t, err, errDataNotFound)

	err = sim.Run()
	assert.ErrorIs(t, err, errAlreadyRunning)

	sim.Stop()
	err = sim.Stop()
	assert.ErrorIs(t, err, errAlreadyStopped)

	sim.Reset()
	err = sim.Stop()
	assert.ErrorIs(t, err, errNotWorking)

	sim.Close()
	err = sim.Close()
	assert.ErrorIs(t, err, errAlreadyClosed)

	err = sim.Run()
	assert.ErrorIs(t, err, errAlreadyClosed)
}
