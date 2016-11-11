package main

import (
	"log"
	"os"
	"github.com/semiversus/nesolution/replay"
	"github.com/semiversus/nesolution/evolve"
)

func main() {
  var replay_path string

  args := os.Args[1:]

  if len(args)<2 {
		log.Fatalln("specify rom file, a command (play, record or evolve) and optional a replay file")
  }

  rom_path := args[0]
  cmd_string := args[1]

  if len(args)>=3 {
    replay_path = args[2]
  }

  switch cmd_string { // check command
  case "play":
    replay.Run(rom_path, replay_path, true)
  case "encode":
    replay.Encode(rom_path, replay_path)
  case "record":
    replay.Run(rom_path, replay_path, false)
  case "evolve":
    evolve.Run(rom_path, replay_path)
  }
}
