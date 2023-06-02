// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build go1.16
// +build go1.16

package fs

import (
	"io/fs"
	"testing"

	fuzz "github.com/google/gofuzz"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

func TestFromBundle(t *testing.T) {
	type args struct {
		b *bundlev1.Bundle
	}
	tests := []struct {
		name    string
		args    args
		want    BundleFS
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name: "empty",
			args: args{
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{},
				},
			},
			wantErr: false,
		},
		{
			name: "valid",
			args: args{
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "application/test",
						},
						{
							Name: "application/production/test",
						},
						{
							Name: "application/staging/test",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FromBundle(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromBundle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFromBundle_Fuzz(t *testing.T) {
	// Making sure the descrption never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		var (
			src bundlev1.Bundle
		)

		// Prepare arguments
		f.Fuzz(&src)

		// Execute
		FromBundle(&src)
	}
}

func mustFromBundle(b *bundlev1.Bundle) BundleFS {
	fs, err := FromBundle(b)
	if err != nil {
		panic(err)
	}
	return fs
}

var testBundle = &bundlev1.Bundle{
	Packages: []*bundlev1.Package{
		{
			Name: "application/test",
		},
		{
			Name: "application/production/test",
		},
		{
			Name: "application/staging/test",
		},
	},
}

func Test_bundleFs_Open(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fs      BundleFS
		args    args
		want    fs.File
		wantErr bool
	}{
		{
			name: "empty",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "",
			},
		},
		{
			name: "directory",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/production",
			},
		},
		{
			name: "directory not exists",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/whatever",
			},
			wantErr: true,
		},
		{
			name: "file",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/production/test",
			},
		},
		{
			name: "file not exists",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/production/whatever",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.fs.Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("bundleFs.Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_bundleFs_Open_Fuzz(t *testing.T) {
	bfs, _ := FromBundle(testBundle)

	// Making sure the descrption never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		var (
			name string
		)

		// Prepare arguments
		f.Fuzz(&name)

		// Execute
		bfs.Open(name)
	}
}

func Test_bundleFs_ReadDir(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fs      BundleFS
		args    args
		want    []fs.DirEntry
		wantErr bool
	}{
		{
			name: "empty",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "",
			},
		},
		{
			name: "directory",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/production",
			},
		},
		{
			name: "directory not exists",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/whatever",
			},
			wantErr: true,
		},
		{
			name: "file",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/production/test",
			},
			wantErr: true,
		},
		{
			name: "file not exists",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/production/whatever",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.fs.ReadDir(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("bundleFs.ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_bundleFs_ReadDir_Fuzz(t *testing.T) {
	bfs, _ := FromBundle(testBundle)

	// Making sure the descrption never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		var (
			name string
		)

		// Prepare arguments
		f.Fuzz(&name)

		// Execute
		bfs.ReadDir(name)
	}
}

func Test_bundleFs_ReadFile(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fs      BundleFS
		args    args
		want    []fs.DirEntry
		wantErr bool
	}{
		{
			name: "empty",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "",
			},
			wantErr: true,
		},
		{
			name: "directory",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/production",
			},
			wantErr: true,
		},
		{
			name: "directory not exists",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/whatever",
			},
			wantErr: true,
		},
		{
			name: "file",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/production/test",
			},
		},
		{
			name: "file not exists",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/production/whatever",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.fs.ReadFile(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("bundleFs.ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_bundleFs_ReadFile_Fuzz(t *testing.T) {
	bfs, _ := FromBundle(testBundle)

	// Making sure the descrption never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		var (
			name string
		)

		// Prepare arguments
		f.Fuzz(&name)

		// Execute
		bfs.ReadFile(name)
	}
}

func Test_bundleFs_Stat(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fs      BundleFS
		args    args
		want    []fs.DirEntry
		wantErr bool
	}{
		{
			name: "empty",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "",
			},
		},
		{
			name: "directory",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/production",
			},
		},
		{
			name: "directory not exists",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/whatever",
			},
			wantErr: true,
		},
		{
			name: "file",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/production/test",
			},
		},
		{
			name: "file not exists",
			fs:   mustFromBundle(testBundle),
			args: args{
				name: "application/production/whatever",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.fs.Stat(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("bundleFs.Stat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_bundleFs_Stat_Fuzz(t *testing.T) {
	bfs, _ := FromBundle(testBundle)

	// Making sure the descrption never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		var (
			name string
		)

		// Prepare arguments
		f.Fuzz(&name)

		// Execute
		bfs.Stat(name)
	}
}
