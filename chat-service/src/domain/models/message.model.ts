import { SchemaFactory } from "@nestjs/mongoose";
import { Prop, Schema } from "@nestjs/mongoose/dist/decorators";
import { Schema as MongooseSchema, Document } from "mongoose";
import { User } from "./user.model";

@Schema({ timestamps: { createdAt: 'sentAt' } })
export class Message extends Document {
    @Prop()
    content: string;

    @Prop({ type: MongooseSchema.Types.ObjectId, ref: 'User', required: true })
    sender: User;

    @Prop({ type: MongooseSchema.Types.ObjectId, ref: 'User' })
    receiver: User;

    @Prop({ type: MongooseSchema.Types.ObjectId, ref: 'Chat' })
    chat: string;

    @Prop({ default: false })
    isRead: boolean;

    // emojis
    @Prop({ type: [{ id: String, emoji: String }], default: [] })
    reacts: { id: string; emoji: string }[];

    @Prop({ type: [{ type: MongooseSchema.Types.ObjectId, ref: 'Message' }], default: [] })
    answers: Message[];
}


export const MessageSchema = SchemaFactory.createForClass(Message);
export const MESSAGE_MODEL = 'MESSAGE_MODEL';