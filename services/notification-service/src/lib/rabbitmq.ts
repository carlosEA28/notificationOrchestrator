import amqp, {ChannelModel} from "amqplib"

export class RabbitmqProvider{
    private connection?: ChannelModel
    private channel?: amqp.Channel

    constructor(private url = process.env.RABBITMQ_URL ?? "amqp://localhost") {}

    private async ensureConnected(){
        if (this.channel) {
            return
        }


        try {
            this.connection = await amqp.connect(this.url)
            this.channel = await this.connection.createChannel()
            console.log("Connected to RabbitMQ")
        } catch(err){
            console.error("Error connecting to RabbitMQ:", err)
            throw err
        }
    }

    async produce(exchange: string, message: unknown, routingKey = exchange, queueName?: string){
        await this.ensureConnected()
        const channel = this.channel
        if (!channel) {
            throw new Error("RabbitMQ channel not initialized")
        }

        await channel.assertExchange(exchange, "topic", { durable: true })
        if (queueName) {
            await channel.assertQueue(queueName, { durable: true })
            await channel.bindQueue(queueName, exchange, routingKey)
        }
        const payload = Buffer.from(typeof message === "string" ? message : JSON.stringify(message))
        channel.publish(exchange, routingKey, payload, { persistent: true })
        console.log(`Message published to exchange ${exchange} with routing key ${routingKey}`)
    }

    async close(){
        if (this.channel) {
            await this.channel.close()
        }
        if (this.connection) {
            await this.connection.close()
        }
        this.channel = undefined
        this.connection = undefined
    }
}
