func publishOrder(ch *amqp.Channel, order Order) error {
    package main
	priority := "standard"
    if order.IsExpress { priority = "express" }
    if order.IsInternational { priority = "international" }

    return ch.Publish(
        DirectExchange,          // exchange
        priority,                // routing key (matches queue binding)
        false,                   // mandatory
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:       order.ToJSON(),
        },
    )
}

// Worker setup
queues := []string{"express", "international", "standard"}
for _, priority := range queues {
    ch.QueueBind(
        priority + "_queue",  // queue name
        priority,            // routing key
        DirectExchange,      // exchange
        false,
        nil,
    )
}