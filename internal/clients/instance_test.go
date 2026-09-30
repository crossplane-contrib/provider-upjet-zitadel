// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package clients

import (
	"context"
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	adminpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	instancepb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/instance"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type instanceServer struct {
	adminpb.UnimplementedAdminServiceServer
	mu     sync.Mutex
	tokens []string
	hosts  []string
}

func (s *instanceServer) GetMyInstance(ctx context.Context, _ *adminpb.GetMyInstanceRequest) (*adminpb.GetMyInstanceResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens = append(s.tokens, strings.Join(md.Get("authorization"), ""))
	s.hosts = append(s.hosts, strings.Join(md.Get("x-zitadel-forwarded-host"), ""))
	return &adminpb.GetMyInstanceResponse{Instance: &instancepb.InstanceDetail{Id: "observed-instance", Name: "installed", Version: "4", State: instancepb.State_STATE_RUNNING}}, nil
}

func TestInstanceClientReadsWithCurrentCredentialsAndHeaders(t *testing.T) {
	var lc net.ListenConfig
	listener, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	api := &instanceServer{}
	adminpb.RegisterAdminServiceServer(server, api)
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)
	for _, token := range []string{"first-credential", "rotated-credential"} {
		data, err := json.Marshal(instanceCredentials{Domain: "127.0.0.1", Port: strconv.Itoa(listener.Addr().(*net.TCPAddr).Port), Insecure: true, AccessToken: token, TransportHeaders: map[string]string{"x-zitadel-forwarded-host": "auth.example.test"}})
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		c, err := instanceClientFromCredentials(ctx, data)
		if err != nil {
			cancel()
			t.Fatal(err)
		}
		got, err := c.Observe(ctx)
		_ = c.Close()
		cancel()
		if err != nil {
			t.Fatal(err)
		}
		if got.ID != "observed-instance" || got.Name != "installed" || got.State != "STATE_RUNNING" {
			t.Fatalf("wrong observation: %+v", got)
		}
	}
	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.tokens) != 2 || api.tokens[0] != "Bearer first-credential" || api.tokens[1] != "Bearer rotated-credential" {
		t.Fatal("credentials were cached or not attached")
	}
	for _, host := range api.hosts {
		if host != "auth.example.test" {
			t.Fatal("transport header missing")
		}
	}
}

func TestInstanceClientRequiresTLSUnlessExplicitlyInsecure(t *testing.T) {
	var lc net.ListenConfig
	listener, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	adminpb.RegisterAdminServiceServer(server, &instanceServer{})
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)
	data, err := json.Marshal(instanceCredentials{Domain: "127.0.0.1", Port: strconv.Itoa(listener.Addr().(*net.TCPAddr).Port), AccessToken: "credential"})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	c, err := instanceClientFromCredentials(ctx, data)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if _, err := c.Observe(ctx); err == nil {
		t.Fatal("default client connected to a plaintext server")
	}
}

func TestInstanceCredentialValidationDoesNotExposeSecrets(t *testing.T) {
	for _, data := range []string{
		`{"domain":"local","access_token":"DO_NOT_LOG", "port":{}}`,
		`{"domain":"local"}`,
		`{"access_token":"DO_NOT_LOG"}`,
		`{"domain":"local","access_token":"DO_NOT_LOG","jwt_profile_json":"DO_NOT_LOG"}`,
		`{"domain":"local","jwt_profile_json":"DO_NOT_LOG"}`,
	} {
		c, err := instanceClientFromCredentials(context.Background(), []byte(data))
		if c != nil {
			c.Close()
		}
		if err == nil {
			t.Fatal("expected invalid credentials")
		}
		if strings.Contains(err.Error(), "DO_NOT_LOG") {
			t.Fatal("credential leaked in error")
		}
	}
}
