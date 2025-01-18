package cache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func getEmptyCacheData() map[string]int {
	cacheData := make(map[string]int)
	return cacheData
}

func getPopulatedCacheData() map[string]int {
	cacheData := make(map[string]int)
	cacheData["key-one"] = 1
	return cacheData
}

func getMultiplePopulatedCacheData() map[string]int {
	cacheData := make(map[string]int)
	cacheData["key-one"] = 1
	cacheData["key-two"] = 2
	return cacheData
}

func TestCache_Add(t *testing.T) {
	type args struct {
		key   string
		value int
	}
	tests := []struct {
		name    string
		data    map[string]int
		args    args
		wantErr bool
		wantLen int
	}{
		{
			name: "Add to empty cache",
			data: getEmptyCacheData(),
			args: args{
				key:   "key-one",
				value: 1,
			},
			wantLen: 1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cache{
				data: tt.data,
			}
			err := c.Add(tt.args.key, tt.args.value)
			if !tt.wantErr {
				assert.NoError(t, err)
			}

			assert.Len(t, tt.data, tt.wantLen)
		})
	}
}

func TestCache_Get(t *testing.T) {
	type args struct {
		key   string
		value int
	}

	tests := []struct {
		name    string
		data    map[string]int
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Find value for key",
			data: getPopulatedCacheData(),
			args: args{
				key: "key-one",
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "Not Found value for key as key does not exist",
			data: getPopulatedCacheData(),
			args: args{
				key: "key-two",
			},
			want: 0,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == NotFoundErr
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cache{
				data: tt.data,
			}
			result, err := c.Get(tt.args.key)
			if tt.wantErr != nil {
				if !tt.wantErr(t, err, fmt.Sprintf("Get(%v)", tt.args.key)) {
					require.Failf(t, "Invalid error message", "result: %v", err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equalf(t, tt.want, result, "Get(%v)", tt.args.key)
		})
	}
}

func TestCache_List(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]int
		want    map[string]int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Get cache entries",
			data:    getMultiplePopulatedCacheData(),
			want:    getMultiplePopulatedCacheData(),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cache{
				data: tt.data,
			}
			result, err := c.List()
			if tt.wantErr != nil {
				if !tt.wantErr(t, err, fmt.Sprintf("List()")) {
					assert.Fail(t, "Invalid error")
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equalf(t, tt.want, result, "List()")
		})
	}
}

func TestCache_Remove(t *testing.T) {
	type args struct {
		key string
	}

	tests := []struct {
		name    string
		data    map[string]int
		args    args
		wantErr assert.ErrorAssertionFunc
		wantLen int
	}{
		{
			name: "Remove key from cache leaves it empty",
			data: getPopulatedCacheData(),
			args: args{
				key: "key-one",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cache{
				data: tt.data,
			}

			err := c.Remove(tt.args.key)

			if tt.wantErr != nil {
				if !tt.wantErr(t, err, fmt.Sprintf("Remove cache error")) {
					assert.Fail(t, "Invalid error")
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Len(t, c.data, tt.wantLen)
		})
	}
}

func TestCache_Update(t *testing.T) {
	type args struct {
		key   string
		value int
	}
	tests := []struct {
		name    string
		data    map[string]int
		args    args
		wantErr assert.ErrorAssertionFunc
		wantLen int
	}{
		{
			name: "Update value for key1 in cache",
			data: getPopulatedCacheData(),
			args: args{
				key:   "key-one",
				value: 2,
			},
			wantErr: nil,
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cache{
				data: tt.data,
			}

			err := c.Update(tt.args.key, tt.args.value)
			if tt.wantErr != nil {
				if !tt.wantErr(t, err, fmt.Sprintf("Update cache error")) {
					assert.Fail(t, "Invalid error")
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Len(t, c.data, tt.wantLen)
		})
	}
}
