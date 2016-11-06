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

type IterateResult struct {
  replay *replay.Replay
  score uint64
}

func Run(rom_path string, replay_path string) {
  var best_score uint64
  var best_replay *replay.Replay

  ch := make(chan IterateResult)

  rand.Seed( time.Now().UnixNano())
  best_replay=replay.Load(replay_path)

  _, best_score=ScoreReplay(rom_path, best_replay)

  for {
    for i:=0; i<10; i++ {
      go func(i int) {
        fmt.Println("Start ", i)
        actual_replay, _, actual_score:=Iterate(rom_path, best_replay)
        result:=IterateResult{actual_replay, actual_score}
        ch <- result
        fmt.Println("Finish ", i, actual_score)
      }(i)
    }
    for i:=0; i<10; i++ {
      actual_score_replay:= <- ch

      if actual_score_replay.score>=best_score {
        best_replay=actual_score_replay.replay
        best_score=actual_score_replay.score
        best_replay.Save("best.mov")
      }
    }
    fmt.Println(best_score)
  }
}

func Iterate(rom_path string, replay_master *replay.Replay) (replay *replay.Replay, state int, score uint64) {
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

  state, score=ScoreReplay(rom_path, replay_actual)
  return replay_actual, state, score
}

func ScoreReplay(rom_path string, replay *replay.Replay) (state int, score uint64) {
  var frame_score uint64;

	console, err := nes.NewConsole(rom_path)
	if err != nil {
		log.Fatalln(err)
	}

  console.Load(replay.GetConsoleState())

  for frame:=0; frame<replay.Len(); frame++ {
    console.StepFrame()
    console.SetButtons1(replay.ReadButtons(frame))
    state, frame_score=GetFrameScore(console)
    score+=frame_score
    if state!=Running {
      break;
    }
  }
  return state, score
}

func GetFrameScore(console *nes.Console) (state int, score uint64) {
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
