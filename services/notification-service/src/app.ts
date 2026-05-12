import fastify from "fastify";
import { ZodError } from "zod";
import { appRoutes } from "./http/routes";
import fastifyApiReference from "@scalar/fastify-api-reference";
import fs from "node:fs";
import path from "node:path";

export const app = fastify();

const openApiPath = path.join(
  process.cwd(),
  "..",
  "..",
  "docs",
  "notification-service",
  "openapi.yaml"
);
const openApiSpec = fs.readFileSync(openApiPath, "utf8");

app.get("/openapi.yaml", async (_request, reply) => {
  return reply.type("text/yaml").send(openApiSpec);
});

app.register(fastifyApiReference, {
  routePrefix: "/docs",
  configuration: {
    url: "/openapi.yaml",
  },
});

app.register(appRoutes);

app.setErrorHandler((err, request, reply) => {
  if (err instanceof ZodError) {
    return reply
      .status(400)
      .send({ message: "Validation error.", issues: err.issues });
  }

  return reply.status(500).send({ message: "Internal server error." });
});