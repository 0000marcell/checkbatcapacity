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

func WriteToFile(text string) {
  file, err := os.OpenFile("/home/mmc/checkbat.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
  if err != nil {
    fmt.Println("Error opening file:", err)
    return
  }
	
	defer file.Close()

  data := []byte(text + "\n")
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

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
  } else if b.Capacity() < 50 && b.Status() != "Charging" { 
    return "mid"
  } else {
    return "high"
  }
}

func main() {
  s := gocron.NewScheduler(time.UTC)
  bat1 := Battery{
    name: "BAT0",
  }
  bat2 := Battery{
    name: "BAT1",
  }

  lowbatcmd := exec.Command("dunstify", "LOW BATTERY!!!")
  midbatcmd := exec.Command("dunstify", "Bat less then 50%!!!")

  s.Every(5).Seconds().Do(func() {
    bat1Cap := bat1.CheckCapacity()
    bat2Cap := bat2.CheckCapacity()
    WriteToFile("5 checking bat capacity, bat1: "+bat1Cap+" bat2: "+bat2Cap)
    if bat1Cap == "low" &&
    bat2Cap == "low" {
      err := lowbatcmd.Run()
      if err != nil {
        fmt.Println("Error executing dunstify command ", err)
      }  
    }
  })

  s.Every(3600).Seconds().Do(func() {
    bat1Cap := bat1.CheckCapacity()
    bat2Cap := bat2.CheckCapacity()
    WriteToFile("3600 checking bat capacity, bat1: "+bat1Cap+" bat2: "+bat2Cap)
    if bat1Cap == "mid" ||
    bat2Cap == "mid" {
      err := midbatcmd.Run()
      if err != nil {
        fmt.Println("Error executing dunstify command ", err)
      }  
    }
  })
  s.StartBlocking()
}
