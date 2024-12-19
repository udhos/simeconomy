package main

import (
	"fmt"
	"math/rand/v2"
)

func newWorld(cfg config) *world {
	w := &world{
		cfg: cfg,
		day: 0,
	}
	for range w.cfg.InitialPopulation {
		w.population = append(w.population, &person{money: cfg.InitialMoney})
	}

	var i int
	for range w.cfg.InitialFarmers {
		w.population[i].job = personJobFarmer
		w.population[i].farmerPrice = cfg.FarmerInitialPrice
		i++
	}
	for range w.cfg.InitialMerchants {
		w.population[i].job = personJobMerchant

		var sellFactor [goodTypeMax]float32
		for j := range goodTypeMax {
			sellFactor[j] = cfg.MerchantInitialSellFactor
		}
		w.population[i].mercant = merchantShop{
			capacity:   100,
			sellFactor: sellFactor,
		}

		i++
	}

	return w
}

type world struct {
	cfg        config
	day        int
	population []*person
	goods      []*good
}

func (w world) populationSize() (population, farmers, merchants int) {
	for _, p := range w.population {
		switch p.job {
		case personJobFarmer:
			farmers++
		case personJobMerchant:
			merchants++
		}
	}
	population = len(w.population)
	return
}

func (w world) amountOfGoods() int {
	return len(w.goods)
}

func (w world) amountOfFood() int {
	var f int
	for _, g := range w.goods {
		if g.kind == goodTypeFood {
			f++
		}
	}
	return f
}

func (w *world) run() {

	w.rot()

	w.eat()

	w.farmersProduce()

	w.merchantsBuy()

	w.day++
}

func (w *world) merchantsBuy() {
	for _, p := range w.population {
		if p.job != personJobMerchant {
			continue // not a merchant
		}

		old := len(p.mercant.items)
		p.merchantBuy(w)
		fmt.Printf("merchant bought %d items\n", len(p.mercant.items)-old)
	}
}

func (p *person) merchantBuy(w *world) {

	buyMoney := p.money / 2
	if buyMoney < 100 {
		return // no money
	}

	have := len(p.mercant.items)
	want := p.mercant.capacity - have
	if want < 1 {
		return // no room for items
	}

	// find farmers
	var farmers []*person
	for _, f := range w.population {
		if f.job != personJobFarmer {
			continue
		}
		farmers = append(farmers, f)
	}
	// find farmers who have food
	var haveFood []*person
	for _, f := range farmers {
		for _, g := range f.goods {
			if g.kind == goodTypeFood {
				haveFood = append(haveFood, f)
				break
			}
		}
	}
	// shuffle
	rand.Shuffle(len(haveFood), func(i, j int) {
		haveFood[i], haveFood[j] = haveFood[j], haveFood[i]
	})
	// pick three
	if len(haveFood) > 2 {
		haveFood = haveFood[:3]
	}

	// scan sellers
	for _, f := range haveFood {
		for _, g := range f.goods {
			if g.kind != goodTypeFood {
				continue
			}

			if buyMoney < f.farmerPrice {
				return // no money
			}

			// do transaction

			buyMoney -= f.farmerPrice
			p.money -= f.farmerPrice
			f.money += f.farmerPrice

			f.removeGood(g)
			p.goods = append(p.goods, g)

			m := p.mercant
			m.items = append(m.items, &merchantItem{
				item:     g,
				buyPrice: f.farmerPrice,
			})
			p.mercant = m

			want--
		}
	}

}

func (w *world) rot() {
	var rotten int
	for _, p := range w.population {
		for _, g := range p.goods {
			if g.validUntilDay > 0 && g.validUntilDay < w.day {
				p.removeGood(g)
				w.removeGood(g)
				rotten++
			}
		}
	}

	fmt.Printf("rotten goods lost: %d\n", rotten)
}

func (w *world) eat() {
	var ate, hungry int
	for _, p := range w.population {
		a, h := p.eat(w)
		ate += a
		hungry += h
	}

	fmt.Printf("population meals: meals eaten %d, meals missed: %d\n", ate, hungry)
}

func (w *world) removeGood(old *good) {
	removeGood(&w.goods, old)
}

func (w *world) farmersProduce() {
	old := len(w.goods)
	for _, p := range w.population {
		if p.job == personJobFarmer {
			// create fruits per farmer
			for range w.cfg.FarmerProduction {
				g := &good{
					kind:          goodTypeFood,
					validUntilDay: w.day + w.cfg.FruitDuration,
				}
				w.goods = append(w.goods, g)
				p.goods = append(p.goods, g)
			}
		}
	}

	fmt.Printf("farmer production: %d\n", len(w.goods)-old)
}

func (w world) money() int {
	var m int
	for _, p := range w.population {
		m += p.money
	}
	return m
}

type personJob int

const (
	personJobNone     personJob = iota
	personJobFarmer   personJob = iota
	personJobMerchant personJob = iota
)

type person struct {
	job         personJob
	money       int
	goods       []*good
	farmerPrice int          // only for farmers
	mercant     merchantShop // only for merchants
}

type merchantShop struct {
	capacity   int
	items      []*merchantItem
	sellFactor [goodTypeMax]float32
}

type merchantItem struct {
	item     *good
	buyPrice int
}

func removeGood(goods *[]*good, old *good) {
	for i, g := range *goods {
		if g == old {
			(*goods)[i] = (*goods)[len(*goods)-1]
			(*goods) = (*goods)[:len(*goods)-1]
			return
		}
	}
	panic("good not found")
}

func (p *person) removeGood(old *good) {
	removeGood(&p.goods, old)

	if p.job == personJobMerchant {
		m := p.mercant
		for i, item := range m.items {
			if item.item == old {
				m.items[i] = m.items[len(m.items)-1]
				m.items = m.items[:len(m.items)-1]
				break
			}
		}
		p.mercant = m
	}
}

func (p *person) eat(w *world) (ate, hungry int) {
	for range w.cfg.DailyMeals {
		if f, found := p.pickFood(); found {
			w.removeGood(f)
			p.removeGood(f)
			ate++
			continue
		}

		// FIXME TODO XXX eat - try to buy
		panic("FIXME TODO XXX eat - try to buy")

		hungry++
	}
	return
}

func (p *person) pickFood() (*good, bool) {
	for _, g := range p.goods {
		// have food?
		if g.kind == goodTypeFood {
			return g, true
		}
	}
	return nil, false
}

type goodType int

const (
	goodTypeNone goodType = iota
	goodTypeFood goodType = iota
)

const goodTypeMax = goodTypeFood

type good struct {
	kind          goodType
	validUntilDay int
}
