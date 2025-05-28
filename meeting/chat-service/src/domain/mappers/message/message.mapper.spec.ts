import { MessageMapper } from "./message.mapper";
import { MessageDocument } from "../../models/message.model";
import { SendMessageCommand } from "../../commands/send-message.command";

describe("MessageMapper", () => {
    let messageMapper: MessageMapper;

    beforeEach(() => {
        messageMapper = new MessageMapper();
    });

    describe("toDto", () => {
        it("should map a Message document to MessageDTO", () => {
            const mockDoc = {
                id: "1",
                chatID: "chat123",
                content: "Hello, world!",
                sender: "user1",
                receiver: "user2",
                reactions: [],
                answers: [],
                isRead: true,
                sentAt: new Date(),
            } as unknown as MessageDocument;

            const result = messageMapper.toDto(mockDoc);

            expect(result).toMatchObject({
                id: "1",
                chatID: "chat123",
                content: "Hello, world!",
                sender: "user1",
                receiver: "user2",
                reactions: [],
                answers: [],
                isRead: true,
                sentAt: mockDoc.sentAt,
            });
        });

        it("should handle undefined reactions and answers gracefully", () => {
            const mockDoc = {
                chatID: "chat123",
                content: "Hello, world!",
                sender: "user1",
                receiver: "user2",
                isRead: false,
                sentAt: new Date(),
                reactions: undefined,
                answers: undefined,
            } as unknown as MessageDocument;

            const result = messageMapper.toDto(mockDoc);

            expect(result.reactions).toEqual([]);
            expect(result.answers).toEqual([]);
        });
    });

    describe("fromCommand", () => {
        it("should map a SendMessageCommand to Partial<Message>", () => {
            const command: SendMessageCommand = {
                content: "Test message",
                meetingID: "meeting123",
                receiverID: "receiver123",
                senderID: "sender123",
            };

            const result = messageMapper.fromCommand(command);

            expect(result).toEqual({
                content: "Test message",
                chatID: "meeting123",
                receiver: "receiver123",
                sender: "sender123",
            });
        });
    });


});
