import { IsString, IsNotEmpty, IsUUID } from 'class-validator';
import { DTO } from './dto';
import { UserDTO } from './user.dto';
import { ApiProperty } from '@nestjs/swagger';

export class ReactionDTO extends DTO {
  @IsUUID()
  readonly id: string;

  @IsNotEmpty()
  readonly user: UserDTO;

  @ApiProperty({
    description: "Emojis code",
    format: '"U+1F44D"|"U+2764"|"U+1F525"',
    type: String,
    examples: ["🔥", "👍", "❤️"]
  })
  @IsString()
  @IsNotEmpty()
  readonly emoji: string;

  @IsString()
  @IsNotEmpty()
  readonly messageID: string;
}
