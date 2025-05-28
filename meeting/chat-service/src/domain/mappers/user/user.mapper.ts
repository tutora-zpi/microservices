import { Document } from "mongoose";
import { UserDTO } from "src/domain/dto/user.dto";
import { User } from "src/domain/models/user.model";
import { IMapper } from "../mapper";

export class UserMapper implements IMapper<UserDTO, User> {
    toDoc(dto: UserDTO): Partial<User> {
        return {
            id: dto.id,
            avatarURL: dto.avatarURL,
            firstName: dto.firstName,
            lastName: dto.lastName,
        };
    }

    toDto(doc: User & Document): UserDTO {
        return {
            id: doc.id,
            avatarURL: doc.avatarURL,
            firstName: doc.firstName,
            lastName: doc.lastName,
        };
    }
}
