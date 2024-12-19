package main

func newWorld(cfg config) *world {
	w := &world{
		cfg: cfg,
		day: 1,
	}
	for range w.cfg.InitialPopulation {
		w.population = append(w.population, &person{money: cfg.InitialMoney})
	}
	return w
}

type world struct {
	cfg        config
	day        int
	population []*person
	goods      []*good
}

func (w world) populationSize() int {
	return len(w.population)
}

func (w world) amountOfGoods() int {
	return len(w.goods)
}

func (w world) money() int {
	var m int
	for _, p := range w.population {
		m += p.money
	}
	return m
}

type person struct {
	money int
}

type good struct {
}
