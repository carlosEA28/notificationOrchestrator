import { z } from "zod";

enum StatusSchema {
  PENDING = "PENDING",
  SENT = "SENT",
  DELIVERED = "DELIVERED",
  FAILED = "FAILED",
}

export const NotificationSchema = z.object({
  userId: z.string().uuid(),
  templateId: z.string().uuid(),
  correlationId: z.string().min(1),
  payload: z.any().optional(),
  priority: z.number().int().min(0).max(10).optional(),
  status: z.nativeEnum(StatusSchema).optional(),
});

export type NotificationDTO = z.infer<typeof NotificationSchema>;
