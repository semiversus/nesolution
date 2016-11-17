package replay

import (
  "os"
  "log"
  "bytes"
  "encoding/gob"
	"github.com/semiversus/nesolution/nes"
)

type Replay struct {
  Controller_data []byte
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
	decoder.Decode(&replay.Controller_data)
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
	encoder.Encode(r.Controller_data)
	encoder.Encode(r.console_state)
  return nil
}

func (r *Replay) Len() int {
  return len(r.Controller_data)
}

func (r *Replay) Copy() *Replay {
  replay := Replay{Controller_data: make([]byte, len(r.Controller_data)), console_state:r.console_state}
  copy(replay.Controller_data, r.Controller_data)
  return &replay
}

func (r *Replay) AppendButtons(buttons [8]bool) {
  value := byte(0)
  for i := uint(0) ; i<8; i++ {
    if buttons[i] {
      value+=1<<i
    }
  }
  r.Controller_data = append(r.Controller_data, value)
}

func (r *Replay) ReadButtons(pos int) (buttons [8]bool) {
  if pos < len(r.Controller_data) {
    value := r.Controller_data[pos]
    for i := uint(0); i<8; i++ {
      buttons[i] = (value>>i)&1==1
    }
  }
  return buttons
}

func (r *Replay) SetButton(pos int, length int, button int) {
  if pos+length>len(r.Controller_data) {
    r.Controller_data=append(r.Controller_data[:], make([]byte, pos+length-len(r.Controller_data))...)
  }

  for i:=pos; i<pos+length; i++ {
    switch button {
    case nes.ButtonLeft:
      r.Controller_data[i]=(r.Controller_data[i]&(^uint8(1<<nes.ButtonRight)))|(1<<nes.ButtonLeft)
    case nes.ButtonRight:
      r.Controller_data[i]=(r.Controller_data[i]&(^uint8(1<<nes.ButtonLeft)))|(1<<nes.ButtonRight)
    case nes.ButtonUp:
      r.Controller_data[i]=(r.Controller_data[i]&(^uint8(1<<nes.ButtonDown)))|(1<<nes.ButtonUp)
    case nes.ButtonDown:
      r.Controller_data[i]=(r.Controller_data[i]&(^uint8(1<<nes.ButtonUp)))|(1<<nes.ButtonDown)
    default:
      r.Controller_data[i]|=1<<uint8(button)
    }
  }
}

func (r *Replay) RemoveButton(pos int, length int, button int) {
  if pos+length>len(r.Controller_data) {
    return
  }
  for i:=pos; i<pos+length; i++ {
    r.Controller_data[i]&=^uint8(1<<uint8(button))
  }
}

func (r *Replay) Cut(pos int, length int) {
  if pos+length>len(r.Controller_data) {
    return
  }
  copy(r.Controller_data[pos:], r.Controller_data[pos+length:])
  for i:=range r.Controller_data[len(r.Controller_data)-length:] {
    r.Controller_data[i]=0
  }
}
