import { RabbitmqProvider } from "../lib/rabbitmq";
import { NotificationRespository } from "../repositories/notification-service";
import { NotificationDTO } from "../types/notification";

export class SendNotificationUseCase {
  constructor(
    private notificationRepository: NotificationRespository,
    private rabbitmqProvider: RabbitmqProvider,
  ) {}

  async execute(params: NotificationDTO) {
    const userExists = await this.notificationRepository.userExists(
      params.userId,
    );
    if (!userExists) {
      throw new Error("User not found");
    }

    const templateExists = await this.notificationRepository.templateExists(
      params.templateId,
    );
    if (!templateExists) {
      throw new Error("Template not found");
    }

    const correlationIdExists =
      await this.notificationRepository.correlationIdExists(
        params.correlationId,
      );
    if (correlationIdExists) {
      throw new Error("Correlation ID already exists");
    }

    const notification = await this.notificationRepository.Create(params);

    await this.rabbitmqProvider.produce("notification-created", notification);

    return notification;
  }
}
