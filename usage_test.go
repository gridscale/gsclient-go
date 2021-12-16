package gsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetGeneralUsage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	mux.HandleFunc(apiProjectLevelUsage, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet(""))
	})
	mux.HandleFunc(apiContractLevelUsage, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet(""))
	})
	type args struct {
		ctx              context.Context
		queryLevel       usageQueryLevel
		fromTime         GSTime
		toTime           *GSTime
		withoutDeleted   bool
		intervalVariable string
	}
	tests := []struct {
		name    string
		args    args
		want    GeneralUsage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successful Project level GetGeneralUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockGeneralUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetGeneralUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockGeneralUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetGeneralUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockGeneralUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetGeneralUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockGeneralUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetGeneralUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockGeneralUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetGeneralUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockGeneralUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetGeneralUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockGeneralUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetGeneralUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockGeneralUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Invalid query level GetGeneralUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       100000,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    GeneralUsage{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetGeneralUsage(tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
			if !tt.wantErr(t, err, fmt.Sprintf("GetGeneralUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetGeneralUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
		})
	}
}

func TestClient_GetServersUsage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	projectServerURI := path.Join(apiProjectLevelUsage, "servers")
	contractServerURI := path.Join(apiContractLevelUsage, "servers")
	mux.HandleFunc(projectServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("servers"))
	})
	mux.HandleFunc(contractServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("servers"))
	})
	type args struct {
		ctx              context.Context
		queryLevel       usageQueryLevel
		fromTime         GSTime
		toTime           *GSTime
		withoutDeleted   bool
		intervalVariable string
	}
	tests := []struct {
		name    string
		args    args
		want    ServersUsage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successful Project level GetServersUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockServersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetServersUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockServersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetServersUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockServersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetServersUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockServersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetServersUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockServersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetServersUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockServersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetServersUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockServersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetServersUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockServersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Invalid query level GetServersUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       100000,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    ServersUsage{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetServersUsage(tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
			if !tt.wantErr(t, err, fmt.Sprintf("GetServersUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetServersUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
		})
	}
}

func TestClient_GetDistributedStoragesUsage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	projectServerURI := path.Join(apiProjectLevelUsage, "distributed_storages")
	contractServerURI := path.Join(apiContractLevelUsage, "distributed_storages")
	mux.HandleFunc(projectServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("distributed_storages"))
	})
	mux.HandleFunc(contractServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("distributed_storages"))
	})
	type args struct {
		ctx              context.Context
		queryLevel       usageQueryLevel
		fromTime         GSTime
		toTime           *GSTime
		withoutDeleted   bool
		intervalVariable string
	}
	tests := []struct {
		name    string
		args    args
		want    DistributedStoragesUsage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successful Project level GetDistributedStoragesUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockDistributedStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetDistributedStoragesUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockDistributedStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetDistributedStoragesUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockDistributedStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetDistributedStoragesUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockDistributedStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetDistributedStoragesUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockDistributedStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetDistributedStoragesUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockDistributedStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetDistributedStoragesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockDistributedStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetDistributedStoragesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockDistributedStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Invalid query level GetDistributedStoragesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       100000,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    DistributedStoragesUsage{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetDistributedStoragesUsage(tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
			if !tt.wantErr(t, err, fmt.Sprintf("GetDistributedStoragesUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetDistributedStoragesUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
		})
	}
}

func TestClient_GetRocketStoragesUsage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	projectServerURI := path.Join(apiProjectLevelUsage, "rocket_storages")
	contractServerURI := path.Join(apiContractLevelUsage, "rocket_storages")
	mux.HandleFunc(projectServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("rocket_storages"))
	})
	mux.HandleFunc(contractServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("rocket_storages"))
	})
	type args struct {
		ctx              context.Context
		queryLevel       usageQueryLevel
		fromTime         GSTime
		toTime           *GSTime
		withoutDeleted   bool
		intervalVariable string
	}
	tests := []struct {
		name    string
		args    args
		want    RocketStoragesUsage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successful Project level GetRocketStoragesUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockRocketStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetRocketStoragesUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockRocketStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetRocketStoragesUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockRocketStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetRocketStoragesUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockRocketStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetRocketStoragesUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockRocketStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetRocketStoragesUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockRocketStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetRocketStoragesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockRocketStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetRocketStoragesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockRocketStoragesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Invalid query level GetRocketStoragesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       100000,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    RocketStoragesUsage{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetRocketStoragesUsage(tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
			if !tt.wantErr(t, err, fmt.Sprintf("GetRocketStoragesUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetRocketStoragesUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
		})
	}
}

func TestClient_GetStorageBackupsUsage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	projectServerURI := path.Join(apiProjectLevelUsage, "storage_backups")
	contractServerURI := path.Join(apiContractLevelUsage, "storage_backups")
	mux.HandleFunc(projectServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("storage_backups"))
	})
	mux.HandleFunc(contractServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("storage_backups"))
	})
	type args struct {
		ctx              context.Context
		queryLevel       usageQueryLevel
		fromTime         GSTime
		toTime           *GSTime
		withoutDeleted   bool
		intervalVariable string
	}
	tests := []struct {
		name    string
		args    args
		want    StorageBackupsUsage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successful Project level GetStorageBackupsUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockStorageBackupsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetStorageBackupsUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockStorageBackupsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetStorageBackupsUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockStorageBackupsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetStorageBackupsUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockStorageBackupsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetStorageBackupsUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockStorageBackupsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetStorageBackupsUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockStorageBackupsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetStorageBackupsUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockStorageBackupsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetStorageBackupsUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockStorageBackupsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Invalid query level GetStorageBackupsUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       100000,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    StorageBackupsUsage{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetStorageBackupsUsage(tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
			if !tt.wantErr(t, err, fmt.Sprintf("GetStorageBackupsUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetStorageBackupsUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
		})
	}
}

func TestClient_GetSnapshotsUsage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	projectServerURI := path.Join(apiProjectLevelUsage, "snapshots")
	contractServerURI := path.Join(apiContractLevelUsage, "snapshots")
	mux.HandleFunc(projectServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("snapshots"))
	})
	mux.HandleFunc(contractServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("snapshots"))
	})
	type args struct {
		ctx              context.Context
		queryLevel       usageQueryLevel
		fromTime         GSTime
		toTime           *GSTime
		withoutDeleted   bool
		intervalVariable string
	}
	tests := []struct {
		name    string
		args    args
		want    SnapshotsUsage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successful Project level GetSnapshotsUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockSnapshotsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetSnapshotsUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockSnapshotsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetSnapshotsUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockSnapshotsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetSnapshotsUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockSnapshotsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetSnapshotsUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockSnapshotsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetSnapshotsUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockSnapshotsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetSnapshotsUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockSnapshotsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetSnapshotsUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockSnapshotsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Invalid query level GetSnapshotsUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       100000,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    SnapshotsUsage{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetSnapshotsUsage(tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
			if !tt.wantErr(t, err, fmt.Sprintf("GetSnapshotsUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetSnapshotsUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
		})
	}
}

func TestClient_GetTemplatesUsage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	projectServerURI := path.Join(apiProjectLevelUsage, "templates")
	contractServerURI := path.Join(apiContractLevelUsage, "templates")
	mux.HandleFunc(projectServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("templates"))
	})
	mux.HandleFunc(contractServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("templates"))
	})
	type args struct {
		ctx              context.Context
		queryLevel       usageQueryLevel
		fromTime         GSTime
		toTime           *GSTime
		withoutDeleted   bool
		intervalVariable string
	}
	tests := []struct {
		name    string
		args    args
		want    TemplatesUsage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successful Project level GetTemplatesUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockTemplatesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetTemplatesUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockTemplatesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetTemplatesUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockTemplatesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetTemplatesUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockTemplatesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetTemplatesUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockTemplatesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetTemplatesUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockTemplatesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetTemplatesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockTemplatesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetTemplatesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockTemplatesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Invalid query level GetTemplatesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       100000,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    TemplatesUsage{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetTemplatesUsage(tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
			if !tt.wantErr(t, err, fmt.Sprintf("GetTemplatesUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetTemplatesUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
		})
	}
}

func TestClient_GetISOImagesUsage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	projectServerURI := path.Join(apiProjectLevelUsage, "iso_images")
	contractServerURI := path.Join(apiContractLevelUsage, "iso_images")
	mux.HandleFunc(projectServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("iso_images"))
	})
	mux.HandleFunc(contractServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("iso_images"))
	})
	type args struct {
		ctx              context.Context
		queryLevel       usageQueryLevel
		fromTime         GSTime
		toTime           *GSTime
		withoutDeleted   bool
		intervalVariable string
	}
	tests := []struct {
		name    string
		args    args
		want    ISOImagesUsage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successful Project level GetISOImagesUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockISOImagesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetISOImagesUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockISOImagesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetISOImagesUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockISOImagesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetISOImagesUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockISOImagesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetISOImagesUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockISOImagesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetISOImagesUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockISOImagesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetISOImagesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockISOImagesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetISOImagesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockISOImagesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Invalid query level GetISOImagesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       100000,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    ISOImagesUsage{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetISOImagesUsage(tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
			if !tt.wantErr(t, err, fmt.Sprintf("GetISOImagesUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetISOImagesUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
		})
	}
}

func TestClient_GetIPsUsage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	projectServerURI := path.Join(apiProjectLevelUsage, "ip_addresses")
	contractServerURI := path.Join(apiContractLevelUsage, "ip_addresses")
	mux.HandleFunc(projectServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("ip_addresses"))
	})
	mux.HandleFunc(contractServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("ip_addresses"))
	})
	type args struct {
		ctx              context.Context
		queryLevel       usageQueryLevel
		fromTime         GSTime
		toTime           *GSTime
		withoutDeleted   bool
		intervalVariable string
	}
	tests := []struct {
		name    string
		args    args
		want    IPsUsage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successful Project level GetIPsUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockIPsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetIPsUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockIPsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetIPsUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockIPsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetIPsUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockIPsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetIPsUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockIPsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetIPsUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockIPsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetIPsUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockIPsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetIPsUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockIPsUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Invalid query level GetIPsUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       100000,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    IPsUsage{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetIPsUsage(tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
			if !tt.wantErr(t, err, fmt.Sprintf("GetIPsUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetIPsUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
		})
	}
}

func TestClient_GetLoadBalancersUsage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	projectServerURI := path.Join(apiProjectLevelUsage, "load_balancers")
	contractServerURI := path.Join(apiContractLevelUsage, "load_balancers")
	mux.HandleFunc(projectServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("load_balancers"))
	})
	mux.HandleFunc(contractServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("load_balancers"))
	})
	type args struct {
		ctx              context.Context
		queryLevel       usageQueryLevel
		fromTime         GSTime
		toTime           *GSTime
		withoutDeleted   bool
		intervalVariable string
	}
	tests := []struct {
		name    string
		args    args
		want    LoadBalancersUsage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successful Project level GetLoadBalancersUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockLoadBalancersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetLoadBalancersUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockLoadBalancersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetLoadBalancersUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockLoadBalancersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetLoadBalancersUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockLoadBalancersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetLoadBalancersUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockLoadBalancersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetLoadBalancersUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockLoadBalancersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetLoadBalancersUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockLoadBalancersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetLoadBalancersUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockLoadBalancersUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Invalid query level GetLoadBalancersUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       100000,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    LoadBalancersUsage{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetLoadBalancersUsage(tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
			if !tt.wantErr(t, err, fmt.Sprintf("GetLoadBalancersUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetLoadBalancersUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
		})
	}
}

func TestClient_GetPaaSServicesUsage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	projectServerURI := path.Join(apiProjectLevelUsage, "paas_services")
	contractServerURI := path.Join(apiContractLevelUsage, "paas_services")
	mux.HandleFunc(projectServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("paas_services"))
	})
	mux.HandleFunc(contractServerURI, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareResourceUsageGet("paas_services"))
	})
	type args struct {
		ctx              context.Context
		queryLevel       usageQueryLevel
		fromTime         GSTime
		toTime           *GSTime
		withoutDeleted   bool
		intervalVariable string
	}
	tests := []struct {
		name    string
		args    args
		want    PaaSServicesUsage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successful Project level GetPaaSServicesUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockPaaSServicesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetPaaSServicesUsage with toTime=nil, withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           nil,
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockPaaSServicesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetPaaSServicesUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockPaaSServicesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetPaaSServicesUsage with withoutDeleted=false, intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   false,
				intervalVariable: "",
			},
			want:    getMockPaaSServicesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetPaaSServicesUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockPaaSServicesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetPaaSServicesUsage with intervalVariable=\"\"",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: "",
			},
			want:    getMockPaaSServicesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Project level GetPaaSServicesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ProjectLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockPaaSServicesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Successful Contract level GetPaaSServicesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       ContractLevelUsage,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    getMockPaaSServicesUsage(),
			wantErr: assert.NoError,
		},
		{
			name: "Invalid query level GetPaaSServicesUsage",
			args: args{
				ctx:              context.Background(),
				queryLevel:       100000,
				fromTime:         GSTime{time.Now().Add(-24 * time.Hour)},
				toTime:           &GSTime{time.Now()},
				withoutDeleted:   true,
				intervalVariable: HourIntervalVariable,
			},
			want:    PaaSServicesUsage{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetPaaSServicesUsage(tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
			if !tt.wantErr(t, err, fmt.Sprintf("GetPaaSServicesUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetPaaSServicesUsage(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.queryLevel, tt.args.fromTime, tt.args.toTime, tt.args.withoutDeleted, tt.args.intervalVariable)
		})
	}
}

var dummyRsCurrentUsagePerMinute = []Usage{{
	ProductNumber: 1,
	Value:         2,
}}

var dummyRsUsagePerInterval = []UsagePerInterval{{
	IntervalStart: dummyTime,
	IntervalEnd:   dummyTime,
	AccumulatedUsage: []Usage{{
		ProductNumber: 1,
		Value:         2,
	}},
}}

func getMockGeneralUsage() GeneralUsage {
	dummyRsUsageInfo := ResourceUsageInfo{
		CurrentUsagePerMinute: dummyRsCurrentUsagePerMinute,
		UsagePerInterval:      dummyRsUsagePerInterval,
	}
	mock := GeneralUsage{
		ResourcesUsage: GeneralUsageProperties{
			Servers:             dummyRsUsageInfo,
			RocketStorages:      dummyRsUsageInfo,
			DistributedStorages: dummyRsUsageInfo,
			StorageBackups:      dummyRsUsageInfo,
			Snapshots:           dummyRsUsageInfo,
			Templates:           dummyRsUsageInfo,
			IsoImages:           dummyRsUsageInfo,
			IPAddresses:         dummyRsUsageInfo,
			LoadBalancers:       dummyRsUsageInfo,
			PaaSServices:        dummyRsUsageInfo,
		},
	}
	return mock
}

func getMockServersUsage() ServersUsage {
	mock := ServersUsage{
		ResourcesUsage: []ServerUsageProperties{{
			ObjectUUID:            dummyUUID,
			Name:                  "test",
			Memory:                2,
			Cores:                 1,
			Power:                 true,
			Labels:                []string{"test"},
			Deleted:               false,
			Status:                "active",
			CurrentUsagePerMinute: dummyRsCurrentUsagePerMinute,
			UsagePerInterval:      dummyRsUsagePerInterval,
		},
		}}
	return mock
}

func getMockDistributedStoragesUsage() DistributedStoragesUsage {
	mock := DistributedStoragesUsage{
		ResourcesUsage: []StorageUsageProperties{{
			ObjectUUID:            dummyUUID,
			ParentUUID:            dummyUUID,
			Name:                  "test",
			Labels:                []string{"test"},
			Deleted:               false,
			Status:                "active",
			StorageType:           string(InsaneStorageType),
			LastUsedTemplate:      "",
			Capacity:              1000,
			CurrentUsagePerMinute: dummyRsCurrentUsagePerMinute,
			UsagePerInterval:      dummyRsUsagePerInterval,
		},
		}}
	return mock
}

func getMockRocketStoragesUsage() RocketStoragesUsage {
	mock := RocketStoragesUsage{
		ResourcesUsage: []StorageUsageProperties{{
			ObjectUUID:            dummyUUID,
			ParentUUID:            dummyUUID,
			Name:                  "test",
			Labels:                []string{"test"},
			Deleted:               false,
			Status:                "active",
			StorageType:           string(InsaneStorageType),
			LastUsedTemplate:      "",
			Capacity:              1000,
			CurrentUsagePerMinute: dummyRsCurrentUsagePerMinute,
			UsagePerInterval:      dummyRsUsagePerInterval,
		},
		}}
	return mock
}

func getMockStorageBackupsUsage() StorageBackupsUsage {
	mock := StorageBackupsUsage{
		ResourcesUsage: []StorageBackupUsageProperties{{
			ObjectUUID:            dummyUUID,
			Name:                  "test",
			CreateTime:            dummyTime,
			ChangeTime:            dummyTime,
			Capacity:              1000,
			CurrentUsagePerMinute: dummyRsCurrentUsagePerMinute,
			UsagePerInterval:      dummyRsUsagePerInterval,
		},
		}}
	return mock
}

func getMockSnapshotsUsage() SnapshotsUsage {
	mock := SnapshotsUsage{
		ResourcesUsage: []SnapshotUsageProperties{{
			ObjectUUID:            dummyUUID,
			Name:                  "test",
			ParentUUID:            dummyUUID,
			ParentName:            "test",
			ProjectUUID:           dummyUUID,
			Labels:                nil,
			Status:                "active",
			CreateTime:            dummyTime,
			ChangeTime:            dummyTime,
			Capacity:              1000,
			Deleted:               false,
			CurrentUsagePerMinute: dummyRsCurrentUsagePerMinute,
			UsagePerInterval:      dummyRsUsagePerInterval,
		},
		}}
	return mock
}

func getMockTemplatesUsage() TemplatesUsage {
	mock := TemplatesUsage{
		ResourcesUsage: []TemplateUsageProperties{{
			ObjectUUID:            dummyUUID,
			Name:                  "test",
			Status:                "active",
			Ostype:                "linux",
			Version:               "0.1",
			CreateTime:            dummyTime,
			ChangeTime:            dummyTime,
			Private:               false,
			LicenseProductNo:      0,
			Capacity:              1000,
			Distro:                "ubuntu",
			Description:           "test",
			Labels:                nil,
			ProjectUUID:           dummyUUID,
			Deleted:               false,
			CurrentUsagePerMinute: dummyRsCurrentUsagePerMinute,
			UsagePerInterval:      dummyRsUsagePerInterval,
		},
		}}
	return mock
}

func getMockISOImagesUsage() ISOImagesUsage {
	mock := ISOImagesUsage{
		ResourcesUsage: []ISOImageUsageProperties{{
			ObjectUUID:            dummyUUID,
			Name:                  "test",
			Description:           "test",
			SourceURL:             "https://example.com",
			Labels:                nil,
			Status:                "active",
			CreateTime:            dummyTime,
			ChangeTime:            dummyTime,
			Version:               "0.1",
			Private:               true,
			Capacity:              1000,
			ProjectUUID:           dummyUUID,
			Deleted:               false,
			CurrentUsagePerMinute: dummyRsCurrentUsagePerMinute,
			UsagePerInterval:      dummyRsUsagePerInterval,
		},
		}}
	return mock
}

func getMockIPsUsage() IPsUsage {
	mock := IPsUsage{
		ResourcesUsage: []IPUsageProperties{{
			ObjectUUID:            dummyUUID,
			Name:                  "test",
			IP:                    "192.168.1.1",
			Family:                0,
			CreateTime:            dummyTime,
			ChangeTime:            dummyTime,
			Status:                "active",
			LocationCountry:       "test",
			LocationName:          "test",
			LocationIata:          "test",
			LocationUUID:          dummyUUID,
			Prefix:                "192.168.1.1",
			DeleteBlock:           false,
			Failover:              false,
			Labels:                nil,
			ReverseDNS:            "192.168.2.1",
			PartnerUUID:           dummyUUID,
			ProjectUUID:           dummyUUID,
			Deleted:               false,
			CurrentUsagePerMinute: dummyRsCurrentUsagePerMinute,
			UsagePerInterval:      dummyRsUsagePerInterval,
		},
		}}
	return mock
}

func getMockLoadBalancersUsage() LoadBalancersUsage {
	mock := LoadBalancersUsage{
		ResourcesUsage: []LoadBalancerUsageProperties{{
			ObjectUUID: dummyUUID,
			Name:       "test",
			ForwardingRules: []ForwardingRule{
				{
					LetsencryptSSL: nil,
					ListenPort:     8080,
					Mode:           "http",
					TargetPort:     8000,
				},
			},
			BackendServers: []BackendServer{
				{
					Weight: 100,
					Host:   "185.201.147.176",
				},
			},
			CreateTime:            dummyTime,
			ChangeTime:            dummyTime,
			Status:                "active",
			RedirectHTTPToHTTPS:   false,
			Algorithm:             "leastconn",
			ListenIPv6UUID:        dummyUUID,
			ListenIPv4UUID:        dummyUUID,
			CurrentUsagePerMinute: dummyRsCurrentUsagePerMinute,
			UsagePerInterval:      dummyRsUsagePerInterval,
		},
		}}
	return mock
}

func getMockPaaSServicesUsage() PaaSServicesUsage {
	mock := PaaSServicesUsage{
		ResourcesUsage: []PaaSServiceUsageProperties{{
			ObjectUUID: dummyUUID,
			Name:       "test",
			Status:     "active",
			Credentials: []Credential{
				{
					Username: "username",
					Password: "password",
					Type:     "type",
				},
			},
			CreateTime:          dummyTime,
			ChangeTime:          dummyTime,
			ServiceTemplateUUID: dummyUUID,
			Parameters: map[string]interface{}{
				"TEST_PARAM": "test value",
			},
			ResourceLimits: []ResourceLimit{
				{
					Resource: "cpu",
					Limit:    2,
				},
			},
			ProjectUUID:           dummyUUID,
			Deleted:               false,
			CurrentUsagePerMinute: dummyRsCurrentUsagePerMinute,
			UsagePerInterval:      dummyRsUsagePerInterval,
		},
		}}
	return mock
}

func prepareResourceUsageGet(resourceName string) string {
	var usage interface{}
	switch resourceName {
	case "servers":
		usage = getMockServersUsage()
	case "distributed_storages":
		usage = getMockDistributedStoragesUsage()
	case "rocket_storages":
		usage = getMockRocketStoragesUsage()
	case "storage_backups":
		usage = getMockStorageBackupsUsage()
	case "snapshots":
		usage = getMockSnapshotsUsage()
	case "templates":
		usage = getMockTemplatesUsage()
	case "iso_images":
		usage = getMockISOImagesUsage()
	case "ip_addresses":
		usage = getMockIPsUsage()
	case "load_balancers":
		usage = getMockLoadBalancersUsage()
	case "paas_services":
		usage = getMockPaaSServicesUsage()
	default:
		usage = getMockGeneralUsage()
	}
	res, _ := json.Marshal(usage)
	return string(res)
}
