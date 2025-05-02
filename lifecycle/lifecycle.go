package lifecycle

var AppLifecycleInit = make(chan *AppLifecycleInitModel)
var AppLifecyclePrepared = make(chan interface{})
var AppLifecycleWorkBegin = make(chan interface{})
var AppLifecycleActived = make(chan interface{})
var AppLifecyclePaused = make(chan interface{})
var AppLifecycleDestroy = make(chan interface{})
