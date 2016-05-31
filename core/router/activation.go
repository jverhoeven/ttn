package router

import (
	"errors"
	"sync"

	pb_broker "github.com/TheThingsNetwork/ttn/api/broker"
	pb_protocol "github.com/TheThingsNetwork/ttn/api/protocol"
	pb_lorawan "github.com/TheThingsNetwork/ttn/api/protocol/lorawan"
	pb "github.com/TheThingsNetwork/ttn/api/router"
	"github.com/TheThingsNetwork/ttn/core/types"
)

func (r *router) HandleActivation(gatewayEUI types.GatewayEUI, activation *pb.DeviceActivationRequest) (*pb.DeviceActivationResponse, error) {

	gateway := r.getGateway(gatewayEUI)

	uplink := &pb.UplinkMessage{
		Payload:          activation.Payload,
		ProtocolMetadata: activation.ProtocolMetadata,
		GatewayMetadata:  activation.GatewayMetadata,
	}

	// Only for LoRaWAN
	gateway.Utilization.AddRx(uplink)

	downlinkOptions := r.buildDownlinkOptions(uplink, true, gateway)

	// Find Broker
	brokers, err := r.brokerDiscovery.All()
	if err != nil {
		return nil, err
	}

	// Prepare request
	request := &pb_broker.DeviceActivationRequest{
		Payload:          activation.Payload,
		DevEui:           activation.DevEui,
		AppEui:           activation.AppEui,
		ProtocolMetadata: activation.ProtocolMetadata,
		GatewayMetadata:  activation.GatewayMetadata,
		ActivationMetadata: &pb_protocol.ActivationMetadata{
			Protocol: &pb_protocol.ActivationMetadata_Lorawan{
				Lorawan: &pb_lorawan.ActivationMetadata{
					AppEui: activation.AppEui,
					DevEui: activation.DevEui,
				},
			},
		},
		DownlinkOptions: downlinkOptions,
	}

	// Prepare LoRaWAN activation
	status, err := gateway.Status.Get()
	if err != nil {
		return nil, err
	}
	band, err := getBand(status.Region)
	if err != nil {
		return nil, err
	}
	lorawan := request.ActivationMetadata.GetLorawan()
	lorawan.Rx1DrOffset = 0
	lorawan.Rx2Dr = uint32(band.RX2DataRate)
	lorawan.RxDelay = uint32(band.ReceiveDelay1.Seconds())
	switch status.Region {
	case "EU_863_870":
		lorawan.CfList = []uint64{867100000, 867300000, 867500000, 867700000, 867900000}
	}

	// Forward to all brokers and collect responses
	var wg sync.WaitGroup
	responses := make(chan *pb_broker.DeviceActivationResponse, len(brokers))
	for _, broker := range brokers {
		broker, err := r.getBroker(broker)
		if err != nil {
			continue
		}

		// Do async request
		wg.Add(1)
		go func() {
			res, err := broker.client.Activate(r.Component.GetContext(), request)
			if err == nil && res != nil {
				responses <- res
			}
			wg.Done()
		}()
	}

	// Make sure to close channel when all requests are done
	go func() {
		wg.Wait()
		close(responses)
	}()

	var gotFirst bool
	for res := range responses {
		if gotFirst {
			// warn for duplicate responses
		} else {
			gotFirst = true
			downlink := &pb_broker.DownlinkMessage{
				Payload:        res.Payload,
				DownlinkOption: res.DownlinkOption,
			}
			err := r.HandleDownlink(downlink)
			if err != nil {
				gotFirst = false // try again
			}
		}
	}

	// Activation not accepted by any broker
	if !gotFirst {
		return nil, errors.New("ttn/router: Activation not accepted at this Gateway")
	}

	// Activation accepted by (at least one) broker
	return &pb.DeviceActivationResponse{}, nil
}
