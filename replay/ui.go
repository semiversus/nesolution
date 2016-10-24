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

func init() {
	// we need a parallel OS thread to avoid audio stuttering
	runtime.GOMAXPROCS(2)

	// we need to keep OpenGL calls on a single thread
	runtime.LockOSThread()
}

type UI struct {
	console  *nes.Console
	window   *glfw.Window
	texture  uint32
  replay   *Replay
}

func NewUI(rom_path string, replay_path string, mode int) *UI {
	console, err := nes.NewConsole(rom_path)
	if err != nil {
		log.Fatalln(err)
	}

	r := UI{console: console, replay: NewReplay(replay_path)}

  if mode==Playing {
    r.replay.Load(console)
  }

	return &r
}

func (r *UI) Run() {
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
  r.window = window
	if err != nil {
		log.Fatalln(err)
	}
	r.window.MakeContextCurrent()

	r.console.SetAudioChannel(audio.channel)
	r.console.SetAudioSampleRate(audio.sampleRate)
	r.window.SetKeyCallback(r.OnKey)

	// initialize gl
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}
	gl.Enable(gl.TEXTURE_2D)
	gl.ClearColor(0, 0, 0, 1)
	r.texture = createTexture()

  old_timestamp := glfw.GetTime()
	for !r.window.ShouldClose() {
    gl.Clear(gl.COLOR_BUFFER_BIT)

    timestamp := glfw.GetTime()
    r.Update(timestamp-old_timestamp)
    old_timestamp = timestamp

		r.window.SwapBuffers()
		glfw.PollEvents()

	}
}

func (r *UI) Update(dt float64) {
	window := r.window
	console := r.console
	cycles := int(nes.CPUFrequency * dt)
	frame := console.PPU.Frame
	for cycles > 0 {
    if console.PPU.Frame>frame {
      if r.replay.GetState()==Playing {
	      console.SetButtons1(r.replay.ReadButtons())
      } else {
        buttons := updateControllers(window, console)
        if r.replay.GetState()==Recording {
          r.replay.AppendButtons(buttons)
        }
      }
      frame++
    }
		cycles -= console.Step()
	}
	gl.BindTexture(gl.TEXTURE_2D, r.texture)
	setTexture(console.Buffer())
	drawBuffer(r.window)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (r *UI) OnKey(window *glfw.Window,
	key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyR:
			r.console.Reset()
		case glfw.KeySpace:
      switch r.replay.GetState() {
      case Idle:
        fmt.Println("Start")
        r.replay.StartRecord(r.console)
      case Recording:
        fmt.Println("Stop")
        r.replay.Save()
      }
		}
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
