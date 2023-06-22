// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package encryption_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/sdk/value/encryption"
	// Register encryption transformers.
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/aead"
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/age"
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/dae"
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/fernet"
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/jwe"
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/paseto"
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/secretbox"
	"zntr.io/harp/v2/pkg/sdk/value/mock"
)

func TestFromKey(t *testing.T) {
	type args struct {
		keyValue string
	}
	tests := []struct {
		name    string
		args    args
		want    value.Transformer
		wantErr bool
	}{
		{
			name: "blank",
			args: args{
				keyValue: "",
			},
			wantErr: true,
		},
		{
			name: "invalid aes-gcm",
			args: args{
				keyValue: "aes-gcm:zQyPnNa-jlQsLW3Ypd87cX88ROMkdgnqv0a3y8",
			},
			wantErr: true,
		},
		{
			name: "invalid secretbox",
			args: args{
				keyValue: "secretbox:gCUODuqhcktiM1USKOfkwVlKhoUyHxXZm6d6",
			},
			wantErr: true,
		},
		{
			name: "invalid fernet",
			args: args{
				keyValue: "fernet:ZER8WwNyw5Dsd65bctxillSrRMX4ObaZsQjaNW1",
			},
			wantErr: true,
		},
		{
			name: "aes-gcm",
			args: args{
				keyValue: "aes-gcm:zQyPnNa-jlQsLW3Ypd87cX88ROMkdgnqv0a3y8LiISg=",
			},
			wantErr: false,
		},
		{
			name: "dae-aes-gcm",
			args: args{
				keyValue: "dae-aes-gcm:zQyPnNa-jlQsLW3Ypd87cX88ROMkdgnqv0a3y8LiISg=",
			},
			wantErr: false,
		},
		{
			name: "dae-aes-gcm",
			args: args{
				keyValue: "dae-aes-gcm:zQyPnNa-jlQsLW3Ypd87cX88ROMkdgnqv0a3y8LiISg=:jc32fV49Vi94NUYPnYR6ShInCD5rAiuMkkK2zb-Up4k=",
			},
			wantErr: false,
		},
		{
			name: "secretbox",
			args: args{
				keyValue: "secretbox:gCUODuqhcktiM1USKOfkwVlKhoUyHxXZm6d64nztCp0=",
			},
			wantErr: false,
		},
		{
			name: "chacha",
			args: args{
				keyValue: "chacha:gCUODuqhcktiM1USKOfkwVlKhoUyHxXZm6d64nztCp0=",
			},
			wantErr: false,
		},
		{
			name: "dae-chacha",
			args: args{
				keyValue: "dae-chacha:gCUODuqhcktiM1USKOfkwVlKhoUyHxXZm6d64nztCp0=",
			},
			wantErr: false,
		},
		{
			name: "dae-chacha with salt",
			args: args{
				keyValue: "dae-chacha:gCUODuqhcktiM1USKOfkwVlKhoUyHxXZm6d64nztCp0=:jc32fV49Vi94NUYPnYR6ShInCD5rAiuMkkK2zb-Up4k=",
			},
			wantErr: false,
		},
		{
			name: "xchacha",
			args: args{
				keyValue: "xchacha:VhfCXaD_QwwwoPCjLJx6vgnaSo0sMPjdCmT0RUUQjBQ=",
			},
			wantErr: false,
		},
		{
			name: "dae-xchacha",
			args: args{
				keyValue: "dae-xchacha:VhfCXaD_QwwwoPCjLJx6vgnaSo0sMPjdCmT0RUUQjBQ=",
			},
			wantErr: false,
		},
		{
			name: "dae-xchacha with salt",
			args: args{
				keyValue: "dae-xchacha:VhfCXaD_QwwwoPCjLJx6vgnaSo0sMPjdCmT0RUUQjBQ=:jc32fV49Vi94NUYPnYR6ShInCD5rAiuMkkK2zb-Up4k=",
			},
			wantErr: false,
		},
		{
			name: "fernet",
			args: args{
				keyValue: "fernet:ZER8WwNyw5Dsd65bctxillSrRMX4ObaZsQjaNW1nBBI=",
			},
			wantErr: false,
		},
		{
			name: "aes-siv",
			args: args{
				keyValue: "aes-siv:2XEKpPbE8T0ghLj8Wr9v6stV0YrUCNSoSbtc69Kh-n7-pVaKmWZ8LSvaJOK9BJHqDWE8vyNSzyNpcTYv3-J9lw==",
			},
			wantErr: false,
		},
		{
			name: "dae-aes-siv",
			args: args{
				keyValue: "dae-aes-siv:2XEKpPbE8T0ghLj8Wr9v6stV0YrUCNSoSbtc69Kh-n7-pVaKmWZ8LSvaJOK9BJHqDWE8vyNSzyNpcTYv3-J9lw==",
			},
			wantErr: false,
		},
		{
			name: "dae-aes-siv with salt",
			args: args{
				keyValue: "dae-aes-siv:2XEKpPbE8T0ghLj8Wr9v6stV0YrUCNSoSbtc69Kh-n7-pVaKmWZ8LSvaJOK9BJHqDWE8vyNSzyNpcTYv3-J9lw==:jc32fV49Vi94NUYPnYR6ShInCD5rAiuMkkK2zb-Up4k=",
			},
			wantErr: false,
		},
		{
			name: "aes-pmac-siv",
			args: args{
				keyValue: "aes-pmac-siv:Brfled4G7okhpCb6T2HMWKgDo1vyqrEdWWVIXfcFUysHaOacXkER5z9GHRuz89scK2TSE962nAFUcScAkihP9w==",
			},
			wantErr: false,
		},
		{
			name: "dae-aes-pmac-siv",
			args: args{
				keyValue: "dae-aes-pmac-siv:Brfled4G7okhpCb6T2HMWKgDo1vyqrEdWWVIXfcFUysHaOacXkER5z9GHRuz89scK2TSE962nAFUcScAkihP9w==",
			},
			wantErr: false,
		},
		{
			name: "dae-aes-pmac-siv with salt",
			args: args{
				keyValue: "dae-aes-pmac-siv:Brfled4G7okhpCb6T2HMWKgDo1vyqrEdWWVIXfcFUysHaOacXkER5z9GHRuz89scK2TSE962nAFUcScAkihP9w==:jc32fV49Vi94NUYPnYR6ShInCD5rAiuMkkK2zb-Up4k=",
			},
			wantErr: false,
		},
		{
			name: "jwe",
			args: args{
				keyValue: "jwe:a256kw:ZER8WwNyw5Dsd65bctxillSrRMX4ObaZsQjaNW1nBBI=",
			},
			wantErr: false,
		},
		{
			name: "paseto",
			args: args{
				keyValue: "paseto:kP1yHnBcOhjowNFXSCyycSuXdUqTlbuE6ES5tTp-I_o=",
			},
			wantErr: false,
		},
		/*{
			name: "age-recipients",
			args: args{
				keyValue: "age-recipients:age1ce20pmz8z0ue97v7rz838v6pcpvzqan30lr40tjlzy40ez8eldrqf2zuxe",
			},
			wantErr: false,
		},
		{
			name: "age-identity",
			args: args{
				keyValue: "age-identity:AGE-SECRET-KEY-1W8E69DQEVASNK68FX7C6QLD99KTG96RHWW0EZ3RD0L29AHV4S84QHUAP4C",
			},
			wantErr: false,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encryption.FromKey(tt.args.keyValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				return
			}

			// Ensure not panic
			assert.NotPanics(t, func() {
				encryption.Must(got, err)
			})

			// Encrypt
			msg := []byte("msg")
			encrypted, err := got.To(context.Background(), msg)
			if err != nil {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Decrypt
			decrypted, err := got.From(context.Background(), encrypted)
			if err != nil {
				t.Errorf("From() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check identity
			if !reflect.DeepEqual(msg, decrypted) {
				t.Errorf("expectd: %v, got: %v", msg, decrypted)
				return
			}
		})
	}
}

func TestMust(t *testing.T) {
	assert.Panics(t, func() {
		encryption.Must(mock.Transformer(nil), errors.New("test"))
	})

	assert.Panics(t, func() {
		encryption.Must(nil, nil)
	})
}
