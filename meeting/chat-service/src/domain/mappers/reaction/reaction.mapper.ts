import { ReactionDTO } from "../../dto/reaction.dto";
import { Reaction, ReactionDocument } from "../../models/reaction.model";
import { IMapper } from "../mapper";
import { User, UserDocument } from "../../models/user.model";
import { ReactMessageOnCommand } from "../../commands/react-on-message.command";
import { UserMapper } from "../user/user.mapper";

export class ReactionMapper implements IMapper<ReactionDTO, Reaction> {
    private readonly userMapper = new UserMapper();

    toDoc(dto: ReactionDTO): Partial<Reaction> {
        const user = this.userMapper.toDoc(dto.user) as User;

        return {
            emoji: dto.emoji,
            user: user.id,
        };
    }

    toDto(doc: ReactionDocument): ReactionDTO {
        const user = doc.user as UserDocument;

        return {
            id: doc.id,
            emoji: doc.emoji,
            user: this.userMapper.toDto(user),
        };
    }

    fromCommand(command: ReactMessageOnCommand): Partial<Reaction> {
        return {
            emoji: command.emoji,
            user: command.userID,
        };
    }
}