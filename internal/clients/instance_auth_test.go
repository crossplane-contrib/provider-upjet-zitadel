// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package clients

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"

	adminpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
)

func jwtProfileTestKey(t *testing.T) (string, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	data, err := json.Marshal(map[string]string{"type": "serviceaccount", "keyId": "test-key", "userId": "test-user", "key": string(keyPEM)})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "key.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return string(data), path
}

func jwtProfileTestCredentials(t *testing.T, endpoint, method, keyJSON, keyPath string) []byte {
	t.Helper()
	host, port, err := net.SplitHostPort(strings.TrimPrefix(endpoint, "http://"))
	if err != nil {
		t.Fatal(err)
	}
	value := keyPath
	if method == "jwt_profile_json" {
		value = keyJSON
	}
	data, err := json.Marshal(map[string]any{"domain": host, "port": port, "insecure": true, method: value})
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestInstanceJWTProfileRespectsReconcileContext(t *testing.T) {
	keyJSON, keyPath := jwtProfileTestKey(t)
	for _, stage := range []string{"discovery", "token"} {
		for _, method := range []string{"jwt_profile_json", "jwt_profile_file", "token"} {
			for _, termination := range []string{"deadline", "cancellation"} {
				t.Run(stage+"/"+method+"/"+termination, func(t *testing.T) {
					entered := make(chan struct{})
					requestCanceled := make(chan struct{})
					release := make(chan struct{})
					var releaseOnce, enteredOnce, canceledOnce sync.Once
					var server *httptest.Server
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if stage == "token" && r.URL.Path == "/.well-known/openid-configuration" {
							w.Header().Set("Content-Type", "application/json")
							_, _ = fmt.Fprintf(w, `{"issuer":%q,"token_endpoint":%q}`, server.URL, server.URL+"/oauth/v2/token")
							return
						}
						path := "/.well-known/openid-configuration"
						if stage == "token" {
							path = "/oauth/v2/token"
						}
						if r.URL.Path != path {
							http.NotFound(w, r)
							return
						}
						// Consume POST bodies so the server can detect client cancellation.
						_, _ = io.Copy(io.Discard, r.Body)
						enteredOnce.Do(func() { close(entered) })
						select {
						case <-r.Context().Done():
							canceledOnce.Do(func() { close(requestCanceled) })
						case <-release:
							http.Error(w, "released stalled endpoint", http.StatusServiceUnavailable)
						}
					}))
					t.Cleanup(server.Close)
					t.Cleanup(func() { releaseOnce.Do(func() { close(release) }) })
					data := jwtProfileTestCredentials(t, server.URL, method, keyJSON, keyPath)
					var ctx context.Context
					var cancel context.CancelFunc
					if termination == "deadline" {
						ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
					} else {
						ctx, cancel = context.WithCancel(context.Background())
					}
					defer cancel()
					done := make(chan error, 1)
					go func() {
						c, err := instanceClientFromCredentials(ctx, data)
						if err == nil {
							_, err = c.Observe(ctx)
							_ = c.Close()
						}
						done <- err
					}()
					select {
					case <-entered:
					case err := <-done:
						t.Fatalf("returned before reaching stalled %s endpoint: %v", stage, err)
					case <-time.After(5 * time.Second):
						t.Fatal("did not reach HTTP endpoint")
					}
					if termination == "cancellation" {
						cancel()
					}
					<-ctx.Done()
					select {
					case err := <-done:
						if err == nil {
							t.Fatal("expected authentication error after context ended")
						}
					case <-time.After(time.Second):
						t.Fatalf("%s request ignored reconcile %s", stage, termination)
					}
					select {
					case <-requestCanceled:
					case <-time.After(time.Second):
						t.Fatal("HTTP request remained active after reconciliation returned")
					}
				})
			}
		}
	}
}

func TestInstanceJWTProfileObservesWithValidCredentials(t *testing.T) {
	keyJSON, keyPath := jwtProfileTestKey(t)
	for _, method := range []string{"jwt_profile_json", "jwt_profile_file", "token"} {
		t.Run(method, func(t *testing.T) {
			grpcServer := grpc.NewServer()
			api := &instanceServer{}
			adminpb.RegisterAdminServiceServer(grpcServer, api)
			t.Cleanup(grpcServer.Stop)
			var server *httptest.Server
			server = httptest.NewServer(h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
					grpcServer.ServeHTTP(w, r)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/.well-known/openid-configuration":
					_, _ = fmt.Fprintf(w, `{"issuer":%q,"token_endpoint":%q}`, server.URL, server.URL+"/oauth/v2/token")
				case "/oauth/v2/token":
					if err := r.ParseForm(); err != nil {
						t.Errorf("cannot parse token request: %v", err)
					}
					if r.Form.Get("grant_type") != "urn:ietf:params:oauth:grant-type:jwt-bearer" || r.Form.Get("assertion") == "" {
						t.Error("missing JWT profile assertion/grant")
					}
					_, _ = io.WriteString(w, `{"access_token":"profile-token","token_type":"Bearer","expires_in":3600}`)
				default:
					http.NotFound(w, r)
				}
			}), &http2.Server{}))
			t.Cleanup(server.Close)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			c, err := instanceClientFromCredentials(ctx, jwtProfileTestCredentials(t, server.URL, method, keyJSON, keyPath))
			if err != nil {
				t.Fatal(err)
			}
			defer c.Close()
			got, err := c.Observe(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if got.ID != "observed-instance" {
				t.Fatalf("wrong observed instance: %+v", got)
			}
			api.mu.Lock()
			defer api.mu.Unlock()
			if len(api.tokens) != 1 || api.tokens[0] != "Bearer profile-token" {
				t.Fatalf("missing exchanged JWT-profile token: %v", api.tokens)
			}
		})
	}
}
