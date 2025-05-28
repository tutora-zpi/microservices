import { IsUUID, IsNotEmpty, IsOptional, IsUrl, IsString } from "class-validator";

export class UserDTO {
    @IsUUID()
    @IsNotEmpty()
    readonly id: string;

    @IsOptional()
    @IsUrl()
    readonly avatarURL?: string;

    @IsString()
    @IsNotEmpty()
    readonly firstName: string;

    @IsString()
    @IsNotEmpty()
    readonly lastName: string;
}