package connector

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"log"
	"bufio"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


type PluginConnector struct {
	grpcConn *grpc.ClientConn
	cmd      *exec.Cmd
	stdin    io.WriteCloser
}

func NewPluginConnector(ctx context.Context, pluginName string) (*PluginConnector, error) {
	pluginPath := fmt.Sprintf("fd-plugin-auth-%s", pluginName)
	cmd := exec.CommandContext(ctx, pluginPath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error capturing stdout: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Warning: could not capture stderr: %v", err)
	} else {
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				log.Printf("PLUGIN_LOG[%s]: %s", pluginName, scanner.Text())
			}
		}()
	}

	addressChan := make(chan string)
	errChan := make(chan error, 1)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting plugin executable: %w", err)
	}

	go func() {
		reader := bufio.NewReader(stdout)
		address, err := reader.ReadString('\n')
		if err != nil {
			errChan <- fmt.Errorf("failed to read address from plugin stdout: %w", err)
			return
		}
		addressChan <- strings.TrimSpace(address)
	}()

	var address string
	select {
	case addr := <-addressChan:
		address = addr
	case err := <-errChan:
		_ = cmd.Wait()
		return nil, err
	case <-ctx.Done():
		_ = cmd.Wait()
		return nil, fmt.Errorf("context cancelled while waiting for plugin address: %w", ctx.Err())
	}

	log.Printf("Connecting to plugin '%s' at address: %s", pluginName, address)

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	var dialAddress string

	switch runtime.GOOS {
	case "linux", "darwin":
		dialAddress = "unix:" + address
	case "windows":
		dialAddress = fmt.Sprintf("127.0.0.1:%s", address)
	default:
		return nil, fmt.Errorf("unsupported OS for plugin connection: %s", runtime.GOOS)
	}

	grpcConn, err := grpc.NewClient(dialAddress, dialOpts...)
	if err != nil {
		_ = cmd.Wait()
		return nil, fmt.Errorf("failed to connect via gRPC: %w", err)
	}

	return &PluginConnector{
		grpcConn: grpcConn,
		cmd:      cmd,
		stdin:    stdin,
	}, nil
}

func (pc *PluginConnector) GRPCConn() *grpc.ClientConn {
	return pc.grpcConn
}

func (pc *PluginConnector) Close() {
	if pc.grpcConn != nil {
		pc.grpcConn.Close()
	}
	if pc.stdin != nil {
		pc.stdin.Close()
	}
	if pc.cmd != nil {
		pc.cmd.Wait()
	}
}
