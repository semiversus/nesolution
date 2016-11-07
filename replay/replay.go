package replay

import (
  "os"
  "log"
  "bytes"
  "encoding/gob"
	"github.com/semiversus/nesolution/nes"
)

type Replay struct {
  controller_data []byte
  console_state []byte
}

func NewReplay(console *nes.Console) *Replay {
  replay := Replay{}
  buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
  console.Save(encoder)
  replay.console_state=buffer.Bytes()
  return &replay
}

func Load(filename string) *Replay {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

  replay := Replay{}
	decoder := gob.NewDecoder(file)
	decoder.Decode(&replay.controller_data)
  decoder.Decode(&replay.console_state)

  return &replay
}

func (r *Replay) GetConsoleState() *gob.Decoder {
  buffer := bytes.NewBuffer(r.console_state)
	decoder := gob.NewDecoder(buffer)
  return decoder
}

func (r *Replay) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	encoder.Encode(r.controller_data)
	encoder.Encode(r.console_state)
  return nil
}

func (r *Replay) Len() int {
  return len(r.controller_data)
}

func (r *Replay) Copy() *Replay {
  replay := Replay{controller_data: make([]byte, len(r.controller_data)), console_state:r.console_state}
  copy(replay.controller_data, r.controller_data)
  return &replay
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

func (r *Replay) ReadButtons(pos int) (buttons [8]bool) {
  if pos < len(r.controller_data) {
    value := r.controller_data[pos]
    for i := uint(0); i<8; i++ {
      buttons[i] = (value>>i)&1==1
    }
  }
  return buttons
}

func (r *Replay) SetButton(pos int, length int, button int) {
  if pos+length>len(r.controller_data) {
    r.controller_data=append(r.controller_data[:], make([]byte, pos+length-len(r.controller_data))...)
  }

  for i:=pos; i<pos+length; i++ {
    switch button {
    case nes.ButtonLeft:
      r.controller_data[i]=(r.controller_data[i]&(^uint8(1<<nes.ButtonRight)))|(1<<nes.ButtonLeft)
    case nes.ButtonRight:
      r.controller_data[i]=(r.controller_data[i]&(^uint8(1<<nes.ButtonLeft)))|(1<<nes.ButtonRight)
    case nes.ButtonUp:
      r.controller_data[i]=(r.controller_data[i]&(^uint8(1<<nes.ButtonDown)))|(1<<nes.ButtonUp)
    case nes.ButtonDown:
      r.controller_data[i]=(r.controller_data[i]&(^uint8(1<<nes.ButtonUp)))|(1<<nes.ButtonDown)
    default:
      r.controller_data[i]|=1<<uint8(button)
    }
  }
}

func (r *Replay) RemoveButton(pos int, length int, button int) {
  for i:=pos; i<pos+length; i++ {
    r.controller_data[i]&=^uint8(1<<uint8(button))
  }
}

func (r *Replay) Cut(pos int, length int) {
  if pos>len(r.controller_data) {
    return 
  }
  copy(r.controller_data[pos:], r.controller_data[pos+length:])
  for i:=range r.controller_data[len(r.controller_data)-length:] {
    r.controller_data[i]=0
  }
}
