package main

import (
	"math"
	"math/rand/v2"
	"strconv"
	"strings"
)

const tail_factor = 10.0

type hedge struct {
	name         string
	delay        float64
	precision    float64
	cancellation bool
}

func newHedge(name string, p95 float64) hedge {
	switch name {
	case "naive_hedge":
		return hedge{
			name:         "naive_hedge",
			delay:        0,
			cancellation: false,
		}
	case "delayed_hedge_p95wc":
		return hedge{
			name:         "delayed_hedge_p95wc",
			delay:        p95,
			cancellation: true,
		}
	case "perfect_hedge":
		return hedge{
			name:         "perfect_hedge",
			delay:        p95,
			cancellation: true,
		}

	case "assisted_hedge_90wc":
		return hedge{
			name:         "assisted_hedge_90wc",
			precision:    0.9,
			cancellation: true,
		}

	case "assisted_hedge_90nc":
		return hedge{
			name:         "assisted_hedge_90nc",
			precision:    0.9,
			cancellation: false,
		}

	default:
		return hedge{}
	}
}

func (h hedge) hedgedRequest(generated_app, generated_func string, ts, service_time, copy_service_time float64) []string {
	delay := h.delay
	isTailLatency := service_time >= tail_factor*copy_service_time

	// perfect hedge always knows when to send the copy
	if h.name == "perfect_hedge" {
		// delaying beyond service_time to copy only if worth
		delay = service_time + 1.0
		if isTailLatency {
			delay = 0.0 // in case of tail, copy sure is faster, send it right away
		}
	}

	// the assistance predictor precision defines acurrancy for catching tail latencies
	// its unprecision is regarding of how much tail latencies are missed
	// false positives are not modeled, since they only cause extra load but no latency improvement
	if strings.HasPrefix(h.name, "assisted_hedge") {
		// delaying beyond service_time to avoid hedging when not worth
		delay = service_time + 1.0
		if isTailLatency && h.precision >= rand.Float64() {
			delay = 0.0 // in case of tail, copy sure is faster, send it right away
		}
	}

	if h.cancellation {
		return h.hedgedWithCancellation(service_time, copy_service_time, delay, ts, h.name, generated_app, generated_func)
	}
	return h.getHedgedNoCancellation(service_time, copy_service_time, delay, ts, h.name, generated_app, generated_func)
}

func (h hedge) hedgedWithCancellation(service_time, copy_service_time, delay, ts float64, name, generated_app, generated_func string) []string {
	response_time := math.Min(service_time, delay+copy_service_time)
	end_timestamp := ts + response_time
	system_load := response_time
	if response_time > delay {
		delta := response_time - delay
		system_load = delay + 2*delta // add additinal time spent running function after delay up to first finish
	}
	return hedgeOutput(name, generated_app, generated_func, end_timestamp, response_time, system_load, service_time, copy_service_time, delay)
}

func (h hedge) getHedgedNoCancellation(service_time, copy_service_time, delay, ts float64, name, generated_app, generated_func string) []string {
	response_time := math.Min(service_time, delay+copy_service_time)
	end_timestamp := ts + response_time
	system_load := response_time
	if response_time > delay {
		system_load = service_time + copy_service_time // if a second is sent, process both completely
	}
	return hedgeOutput(name, generated_app, generated_func, end_timestamp, response_time, system_load, service_time, copy_service_time, delay)
}

func hedgeOutput(name, app, function string, end_ts, response_time, system_load, service_time, copy_service_time, delay float64) []string {
	return []string{
		name,
		app,
		function,
		strconv.FormatFloat(end_ts, 'f', -1, 64),
		strconv.FormatFloat(response_time, 'f', -1, 64),
		strconv.FormatFloat(system_load, 'f', -1, 64),
		strconv.FormatFloat(service_time, 'f', -1, 64),
		strconv.FormatFloat(copy_service_time, 'f', -1, 64),
		strconv.FormatFloat(delay, 'f', -1, 64),
	}
}
