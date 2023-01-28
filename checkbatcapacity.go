package main

import (
 "fmt"
 "time"
 "github.com/go-co-op/gocron"
 "os"
 "strconv"
 "strings"
 "os/exec"
)

func ReadFile(path string) string {
  data, err := os.ReadFile(path)
  if err != nil {
    fmt.Println("error reading file: ", path)
    return "error!"
  }
  return strings.TrimSpace(string(data))
}

type Battery struct {
  name string
}

func (b *Battery) CapacityPath() string {
  return "/sys/class/power_supply/" + b.name + "/capacity"
}

func (b *Battery) StatusPath() string {
  return "/sys/class/power_supply/" + b.name + "/status"
}

func (b *Battery) Capacity() int {
  res, err := strconv.Atoi(ReadFile(b.CapacityPath()))
  if err != nil {
    fmt.Println("Error during conversion: ", "capacity")
    return -1
  }
  return res
}

func (b *Battery) Status() string {
  return ReadFile(b.StatusPath())
}

func (b *Battery) CheckCapacity() string {
  if b.Capacity() < 30 && b.Status() != "Charging" {
    return "low"
  } else {
    return "high"
  }
}

func main() {
 //runCronJobs()
 s := gocron.NewScheduler(time.UTC)
 bat1 := Battery{
   name: "BAT0",
 }
 bat2 := Battery{
   name: "BAT1",
 }

 cmd := exec.Command("dunstify", "LOW BATTERY!!!")

 s.Every(5).Seconds().Do(func() {
  if bat1.CheckCapacity() == "low" &&
     bat2.CheckCapacity() == "low" {
    err := cmd.Run()
    if err != nil {
      fmt.Println("Error executing dunstify command ", err)
    }  
  }
 })
 s.StartBlocking()
}
