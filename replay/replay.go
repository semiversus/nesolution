package replay

import (
  "os"
  "encoding/gob"
	"github.com/semiversus/nesolution/nes"
)

const (
  Idle = iota
  Playing
  Recording
)

type Replay struct {
  controller_data []byte
  controller_index int
  filename string
  file *os.File
  state int
  encoder *gob.Encoder
}

func NewReplay(filename string) *Replay {
  replay := Replay{controller_data: make([]byte, 0, 50*60*3), state: Idle, filename: filename}
  return &replay
}

func (r *Replay) GetState() int {
  return r.state
}

func (r *Replay) StartRecord(console *nes.Console) error {
	file, err := os.Create(r.filename)
  r.file = file
	if err != nil {
		return err
	}

  r.state = Recording
  r.encoder = gob.NewEncoder(r.file)
  console.Save(r.encoder)
  return nil
}

func (r *Replay) Save() {
  if r.file==nil {
    return
  }
	r.encoder.Encode(r.controller_data)
  r.state = Idle
  r.file.Close()
}

func (r *Replay) Load(console *nes.Console) error {
	file, err := os.Open(r.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	console.Load(decoder)
	decoder.Decode(&r.controller_data)
  r.controller_index=0
  r.state = Playing
  return nil
}

func (r *Replay) AppendButtons(buttons [8]bool) {
  value := byte(0)
  for i := uint(0) ; i<8; i++ {
    if buttons[i] {
      value+=1<<i
    }
  }
  r.controller_data = append(r.controller_data, value)
}

func (r *Replay) ReadButtons() (buttons [8]bool) {
  if r.controller_index < len(r.controller_data) {
    value := r.controller_data[r.controller_index]
    r.controller_index++
    for i := uint(0); i<8; i++ {
      buttons[i] = (value>>i)&1==1
    }
  }
  return buttons
}

func (r *Replay) PlayFinished() bool {
  return r.state==Playing && r.controller_index>=len(r.controller_data)
}
