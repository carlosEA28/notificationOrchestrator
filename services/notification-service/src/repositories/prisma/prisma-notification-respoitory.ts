import { Notification } from "@prisma/client";
import { NotificationDTO } from "../../types/notification";
import { NotificationRespository } from "../notification-service";
import { prisma } from "../../lib/prisma";

export class PrismaNotificationRepository implements NotificationRespository {
  async Create(params: NotificationDTO): Promise<Notification> {
    const notification = await prisma.notification.create({
      data: {
        userId: params.userId,
        templateId: params.templateId,
        correlationId: params.correlationId,
        priority: params.priority,
        payload: params.payload as any,
        status: params.status || "PENDING",
      },
    });

    return notification;
  }

  async userExists(userId: string): Promise<boolean> {
    const user = await prisma.user.findUnique({
      where: { id: userId },
      select: { id: true },
    });
    return !!user;
  }

  async templateExists(templateId: string): Promise<boolean> {
    const template = await prisma.notificationTemplate.findUnique({
      where: { id: templateId },
      select: { id: true },
    });
    return !!template;
  }

  async correlationIdExists(correlationId: string): Promise<boolean> {
    const notification = await prisma.notification.findUnique({
      where: { correlationId },
      select: { id: true },
    });
    return !!notification;
  }
}
