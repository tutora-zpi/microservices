import { Prop, Schema, SchemaFactory } from "@nestjs/mongoose";
import { Document, Schema as MongooseSchema } from 'mongoose';
import { Message } from "./message.model";
import { User } from "./user.model";
import { Model } from "./model";

export type ChatDocument = Chat & Document;

@Schema({ timestamps: true })
export class Chat extends Model {
    /// generated from entrypoint which starts meeting
    @Prop({ required: true })
    _id: string; // meetingID

    @Prop({ type: [{ type: String, ref: 'User' }], default: [] })
    members: User[] | string[];

    @Prop({ type: [{ type: MongooseSchema.Types.ObjectId, ref: 'Message' }], default: [] })
    messages: Message[] | MongooseSchema.Types.ObjectId[];
}

export const ChatSchema = SchemaFactory.createForClass(Chat);

export const CHAT_MODEL = 'CHAT_MODEL';