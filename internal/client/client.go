package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/klausborkowski/wordofwisdom/config"
	"github.com/klausborkowski/wordofwisdom/internal/pow"
	"github.com/klausborkowski/wordofwisdom/internal/protocol"
)

func Run(ctx context.Context, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	fmt.Println("connected to", address)
	defer conn.Close()

	for {
		message, err := HandleConnection(ctx, conn, conn)
		if err != nil {
			return err
		}
		fmt.Println("quote result:", message)
		time.Sleep(5 * time.Second)
	}
}

func HandleConnection(ctx context.Context, readerConn io.Reader, writerConn io.Writer) (string, error) {
	reader := bufio.NewReader(readerConn)

	err := sendMsg(protocol.Message{
		Header: protocol.RequestChallenge,
	}, writerConn)
	if err != nil {
		return "", fmt.Errorf("err send request: %w", err)
	}

	msgStr, err := readConnMsg(reader)
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}
	msg, err := protocol.ParseMessage(msgStr)
	if err != nil {
		return "", fmt.Errorf("err parse msg: %w", err)
	}
	var hashcash pow.HashcashData
	err = json.Unmarshal([]byte(msg.Payload), &hashcash)
	if err != nil {
		return "", fmt.Errorf("err parse hashcash: %w", err)
	}
	fmt.Println("got hashcash:", hashcash)

	conf := ctx.Value(config.ConfigCtxKey).(*config.Configuration)
	hashcash, err = hashcash.ComputeHashcash(conf.HashcashConfig.MaxIteration)
	if err != nil {
		return "", fmt.Errorf("err compute hashcash: %w", err)
	}
	fmt.Println("hashcash computed:", hashcash)
	byteData, err := json.Marshal(hashcash)
	if err != nil {
		return "", fmt.Errorf("err marshal hashcash: %w", err)
	}

	err = sendMsg(protocol.Message{
		Header:  protocol.RequestResource,
		Payload: string(byteData),
	}, writerConn)
	if err != nil {
		return "", fmt.Errorf("err send request: %w", err)
	}
	fmt.Println("challenge sent to server")

	msgStr, err = readConnMsg(reader)
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}
	msg, err = protocol.ParseMessage(msgStr)
	if err != nil {
		return "", fmt.Errorf("err parse msg: %w", err)
	}
	return msg.Payload, nil
}

func readConnMsg(reader *bufio.Reader) (string, error) {
	return reader.ReadString('\n')
}

func sendMsg(msg protocol.Message, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := conn.Write([]byte(msgStr))
	return err
}
