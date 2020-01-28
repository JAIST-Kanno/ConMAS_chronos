package main

import (
    "flag"
    "github.com/nats-io/nats.go"
    "log"
    "os"
    "runtime"
    "time"
)

func main() {
    var (
        max_step int
        agents int
        nats_server string
    )
    flag.IntVar(&max_step, "max_step", 2000, "for max step")
    flag.IntVar(&agents, "agents", 30, "init agents number")
    flag.StringVar(&nats_server, "server", "nats:4222", "NATS messaging server")
    flag.Parse()
    step := 0
    readyAgents := 0
    t1 := time.Now()
    opt := []nats.Option{}
    opt = setupConnOptions(opt)
    nc, err := nats.Connect("nats://" + nats_server, opt...)
    n, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
    if err != nil {
        panic(err)
    }

    n.Subscribe("agents.init", func(msg string) {
        readyAgents = readyAgents + 1
        println("readyAgents:", readyAgents)
        if readyAgents >= agents {
            readyAgents = 0
            t1 = time.Now()
            n.Publish("api.next", "Go Next Step!")
            println("START!")
        }
    })
    n.Subscribe("api.next", func(msg string) {
        if max_step > step {
            step = step + 1
            //fmt.Printf("simulation loop: %d\n", step)
        } else {
            t2 := time.Now()
            n.Publish("api.exit", "All done! EXIT!")
            timeDiff := t2.Sub(t1)
            sps := float64(step) / timeDiff.Seconds()
            println("Done ", step, "steps with", agents, "agents")
            println("in", timeDiff.Seconds(), "seconds")
            println("Steps per Seconds:", sps)
            println("Simulation terminating...")
            os.Exit(0)
        }
    })
    nc.Subscribe("agents.report", func(msg *nats.Msg){
        readyAgents = readyAgents + 1
        if readyAgents >= agents {
            readyAgents = 0
            n.Publish("api.move", "Move agents!")
        }
    })
    nc.Subscribe("agents.moved", func(msg *nats.Msg) {
        readyAgents = readyAgents + 1
        if readyAgents >= agents {
            readyAgents = 0
            n.Publish("api.next", "Go Next Step!")
        }
    })
    if err != nil {
        panic(err)
    }
    nc.Flush()
    n.Flush()
    runtime.Goexit()
}

func setupConnOptions(opts []nats.Option) []nats.Option {
    totalWait := 10 * time.Minute
    reconnectDelay := time.Second

    opts = append(opts, nats.ReconnectWait(reconnectDelay))
    opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
    opts = append(opts, nats.DisconnectHandler(func(nc *nats.Conn) {
        log.Printf("Disconnected: will attempt reconnects for %.0fm", totalWait.Minutes())
    }))
    opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
        log.Printf("Reconnected [%s]", nc.ConnectedUrl())
    }))
    opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
        log.Fatalf("Exiting: %v", nc.LastError())
    }))
    return opts
}
