package main

import (
	"fmt"
	"encoding/json"
	"bufio"
	"os"
	"strconv"
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
			return fmt.Errorf("Key %s is used more than once", val)
		}
		(*players)[n]=player
	}
	return nil
}
