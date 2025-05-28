import { ChatMapper } from "./chat.mapper";
import { ChatDocument } from "../../models/chat.model";
import { MeetingStartedEvent } from "../../events/meeting-started.event";
import { UserDocument } from "../../models/user.model";
import { MessageDocument } from "src/domain/models/message.model";

describe("ChatMapper", () => {
    let chatMapper: ChatMapper;
    let user1: UserDocument;
    let user2: UserDocument;
    let msg: MessageDocument;

    beforeEach(() => {
        chatMapper = new ChatMapper();

        user1 = {
            id: "1",
            firstName: "User1",
            lastName: "One",
            avatarURL: "http://example.com/avatar1.jpg",
        } as UserDocument;

        user2 = {
            id: "2",
            firstName: "User2",
            lastName: "Two",
            avatarURL: "http://example.com/avatar2.jpg",
        } as UserDocument;

        msg = {
            id: "m1",
            content: "Hello",
            sender: user1,
            receiver: user2,
            chatID: "chat123",
            sentAt: new Date(),
            isRead: false,
            reactions: [],
            answers: [],
        } as unknown as MessageDocument;
    });

    describe("toDto", () => {
        it("should map Chat document to ChatDTO", () => {
            const mockChatDoc = {
                id: "123",
                members: [user1, user2],
                messages: [msg],
            } as unknown as ChatDocument;

            const result = chatMapper.toDto(mockChatDoc);

            expect(result).toMatchObject({
                id: "123",
                members: [
                    {
                        id: "1",
                        firstName: "User1",
                        lastName: "One",
                        avatarURL: "http://example.com/avatar1.jpg",
                    },
                    {
                        id: "2",
                        firstName: "User2",
                        lastName: "Two",
                        avatarURL: "http://example.com/avatar2.jpg",
                    },
                ],
                messages: [
                    {
                        id: "m1",
                        content: "Hello",
                        chatID: "chat123",
                        isRead: false,
                        reactions: [],
                        answers: [],
                    },
                ],
            });
        });
    });

    describe("fromEvent", () => {
        it("should map MeetingStartedEvent to Partial<Chat>", () => {
            const mockEvent: MeetingStartedEvent = {
                meetingID: "456",
                members: [user1, user2],
                startedTime: new Date(),
            };

            const result = chatMapper.fromEvent(mockEvent);

            expect(result).toMatchObject({
                id: "456",
                members: [
                    {
                        id: "1",
                        firstName: "User1",
                        lastName: "One",
                        avatarURL: "http://example.com/avatar1.jpg",
                    },
                    {
                        id: "2",
                        firstName: "User2",
                        lastName: "Two",
                        avatarURL: "http://example.com/avatar2.jpg",
                    },
                ],
            });
        });
    });
});
