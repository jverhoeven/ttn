package router

import (
	"testing"

	pb_gateway "github.com/TheThingsNetwork/ttn/api/gateway"
	"github.com/TheThingsNetwork/ttn/core/router/gateway"
	"github.com/TheThingsNetwork/ttn/core/types"
	. "github.com/smartystreets/assertions"
)

func TestHandleGatewayStatus(t *testing.T) {
	a := New(t)
	eui := types.GatewayEUI{0, 0, 0, 0, 0, 0, 0, 2}

	router := &router{
		gateways: map[types.GatewayEUI]*gateway.Gateway{},
	}

	// Handle
	statusMessage := &pb_gateway.StatusMessage{Description: "Fake Gateway"}
	err := router.HandleGatewayStatus(eui, statusMessage)
	a.So(err, ShouldBeNil)

	// Check storage
	status, err := router.getGateway(eui).Status.Get()
	a.So(err, ShouldBeNil)
	a.So(status, ShouldNotBeNil)
	a.So(*status, ShouldResemble, *statusMessage)
}
