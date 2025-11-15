package distuv

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type Distribution struct {
	name         string
	dist         string
	latency      float64
	tail_latency float64
	b            distuv.Bernoulli
	ln           distuv.LogNormal
	ps           distuv.Poisson
	wb           distuv.Weibull
}

func NewDistribution(name, dist string) *Distribution {
	src := rand.NewSource(uint64(time.Now().Nanosecond()))
	switch dist {
	case "constant":
		return &Distribution{
			name:    name,
			dist:    dist,
			latency: getDistMetric(name, dist, "latency"),
		}
	case "poisson":
		return &Distribution{
			dist: dist,
			ps: distuv.Poisson{
				Lambda: getDistMetric(name, dist, "lambda"),
				Src:    src,
			},
		}
	case "weibull":
		return &Distribution{
			dist: dist,
			wb: distuv.Weibull{
				K:      getDistMetric(name, dist, "k"),
				Lambda: getDistMetric(name, dist, "lambda"),
				Src:    src,
			},
		}
	case "lognormal":
		return &Distribution{
			dist: dist,
			ln: distuv.LogNormal{
				Mu:    getDistMetric(name, dist, "mu"),
				Sigma: getDistMetric(name, dist, "sigma"),
				Src:   src,
			},
		}
	default:
		return nil
	}
}

func getDistMetric(name, dist, metric string) float64 {
	return viper.GetFloat64(fmt.Sprintf("%s.distributions.%s.%s", name, dist, metric))
}

func (d *Distribution) NextValue() float64 {
	switch d.dist {
	case "poisson":
		return d.ps.Rand()
	case "weibull":
		return d.wb.Rand()
	case "lognormal":
		return d.ln.Rand()
	default:
		return d.latency
	}
}

func (d *Distribution) GetPercentile(p float64) float64 {
	switch d.dist {
	case "poisson":
		return d.ps.Quantile(p)
	case "weibull":
		return d.wb.Quantile(p)
	case "lognormal":
		return d.ln.Quantile(p)
	default:
		return d.latency
	}
}
