import {z} from "zod";

export const NotificationSchema= z.object({
    notificationId: z.uuid(),
    userId: z.string(),
    payload: z.string(),
    correlationId:z.string(),
    priority:z.number(),
})

export type Notification = z.infer<typeof NotificationSchema>