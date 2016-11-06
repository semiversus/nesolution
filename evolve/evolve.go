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
  Running = iota
  Timeout
  BadEnd
  GoodEnd
)

func Run(rom_path string, replay_path string) {
  var best_score, actual_score uint64
  var best_replay, actual_replay *replay.Replay
  var state int

  rand.Seed( time.Now().UnixNano())
  best_replay=replay.Load(replay_path)

  for {
    actual_replay, state, actual_score=Iterate(rom_path, best_replay)
    fmt.Println(best_score, actual_score, state)
    if actual_score>=best_score {
      best_replay=actual_replay
      best_score=actual_score
      best_replay.Save("best.mov")
    }
  }
}

func Iterate(rom_path string, replay_master *replay.Replay) (replay *replay.Replay, state int, score uint64) {
  var frame_score uint64;
	console, err := nes.NewConsole(rom_path)
	if err != nil {
		log.Fatalln(err)
	}

  console.Load(replay_master.GetConsoleState())
  replay_actual := replay_master.Copy()

  changes := int(math.Pow(10, rand.Float64()))
  for i:=0; i<changes; i++ {
    pos := int(math.Pow(float64(replay_actual.Len()), rand.Float64()))
    length := int(math.Pow(400.0, rand.Float64()))
    if pos+length>replay_actual.Len() {
      length = replay_actual.Len()-pos
    }
    mode := rand.Float64()
    button := rand.Intn(6)
    if button>=2 { // skip start and select button
      button+=2
    }

    switch {
    case mode < 0.5: // add button
      replay_actual.SetButton(pos, length, button)
    case mode < 0.8: // remove slice
      replay_actual.Cut(pos, length)
    default: // remove button
      replay_actual.RemoveButton(pos, length, button)
    }
  }

  for frame:=0; frame<replay_actual.Len(); frame++ {
    console.StepFrame()
    console.SetButtons1(replay_actual.ReadButtons(frame))
    state, frame_score=GetScore(console)
    score+=frame_score
    if state!=Running {
      break;
    }
  }
  return replay_actual, state, score
}

func GetScore(console *nes.Console) (state int, score uint64) {
  score=uint64(console.RAM[0x6D])*256+uint64(console.RAM[0x86]) // x pos
  state=Running
  if console.RAM[0x0E]==0x06 || console.RAM[0x0E]==0x0B || console.RAM[0xB5]==255 {
    state=BadEnd
  }
  if console.RAM[0x70f]!=0 && console.RAM[0x70f]!=255 {
    state=GoodEnd
  }
  return state, score
}
