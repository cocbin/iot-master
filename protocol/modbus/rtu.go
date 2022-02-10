package modbus

import (
	"errors"
	"fmt"
	"github.com/zgwit/iot-master/protocol/helper"
	"github.com/zgwit/iot-master/service"
)

type response struct {
	buf []byte
	err error
}

type request struct {
	cmd  []byte
	resp chan response //out
}

type RTU struct {
	link  service.Link
	queue chan request //in
}

func NewModbusRtu(link service.Link) *RTU {
	rtu := &RTU{
		link:  link,
		queue: make(chan request),
	}
	link.On("data", func(data []byte) {
		rtu.OnData(data)
	})

	return rtu
}

func (m *RTU) execute(cmd []byte) ([]byte, error) {
	req := request{
		cmd: cmd,
		resp: make(chan response),
	}
	//排队等待
	m.queue <- req

	//下发指令
	err := m.link.Write(cmd)
	if err != nil {
		//释放队列
		<- m.queue
		return nil, err
	}

	//等待结果
	resp := <-req.resp
	return resp.buf, resp.err
}

func (m *RTU) OnData(buf []byte) {
	req := <- m.queue

	//解析数据
	l := len(buf)
	crc := helper.ParseUint16(buf[l-2:])

	if crc != CRC16(buf[:l-2]) {
		//检验错误
		req.resp <- response{err: errors.New("校验错误")}
		return
	}

	//解析错误码
	if buf[1]&0x80 > 0 {
		req.resp <- response{err: fmt.Errorf("错误码：%d", buf[2])}
		return
	}

	//解析数据
	length := 4
	count := int(helper.ParseUint16(buf[1:]))
	switch buf[1] {
	case FuncCodeReadDiscreteInputs,
		FuncCodeReadCoils:
		length += 1 + count/8
		if count%8 != 0 {
			length++
		}

		if l < length {
			//长度不够
			req.resp <- response{err: errors.New("长度不够")}
			return
		}
		b := buf[2 : l-2]
		//数组解压
		bb := helper.ExpandBool(b, count)
		req.resp <- response{buf: bb}
	case FuncCodeReadInputRegisters,
		FuncCodeReadHoldingRegisters,
		FuncCodeReadWriteMultipleRegisters:
		length += 1 + count*2
		if l < length {
			//长度不够
			req.resp <- response{err: errors.New("长度不够")}
			return
		}
		b := buf[2 : l-2]
		req.resp <- response{buf: b}
	default:
		req.resp <- response{}
	}
}

func (m *RTU) Read(slave uint8, code uint8, offset uint16, size uint16) ([]byte, error) {
	b := make([]byte, 8)
	b[0] = slave
	b[1] = code
	helper.WriteUint16(b[2:], offset)
	helper.WriteUint16(b[4:], size)
	helper.WriteUint16(b[6:], CRC16(b[:6]))

	return m.execute(b)
}

func (m *RTU) Write(slave uint8, code uint8, offset uint16, buf []byte) error {
	length := len(buf)
	//如果是线圈，需要Shrink
	if code == 1 {
		switch code {
		case FuncCodeReadCoils:
			if length == 1 {
				code = 5
				//数据 转成 0x0000 0xFF00
				if buf[1] > 0 {
					buf = []byte{0xFF, 0}
				} else {
					buf = []byte{0, 0}
				}
			} else {
				code = 15 //0x0F
				//数组压缩
				b := helper.ShrinkBool(buf)
				count := len(b)
				buf = make([]byte, 3+count)
				helper.WriteUint16(buf, uint16(length))
				buf[2] = uint8(count)
				copy(buf[3:], b)
			}
		case FuncCodeReadHoldingRegisters:
			if length == 2 {
				code = 6
			} else {
				code = 16 //0x10
				b := make([]byte, 3+length)
				helper.WriteUint16(b, uint16(length/2))
				b[2] = uint8(length)
				copy(b[3:], buf)
				buf = b
			}
		default:
			return errors.New("功能码不支持")
		}
	}

	l := 6 + len(buf)
	b := make([]byte, l)
	b[0] = slave
	b[1] = code
	helper.WriteUint16(b[2:], offset)
	copy(b[4:], buf)
	helper.WriteUint16(b[l-2:], CRC16(b[:l-2]))

	_, err := m.execute(b)
	return err
}
