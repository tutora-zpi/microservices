import { SchemaFactory } from "@nestjs/mongoose";
import { Prop, Schema } from "@nestjs/mongoose/dist/decorators";
import { Document, Schema as MongooseSchema } from 'mongoose';
import { Message } from "./message.model";


@Schema({ timestamps: true })
export class Chat extends Document {
    @Prop({ type: [String], default: [] })
    members: string[]

    @Prop({ type: [{ type: MongooseSchema.Types.ObjectId, ref: 'Message' }] })
    messages: Message[];
}


export const ChatSchema = SchemaFactory.createForClass(Chat);