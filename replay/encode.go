package replay

import (
  "os"
  "log"
  "image/png"
  "encoding/binary"
  "fmt"

  "github.com/semiversus/nesolution/nes"
)

func Encode(rom_path string, replay_path string) {

  console, err := nes.NewConsole(rom_path)
  if err != nil {
    log.Fatalln(err)
  }

  audio := make(chan float32, 44100)

  audio_file,_:=os.Create("audio.raw")
  defer audio_file.Close()

  console.SetAudioChannel(audio)
  console.SetAudioSampleRate(44100)

  replay := Load(replay_path)
  console.Load(replay.GetConsoleState())

  replay_pos := 0

  frame := console.PPU.Frame
  for replay_pos<=replay.Len() {
    if console.PPU.Frame>frame {
      console.SetButtons1(replay.ReadButtons(replay_pos))
      replay_pos++
      frame++
      file,_:=os.Create(fmt.Sprintf("img%06d.png", replay_pos))
      png.Encode(file, console.Buffer())
      file.Close()
    audio_samples:
      for {
        select {
        case val := <-audio:
          binary.Write(audio_file, binary.LittleEndian, &val)
        default:
          break audio_samples
        }
      }

    }
    console.Step()
  }
}
