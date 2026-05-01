import { Notification } from "@prisma/client";
import { NotificationDTO } from "../types/notification";

export interface NotificationRespository {
  Create(params: NotificationDTO): Promise<Notification>;
  userExists(userId: string): Promise<boolean>;
  templateExists(templateId: string): Promise<boolean>;
  correlationIdExists(correlationId: string): Promise<boolean>;
}
