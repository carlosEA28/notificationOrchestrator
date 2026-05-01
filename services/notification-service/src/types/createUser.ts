import {z} from "zod";

const isValidPtBrPhone = (value: string) => {
    let digits = value.replace(/\D/g, "")
    if (digits.startsWith("55") && digits.length > 11) {
        digits = digits.slice(2)
    }

    if (digits.length !== 10 && digits.length !== 11) {
        return false
    }
    if (digits.startsWith("0")) {
        return false
    }

    const ddd = Number(digits.slice(0, 2))
    if (Number.isNaN(ddd) || ddd < 11 || ddd > 99) {
        return false
    }

    const subscriberFirst = digits[2]
    if (digits.length === 11) {
        return subscriberFirst === "9"
    }

    return subscriberFirst >= "2" && subscriberFirst <= "8"
}

export const userPreferences = z.object({
    eventType: z.string(),
    channel: z.string(),
    enabled: z.boolean(),
})

export const createUserSchema = z.object({
    email: z.string().email(),
    phone: z.string().refine(isValidPtBrPhone, { message: "Telefone pt-BR invalido" }),
    pushToken: z.string().optional(),
    preferences: userPreferences
})

export type CreateUserDTO = z.infer<typeof createUserSchema>