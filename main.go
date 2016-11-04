package main

import (
	"log"
	"os"
	"github.com/semiversus/nesolution/replay"
	"github.com/semiversus/nesolution/evolve"
)

func main() {
	log.SetFlags(0)

  args := os.Args[1:]

  if len(args)<3 {
		log.Fatalln("specify rom file, a command (play, record or evolve) and a replay file")
  }

  rom_path := args[0]
  cmd_string := args[1]
  replay_path := args[2]

  switch cmd_string { // check command
  case "play":
    replay.Run(rom_path, replay_path, true)
  case "record":
    replay.Run(rom_path, replay_path, false)
  case "evolve":
    evolve.Run(rom_path, replay_path)
  }
}
