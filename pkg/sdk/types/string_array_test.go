// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package types

import (
	"testing"
)

func TestStringArray_Contains(t *testing.T) {
	type args struct {
		item string
	}
	tests := []struct {
		name string
		s    StringArray
		args args
		want bool
	}{
		{
			name: "empty",
			s:    StringArray{},
			args: args{
				item: "",
			},
			want: false,
		},
		{
			name: "not same case",
			s:    StringArray{"fOo"},
			args: args{
				item: "foo",
			},
			want: true,
		},
		{
			name: "same case",
			s:    StringArray{"foo"},
			args: args{
				item: "foo",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Contains(tt.args.item); got != tt.want {
				t.Errorf("StringArray.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringArray_AddIfNotContains(t *testing.T) {
	type args struct {
		item string
	}
	tests := []struct {
		name string
		s    *StringArray
		args args
		want bool
	}{
		{
			name: "empty",
			s:    &StringArray{},
			args: args{
				item: "1",
			},
			want: true,
		},
		{
			name: "already contains",
			s:    &StringArray{"1"},
			args: args{
				item: "1",
			},
			want: false,
		},
		{
			name: "already contains with different case",
			s:    &StringArray{"fOo"},
			args: args{
				item: "foo",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.AddIfNotContains(tt.args.item); got != tt.want {
				t.Errorf("StringArray.AddIfNotContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringArray_Remove(t *testing.T) {
	type args struct {
		item string
	}
	tests := []struct {
		name string
		s    *StringArray
		args args
		want bool
	}{
		{
			name: "empty",
			s:    &StringArray{},
			args: args{
				item: "1",
			},
			want: false,
		},
		{
			name: "contains",
			s:    &StringArray{"1"},
			args: args{
				item: "1",
			},
			want: true,
		},
		{
			name: "contains with different case",
			s:    &StringArray{"fOo"},
			args: args{
				item: "foo",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Remove(tt.args.item); got != tt.want {
				t.Errorf("StringArray.Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringArray_HasOneOf(t *testing.T) {
	type args struct {
		items []string
	}
	tests := []struct {
		name string
		s    *StringArray
		args args
		want bool
	}{
		{
			name: "empty",
			s:    &StringArray{},
			args: args{
				items: []string{"1"},
			},
			want: false,
		},
		{
			name: "contains",
			s:    &StringArray{"1"},
			args: args{
				items: []string{"1", "2", "3"},
			},
			want: true,
		},
		{
			name: "contains with different case",
			s:    &StringArray{"fOo"},
			args: args{
				items: []string{"foo", "2", "3"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.HasOneOf(tt.args.items...); got != tt.want {
				t.Errorf("StringArray.HasOneOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
