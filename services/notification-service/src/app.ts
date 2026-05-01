import fastify from "fastify";
import {ZodError} from "zod";
import {appRoutes} from "./http/routes";

export const app = fastify();

app.register(appRoutes)

app.setErrorHandler((err, request, reply) => {
    if (err instanceof ZodError) {
        return reply
            .status(400)
            .send({ message: "Validation error.", issues: err.issues });
    }

    return reply.status(500).send({ message: "Internal server error." });
});