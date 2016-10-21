### Summary

This project extends the NES emulator from Michael Fogleman to let computers play NES. This repository is forked from http://github.com/fogleman/nes .

Evolutionary algorithm iterates through following steps:
- Use a given button pattern and make random changes (add/remove button presses, cut pieces out of the pattern
- Let the pattern run in the emulator
- Rate the success of the button pattern
- If the actual is better than the original use this pattern for the next iteration

### Dependencies

    github.com/go-gl/gl/v2.1/gl
    github.com/go-gl/glfw/v3.1/glfw
    github.com/gordonklaus/portaudio

The portaudio-go dependency requires PortAudio on your system:

> To build portaudio-go, you must first have the PortAudio development headers
> and libraries installed. Some systems provide a package for this; e.g., on
> Ubuntu you would want to run apt-get install portaudio19-dev. On other systems
> you might have to install from source.

On Mac, you can use homebrew:

    brew install portaudio

### Installation

The `go get` command will automatically fetch the dependencies listed above,
compile the binary and place it in your `$GOPATH/bin` directory.

    go get github.com/semiversus/nesolution

### Usage

    nesolution rom_file cmd replay_file
    
Possible commands(`cmd`) are:
- `record`: play the game and use spacebar to start and end recording to the given `replay_file`
- `play`: replay the button pattern stored in `replay_file`
- `evolve`: optimize the given button pattern

### Controls

Joysticks are supported, although the button mapping is currently hard-coded.
Keyboard controls are indicated below.

| Nintendo              | Emulator    |
| --------------------- | ----------- |
| Up, Down, Left, Right | Arrow Keys  |
| Start                 | Enter       |
| Select                | Right Shift |
| A                     | Z           |
| B                     | X           |
| A (Turbo)             | A           |
| B (Turbo)             | S           |
| Reset                 | R           |
| Start/Stop Recording  | Space       |

### Mappers

The following mappers have been implemented:

* NROM (0)
* MMC1 (1)
* UNROM (2)
* CNROM (3)
* MMC3 (4)
* AOROM (7)

These mappers cover about 85% of all NES games. I hope to implement more
mappers soon. To see what games should work, consult this list:

[NES Mapper List](http://tuxnes.sourceforge.net/nesmapper.txt)

### Known Issues

* there are some minor issues with PPU timing, but most games work OK anyway
* the APU emulation isn't quite perfect, but not far off

### Documentation

Interested in writing your own emulator? Curious about the NES internals? Here
are some good resources:

* [NES Documentation (PDF)](http://nesdev.com/NESDoc.pdf)
* [NES Reference Guide (Wiki)](http://wiki.nesdev.com/w/index.php/NES_reference_guide)
* [6502 CPU Reference](http://www.obelisk.demon.co.uk/6502/)
