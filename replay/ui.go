package replay

import (
	"log"
  "runtime"
  "fmt"

	"github.com/semiversus/nesolution/nes"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/gordonklaus/portaudio"
)

const (
	width  = 256
	height = 240
	scale  = 3
	title  = "NES"
)

const (
  Idle = iota
  Playing
  Recording
)

func init() {
	// we need a parallel OS thread to avoid audio stuttering
	runtime.GOMAXPROCS(2)

	// we need to keep OpenGL calls on a single thread
	runtime.LockOSThread()
}

func Run(rom_path string, replay_path string, replay_mode bool) {
	console, err := nes.NewConsole(rom_path)
	if err != nil {
		log.Fatalln(err)
	}

	portaudio.Initialize()
	defer portaudio.Terminate()

	audio := NewAudio()
	if err := audio.Start(); err != nil {
		log.Fatalln(err)
	}
	defer audio.Stop()

	// initialize glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln(err)
	}
	defer glfw.Terminate()

	// create window
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(width*scale, height*scale, title, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}
	window.MakeContextCurrent()

	console.SetAudioChannel(audio.channel)
	console.SetAudioSampleRate(audio.sampleRate)

  mode := Idle
  replay := NewReplay(console)

  if replay_mode==true {
    mode = Playing
    replay = Load(replay_path)
    console.Load(replay.GetConsoleState())
  }

  onKey := func (window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
    if action != glfw.Press {
      return
    }

    switch key {
    case glfw.KeyR:
      console.Reset()
    case glfw.KeySpace:
      switch mode {
      case Idle:
        fmt.Println("Start")
        replay = NewReplay(console)
        mode = Recording
      case Recording:
        fmt.Println("Stop")
        replay.Save(replay_path)
      }
    }
  }

	window.SetKeyCallback(onKey)

	// initialize gl
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}
	gl.Enable(gl.TEXTURE_2D)
	gl.ClearColor(0, 0, 0, 1)
	texture := createTexture()

  old_timestamp := glfw.GetTime()
  timestamp := old_timestamp
  cycles := int(0)
  frame := uint64(0)

  replay_pos := 0

	for !(window.ShouldClose() || (mode==Playing && replay_pos>replay.Len())) {
    gl.Clear(gl.COLOR_BUFFER_BIT)

    timestamp = glfw.GetTime()
    cycles = int(nes.CPUFrequency * (timestamp-old_timestamp))
    frame = console.PPU.Frame
    for cycles > 0 {
      if console.PPU.Frame>frame {
        switch mode {
        case Playing:
          console.SetButtons1(replay.ReadButtons(replay_pos))
          replay_pos++
        case Recording:
          replay.AppendButtons(updateControllers(window, console))
        case Idle:
          updateControllers(window, console)
        }
        frame++
      }
      cycles -= console.Step()
    }
    gl.BindTexture(gl.TEXTURE_2D, texture)
    setTexture(console.Buffer())
    drawBuffer(window)
    gl.BindTexture(gl.TEXTURE_2D, 0)
    old_timestamp = timestamp

		window.SwapBuffers()
		glfw.PollEvents()
	}
}


func drawBuffer(window *glfw.Window) {
	w, h := window.GetFramebufferSize()
	s1 := float32(w) / 256
	s2 := float32(h) / 240
	f := float32(1)
	var x, y float32
	if s1 >= s2 {
		x = f * s2 / s1
		y = f
	} else {
		x = f
		y = f * s1 / s2
	}
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-x, -y)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()
}

func updateControllers(window *glfw.Window, console *nes.Console) [8]bool {
	turbo := console.PPU.Frame%6 < 3
	k1 := readKeys(window, turbo)
	j1 := readJoystick(glfw.Joystick1, turbo)
	j2 := readJoystick(glfw.Joystick2, turbo)
  buttons := combineButtons(k1, j1)
	console.SetButtons1(buttons)
	console.SetButtons2(j2)
  return buttons
}
