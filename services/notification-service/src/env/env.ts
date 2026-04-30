import {z} from "zod";

const envschema = z.object({
    PORT: z.string()
})

const _env = envschema.safeParse(process.env)

if (!_env.success) {
    console.error("Invalid environment variables", _env.error);

    throw new Error("Invalid environment variables");
}

export const env = _env.data;