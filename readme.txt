Connection manager
    TODO: Peer dicovery with DHT
    TODO: manage connection

Efficient broadcast (simple, FEC-based Unicast)
    what if receive fake message (msgID) -> using bloom filter for each chunk hash
    what if receive same chunk in different sender -> DiscardDuplicateMessageValidator -> same bloom filter will be cache with same key

Pubsub
    how to strictly allow subscribe
        must include signature
        monitor and add to blacklist
    how to avoid duplicate Pubsub -> DiscardDuplicateMessageValidator




