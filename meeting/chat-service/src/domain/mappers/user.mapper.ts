import { UserDTO } from "../dto/user.dto";
import { User } from "../models/user.model";
import { IMapper } from "./mapper";
import { Document } from "mongoose";

export class UserMapper implements IMapper<UserDTO, User> {
    toDoc(dto: UserDTO): Partial<User> {
        return {
            _id: dto.id,
            avatarURL: dto.avatarURL,
            firstName: dto.firstName,
            lastName: dto.lastName,
        };
    }

    toDto(doc: User & Document): UserDTO {
        return {
            id: doc._id,
            avatarURL: doc.avatarURL,
            firstName: doc.firstName,
            lastName: doc.lastName,
        };
    }
}
