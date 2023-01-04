package dishy

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"olowe.co/dishy/device"
)

//go:generate ./protoc.sh 127.0.0.1:9200

const (
	DefaultDishyAddr = "192.168.100.1:9200"
	DefaultWifiAddr  = "192.168.1.1:9000"
)

// A Client is a high-level client to communicate with dishy over the network.
// A new Client must be created with Dial.
type Client struct {
	// Timeout specifies a time limit for requests made by the
	// client. A timeout of zero means no timeout.
	Timeout time.Duration
	dc      device.DeviceClient
	conn    *grpc.ClientConn
}

// Dial returns a new Client connected to the dishy device at addr.
// Most callers should specify DefaultDishyAddr.
func Dial(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	return &Client{
		conn: conn,
		dc:   device.NewDeviceClient(conn),
	}, err
}

func (c *Client) Unstow() error {
	req := &device.Request{
		Request: &device.Request_DishStow{
			DishStow: &device.DishStowRequest{
				Unstow: true,
			},
		},
	}
	_, err := c.do(req)
	return err
}

func (c *Client) Stow() error {
	req := &device.Request{
		Request: &device.Request_DishStow{
			DishStow: &device.DishStowRequest{
				Unstow: false,
			},
		},
	}
	_, err := c.do(req)
	return err
}

func (c *Client) Reboot() error {
	req := &device.Request{
		Request: &device.Request_Reboot{
			Reboot: &device.RebootRequest{},
		},
	}
	_, err := c.do(req)
	return err
}

func (c *Client) Status() (*device.DishGetStatusResponse, error) {
	req := &device.Request{
		Request: &device.Request_GetStatus{
			GetStatus: &device.GetStatusRequest{},
		},
	}
	resp, err := c.do(req)
	return resp.GetDishGetStatus(), err
}

func (c *Client) Interfaces() ([]device.NetworkInterface, error) {
	req := &device.Request{
		Request: &device.Request_GetNetworkInterfaces{
			GetNetworkInterfaces: &device.GetNetworkInterfacesRequest{},
		},
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	if resp.GetGetNetworkInterfaces() == nil {
		return nil, fmt.Errorf("no interfaces in response")
	}
	var ifaces []device.NetworkInterface
	for _, iface := range resp.GetGetNetworkInterfaces().NetworkInterfaces {
		if iface == nil {
			continue
		}
		ifaces = append(ifaces, *iface)
	}
	return ifaces, nil
}

func (c *Client) TransceiverTelemetry() (*device.TransceiverGetTelemetryResponse, error) {
	req := &device.Request{
		Request: &device.Request_TransceiverGetTelemetry{
			TransceiverGetTelemetry: &device.TransceiverGetTelemetryRequest{},
		},
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	if resp.GetTransceiverGetTelemetry() == nil {
		return nil, fmt.Errorf("no telemetry in response")
	}
	return resp.GetTransceiverGetTelemetry(), nil
}

func (c *Client) TransceiverStat() (*device.TransceiverGetStatusResponse, error) {
	req := &device.Request{
		Request: &device.Request_TransceiverGetStatus{
			TransceiverGetStatus: &device.TransceiverGetStatusRequest{},
		},
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	if resp.GetTransceiverGetStatus() == nil {
		return nil, fmt.Errorf("no telemetry in response")
	}
	return resp.GetTransceiverGetStatus(), nil
}

func (c *Client) do(req *device.Request) (*device.Response, error) {
	ctx := context.Background()
	if c.Timeout > 0 {
		ctx, _ = context.WithTimeout(context.Background(), c.Timeout)
	}
	return c.dc.Handle(ctx, req)
}
