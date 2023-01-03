package dishy

import (
	"context"
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

func (c *Client) do(req *device.Request) (*device.Response, error) {
	ctx := context.Background()
	if c.Timeout > 0 {
		ctx, _ = context.WithTimeout(context.Background(), c.Timeout)
	}
	return c.dc.Handle(ctx, req)
}
