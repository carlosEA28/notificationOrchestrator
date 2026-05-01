import { RabbitmqProvider } from "../../lib/rabbitmq";
import { PrismaNotificationRepository } from "../../repositories/prisma/prisma-notification-respoitory";
import { SendNotificationUseCase } from "../sendNotification";

export function makeSendNotificationUseCase() {
  const prismaNotificationRepository = new PrismaNotificationRepository();
  const rabbitmqProvider = new RabbitmqProvider();
  const sendNotificationUseCase = new SendNotificationUseCase(
    prismaNotificationRepository,
    rabbitmqProvider,
  );
  return sendNotificationUseCase;
}
