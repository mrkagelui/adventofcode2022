package main

import (
	"fmt"
	"strconv"
	"strings"
)

type valve struct {
	flow int
	next []string
}

func parse(ss []string) (world, error) {
	r := make(map[string]valve, len(ss))
	for _, s := range ss {
		name, v, err := parseValve(s)
		if err != nil {
			return nil, err
		}
		r[name] = v
	}
	return r, nil
}

func parseValve(s string) (string, valve, error) {

	// s is expected to be of format "Valve IO has flow rate=20; tunnels lead to valves YT, TX"
	s = s[len("Valve "):]

	valveName, rest, found := strings.Cut(s, " ")
	if !found {
		return "", valve{}, fmt.Errorf("no valve name in [%s]", s)
	}

	s = rest[len("has flow rate="):]

	var ss []string
	ss = strings.Split(s, "; tunnels lead to valves ")
	if len(ss) != 2 {
		ss = strings.Split(s, "; tunnel leads to valve ")
		if len(ss) != 2 {
			return "", valve{}, fmt.Errorf("no next value in [%s]", s)
		}
	}
	flow, err := strconv.Atoi(ss[0])
	if err != nil {
		return "", valve{}, fmt.Errorf("invalid flow in [%s]: %v", s, err)
	}

	next := strings.Split(ss[1], ", ")

	return valveName, valve{
		flow: flow,
		next: next,
	}, nil
}
