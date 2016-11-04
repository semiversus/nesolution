package evolve

import (
	"log"
  "math"
  "time"
  "math/rand"
	"github.com/semiversus/nesolution/nes"
	"github.com/semiversus/nesolution/replay"
  "fmt"
)

const (
  Timeout = iota
  BadEnd
  GoodEnd
)

func Run(rom_path string, replay_path string) {
  rand.Seed( time.Now().UnixNano())
  Iterate(rom_path, replay.Load(replay_path))
}

func Iterate(rom_path string, replay_master *replay.Replay) (state int, score uint64) {
	console, err := nes.NewConsole(rom_path)
	if err != nil {
		log.Fatalln(err)
	}

  console.Load(replay_master.GetConsoleState())
  replay := replay_master.Copy()

  changes := int(math.Pow(10, rand.Float64()))
  for i:=0; i<changes; i++ {
    pos := int(math.Pow(float64(replay.Len()), rand.Float64()))
    length := int(math.Pow(400.0, rand.Float64()))
    if pos+length>replay.Len() {
      length = replay.Len()-pos
    }
    mode := rand.Float64()
    button := rand.Intn(6)
    if button>=2 { // skip start and select button
      button+=2
    }

    switch {
    case mode < 0.8: // add button
      replay.SetButton(pos, length, button)
    case mode < 0.9: // remove slice
      replay.Cut(pos, length)
    default: // remove button
      replay.RemoveButton(pos, length, button)
    }
  }

  for frame:=0; frame<replay.Len(); frame++ {
    console.StepFrame()
    fmt.Println(replay.ReadButtons(frame))
    console.SetButtons1(replay.ReadButtons(frame))
    score+=GetScore(console)
  }
  replay.Save("test.mov")
  return state, score
}

func GetScore(console *nes.Console) (score uint64) {
  return 0
}
