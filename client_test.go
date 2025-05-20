// Copyright 2022-2025 The sacloud/api-client-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	client "github.com/sacloud/api-client-go"
	"github.com/stretchr/testify/require"
)

// TestClient clientパッケージの利用例
//
// 実行には環境変数やプロファイルが必要なためExampleテストではなく通常のテストとしている
func TestClient(t *testing.T) {
	// 環境変数/プロファイルを読み込んでオプションを組み立てる
	opt, _ := client.DefaultOption()
	if opt.AccessToken == "" || opt.AccessTokenSecret == "" {
		t.Skip("required: SAKURACLOUD_ACCESS_TOKEN, SAKURACLOUD_ACCESS_TOKEN_SECRET")
	}

	// オプションからファクトリー生成
	clientFactory := client.NewFactory(opt)

	// ファクトリーからHttpRequestDoerを生成
	doer := clientFactory.NewHttpRequestDoer()

	// doerを用いてHttpリクエスト実施
	url := "https://secure.sakura.ad.jp/cloud/zone/is1a/api/cloud/1.1/zone/is1a"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := doer.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var responseData map[string]interface{}
	json.Unmarshal(data, &responseData) //nolint

	zoneInfo := responseData["Zone"].(map[string]interface{})

	require.EqualValues(t, "is1a", zoneInfo["Name"])
}

func TestNewClient(t *testing.T) {
	c, err := client.NewClient("https://secure.sakura.ad.jp/cloud/zone/is1a/api/cloud/1.1/zone/is1a")
	if err != nil {
		t.Skip(err)
	}
	if c.Options().AccessToken == "" || c.Options().AccessTokenSecret == "" {
		t.Skip("required: SAKURACLOUD_ACCESS_TOKEN, SAKURACLOUD_ACCESS_TOKEN_SECRET")
	}

	// クライアントからHttpRequestDoerを生成
	doer := c.NewHttpRequestDoer()
	req, _ := http.NewRequest(http.MethodGet, c.ServerURL(), nil)
	resp, err := doer.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var responseData map[string]interface{}
	json.Unmarshal(data, &responseData) //nolint

	zoneInfo := responseData["Zone"].(map[string]interface{})
	require.EqualValues(t, "is1a", zoneInfo["Name"])
}

func TestNewClientWithArgs(t *testing.T) {
	c, err := client.NewClient("http://127.0.0.1/",
		client.WithApiKeys("foo", "bar"),
		client.WithUserAgent("TestAgent"),
		client.WithProfile("profileName"), // これは実際には呼ばれないがWithProfileが動いてるのをチェック
		client.WithDisableEnv(true),
		client.WithDisableProfile(true),
		client.WithOptions(&client.Options{HttpRequestTimeout: 100}),
	)
	if err != nil {
		t.Skip(err)
	}
	opts := c.Options()

	require.EqualValues(t, "http://127.0.0.1/", c.APIRootURL)
	require.EqualValues(t, "foo", opts.AccessToken)
	require.EqualValues(t, "bar", opts.AccessTokenSecret)
	require.EqualValues(t, "TestAgent", opts.UserAgent)
	require.EqualValues(t, 100, opts.HttpRequestTimeout)
}
