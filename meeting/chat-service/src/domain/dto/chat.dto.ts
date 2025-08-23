import { DTO } from './dto';
import { MessageDTO } from './message.dto';
import { UserDTO } from './user.dto';
import { IsArray, IsDate, IsNotEmpty, IsOptional, IsUUID } from 'class-validator';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

export class ChatDTO extends DTO {
  @ApiProperty({
    description: 'Unique id of the chat (meetingID/roomID)',
    type: String,
    format: 'uuid',
    example: '4fa0c4f2-3b52-4e61-91a5-bbbd1b2e0a0f',
  })
  @IsUUID()
  @IsNotEmpty()
  readonly id: string;

  @ApiProperty({
    description: 'List of chat members',
    type: [UserDTO],
    example: [
      { id: '9b7d4d3e-2c36-419c-b90c-d51a5f038bce', firstName: 'John', lastName: 'Doe', avatarURL: 'https://example.com/avatar1.png' },
      { id: 'f3ab2bc9-8e22-4f52-b2b4-1b3d73fd6c1c', firstName: 'Jane', lastName: 'Doe', avatarURL: 'https://example.com/avatar2.png' }
    ],
  })
  @IsArray()
  readonly members: UserDTO[];

  @ApiProperty({
    description: 'Messages sent in the chat',
    type: [MessageDTO],
  })
  @IsArray()
  readonly messages: MessageDTO[];

  @ApiPropertyOptional({
    description: 'Timestamp when the chat was created',
    type: String,
    format: 'date-time',
    example: '2025-05-11T20:26:57.023Z',
  })
  @IsOptional()
  @IsDate()
  readonly createdAt?: Date;

  @ApiPropertyOptional({
    description: 'Timestamp when the chat was last updated',
    type: String,
    format: 'date-time',
    example: '2025-05-27T10:21:07.716Z',
  })
  @IsOptional()
  @IsDate()
  readonly updatedAt?: Date;
}
