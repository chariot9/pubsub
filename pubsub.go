package pubsub

type chanMapTopics map[chan interface{}][]string
type topicMapChans map[string][]chan interface{}

type Pubsub struct {
	capacity      int
	chanMapTopics chanMapTopics
	topicMapChans topicMapChans
}

func NewPubsub(capacity int) *Pubsub {
	server := Pubsub{capacity, make(chanMapTopics), make(topicMapChans)}
	return &server
}

func (p *Pubsub) Publish(content interface{}, topics ...string) {
	for _, topic := range topics {
		if chanList, ok := p.topicMapChans[topic]; ok {
			for _, channel := range chanList {
				channel <- content
			}
		}
	}
}

func (p *Pubsub) Subscribe(topics ...string) chan interface{} {
	chanClient := make(chan interface{}, p.capacity)
	p.updateTopicMapClient(chanClient, topics)
	return chanClient
}

func (p *Pubsub) updateTopicMapClient(clientChan chan interface{}, topics []string) {
	var updateChanList []chan interface{}
	for _, topic := range topics {
		updateChanList, _ = p.topicMapChans[topic]
		updateChanList = append(updateChanList, clientChan)
		p.topicMapChans[topic] = updateChanList
	}
	p.chanMapTopics[clientChan] = topics
}

func (p *Pubsub) AddSubscription(clientChan chan interface{}, topics ...string) {
	p.updateTopicMapClient(clientChan, topics)
}

func (p *Pubsub) RemoveSubscription(clientChan chan interface{}, topics ...string) {
	for _, topic := range topics {
		if chanList, ok := p.topicMapChans[topic]; ok {
			var updateChanList []chan interface{}
			for _, client := range chanList {
				if client != clientChan {
					updateChanList = append(updateChanList, client)
				}
			}
			p.topicMapChans[topic] = updateChanList
		}

		if topicList, ok := p.chanMapTopics[clientChan]; ok {
			var updateTopicList []string
			for _, updateTopic := range topicList {
				if updateTopic != topic {
					updateTopicList = append(updateTopicList, topic)
				}
			}
			p.chanMapTopics[clientChan] = updateTopicList
		}
	}
}
