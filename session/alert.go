package session

import (
	"bytes"
	"logAnalysis/conf"
	"mmx/mediaservice/log"
	"strconv"
	"strings"
	"time"

	g "github.com/soniah/gosnmp"
)

//为了确保一次通话只有一个同类型的告警信息
//需要记录通话的ssrc，24小时之后自动删除
type Call struct {
	Date  time.Time //当前通话时间
	Ssrc  string    // 当前ssrc
	Types string    //是rtt还是lost_rate
}
type Alert struct {
	Calls []Call
}

/*实时告警相关功能*/
func parseRTT(line []byte) []byte {
	index := bytes.Index(line, []byte("rtt:"))
	if index < 0 {
		return nil
	}
	return line[index+4:]
}
func parseLostRate(line []byte) []byte {
	index := bytes.Index(line, []byte("lost_rate:"))
	if index < 0 {
		return nil
	}
	line = line[index+10:]
	index = bytes.Index(line, []byte("%"))
	if index < 0 {
		return nil
	}
	return line[:index]
}
func parsessrc(line []byte) []byte {
	index := bytes.Index(line, []byte("ssrc:"))
	if index < 0 {
		return nil
	}
	line = line[index+5:]
	index = bytes.Index(line, []byte(" "))
	if index < 0 {
		return nil
	}
	return line[:index]
}
func parseSSRC(line []byte) []byte {
	index := bytes.Index(line, []byte("SSRC:"))
	if index < 0 {
		return nil
	}
	line = line[index+5:]
	index = bytes.Index(line, []byte("]"))
	if index < 0 {
		return nil
	}
	return line[:index]
}

func (a *Alert) inCalls(types string, ssrc string) (in bool) {
	now := time.Now()
	var newCalls []Call
	for _, call := range a.Calls {
		if call.Types == types && call.Ssrc == ssrc {
			in = true
		}
		//新切只保存没有超过24小时的会话
		if now.Sub(call.Date) >= 24*time.Hour {
			continue
		} else {
			newCalls = append(newCalls, call)
		}
	}
	//如果不存在，则加入
	if !in {
		a.Calls = append(a.Calls, Call{
			Date:  time.Now(),
			Ssrc:  ssrc,
			Types: types,
		})
	} else {
		//更新为去除了存在24小时的会话切片
		a.Calls = newCalls
	}
	return
}

func (a *Alert) ParseLineAndAlert(lines []string) {

	for _, line := range lines {
		if strings.Contains(line, "rtt") {
			//获取ssrc
			currentSsrc := parsessrc([]byte(line))
			//检查当前ssrc是否存在， 存在则不再发告警信息
			if !a.inCalls("rtt", string(currentSsrc)) {
				currentRtt := string(parseRTT([]byte(line)))
				currentRttInt, _ := strconv.Atoi(currentRtt)
				//比较字符串长度大则说明超过阈值，因为rtt的值很大，转int会失败
				if len(currentRtt) > len(strconv.Itoa(conf.Cf.Alarm.Rtt)) {
					log.Info("send trap rtt: ", line)
					a.sendTrap([]byte(line))

				} else if currentRttInt > conf.Cf.Alarm.Rtt {
					//长度相同，则比较数值大小
					log.Info("send trap rtt: ", line)
					a.sendTrap([]byte(line))
				}
			}
		}

		if strings.Contains(line, "lost_rate") {
			//获取ssrc
			currentSsrc := parseSSRC([]byte(line))
			if !a.inCalls("lost_rates", string(currentSsrc)) {
				//获取到rate
				currentRate, _ := strconv.Atoi(string(parseLostRate([]byte(line))))
				if currentRate > conf.Cf.Alarm.LostRate {
					log.Info("send trap lost_rate: ", line)
					a.sendTrap([]byte(line))
				}
			}
		}
	}

}

func (a *Alert) sendTrap(alertlog []byte) {
	portInt, _ := strconv.Atoi(conf.Cf.Alarm.Port)
	g.Default.Target = conf.Cf.Alarm.SendTo
	g.Default.Port = uint16(portInt)
	g.Default.Version = g.Version2c
	g.Default.Community = "public"
	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
		return
	}
	defer g.Default.Conn.Close()
	pdu := g.SnmpPDU{
		Name:  conf.Cf.Alarm.Oid,
		Type:  g.OctetString,
		Value: alertlog,
	}
	trap := g.SnmpTrap{
		Variables: []g.SnmpPDU{pdu},
	}
	_, err = g.Default.SendTrap(trap)
	if err != nil {
		log.Fatalf("SendTrap() err: %v", err)
		return
	}
}
