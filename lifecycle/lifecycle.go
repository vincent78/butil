package lifecycle

import "context"

var AppLifecycleInit = make(chan *AppLifecycleInitModel)
var AppLifecyclePrepared = make(chan context.Context)
var AppLifecycleWorkBegin = make(chan context.Context)
var AppLifecycleActived = make(chan context.Context)
var AppLifecyclePaused = make(chan context.Context)
var AppLifecycleDestroy = make(chan context.Context)
