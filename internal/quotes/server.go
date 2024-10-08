package quotes

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"

	"github.com/klausborkowski/wordofwisdom/config"
	"github.com/klausborkowski/wordofwisdom/internal/cache"
	"github.com/klausborkowski/wordofwisdom/internal/clock"
	"github.com/klausborkowski/wordofwisdom/internal/pow"
	"github.com/klausborkowski/wordofwisdom/internal/protocol"
)

var ErrQuit = errors.New("client requests to close connection")

func init() {
	reader := bytes.NewReader(quotes)
	s := bufio.NewScanner(reader)

	for s.Scan() {
		if q := strings.TrimSpace(s.Text()); q != "" {
			Quotes = append(Quotes, q)
		}
	}
}

func RunServer(ctx context.Context, address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Println("listening", listener.Addr())
	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accept connection: %w", err)
		}

		go handleConnection(ctx, conn)
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	fmt.Println("new client:", conn.RemoteAddr())
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("err read connection:", err)
			return
		}
		msg, err := ProcessRequest(ctx, req, conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("err process request:", err)
			return
		}
		if msg != nil {
			err := sendMsg(*msg, conn)
			if err != nil {
				fmt.Println("err send message:", err)
			}
		}
	}
}

func ProcessRequest(ctx context.Context, msgStr string, clientInfo string) (*protocol.Message, error) {
	msg, err := protocol.ParseMessage(msgStr)
	if err != nil {
		return nil, err
	}

	switch msg.Header {
	case protocol.Quit:
		return nil, ErrQuit
	case protocol.RequestChallenge:
		fmt.Printf("client %s requests challenge\n", clientInfo)
		conf := ctx.Value("config").(*config.Configuration)
		clock := ctx.Value("clock").(clock.Clock)
		cache := ctx.Value("cache").(cache.Cache)
		date := clock.Now()

		randValue := rand.Intn(100000)
		err := cache.Add(randValue, int64(conf.HashcashConfig.Duration))
		if err != nil {
			return nil, fmt.Errorf("err add rand to cache: %w", err)
		}

		hashcash := pow.HashcashData{
			Version:    1,
			ZerosCount: conf.HashcashConfig.ZerosCount,
			Date:       date.Unix(),
			Resource:   clientInfo,
			Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", randValue))),
			Counter:    0,
		}
		hashcashMarshaled, err := json.Marshal(hashcash)
		if err != nil {
			return nil, fmt.Errorf("err marshal hashcash: %v", err)
		}
		msg := protocol.Message{
			Header:  protocol.ResponseChallenge,
			Payload: string(hashcashMarshaled),
		}
		return &msg, nil
	case protocol.RequestResource:
		fmt.Printf("client %s requests resource with payload %s\n", clientInfo, msg.Payload)
		var hashcash pow.HashcashData
		err := json.Unmarshal([]byte(msg.Payload), &hashcash)
		if err != nil {
			return nil, fmt.Errorf("err unmarshal hashcash: %w", err)
		}

		if hashcash.Resource != clientInfo {
			return nil, fmt.Errorf("invalid hashcash resource")
		}
		conf := ctx.Value("config").(*config.Configuration)
		clock := ctx.Value("clock").(clock.Clock)
		cache := ctx.Value("cache").(cache.Cache)

		randValueBytes, err := base64.StdEncoding.DecodeString(hashcash.Rand)
		if err != nil {
			return nil, fmt.Errorf("err decode rand: %w", err)
		}
		randValue, err := strconv.Atoi(string(randValueBytes))
		if err != nil {
			return nil, fmt.Errorf("err decode rand: %w", err)
		}

		exists, err := cache.Get(randValue)
		if err != nil {
			return nil, fmt.Errorf("err get rand from cache: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("challenge expired or not sent")
		}

		if clock.Now().Unix()-hashcash.Date > int64(conf.HashcashConfig.Duration) {
			return nil, fmt.Errorf("challenge expired")
		}

		maxIter := hashcash.Counter
		if maxIter == 0 {
			maxIter = 1
		}
		_, err = hashcash.ComputeHashcash(maxIter)
		if err != nil {
			return nil, fmt.Errorf("invalid hashcash")
		}

		fmt.Printf("client %s succesfully computed hashcash %s\n", clientInfo, msg.Payload)
		msg := protocol.Message{
			Header:  protocol.ResponseResource,
			Payload: Quotes[rand.Intn(4)],
		}

		cache.Delete(randValue)
		return &msg, nil
	default:
		return nil, fmt.Errorf("unknown header")
	}
}

func sendMsg(msg protocol.Message, conn net.Conn) error {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := conn.Write([]byte(msgStr))
	return err
}
