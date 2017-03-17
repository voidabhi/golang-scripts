package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/cloud"
	"google.golang.org/cloud/bigtable"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	Project = "fs-staging"
	Zone    = "us-central1-b"
	Cluster = "cluster1"

	table  = "timeout-test"
	cf     = "horses"
	column = "silver" // as in, Hi Ho
)

func openTable(ctx context.Context, addr string) *bigtable.Table {
	var opts []cloud.ClientOption
	if addr != "" {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("failed to connect to in-memory bttest server via grpc: %s", err)
		}
		opts = append(opts, cloud.WithBaseGRPC(conn))
		log.Printf("yeah!")
	}

	client, err := bigtable.NewClient(ctx, Project, Zone, Cluster, opts...)
	if err != nil {
		log.Fatalf("Failed to create bigtable client: %v", err)
	}

	return client.Open(table)
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	addr := flag.String("addr", "", "address to connect to")
	flag.Parse()

	row := fmt.Sprintf("wat-%d", time.Now().Unix())

	log.Printf("using row %q in cf %q in table %q of project %q", row, cf, table, Project)
	time.Sleep(time.Second)

	ctx := context.Background()
	tbl := openTable(ctx, *addr)

	start := time.Now()

	// not necessary:
	// m := bigtable.NewMutation()
	// m.Set(cf, column, bigtable.Now(), []byte("foobar"))
	// if err := tbl.Apply(ctx, row, m); err != nil {
	// 	log.Fatalf("failed to Set: %s", err)
	// }

	readOpt := bigtable.RowFilter(bigtable.LatestNFilter(1))

	v, err := tbl.ReadRow(ctx, row, readOpt)
	if err != nil {
		log.Fatalf("failed to read row: %s", err)
	} else {
		log.Printf("row output: %#v", v)
	}

	var wg sync.WaitGroup

	ctx2, _ := context.WithTimeout(ctx, 5*time.Millisecond)

	/*
	   3000 repros
	   2000 repros
	   1000 repros
	*/
	for i := 0; i < 3000; i++ {
		wg.Add(1)
		go func() {
			if _, err := tbl.ReadRow(ctx2, row, readOpt); err != nil {
				if err != context.DeadlineExceeded && grpc.Code(err) != codes.DeadlineExceeded {
					log.Printf("failed to read row: %v", err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	log.Println("1st pass done")

	log.Println("starting last call on ctx")

	go func() {
		time.Sleep(15 * time.Second)
		log.Println("taking too damn long!  let's try a new connection")
		tbl2 := openTable(ctx, *addr)
		_, err := tbl2.ReadRow(ctx, row, readOpt)
		if err != nil {
			log.Println("connection #2 ReadRow failed: %s", err)
		} else {
			log.Printf("connection #2 worked!")
		}
		os.Exit(1)
	}()

	if _, err = tbl.ReadRow(ctx, row, readOpt); err != nil {
		log.Println("last", err)
	} else {
		log.Println("last ok")
	}

	if time.Since(start) > 15*time.Second {
		log.Printf("runtime: %s", time.Since(start))
		os.Exit(1) // exit with non-0 code so that beast will see this as a failure
	}
}
