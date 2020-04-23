package middleware

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
)

//Response 处理api返回值
func Response(conf *conf.MetadataConf) dispatcher.HandlerFunc {
	return func(ctx *dispatcher.Context) {
		ctx.Next()
		nctx := getCTX(ctx)
		if nctx == nil {
			return
		}
		defer nctx.Close()
		if err := nctx.Response.GetError(); err != nil {
			getLogger(ctx).Errorf("err:%v", err)
			if !servers.IsDebug {
				nctx.Response.ShouldContent(errors.New("请求发生错误"))
			}
		}
		if ctx.Writer.Written() {
			return
		}
		tp, content, err := nctx.Response.GetJSONRenderContent()
		writeTrace(getTrace(conf), tp, ctx, content)
		if err != nil {
			getLogger(ctx).Error(err)
			ctx.JSON(nctx.Response.GetStatus(), map[string]interface{}{"err": err})
			return
		}
		ctx.JSON(nctx.Response.GetStatus(), content)
	}
}
func writeTrace(b bool, tp int, ctx *dispatcher.Context, c interface{}) {
	if !b {
		return
	}
	switch v := c.(type) {
	case []byte:
		setResponseRaw(ctx, string(v))
	case string:
		setResponseRaw(ctx, v)
	default:
		var buff = bytes.NewBufferString("")
		switch tp {
		case context.CT_XML:
			xml.NewEncoder(buff).Encode(c)
		default:
			json.NewEncoder(buff).Encode(c)
		}
		setResponseRaw(ctx, strings.Trim(buff.String(), "\n"))
		buff.Reset()
	}
}
