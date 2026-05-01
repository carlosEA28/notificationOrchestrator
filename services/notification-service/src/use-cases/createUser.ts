import {CreateUserDTO} from "../types/createUser";
import {UsersRepository} from "../repositories/users-repository";

export class CreateUserUseCase{
    constructor(private userRepository:UsersRepository ) {
    }
    async execute(params: CreateUserDTO) {
        const userExists = await this.userRepository.FindByEmail(params.email)

        if(userExists){
            throw new Error("User already exists")
        }

        return await this.userRepository.Create(params);
    }
}