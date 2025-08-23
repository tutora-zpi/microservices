import { ICommand } from '@nestjs/cqrs';
import { IsNotEmpty, IsString, Length, IsUUID } from 'class-validator';

export class SendMessageCommand implements ICommand {
  @IsString()
  @IsNotEmpty()
  @Length(1, 400)
  readonly content: string;

  @IsUUID()
  @IsNotEmpty()
  readonly senderID: string;

  @IsUUID()
  @IsNotEmpty()
  readonly meetingID: string; // chatID !!!

  constructor(
    senderID: string,
    chatID: string,
    content: string,
  ) {
    this.content = content;
    this.senderID = senderID;
    this.meetingID = chatID;
  }
}
