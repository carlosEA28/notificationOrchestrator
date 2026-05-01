import { FastifyReply, FastifyRequest } from "fastify";
import { NotificationSchema } from "../../../types/notification";
import { makeSendNotificationUseCase } from "../../../use-cases/factories/make-send-notification";

export async function sendNotification(request: FastifyRequest, reply: FastifyReply) {
  const params = await NotificationSchema.parseAsync(request.body);

  try {
    const sendNotificationUseCase = makeSendNotificationUseCase();
    await sendNotificationUseCase.execute(params);
  } catch (error) {
    if (error instanceof Error) {
      return reply.status(400).send({ message: error.message });
    }
    return reply.status(500).send({ message: "Internal server error" });
  }

  return reply.status(201).send();
}