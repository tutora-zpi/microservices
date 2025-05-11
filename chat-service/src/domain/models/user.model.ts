import { SchemaFactory } from "@nestjs/mongoose";
import { Prop, Schema } from "@nestjs/mongoose/dist/decorators";

@Schema()
export class User extends Document {
    @Prop({ required: false })
    avatarURL?: string;

    @Prop({ required: true })
    firstName: string;

    @Prop({ required: true })
    lastName: string;
}

export const UserSchema = SchemaFactory.createForClass(User);