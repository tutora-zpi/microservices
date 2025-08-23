import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import {
  IsUUID,
  IsNotEmpty,
  IsOptional,
  IsUrl,
  IsString,
} from 'class-validator';

export class UserDTO {
  @ApiProperty({
    description: 'Unique id of the user.',
    type: String,
    format: "uuid",
    example: "4c661702-2eb2-4d11-bb39-75b2a324d91e"
  })
  @IsUUID()
  @IsNotEmpty()
  readonly id: string;

  @ApiPropertyOptional({
    description: 'Users avatar',
    type: String,
    format: "url",
    example: "https://exmaple.com/avatar/"
  })
  @IsOptional()
  @IsUrl()
  readonly avatarURL?: string;

  @ApiProperty({
    description: 'Users first name',
    type: String,
    example: "Joe"
  })
  @IsString()
  @IsNotEmpty()
  readonly firstName: string;

  @ApiProperty({
    description: 'Users surname',
    type: String,
    example: "Doe"
  })
  @IsString()
  @IsNotEmpty()
  readonly lastName: string;
}
