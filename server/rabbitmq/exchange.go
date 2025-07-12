package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

const (
    DirectExchange  = "orders_direct"
    TopicExchange   = "orders_topic"
    FanoutExchange  = "orders_fanout"
)

func SetupExchanges(ch *amqp.Channel) error {
    // Direct Exchange - for priority-based routing
    err := ch.ExchangeDeclare(
        DirectExchange,
        "direct", // type
        true,      // durable
        false,     // auto-deleted
        false,     // internal
        false,     // no-wait
        nil,
    )
    if err != nil { return err }

    // Topic Exchange - for regional routing (e.g. "orders.europe.#")
    err = ch.ExchangeDeclare(
        TopicExchange,
        "topic",  // type
        true,     // durable
        false,
        false,
        false,
        nil,
    )
    
    // Fanout Exchange - for broadcasting to all services
    err = ch.ExchangeDeclare(
        FanoutExchange,
        "fanout", // type
        true,
        false,
        false,
        false,
        nil,
    )

    return err
}