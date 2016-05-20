package ohayou

var (
	eventFuncs = map[string]func(){
		"catEvent":          catEvent,
		"doubleOhayouEvent": doubleOhayouEvent}
)

func startEvents() {
	for e := range eventFuncs {
		startEvent := eventFuncs[e]
		go startEvent()
	}
}
