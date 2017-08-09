
///////
// Pool

type ThingPool struct {
    factory *ThingerFactory
    unused chan *Thinger
}

func NewThingPool() *ThingPool {
    tp := &ThingPool{}
    tp.factory = &RealThingFactory{false}
    tp.unused = make( chan Thinger )
    return tp
}

func (tp *ThingPool) Get() *Thinger {
    select {
        case t := <- tp.thing:
            return t;
        default: 
            return tp.CreateThing()
    }
}

func (tp *ThingPool) Put( t *Thinger)  {
    tp.unused <- t
}
