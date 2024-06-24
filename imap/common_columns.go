package imap

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func commonColumns(c []*plugin.Column) []*plugin.Column {
	return append([]*plugin.Column{
		{
			Name:        "login",
			Description: "The login name.",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getLoginName,
			Transform:   transform.FromValue(),
		},
	}, c...)
}

// if the caching is required other than per connection, build a cache key for the call and use it in Memoize.
var getLoginNameMemoized = plugin.HydrateFunc(getLoginNameUncached).Memoize(memoize.WithCacheKeyFunction(getLoginNameCacheKey))

// declare a wrapper hydrate function to call the memoized function
// - this is required when a memoized function is used for a column definition
func getLoginName(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getLoginNameMemoized(ctx, d, h)
}

// Build a cache key for the call to getLoginNameCacheKey.
func getLoginNameCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := "getLoginName"
	return key, nil
}

func getLoginNameUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cfg := GetConfig(d.Connection)
	
	return cfg.Login, nil
}
