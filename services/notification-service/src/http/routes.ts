import {FastifyInstance} from "fastify";
import {createUser} from "./controller/user/create-user";
import {sendNotification} from "./controller/notification/send-notification";

export async function appRoutes(app:FastifyInstance){
    app.post("/users", createUser)
    app.post("/notifications", sendNotification)
}