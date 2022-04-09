package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Broker interface {
	publish(topic string, msg interface{}) error
	close()         //关闭
	setCap(cap int) //容量
	sub(topic string) (<-chan interface{}, error)
	unsub(topic string, sub <-chan interface{}) error
	brocast(msg interface{}, subers []chan interface{})
}

type BrokerImpl struct {
	exit   chan bool //退出
	cap    int       //大小
	topics map[string][]chan interface{}
	sync.RWMutex
}

func (b *BrokerImpl) publish(topic string, msg interface{}) error {
	select {
	case <-b.exit:
		return errors.New("broker closed")
	default:
	}

	b.RLock()
	defer b.RUnlock()

	subers, ok := b.topics[topic]
	if !ok {
		fmt.Println("subers len 0")
		return nil
	}

	b.brocast(msg, subers)
	return nil
}

func (b *BrokerImpl) brocast(msg interface{}, subers []chan interface{}) {
	count := len(subers)
	fmt.Println("brocast subers len ", count)
	go func() {
		for i := 0; i < count; i++ {
			select {
			case subers[i] <- msg:
			case <-time.After(time.Millisecond * 10):
			case <-b.exit:
				return
			}

		}
	}()
}

func (b *BrokerImpl) sub(topic string) (<-chan interface{}, error) {
	select {
	case <-b.exit:
		return nil, errors.New("broker end")
	default:
	}

	ch := make(chan interface{}, b.cap)
	b.Lock()
	defer b.Unlock()

	b.topics[topic] = append(b.topics[topic], ch)
	fmt.Println("broker sub topic : ", topic)

	go func() {
		for {
			select {
			case <-b.exit:
				fmt.Println("broker sub go func finish")
				return
			case msg := <-ch:
				fmt.Println("get msg ", msg)
			default:
			}
		}
	}()
	return ch, nil
}

func (b *BrokerImpl) unsub(topic string, sub <-chan interface{}) error {
	select {
	case <-b.exit:
		return errors.New("broker end")
	default:
	}

	b.RUnlock()
	subs, ok := b.topics[topic]
	b.Unlock()

	if !ok {
		return nil
	}

	var newSubs []chan interface{}
	for _, suber := range subs {
		if suber == sub {
			continue
		}

		newSubs = append(newSubs, suber)
	}

	b.Lock()
	b.topics[topic] = newSubs
	b.Unlock()
	return nil
}

func (b *BrokerImpl) close() {
	select {
	case <-b.exit:
		fmt.Println("brocker close exit")
		return
	default:
		close(b.exit)
		b.Lock()
		b.topics = make(map[string][]chan interface{})
		b.Unlock()
	}

	return
}

func (b *BrokerImpl) setCap(cap int) {
	b.cap = cap
}

type Client struct {
	bro *BrokerImpl
}

func NewBroImp() *BrokerImpl {
	return &BrokerImpl{topics: make(map[string][]chan interface{}), exit: make(chan bool, 1)}
}
func NewClient() *Client {
	return &Client{bro: NewBroImp()}
}

func main() {
	b := NewClient()
	b.bro.setCap(100)
	topic := "test"
	ch, err := b.bro.sub(topic)
	if err != nil {
		fmt.Println("error")
	}

	payload := fmt.Sprintf("asong%d", 10)
	if err := b.bro.publish(topic, payload); err != nil {
		fmt.Println("err 1")
	}

	fmt.Println(ch)
	//b.bro.unsub(topic, ch)

	time.Sleep(time.Second * 10)
	fmt.Println("time sleep 10s finish")
	b.bro.exit <- true
	time.Sleep(time.Second * 10)
	fmt.Println("main fun exit")
}
