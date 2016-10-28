package evolve

import (
	"log"
  "math"
  "time"
  "math/rand"
	"github.com/semiversus/nesolution/nes"
	"github.com/semiversus/nesolution/replay"
)

const (
  Timeout = iota
  BadEnd
  GoodEnd
)

func Run(rom_path string, replay_path string) {
  rand.Seed( time.Now().UnixNano())
  Iterate(rom_path, replay.NewReplay(replay_path))
}

func Iterate(rom_path string, replay_master *replay.Replay) (state int, score uint64) {
	console, err := nes.NewConsole(rom_path)
	if err != nil {
		log.Fatalln(err)
	}

  replay_master.Load(console)
  replay := replay_master.Copy()
  replay.StartRecord(console)

  for i:=0; i<50; i++ {
    pos := int(math.Pow(float64(replay.Len()), rand.Float64()))
    length := int(math.Pow(400.0, rand.Float64()))
    if pos+length>replay.Len() {
      length = replay.Len()-pos
    }
    mode := rand.Intn(3)
    button := rand.Intn(6)
    if button>=2 { // skip start and select button
      button+=2
    }
    switch mode {
    case 0: // add button
      replay.SetButton(pos, length, button)
    case 1: // remove slice
      replay.Cut(pos, length)
    case 2: // remove button
      replay.RemoveButton(pos, length, button)
    }
  }

  for frame:=0; frame<replay.Len(); frame++ {
    console.StepFrame()
    console.SetButtons1(replay.ReadButtons())
    log.Println(console.PPU.Frame)
  }
  replay.Save()
  return 0, 0
}
