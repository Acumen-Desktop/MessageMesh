package backend

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/raft"
	consensus "github.com/libp2p/go-libp2p-consensus"
	libp2praft "github.com/libp2p/go-libp2p-raft"
)

type raftState struct {
	Now string
}

type raftOP struct {
	Op string
}

func (o *raftOP) ApplyTo(state consensus.State) (consensus.State, error) {
	fmt.Println("Applying OP: ", o.Op)
	return state, nil
}

func StartRaft(network *Network) {
	pids := network.ChatRoom.PeerList()
	// pids := network.P2p.Host.Peerstore().Peers()
	// -- Create the consensus with no actor attached
	raftconsensus := libp2praft.NewOpLog(&raftState{}, &raftOP{})
	// raftconsensus = libp2praft.NewConsensus(&raftState{"i am not consensuated"})
	// --

	// -- Create Raft servers configuration
	pids = append(pids, network.P2p.Host.ID())
	servers := make([]raft.Server, len(pids))
	for i, pid := range pids {
		servers[i] = raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(pid.String()),
			Address:  raft.ServerAddress(pid.String()),
		}
		fmt.Println("Server: ", servers[i])
	}
	serverConfig := raft.Configuration{
		Servers: servers,
	}
	// --

	// -- Create LibP2P transports Raft
	transport, err := libp2praft.NewLibp2pTransport(network.P2p.Host, 2*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	// --

	// -- Configuration
	raftQuiet := false
	config := raft.DefaultConfig()
	if raftQuiet {
		config.LogOutput = io.Discard
		config.Logger = nil
	}
	config.LocalID = raft.ServerID(network.P2p.Host.ID().String())
	// --

	// -- SnapshotStore
	var raftTmpFolder = "db/raft_testing_tmp"
	snapshots, err := raft.NewFileSnapshotStore(raftTmpFolder, 3, nil)
	if err != nil {
		fmt.Println(err)
	}

	// -- Log store and stable store: we use inmem.
	logStore := raft.NewInmemStore()
	// logStore, _ := raftboltdb.NewBoltStore("db/raft.db")
	// --

	// -- Boostrap everything if necessary
	bootstrapped, err := raft.HasExistingState(logStore, logStore, snapshots)
	if err != nil {
		fmt.Println(err)
	}

	if !bootstrapped {
		// Bootstrap cluster first
		raft.BootstrapCluster(config, logStore, logStore, snapshots, transport, serverConfig)
	} else {
		fmt.Println("Already initialized!!")
	}

	raft, err := raft.NewRaft(config, raftconsensus.FSM(), logStore, logStore, snapshots, transport)
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		actor := libp2praft.NewActor(raft)
		raftconsensus.SetActor(actor)

		waitForLeader(raft)

		for {
			if actor.IsLeader() {
				fmt.Println("I am the leader")
				fmt.Println("Raft State: " + raft.State().String())
				updateState(raftconsensus)
				getState(raftconsensus)
			} else {
				fmt.Println("I am not the leader")
				fmt.Println("Raft State: " + raft.State().String())
				getState(raftconsensus)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		for {
			select {
			case <-raftconsensus.Subscribe():
				newState, _ := raftconsensus.GetCurrentState()
				fmt.Println("New state is: ", newState.(*raftState).Now)
			}
		}
	}()

}

func updateState(c *libp2praft.Consensus) {
	loc, _ := time.LoadLocation("UTC")
	newState := &raftState{Now: time.Now().In(loc).String()}

	// CommitState() blocks until the state has been
	// agreed upon by everyone
	agreedState, err := c.CommitState(newState)
	if err != nil {
		fmt.Println(err)
	}
	if agreedState == nil {
		fmt.Println("agreedState is nil: commited on a non-leader?")
	}
}

func getState(c *libp2praft.Consensus) {
	state, err := c.GetCurrentState()
	if err != nil {
		fmt.Println(err)
	}
	if state == nil {
		fmt.Println("state is nil: commited on a non-leader?")
		return
	}
	fmt.Printf("Current state: %d\n", state)
}

func waitForLeader(r *raft.Raft) {
	obsCh := make(chan raft.Observation, 1)
	observer := raft.NewObserver(obsCh, true, nil)
	r.RegisterObserver(observer)
	defer r.DeregisterObserver(observer)

	// New Raft does not allow leader observation directy
	// What's worse, there will be no notification that a new
	// leader was elected because observations are set before
	// setting the Leader and only when the RaftState has changed.
	// Therefore, we need a ticker.

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	ticker := time.NewTicker(time.Second / 2)
	defer ticker.Stop()
	for {
		select {
		case obs := <-obsCh:
			switch obs.Data.(type) {
			case raft.RaftState:
				if r.Leader() != "" {
					return
				}
			}
		case <-ticker.C:
			if r.Leader() != "" {
				return
			}
		case <-ctx.Done():
			fmt.Println("timed out waiting for Leader")
			fmt.Println("Current Raft State: ", r.State())
			fmt.Println("Current Leader: ", r.Leader())
			return
		}
	}
}

// func AppendServer(network *Network, pid string) {
// 	// -- Create Raft servers configuration
// 	server := raft.Server{
// 		Suffrage: raft.Voter,
// 		ID:       raft.ServerID(pid),
// 		Address:  raft.ServerAddress(pid),
// 	}

// 	serverConfig := raft.Configuration{
// 		Servers: []raft.Server{server},
// 	}
// }
