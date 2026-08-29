package cli

import (
	"fmt"
	"strconv"
)

// stringValue implements lib.CliValue for a string.
type stringValue string

func (v stringValue) String() string  { return string(v) }
func (v stringValue) Int() int        { n, _ := strconv.Atoi(string(v)); return n }
func (v stringValue) Float() float64  { f, _ := strconv.ParseFloat(string(v), 64); return f }
func (v stringValue) Bool() bool      { return string(v) == "true" || string(v) == "1" || string(v) == "yes" }

// intValue implements lib.CliValue for an int.
type intValue int

func (v intValue) String() string  { return fmt.Sprintf("%d", int(v)) }
func (v intValue) Int() int        { return int(v) }
func (v intValue) Float() float64  { return float64(v) }
func (v intValue) Bool() bool      { return int(v) != 0 }

// floatValue implements lib.CliValue for a float64.
type floatValue float64

func (v floatValue) String() string  { return fmt.Sprintf("%g", float64(v)) }
func (v floatValue) Int() int        { return int(v) }
func (v floatValue) Float() float64  { return float64(v) }
func (v floatValue) Bool() bool      { return float64(v) != 0 }

// boolValue implements lib.CliValue for a bool.
type boolValue bool

func (v boolValue) String() string  { if bool(v) { return "true" }; return "false" }
func (v boolValue) Int() int        { if bool(v) { return 1 }; return 0 }
func (v boolValue) Float() float64  { if bool(v) { return 1.0 }; return 0.0 }
func (v boolValue) Bool() bool      { return bool(v) }
