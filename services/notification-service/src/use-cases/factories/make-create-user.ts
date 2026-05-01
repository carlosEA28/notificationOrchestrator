import {PrismaUserRepository} from "../../repositories/prisma-user-repository";
import {CreateUserUseCase} from "../createUser";

export function makeCreateUserUseCase(){
    const prismaUserRepository = new PrismaUserRepository()
    const useCase = new CreateUserUseCase(prismaUserRepository)

    return useCase;
}