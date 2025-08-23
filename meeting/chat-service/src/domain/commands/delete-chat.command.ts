import { ICommand } from "@nestjs/cqrs";
import { ApiProperty } from "@nestjs/swagger";
import { IsNotEmpty, IsUUID } from "class-validator";

export class DeleteChatCommand implements ICommand {

    @ApiProperty({
        description: 'ID of chat to delete - general or meeting.',
        type: String,
        format: 'uuid',
        example: [
            { chatID: '9b7d4d3e-2c36-419c-b90c-d51a5f038bce' },
        ],
    })
    @IsUUID()
    @IsNotEmpty()
    readonly chatID: string;

    constructor(chatID: string) {
        this.chatID = chatID;
    }
}