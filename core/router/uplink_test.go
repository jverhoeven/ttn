package router

import (
	"testing"

	"github.com/TheThingsNetwork/ttn/api/discovery"
	pb_gateway "github.com/TheThingsNetwork/ttn/api/gateway"
	pb_semtech "github.com/TheThingsNetwork/ttn/api/gateway/semtech"
	pb_protocol "github.com/TheThingsNetwork/ttn/api/protocol"
	pb_lorawan "github.com/TheThingsNetwork/ttn/api/protocol/lorawan"
	pb "github.com/TheThingsNetwork/ttn/api/router"
	"github.com/TheThingsNetwork/ttn/core/router/gateway"
	"github.com/TheThingsNetwork/ttn/core/types"
	. "github.com/smartystreets/assertions"
)

// newReferenceGateway returns a default gateway
func newReferenceGateway(region string) *gateway.Gateway {
	gtw := gateway.NewGateway(types.GatewayEUI{0, 1, 2, 3, 4, 5, 6, 7})
	gtw.Status.Update(&pb_gateway.Status{
		Region: region,
	})
	return gtw
}

// newReferenceUplink returns a default uplink message
func newReferenceUplink() *pb.UplinkMessage {
	up := &pb.UplinkMessage{
		Payload: make([]byte, 20),
		ProtocolMetadata: &pb_protocol.RxMetadata{Protocol: &pb_protocol.RxMetadata_Lorawan{Lorawan: &pb_lorawan.Metadata{
			CodingRate: "4/5",
			DataRate:   "SF7BW125",
			Modulation: pb_lorawan.Modulation_LORA,
		}}},
		GatewayMetadata: &pb_gateway.RxMetadata{Gateway: &pb_gateway.RxMetadata_Semtech{Semtech: &pb_semtech.RxMetadata{
			Timestamp: 100,
		}},
			Frequency: 868100000,
			Rssi:      -25.0,
			Snr:       5.0,
		},
	}
	return up
}

type mockBrokerDiscovery struct{}

func (d *mockBrokerDiscovery) Discover(devAddr types.DevAddr) ([]*discovery.Announcement, error) {
	return []*discovery.Announcement{}, nil
}

func (d *mockBrokerDiscovery) All() ([]*discovery.Announcement, error) {
	return []*discovery.Announcement{}, nil
}

func TestHandleUplink(t *testing.T) {
	a := New(t)

	r := &router{
		gateways:        map[types.GatewayEUI]*gateway.Gateway{},
		brokerDiscovery: &mockBrokerDiscovery{},
	}

	uplink := newReferenceUplink()
	gtwEUI := types.GatewayEUI{0, 1, 2, 3, 4, 5, 6, 7}

	err := r.HandleUplink(gtwEUI, uplink)
	a.So(err, ShouldBeNil)
	utilization := r.getGateway(gtwEUI).Utilization
	utilization.Tick()
	rx, _ := utilization.Get()
	a.So(rx, ShouldBeGreaterThan, 0)

	// TODO: Integration test that checks broker forward
}
