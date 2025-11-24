package main

import (
	"math/rand"

	"github.com/dfquaresma/hedge/common/distuv"
)

type traffic struct {
	name              string
	interarrival_dist *distuv.Distribution
	servicetime_dist  []*distuv.Distribution
}

func newTraffic(name string, interarrival_dist *distuv.Distribution, servicetime_dist ...*distuv.Distribution) traffic {
	return traffic{
		name:              name,
		interarrival_dist: interarrival_dist,
		servicetime_dist:  servicetime_dist,
	}
}

func (t traffic) NextServiceValue() float64 {
	i := rand.Intn(len(t.servicetime_dist))
	return t.servicetime_dist[i].NextValue()
}

func (t traffic) NextArrivalValue() float64 {
	return t.interarrival_dist.NextValue()
}
