package http

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// Kratos Request 解码器个性化处理
func RequestDecoder(r *http.Request, v interface{}) error {
	// 从Request Header的Content-Type中提取出对应的解码器
	codec, ok := khttp.CodecForRequest(r, "Content-Type")
	// 如果找不到对应的解码器此时会报错
	if !ok {
		return errors.BadRequest("CODEC", r.Header.Get("Content-Type"))
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}

	r.Body.Close()
	// 将读取的信息塞回 body 中
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	if err = codec.Unmarshal(data, v); err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}

	return nil
}
