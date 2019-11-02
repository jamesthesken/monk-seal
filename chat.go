/*

Functions which drive the bootstrapping and rendezvous for the peers, see for more info:
https://github.com/libp2p/go-libp2p-examples/tree/master/chat-with-rendezvous

*/

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/rivo/tview"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-discovery"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	multiaddr "github.com/multiformats/go-multiaddr"

	"github.com/gdamore/tcell"

)


func handleStream(stream network.Stream) {
	//fmt.Println("Got a new stream!\n")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	// Formatting for the rendezvous printouts
	format := tview.NewTextView().
		SetTextColor(tcell.ColorYellow).
		SetScrollable(true).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetChangedFunc(func() {
			app.Draw()
		})

	go readData(rw, format)
	go writeData(rw, format)

	// TODO: Is there a way to change this?
	// 'stream' will stay open until you close it (or the other side closes it).
}

// Incoming data
func readData(rw *bufio.ReadWriter, format *tview.TextView) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}
// Outgoing data
func writeData(rw *bufio.ReadWriter, format *tview.TextView) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}


// Add a second argument which points to the "Message field" to contain the read/write messages
// Function that bootstraps the network and connects to the desired peer
func rendezvousChat(format *tview.TextView, format2 *tview.TextView) {

	config, err := ParseFlags()
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(format, "Starting MonkSeal chat \n")


	ctx := context.Background()

	// libp2p.New constructs a new libp2p Host. Other options can be added
	// here.
	host, err := libp2p.New(ctx,
		libp2p.ListenAddrs([]multiaddr.Multiaddr(config.ListenAddresses)...),
	)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(format, "Host created. We are:", host.ID() + "\n")
	//fmt.Fprintf(format, host.Addrs())

	// Set a function as stream handler. This function is called when a peer
	// initiates a connection and starts a stream with this peer.
	host.SetStreamHandler(protocol.ID(config.ProtocolID), handleStream)

	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	fmt.Fprintf(format, "Bootstrapping the DHT\n")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}

	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup
	for _, peerAddr := range config.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				//logger.Warning(err)
				fmt.Fprintf(format, "Warning: ", err, "\n")
			} else {
				fmt.Fprintf(format, "Connection established with bootstrap node:", *peerinfo, "\n")
			}
		}()
	}
	wg.Wait()

	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	fmt.Fprintf(format, "Announcing ourselves...\n")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(ctx, routingDiscovery, config.RendezvousString)
	fmt.Fprintf(format, "Successfully announced!\n")

	// Now, look for others who have announced
	// This is like your friend telling you the location to meet you.
	fmt.Fprintf(format, "Searching for other peers...\n")
	peerChan, err := routingDiscovery.FindPeers(ctx, config.RendezvousString)
	if err != nil {
		panic(err)
	}

	for peer := range peerChan {
		if peer.ID == host.ID() {
			continue
		}
		fmt.Fprintf(format, "Found peer:", peer, "\n")

		fmt.Fprintf(format, "Connecting to:", peer, "\n")
		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(config.ProtocolID))

		if err != nil {
			fmt.Fprintf(format, "Connection failed:", err, "\n")
			continue
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go writeData(rw, format2)
			go readData(rw, format2)
		}

		fmt.Fprintf(format, "Connected to:", peer, "\n")
	}

	select {}

}

