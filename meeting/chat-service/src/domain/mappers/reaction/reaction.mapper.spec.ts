import { ReactionMapper } from "./reaction.mapper";
import { ReactionDTO } from "../../dto/reaction.dto";
import { ReactionDocument } from "../../models/reaction.model";
import { ReactMessageOnCommand } from "../../commands/react-on-message.command";
import { UserMapper } from "../user/user.mapper";
import { UserDTO } from "src/domain/dto/user.dto";

describe("ReactionMapper", () => {
    let reactionMapper: ReactionMapper;
    let userDTO: UserDTO;

    beforeEach(() => {
        reactionMapper = new ReactionMapper();
        userDTO = {
            id: "user123",
            lastName: "Doe",
            firstName: "John",
            avatarURL: "http://example.com/avatar.jpg",
        };
    });

    describe("toDoc", () => {
        it("should map ReactionDTO to Reaction document", () => {
            const mockUserMapper = jest.spyOn(UserMapper.prototype, "toDoc").mockReturnValue(userDTO);

            const dto: ReactionDTO = {
                id: "reaction123",
                emoji: "👍",
                user: userDTO,
            };

            const result = reactionMapper.toDoc(dto);

            expect(result).toEqual({
                emoji: "👍",
                user: userDTO.id,
            });
            expect(mockUserMapper).toHaveBeenCalledWith(dto.user);
        });
    });

    describe("toDto", () => {
        it("should map ReactionDocument to ReactionDTO", () => {
            const mockUserMapper = jest.spyOn(UserMapper.prototype, "toDto").mockReturnValue(userDTO);
            const doc: ReactionDocument = {
                id: "reaction123",
                emoji: "👍",
                user: { id: "user123" } as any,
            } as ReactionDocument;

            const result = reactionMapper.toDto(doc);

            expect(result).toEqual({
                id: "reaction123",
                emoji: "👍",
                user: userDTO,
            });
            expect(mockUserMapper).toHaveBeenCalledWith(doc.user);
        });
    });

    describe("fromCommand", () => {
        it("should map ReactMessageOnCommand to Reaction document", () => {
            const command: ReactMessageOnCommand = {
                messageID: "message123",
                userID: userDTO.id,
                emoji: "👍",
                chatID: "chat123",
            };

            const result = reactionMapper.fromCommand(command);

            expect(result).toEqual({
                emoji: "👍",
                user: "user123",
            });
        });
    });
});