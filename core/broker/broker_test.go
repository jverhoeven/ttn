package broker

import (
	"sync"
	"testing"

	pb "github.com/TheThingsNetwork/ttn/api/broker"
	pb_handler "github.com/TheThingsNetwork/ttn/api/handler"
	pb_networkserver "github.com/TheThingsNetwork/ttn/api/networkserver"
	. "github.com/smartystreets/assertions"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type mockNetworkServer struct {
	devices []*pb_networkserver.DevicesResponse_Device
}

func (s *mockNetworkServer) GetDevices(ctx context.Context, req *pb_networkserver.DevicesRequest, options ...grpc.CallOption) (*pb_networkserver.DevicesResponse, error) {
	return &pb_networkserver.DevicesResponse{
		Results: s.devices,
	}, nil
}

func (s *mockNetworkServer) PrepareActivation(ctx context.Context, activation *pb.DeduplicatedDeviceActivationRequest, options ...grpc.CallOption) (*pb.DeduplicatedDeviceActivationRequest, error) {
	return activation, nil
}

func (s *mockNetworkServer) Activate(ctx context.Context, activation *pb_handler.DeviceActivationResponse, options ...grpc.CallOption) (*pb_handler.DeviceActivationResponse, error) {
	return activation, nil
}

func (s *mockNetworkServer) Uplink(ctx context.Context, message *pb.DeduplicatedUplinkMessage, options ...grpc.CallOption) (*pb.DeduplicatedUplinkMessage, error) {
	return message, nil
}

func (s *mockNetworkServer) Downlink(ctx context.Context, message *pb.DownlinkMessage, options ...grpc.CallOption) (*pb.DownlinkMessage, error) {
	return message, nil
}

func TestActivateDeactivateRouter(t *testing.T) {
	a := New(t)

	b := &broker{
		routers: make(map[string]chan *pb.DownlinkMessage),
	}

	err := b.DeactivateRouter("RouterID")
	a.So(err, ShouldNotBeNil)

	ch, err := b.ActivateRouter("RouterID")
	a.So(err, ShouldBeNil)
	a.So(ch, ShouldNotBeNil)

	_, err = b.ActivateRouter("RouterID")
	a.So(err, ShouldNotBeNil)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for range ch {
		}
		wg.Done()
	}()

	err = b.DeactivateRouter("RouterID")
	a.So(err, ShouldBeNil)

	wg.Wait()
}

func TestActivateDeactivateHandler(t *testing.T) {
	a := New(t)

	b := &broker{
		handlers: make(map[string]chan *pb.DeduplicatedUplinkMessage),
	}

	err := b.DeactivateHandler("HandlerID")
	a.So(err, ShouldNotBeNil)

	ch, err := b.ActivateHandler("HandlerID")
	a.So(err, ShouldBeNil)
	a.So(ch, ShouldNotBeNil)

	_, err = b.ActivateHandler("HandlerID")
	a.So(err, ShouldNotBeNil)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for range ch {
		}
		wg.Done()
	}()

	err = b.DeactivateHandler("HandlerID")
	a.So(err, ShouldBeNil)

	wg.Wait()
}
