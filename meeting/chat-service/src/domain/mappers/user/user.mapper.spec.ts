import { UserDTO } from "src/domain/dto/user.dto";
import { UserMapper } from "./user.mapper";
import { UserDocument } from "src/domain/models/user.model";

describe("UserMapper", () => {
    let userMapper: UserMapper;

    beforeEach(() => {
        userMapper = new UserMapper();
    });

    describe("toDoc", () => {
        it("should map UserDTO to Partial<User>", () => {
            const userDTO: UserDTO = {
                id: "123",
                avatarURL: "http://example.com/avatar.jpg",
                firstName: "John",
                lastName: "Doe",
            };

            const result = userMapper.toDoc(userDTO);

            expect(result).toEqual({
                id: "123",
                avatarURL: "http://example.com/avatar.jpg",
                firstName: "John",
                lastName: "Doe",
            });
        });
    });

    describe("toDto", () => {
        it("should map User & Document to UserDTO", () => {
            const userDoc = {
                id: "123",
                avatarURL: "http://example.com/avatar.jpg",
                firstName: "John",
                lastName: "Doe",
            } as UserDocument;

            const result = userMapper.toDto(userDoc);

            expect(result).toEqual({
                id: "123",
                avatarURL: "http://example.com/avatar.jpg",
                firstName: "John",
                lastName: "Doe",
            });
        });

        it("should handle missing fields gracefully", () => {
            const userDoc = {
                id: "123",
                firstName: "John",
            } as UserDocument;

            const result = userMapper.toDto(userDoc);

            expect(result).toEqual({
                id: "123",
                avatarURL: undefined,
                firstName: "John",
                lastName: undefined,
            });
        });
    });
});