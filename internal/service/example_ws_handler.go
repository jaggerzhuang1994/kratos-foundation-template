package service

import (
	"errors"
	"fmt"

	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/server"
)

type ExampleWsHandler struct {
}

func NewExampleWsHandler(
// todo biz
) *ExampleWsHandler {
	return &ExampleWsHandler{}
}

// OnHandshake 建立连接前校验，返回错误表示拒绝连接
func (h *ExampleWsHandler) OnHandshake(request *http.Request) error {
	fmt.Println("OnHandshake ", request.RemoteAddr)
	return errors.New("handshake fail")
}

// OnConnect 建立连接后调用
func (h *ExampleWsHandler) OnConnect(client server.WebsocketClient) {
	fmt.Println("OnConnect ", client.Request().RemoteAddr)
}

// OnClose 连接关闭后回调
func (h *ExampleWsHandler) OnClose(client server.WebsocketClient) {
	fmt.Println("OnClose ", client.Request().RemoteAddr)
}

// OnError 读取消息过程发生错误调用
func (h *ExampleWsHandler) OnError(client server.WebsocketClient, err error) {
	fmt.Println("OnError ", client.Request().RemoteAddr, err)
}

// OnMessage 收到消息后调用
func (h *ExampleWsHandler) OnMessage(client server.WebsocketClient, message []byte, messageType server.MessageType) {
	fmt.Println("OnMessage ", client.Request().RemoteAddr, messageType, message)
	err := client.SendText("echo")
	if err != nil {
		fmt.Println(err)
	}
	_ = client.SendJSON(map[string]interface{}{
		"message_type": messageType,
		"message":      string(message),
	})
	// 主动关闭连接
	// client.Close()

	// todo 在 service 解析ws的参数，并调用biz的业务接口
}
