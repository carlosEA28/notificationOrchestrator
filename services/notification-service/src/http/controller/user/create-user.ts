import {FastifyReply, FastifyRequest} from "fastify";
import {createUserSchema} from "../../../types/createUser";
import {makeCreateUserUseCase} from "../../../use-cases/factories/make-create-user";

export async function createUser(request: FastifyRequest,reply: FastifyReply){
    const params = await  createUserSchema.parseAsync(request.body)

    try {
        const registerUseCase = makeCreateUserUseCase()

        await registerUseCase.execute({
            email: params.email,
            phone: params.phone,
            pushToken: params.pushToken,
            preferences: params.preferences
        })
    }catch (error){
        if (error instanceof Error) {
            return reply.status(400).send({message: error.message})
        }
        return reply.status(500).send({message: "Internal server error"})
    }

    return reply.status(201).send()
}