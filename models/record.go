package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strconv"
)

// 将[]byte经md5计算哈希后转为hex字符串，取前8字符
func data2hashhex8(data []byte) string {
	hash := md5.Sum(data)
	hashhex := hex.EncodeToString(hash[:])
	hashhex8 := hashhex[:8]
	return hashhex8
}

type Record interface {
	// 转为Json格式输出，转成json之前一定加上消息类型
	Json() ([]byte, error)
	// 消息的唯一标识，应当是由 内容哈希拼上时间戳的
	ID() string
	// 消息类型
	Type() string
}

///////////////////////////// 各类Record ////////////////////////////

type PotNodeStateSwitchRecord struct {
	// 生成Json时自动生成
	RecordId string	// Record哈希用作标识，除RecordId以外全部内容
	RecordType string

	// 人工填入
	Time int64	// 时间
	NodeId string	// 节点id
	NewState string	// 新状态
}

// 正常来说，这些API都只会使用一次，所以无所谓会不会浪费些资源了
func (r PotNodeStateSwitchRecord) Json() ([]byte, error) {
	r.RecordType = "PotNodeStateSwitch"
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	r.RecordId = strconv.Itoa(int(r.Time)) + data2hashhex8(data)

	return json.Marshal(r)
}

func (r PotNodeStateSwitchRecord) ID() string {
	r.RecordType = "PotNodeStateSwitch"
	data, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	r.RecordId = strconv.Itoa(int(r.Time)) + data2hashhex8(data)
	return r.RecordId
}

func (r PotNodeStateSwitchRecord) Type() string {
	return "PotNodeStateSwitch"
}

/////////////////////////////////////////////////////////////////

type PotNodeSendMsgRecord struct {

}

////////////////////////////////////////////////////////////////

type FakeNodeFaultRecord struct {
	// 生成Json时自动生成
	RecordId string	// Record哈希用作标识，除RecordId以外全部内容
	RecordType string

	// 人工填入
	Time int64	// 时间
	NodeId string	// 节点id
}

func (r FakeNodeFaultRecord) Json() ([]byte, error) {
	r.RecordType = "FakeNodeFault"
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	r.RecordId = strconv.Itoa(int(r.Time)) + data2hashhex8(data)

	return json.Marshal(r)
}

func (r FakeNodeFaultRecord) ID() string {
	r.RecordType = "FakeNodeFault"
	data, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	r.RecordId = strconv.Itoa(int(r.Time)) + data2hashhex8(data)
	return r.RecordId
}

func (r FakeNodeFaultRecord) Type() string {
	return "FakeNodeFault"
}