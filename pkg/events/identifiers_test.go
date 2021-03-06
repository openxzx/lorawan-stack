// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package events_test

import (
	"fmt"
	"testing"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestIdentifierFilter(t *testing.T) {
	a := assertions.New(t)

	ch := make(events.Channel, 10)

	filter := events.NewIdentifierFilter()

	evtAppFoo := events.New(test.Context(), "test", "test", events.WithIdentifiers(ttnpb.ApplicationIdentifiers{ApplicationID: "foo"}))
	evtAppBar := events.New(test.Context(), "test", "test", events.WithIdentifiers(ttnpb.ApplicationIdentifiers{ApplicationID: "bar"}))

	evtCliFoo := events.New(test.Context(), "test", "test", events.WithIdentifiers(ttnpb.ClientIdentifiers{ClientID: "foo"}))
	evtCliBar := events.New(test.Context(), "test", "test", events.WithIdentifiers(ttnpb.ClientIdentifiers{ClientID: "bar"}))

	evtDevFoo := events.New(test.Context(), "test", "test", events.WithIdentifiers(ttnpb.EndDeviceIdentifiers{
		DeviceID:               "foo",
		ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "foo"},
	}))
	evtDevBar := events.New(test.Context(), "test", "test", events.WithIdentifiers(ttnpb.EndDeviceIdentifiers{
		DeviceID:               "bar",
		ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "bar"},
	}))

	evtGtwFoo := events.New(test.Context(), "test", "test", events.WithIdentifiers(ttnpb.GatewayIdentifiers{GatewayID: "foo"}))
	evtGtwBar := events.New(test.Context(), "test", "test", events.WithIdentifiers(ttnpb.GatewayIdentifiers{GatewayID: "bar"}))

	evtOrgFoo := events.New(test.Context(), "test", "test", events.WithIdentifiers(ttnpb.OrganizationIdentifiers{OrganizationID: "foo"}))
	evtOrgBar := events.New(test.Context(), "test", "test", events.WithIdentifiers(ttnpb.OrganizationIdentifiers{OrganizationID: "bar"}))

	evtUsrFoo := events.New(test.Context(), "test", "test", events.WithIdentifiers(ttnpb.UserIdentifiers{UserID: "foo"}))
	evtUsrBar := events.New(test.Context(), "test", "test", events.WithIdentifiers(ttnpb.UserIdentifiers{UserID: "bar"}))

	fooIDs := &ttnpb.CombinedIdentifiers{
		EntityIdentifiers: []*ttnpb.EntityIdentifiers{
			ttnpb.ApplicationIdentifiers{ApplicationID: "foo"}.EntityIdentifiers(),
			ttnpb.ClientIdentifiers{ClientID: "foo"}.EntityIdentifiers(),
			ttnpb.EndDeviceIdentifiers{DeviceID: "foo", ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "foo"}}.EntityIdentifiers(),
			ttnpb.GatewayIdentifiers{GatewayID: "foo"}.EntityIdentifiers(),
			ttnpb.OrganizationIdentifiers{OrganizationID: "foo"}.EntityIdentifiers(),
			ttnpb.UserIdentifiers{UserID: "foo"}.EntityIdentifiers(),
		},
	}

	filter.Subscribe(test.Context(), fooIDs, ch)

	filter.Notify(evtAppBar)
	filter.Notify(evtAppFoo)

	a.So(<-ch, should.Equal, evtAppFoo)
	a.So(ch, should.BeEmpty)

	filter.Notify(evtCliBar)
	filter.Notify(evtCliFoo)

	a.So(<-ch, should.Equal, evtCliFoo)
	a.So(ch, should.BeEmpty)

	filter.Notify(evtDevBar)
	filter.Notify(evtDevFoo)

	a.So(<-ch, should.Equal, evtDevFoo)
	a.So(ch, should.BeEmpty)

	filter.Notify(evtGtwBar)
	filter.Notify(evtGtwFoo)

	a.So(<-ch, should.Equal, evtGtwFoo)
	a.So(ch, should.BeEmpty)

	filter.Notify(evtOrgBar)
	filter.Notify(evtOrgFoo)

	a.So(<-ch, should.Equal, evtOrgFoo)
	a.So(ch, should.BeEmpty)

	filter.Notify(evtUsrBar)
	filter.Notify(evtUsrFoo)

	a.So(<-ch, should.Equal, evtUsrFoo)
	a.So(ch, should.BeEmpty)

	filter.Unsubscribe(test.Context(), fooIDs, ch)

	filter.Notify(evtAppFoo)
	filter.Notify(evtCliFoo)
	filter.Notify(evtDevFoo)
	filter.Notify(evtGtwFoo)
	filter.Notify(evtOrgFoo)
	filter.Notify(evtUsrFoo)

	if !a.So(ch, should.BeEmpty) {
	loop:
		for {
			select {
			case evt := <-ch:
				fmt.Println(evt)
			default:
				break loop
			}
		}
	}
}
