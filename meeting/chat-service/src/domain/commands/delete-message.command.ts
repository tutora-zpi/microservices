import { ICommand } from '@nestjs/cqrs';
import { ApiProperty } from '@nestjs/swagger';
import { IsNotEmpty, IsUUID } from 'class-validator';

export class DeleteMessageCommand implements ICommand {
  @ApiProperty({
    description: 'ID of message to delete.',
    type: String,
    format: 'uuid',
    example: [
      { messageID: '9b7d4d3e-2c36-419c-b90c-d51a5f038bce' },
    ],
  })
  @IsUUID()
  @IsNotEmpty()
  readonly messageID: string;

  constructor(messageID: string) {
    this.messageID = messageID;
  }
}
