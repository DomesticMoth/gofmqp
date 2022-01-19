/*
    This file is part of gofmqp.

    gofmqp is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    gofmqp is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with gofmqp.  If not, see <https://www.gnu.org/licenses/>.
*/
package gofmqp

const DEFAULT_COLLECTION_FREQUENCY = 100
const NO_COLLECT = 9223372036854775807

type Id int

type topicInfo struct{
	subs map[Id]bool
	topics map[MsgTopic]bool
}

type subInfo struct{
	stopics map[MsgTopic]bool
	ptopics map[MsgTopic]bool
}

type Router struct{
	subs map[Id]subInfo
	// Subscription topic
	stopics map[MsgTopic]topicInfo
	// Publication topic
	ptopics map[MsgTopic]topicInfo
	// TODO Come up with adequate variable names
	tflc uint64 // Time From Last Collection
	cf uint64 // Collection Frequency
}

func newTopicInfo() topicInfo {
	return topicInfo{
		make(map[Id]bool),
		make(map[MsgTopic]bool),
	}
}

func newSubInfo() subInfo {
	return subInfo{
		make(map[MsgTopic]bool),
		make(map[MsgTopic]bool),
	}
}

func NewRouter() Router {
	return Router{
		make(map[Id]subInfo),
		make(map[MsgTopic]topicInfo),
		make(map[MsgTopic]topicInfo),
		0,
		DEFAULT_COLLECTION_FREQUENCY,
	}
}

func NewRouterFrq(cf uint64) Router {
	return Router{
		make(map[Id]subInfo),
		make(map[MsgTopic]topicInfo),
		make(map[MsgTopic]topicInfo),
		0,
		cf,
	}
}


func (router *Router) Sub(sub Id, topic MsgTopic){
	// If subscriber does not exist, create it
	if _, exists := router.subs[sub]; !exists {
		router.subs[sub] = newSubInfo()
	}
	// If subscrioption topic does not exist, create it
	if _, exists := router.stopics[topic]; !exists {
		router.stopics[topic] = newTopicInfo()
	}
	// Link subscription topic and subscriber
	router.subs[sub].stopics[topic] = true
	router.stopics[topic].subs[sub] = true
	// For each publication topic
	for ptopic, ptopicInfo := range router.ptopics {
		// If it is equivalent to the subscription topic
		if CompareTopics(ptopic, topic){
			// Link subscription topic and publication topic
			router.stopics[topic].topics[ptopic] = true
			ptopicInfo.topics[topic] = true
			// Link publication topic and subscriber
			ptopicInfo.subs[sub] = true
			router.subs[sub].ptopics[ptopic] = true
		}
	}
}

func (router *Router) Collect(){
	// DRY principle is violated
	// TODO Refactor
	// For each publication topic
	for ptopic, ptopicInfo := range router.ptopics {
		// If there are no subscribers
		if len(ptopicInfo.subs) < 1 {
			// Unlink publication topic from all linked subscription topics
			for stopic, _ := range ptopicInfo.topics {
				delete(router.stopics[stopic].topics, ptopic)
			}
			// Delete publication topic
			delete(router.ptopics, ptopic)
		}
	}
	// For each subscription topic
	for stopic, stopicInfo := range router.stopics {
		// If there are no subscribers
		if len(stopicInfo.subs) < 1 {
			// Unlink subscription topic from all linked publication topics
			for ptopic, _ := range stopicInfo.topics {
				delete(router.ptopics[ptopic].topics, stopic)
			}
			// Delete subscription topic
			delete(router.stopics, stopic)
		}
	}
}

func (router *Router) Unsub(sub Id, topic MsgTopic){
	if subInf, exists := router.subs[sub]; exists {
		if _, exists := subInf.stopics[topic]; exists {
			// Unlink subscription topic and subscriber
			delete(router.stopics[topic].subs, sub)
			delete(router.subs[sub].stopics, topic)
			for ptopic, _ := range router.stopics[topic].topics {
				toremove := true
				ptopicInf := router.ptopics[ptopic]
				for stopic, _ := range ptopicInf.topics {
					if router.stopics[stopic].subs[sub] {
						toremove = false
					}
				}
				if toremove {
					// Unlink publication topic and subscriber
					delete(ptopicInf.subs, sub)
					delete(router.subs[sub].ptopics, ptopic)
				}
			}
			// If subscriber does not have subscriptions remove it
			if len(router.subs[sub].stopics) < 1 && len(router.subs[sub].ptopics) < 1 {
				delete(router.subs, sub)
			}
			// Collect garbage
			router.tflc += 1
			if router.tflc >= router.cf {
				router.Collect()
			}
		}
	}
}

func (router *Router) UnsubAll(sub Id){
	_, exist := router.subs[sub]
	if exist {
		stopics := router.GetSubscriptions(sub)
		for _, stopic := range stopics {
			router.Unsub(sub, stopic)
		}
	}
}

func (router *Router) GetSubscriptions(sub Id) []MsgTopic{
	subInf, exist := router.subs[sub]
	if exist {
		stopics:= make([]MsgTopic, 0, len(subInf.stopics))
    	for k := range subInf.stopics {
        	stopics = append(stopics, k)
    	}
    	return stopics
    }
    return []MsgTopic{}
}

func (router *Router) SubscribersCount(stopic MsgTopic) int {
	stopicInf, exist := router.stopics[stopic]
	if exist {
		return len(stopicInf.subs)
	}
	return 0
}

func (router *Router) Route(ptopic MsgTopic) map[Id]bool {
	_, exist := router.ptopics[ptopic]
	// If there is no publication topic in cashe
	// Add it and
	if !exist {
		router.ptopics[ptopic] = newTopicInfo()
		// For all subscription topics
		for stopic, _ := range router.stopics{
			// If eqal with publication topic
			if CompareTopics(ptopic, stopic){
				// Link subscription topic and publication topic
				router.stopics[stopic].topics[ptopic] = true
				router.ptopics[ptopic].topics[stopic] = true
				// For all related subscribers
				for sub, _ := range router.stopics[stopic].subs {
					// Link publication topic and subscriber
					router.subs[sub].ptopics[ptopic] = true
					router.ptopics[ptopic].subs[sub] = true
				}
			}
		}
	}
	return router.ptopics[ptopic].subs
}
