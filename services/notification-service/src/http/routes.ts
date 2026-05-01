import {FastifyInstance} from "fastify";
import {createUser} from "./controller/user/create-user";

export async function appRoutes(app:FastifyInstance){
    app.post("/users", createUser)
}