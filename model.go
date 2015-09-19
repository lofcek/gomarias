package main

import (
	"fmt"
	"encoding/json"
	"bufio"
	"os"
	"strconv"
	"bytes"
)

type Basic struct {
	Name string
}

type Player struct {
	Name string
	Club string	 `json:",omitempty`
}

type Players map[int] Player

type Tournament struct {
	Basic Basic
	Players Players
}

func load_tournament(filename string) (*Tournament,error) {
	f, err := os.Open("data/" + filename + ".json")
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)
	decoder := json.NewDecoder(reader)
	var t Tournament = Tournament{
		Basic: Basic{Name: filename},
		Players: make(Players),
	}
	if err := decoder.Decode(&t); err != nil {
		return nil, err
	}
	return  &t, nil
}

func (players *Players)UnmarshalJSON(b []byte) error {
	var tmp map[string]Player
	if err := json.Unmarshal(b, &tmp); err!=nil {
		return err
	}
	for str,player := range(tmp) {
		n, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		if val,ok :=(*players)[n]; ok {
			return fmt.Errorf(gettext(`Kľúč "%s" bol použitý viac krát`), val)
		}
		(*players)[n]=player
	}
	return nil
}

func (self Player)empty() bool {
	return self.Name=="" && self.Club==""
}

func (self Players)MarshalJSON() ([]byte, error) {
	if len(self) == 0 {
		return []byte("{}"), nil
	}
	var start rune = '{'
	var b bytes.Buffer
	for i:=0; i<len(self); i++ {
		txt,err := json.Marshal((self)[i])
		if err!=nil {
			return nil,err
		}
		b.WriteString(fmt.Sprintf(`%c"%d":%s`, start, i, txt))
		start=','
	}
	b.WriteRune('}')
	return b.Bytes(), nil
}
