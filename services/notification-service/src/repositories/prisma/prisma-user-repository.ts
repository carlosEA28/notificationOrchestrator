import { prisma } from "../../lib/prisma";
import { CreateUserDTO } from "../../types/createUser";
import { UsersRepository } from "../users-repository";

export class PrismaUserRepository implements UsersRepository {
  async Create(params: CreateUserDTO) {
    return prisma.user.create({
      data: {
        email: params.email,
        phone: params.phone,
        preferences: {
          create: {
            channel: params.preferences.channel,
            eventType: params.preferences.eventType,
            enabled: params.preferences.enabled,
          },
        },
        pushToken: params.pushToken,
      },
    });
  }

  async FindByEmail(email: string) {
    return prisma.user.findUnique({
      where: { email },
    });
  }
}
