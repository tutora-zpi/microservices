import { DTO } from './dto';
import { ReactionDTO } from './reaction.dto';
import { UserDTO } from './user.dto';
import { IsArray, IsBoolean, IsDate, IsNotEmpty, IsOptional, IsString, IsUUID } from 'class-validator';
import { ApiProperty, ApiPropertyOptional, getSchemaPath } from '@nestjs/swagger';

export class MessageDTO extends DTO {
  @IsUUID()
  readonly id: string;

  @IsDate()
  readonly sentAt: Date;

  @IsString()
  @IsNotEmpty()
  readonly content: string;

  readonly sender: UserDTO | string;

  @IsUUID()
  readonly chatID: string;

  @IsBoolean()
  readonly isRead: boolean;

  @IsArray()
  readonly reactions: ReactionDTO[];

  @IsArray()
  readonly answers: MessageDTO[];
}
