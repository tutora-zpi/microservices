import { ICommand } from '@nestjs/cqrs';
import { Type } from 'class-transformer';
import { ArrayMinSize, IsArray, IsNotEmpty, IsUUID, ValidateNested } from 'class-validator';
import { UserDTO } from '../dto/user.dto';
import { ApiProperty } from '@nestjs/swagger';

export class CreateGeneralChatCommand implements ICommand {
    @ApiProperty({
        description: 'Unique id of the room',
        type: String,
        format: 'uuid',
        example: '4fa0c4f2-3b52-4e61-91a5-bbbd1b2e0a0f',
    })
    @IsUUID()
    @IsNotEmpty()
    readonly roomID: string;

    @ApiProperty({
        description: 'List with class members',
        type: [UserDTO],
        minItems: 2,
        example: [
            { id: '9b7d4d3e-2c36-419c-b90c-d51a5f038bce', firstName: 'John', lastName: 'Doe', avatarURL: 'https://example.com/avatar1.png' },
            { id: 'f3ab2bc9-8e22-4f52-b2b4-1b3d73fd6c1c', firstName: 'Jane', lastName: 'Doe', avatarURL: 'https://example.com/avatar2.png' }
        ],
    })
    @IsArray()
    @ArrayMinSize(2)
    @ValidateNested({ each: true })
    @Type(() => UserDTO)
    @IsNotEmpty({ each: true })
    readonly members: UserDTO[];

    constructor(roomID: string, members: UserDTO[]) {
        this.roomID = roomID;
        this.members = members;
    }
}
