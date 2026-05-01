import {CreateUserDTO} from "../types/createUser";
import {User} from "@prisma/client";

export interface UsersRepository {
    Create(params:CreateUserDTO):Promise<User>
    FindByEmail(email:string):Promise<User>
}