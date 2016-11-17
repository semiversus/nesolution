package evolve

import (
	"log"
  "math"
  "time"
  "math/rand"
	"github.com/semiversus/nesolution/nes"
	"github.com/semiversus/nesolution/replay"
  "fmt"
  "regexp"
  "strconv"
)

const (
  Running = iota
  Timeout
  BadEnd
  GoodEnd
)

const (
  replay_length = 600
)

type IterateResult struct {
  replay *replay.Replay
  score uint64
}

func Run(rom_path string, replay_path string) {
  var prefix_path, postfix_path string
  var number_path int
  var best_score uint64
  var best_replay *replay.Replay
  var complete_replay *replay.Replay

  ch := make(chan IterateResult)

  rand.Seed( time.Now().UnixNano())

  if replay_path!="" {
    re := regexp.MustCompile("([^0-9]+)([0-9]+)([^0-9]+)")
    path_slice:=re.FindStringSubmatch(replay_path)
    prefix_path, postfix_path=path_slice[1],  path_slice[3]
    number_path,_=strconv.Atoi(path_slice[2])
    complete_replay=replay.Load(replay_path)
  } else {
    complete_replay=new(replay.Replay)
    prefix_path="best"
    number_path=0
    postfix_path=".mov"
  }

  if complete_replay.Len()>replay_length {
    console, _ := nes.NewConsole(rom_path)
    console.Load(complete_replay.GetConsoleState())
    for frame:=0; frame<complete_replay.Len()-replay_length; frame++ {
      console.StepFrame()
      console.SetButtons1(complete_replay.ReadButtons(frame))
    }
    best_replay=replay.NewReplay(console)
    copy(best_replay.Controller_data, complete_replay.Controller_data[complete_replay.Len()-replay_length:])
  } else {
    best_replay=complete_replay
  }

  for {
    _, best_score=ScoreReplay(rom_path, best_replay)

    for j:=0; j<1000; j++ {
      for i:=0; i<10; i++ {
        go func(i int) {
          actual_replay, _, actual_score:=Iterate(rom_path, best_replay)
          result:=IterateResult{actual_replay, actual_score}
          ch <- result
        }(i)
      }
      for i:=0; i<10; i++ {
        actual_score_replay:= <-ch

        if actual_score_replay.score>=best_score {
          best_replay=actual_score_replay.replay
          best_score=actual_score_replay.score
        }
      }
    }
    complete_replay.Controller_data=append(complete_replay.Controller_data, best_replay.Controller_data[:replay_length/2]...)
    number_path++
    filename:=prefix_path+strconv.Itoa(number_path)+postfix_path
    complete_replay.Save(filename)
    fmt.Println(filename, best_score)
    console, _ := nes.NewConsole(rom_path)
    console.Load(best_replay.GetConsoleState())
    for frame:=0; frame<replay_length/2; frame++ {
      console.StepFrame()
      console.SetButtons1(best_replay.ReadButtons(frame))
    }
    second_half:=best_replay.Controller_data[replay_length/2:]
    best_replay=replay.NewReplay(console)
    copy(best_replay.Controller_data, second_half)
  }
}

func Iterate(rom_path string, replay_master *replay.Replay) (replay *replay.Replay, state int, score uint64) {
  replay_actual := replay_master.Copy()

  changes := int(10*rand.Float64())
  for i:=0; i<changes; i++ {
    pos := int(replay_length*rand.Float64())
    length := int(math.Pow(400.0, rand.Float64()))
    if pos+length>replay_length {
      length = replay_length-pos
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

  for frame:=0; frame<replay_length; frame++ {
    console.StepFrame()
    console.SetButtons1(replay.ReadButtons(frame))
    state, frame_score=GetFrameScore(console)
    score+=frame_score
    if state!=Running {
      break;
    }
  }
  score+=uint64(console.RAM[0x6D])*256+uint64(console.RAM[0x86]) // x pos
  score+=uint64(console.RAM[0x7dd])*100000+uint64(console.RAM[0x7de])*10000+uint64(console.RAM[0x7df])*1000+uint64(console.RAM[0x7e0])*100+uint64(console.RAM[0x7e1])*10+uint64(console.RAM[0x7e2]) // points
  score+=uint64(console.RAM[0x756])*100 // big mario or firey
  score+=uint64(console.RAM[0x75f])*200000+uint64(console.RAM[0x75c])*50000

  // fmt.Println(uint64(console.RAM[0x6D])*256+uint64(console.RAM[0x86]), uint64(console.RAM[0x7dd])*100000+uint64(console.RAM[0x7de])*10000+uint64(console.RAM[0x7df])*1000+uint64(console.RAM[0x7e0])*100+uint64(console.RAM[0x7e1])*10+uint64(console.RAM[0x7e2]), (1+float32((uint64(console.RAM[0x7f8])*100+uint64(console.RAM[0x7f9])*10+uint64(console.RAM[0x7fa])))/400.0), score, state)
  return state, score
}

func GetFrameScore(console *nes.Console) (state int, score uint64) {
  state=Running

  if console.RAM[0x0E]==0x06 || console.RAM[0x0E]==0x0B || console.RAM[0x7b1]==1 {
    state=BadEnd
  }

  return state, score
}
