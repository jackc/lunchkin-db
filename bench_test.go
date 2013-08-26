package main

import (
	"fmt"
	"strconv"
	"testing"
	"unsafe"
)

func TestSelectAllGamesWithDetails(t *testing.T) {
	games, err := SelectAllGamesWithDetails()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("%%v", games)
}

func BenchmarkMapStringEmptyInterface(b *testing.B) {
	var output string

	for i := 0; i < b.N; i++ {
		m := make(map[string]interface{})
		m["name"] = "Bob Marley"
		m["status"] = "ghost"
		m["disposition"] = "sad"
		m["level"] = int64(9)
		m["effective_level"] = int64(15)
		m["winner"] = true

		output = m["name"].(string)
		output = m["status"].(string)
		output = m["disposition"].(string)
		output = strconv.FormatInt(m["level"].(int64), 10)
		output = strconv.FormatInt(m["effective_level"].(int64), 10)
		output = strconv.FormatBool(m["winner"].(bool))
	}

	fmt.Println(output)
}

func BenchmarkMapStringString(b *testing.B) {
	var output string

	for i := 0; i < b.N; i++ {
		m := make(map[string]string)
		m["name"] = "Bob Marley"
		m["status"] = "ghost"
		m["disposition"] = "sad"
		m["level"] = "9"
		m["effective_level"] = "15"
		m["winner"] = "true"

		output = m["name"]
		output = m["status"]
		output = m["disposition"]
		output = m["level"]
		output = m["effective_level"]
		output = m["winner"]
	}

	fmt.Println(output)
}

func BenchmarkStruct(b *testing.B) {
	type participation struct {
		name            string
		status          string
		disposition     string
		level           int64
		effective_level int64
		winner          bool
	}

	var output string

	for i := 0; i < b.N; i++ {
		var p participation
		p.name = "Bob Marley"
		p.status = "ghost"
		p.disposition = "sad"
		p.level = 9
		p.effective_level = 15
		p.winner = true

		output = p.name
		output = p.status
		output = p.disposition
		output = strconv.FormatInt(p.level, 10)
		output = strconv.FormatInt(p.effective_level, 10)
		output = strconv.FormatBool(p.winner)
	}

	fmt.Println(output)
}

type nullString struct {
	value string
	null  bool
}

type nullInt64 struct {
	value int64
	null  bool
}

type nullBool struct {
	value bool
	null  bool
}

type nullStruct struct {
	name            nullString
	status          nullString
	disposition     nullString
	level           nullInt64
	effective_level nullInt64
	winner          nullBool
}

func BenchmarkNullStruct(b *testing.B) {
	var output string

	for i := 0; i < b.N; i++ {
		var ns nullStruct
		ns.name = nullString{value: "Bob Marley"}
		ns.status = nullString{value: "ghost"}
		ns.disposition = nullString{null: true}
		ns.level = nullInt64{value: 9}
		ns.effective_level = nullInt64{value: 15}
		ns.winner = nullBool{null: true}

		if !ns.name.null {
			output = ns.name.value
		}
		if !ns.status.null {
			output = ns.status.value
		}
		if !ns.disposition.null {
			output = ns.disposition.value
		}
		if !ns.level.null {
			output = strconv.FormatInt(ns.level.value, 10)
		}
		if !ns.effective_level.null {
			output = strconv.FormatInt(ns.effective_level.value, 10)
		}
		if !ns.winner.null {
			output = strconv.FormatBool(ns.winner.value)
		}
	}

	fmt.Println(output)
}

type bitStruct struct {
	nulls           int64
	name            string
	status          string
	disposition     string
	level           int64
	effective_level int64
	winner          bool
}

func (bs *bitStruct) isNameNull() bool {
	return bs.nulls&0x1 != 0
}

func (bs *bitStruct) isStatusNull() bool {
	return bs.nulls&0x2 != 0
}

func (bs *bitStruct) isDispositionNull() bool {
	return bs.nulls&0x4 != 0
}

func (bs *bitStruct) isLevelNull() bool {
	return bs.nulls&0x8 != 0
}

func (bs *bitStruct) isEffectiveLevelNull() bool {
	return bs.nulls&0x10 != 0
}

func (bs *bitStruct) isWinnerNull() bool {
	return bs.nulls&0x20 != 0
}

func BenchmarkBitstringStruct(b *testing.B) {

	var output string

	for i := 0; i < b.N; i++ {
		var bs bitStruct
		bs.name = "Bob Marley"
		bs.status = "ghost"
		bs.nulls = bs.nulls | 0x4 // disposition is null
		bs.level = 9
		bs.effective_level = 15
		bs.nulls = bs.nulls | 0x20 // winner is null

		if !bs.isNameNull() {
			output = bs.name
		}
		if !bs.isStatusNull() {
			output = bs.status
		}
		if !bs.isDispositionNull() {
			output = bs.disposition
		}
		if !bs.isLevelNull() {
			output = strconv.FormatInt(bs.level, 10)
		}
		if !bs.isEffectiveLevelNull() {
			output = strconv.FormatInt(bs.effective_level, 10)
		}
		if !bs.isWinnerNull() {
			output = strconv.FormatBool(bs.winner)
		}
	}

	fmt.Println(output)
}

func TestStructSizes(t *testing.T) {
	var ns nullStruct
	var bs bitStruct

	fmt.Printf("nullStruct: %d\n", unsafe.Sizeof(ns))
	fmt.Printf("bitStruct: %d\n", unsafe.Sizeof(bs))
}
