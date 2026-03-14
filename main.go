package main

import (
	"github.com/ibesuperv/Distributed-File-Storage-System/p2p"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := FileServerOpts{
		ID:                strings.ReplaceAll(listenAddr, ":", "") + "_node",
		EncKey:            []byte("mysecretkey123456789012345678901"), // 32 bytes static testing key
		StorageRoot:       strings.ReplaceAll(listenAddr, ":", "") + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	uploadFlag := flag.String("u", "", "File path to upload")
	downloadFlag := flag.String("d", "", "Filename to download")
	flag.Parse()

	if *uploadFlag == "" && *downloadFlag == "" {
		log.Fatal("Please provide valid flags. Example: \nmake run ARGS=\"-u test_files/image.png\"\nmake run ARGS=\"-d image.png\"")
	}

	s1 := makeServer(":3000", "")
	s2 := makeServer(":7000", "")
	s3 := makeServer(":5000", ":3000", ":7000")

	go func() { log.Fatal(s1.Start()) }()
	time.Sleep(500 * time.Millisecond)
	go func() { log.Fatal(s2.Start()) }()

	time.Sleep(2 * time.Second)

	go s3.Start()
	time.Sleep(2 * time.Second)

	if *uploadFlag != "" {
		uploadFile(s3, *uploadFlag)
	}

	if *downloadFlag != "" {
		downloadFile(s3, *downloadFlag)
	}

	fmt.Println("Operations completed. Terminating...")
}

func uploadFile(s3 *FileServer, filePath string) {
	fmt.Printf("--- UPLOADING FILE: %s ---\n", filePath)

	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("failed to open file %s: %v", filePath, err)
		return
	}
	defer f.Close()

	key := filepath.Base(filePath)

	if err := s3.Store(key, f); err != nil {
		log.Printf("failed to store file %s: %v", key, err)
		return
	}

	// Wait for network peers to finish storing the file locally before returning
	time.Sleep(3 * time.Second)
	fmt.Printf("Successfully uploaded file %s\n", key)
}

func downloadFile(s3 *FileServer, key string) {
	fmt.Printf("--- DOWNLOADING FILE: %s ---\n", key)

	// Explicitly target the network to simulate proper retrieval by deleting local copy first
	if err := s3.store.Delete(s3.ID, key); err != nil {
		log.Printf("failed to delete file locally %s: %v", key, err)
	}

	r, err := s3.Get(key)
	if err != nil {
		log.Printf("failed to get file %s: %v", key, err)
		return
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Printf("failed to read retrieved file %s: %v", key, err)
		return
	}

	// Save the retrieved data to the recovred_files directory
	recoveredPath := filepath.Join("recovred_files", key)
	if err := os.MkdirAll("recovred_files", os.ModePerm); err != nil {
		log.Printf("failed to create recovred_files directory: %v", err)
	} else {
		if err := ioutil.WriteFile(recoveredPath, b, 0644); err != nil {
			log.Printf("failed to write recovered file %s to disk: %v", key, err)
		} else {
			fmt.Printf(">> SUCCESS: Recovered file written permanently to %s\n", recoveredPath)
		}
	}

	fmt.Printf("Successfully retrieved file %s (size: %d bytes)\n\n", key, len(b))
}
