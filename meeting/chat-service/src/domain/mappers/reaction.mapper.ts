import { Document } from "mongoose";
import { ReactionDTO } from "../dto/reaction.dto";
import { Reaction } from "../models/reaction.model";
import { IMapper } from "./mapper";
import { User, UserDocument } from "../models/user.model";
import { UserMapper } from "./user.mapper";

export class ReactionMapper implements IMapper<ReactionDTO, Reaction> {
    private readonly userMapper = new UserMapper();

    toDoc(dto: ReactionDTO): Partial<Reaction> {
        const user = this.userMapper.toDoc(dto.user) as User;

        return {
            emoji: dto.emoji,
            user: user,
        };
    }

    toDto(doc: Reaction & Document): ReactionDTO {
        const user = doc.user as UserDocument;

        return {
            id: doc.id,
            emoji: doc.emoji,
            user: this.userMapper.toDto(user),
        };
    }
}