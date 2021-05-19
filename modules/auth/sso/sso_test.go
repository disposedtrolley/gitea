// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sso

import (
	"net/http"
	"testing"

	"code.gitea.io/gitea/modules/setting"
)

func Test_isGitOrLFSPath(t *testing.T) {

	tests := []struct {
		path string

		want bool
	}{
		{
			"/owner/repo/git-upload-pack",
			true,
		},
		{
			"/owner/repo/git-receive-pack",
			true,
		},
		{
			"/owner/repo/info/refs",
			true,
		},
		{
			"/owner/repo/HEAD",
			true,
		},
		{
			"/owner/repo/objects/info/alternates",
			true,
		},
		{
			"/owner/repo/objects/info/http-alternates",
			true,
		},
		{
			"/owner/repo/objects/info/packs",
			true,
		},
		{
			"/owner/repo/objects/info/blahahsdhsdkla",
			true,
		},
		{
			"/owner/repo/objects/01/23456789abcdef0123456789abcdef01234567",
			true,
		},
		{
			"/owner/repo/objects/pack/pack-123456789012345678921234567893124567894.pack",
			true,
		},
		{
			"/owner/repo/objects/pack/pack-0123456789abcdef0123456789abcdef0123456.idx",
			true,
		},
		{
			"/owner/repo/stars",
			false,
		},
		{
			"/notowner",
			false,
		},
		{
			"/owner/repo",
			false,
		},
		{
			"/owner/repo/commit/123456789012345678921234567893124567894",
			false,
		},
	}
	lfsTests := []string{
		"/owner/repo/info/lfs/",
		"/owner/repo/info/lfs/objects/batch",
		"/owner/repo/info/lfs/objects/oid/filename",
		"/owner/repo/info/lfs/objects/oid",
		"/owner/repo/info/lfs/objects",
		"/owner/repo/info/lfs/verify",
		"/owner/repo/info/lfs/locks",
		"/owner/repo/info/lfs/locks/verify",
		"/owner/repo/info/lfs/locks/123/unlock",
	}

	origLFSStartServer := setting.LFS.StartServer

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "http://localhost"+tt.path, nil)
			setting.LFS.StartServer = false
			if got := isGitOrLFSPath(req); got != tt.want {
				t.Errorf("isGitOrLFSPath() = %v, want %v", got, tt.want)
			}
			setting.LFS.StartServer = true
			if got := isGitOrLFSPath(req); got != tt.want {
				t.Errorf("isGitOrLFSPath() = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range lfsTests {
		t.Run(tt, func(t *testing.T) {
			req, _ := http.NewRequest("POST", tt, nil)
			setting.LFS.StartServer = false
			if got := isGitOrLFSPath(req); got != setting.LFS.StartServer {
				t.Errorf("isGitOrLFSPath(%q) = %v, want %v, %v", tt, got, setting.LFS.StartServer, gitPathRe.MatchString(tt))
			}
			setting.LFS.StartServer = true
			if got := isGitOrLFSPath(req); got != setting.LFS.StartServer {
				t.Errorf("isGitOrLFSPath(%q) = %v, want %v", tt, got, setting.LFS.StartServer)
			}
		})
	}
	setting.LFS.StartServer = origLFSStartServer
}