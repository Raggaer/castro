package lua

import (
	"bufio"
	"io"
	"io/ioutil"
	"net"
	"time"

	glua "github.com/yuin/gopher-lua"
)

// SetSocketProtocolMetaTable sets the socket metatable on the given lua state
func SetSocketMetaTable(luaState *glua.LState) {
	// Create and set socket metatable
	socketMetaTable := luaState.NewTypeMetatable(SocketMetatableName)
	luaState.SetGlobal(SocketMetatableName, socketMetaTable)

	// Set all socket metatable functions
	luaState.SetFuncs(socketMetaTable, socketMethods)
}

// GetSocket opens a socket to the given address and returns the response
func GetSocket(L *glua.LState) int {
	// Get protocol string
	protocol := L.Get(2)
	if protocol.Type() != glua.LTString {
		L.ArgError(1, "Invalid protocol type. Expected string")
		return 0
	}

	// Get address string
	address := L.Get(3)
	if address.Type() != glua.LTString {
		L.ArgError(2, "Invalid address type. Expected string")
		return 0
	}

	// Get message string
	message := L.Get(4)
	if address.Type() != glua.LTString {
		L.ArgError(3, "Invalid message type. Expected string")
		return 0
	}

	// Open socket connection
	sock, err := net.DialTimeout(protocol.String(), address.String(), 5*time.Second)
	if err != nil {
		L.Push(glua.LNil)
		L.Push(glua.LString(string("Failed to connect socket: " + err.Error())))
		return 2
	}

	defer sock.Close()

	// Write message to socket
	if _, err := sock.Write([]byte(message.String())); err != nil {
		L.Push(glua.LNil)
		L.Push(glua.LString(string("Failed to write to socket: " + err.Error())))
		return 2
	}

	reader := bufio.NewReader(sock)

	// Read response
	data, err := ioutil.ReadAll(reader)
	if err != nil && err != io.EOF {
		L.Push(glua.LNil)
		L.Push(glua.LString(string("Failed to read socket: " + err.Error())))
		return 2
	}

	// Push response as string
	L.Push(glua.LString(string(data)))

	return 1
}
