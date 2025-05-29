import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { Document } from 'mongoose';

export type BoardDocument = BoardEntity & Document;

@Schema()
export class BoardEntity {
    @Prop({ required: true, index: true })
    sessionId: string;

    @Prop({ type: Object, required: true })
    excalidrawData: any;

    @Prop({ default: Date.now })
    updatedAt: Date;
}

export const BoardSchema = SchemaFactory.createForClass(BoardEntity);