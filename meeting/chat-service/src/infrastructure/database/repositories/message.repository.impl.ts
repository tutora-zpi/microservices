import { Inject, Injectable, Logger } from '@nestjs/common';
import mongoose, { Model, Types } from 'mongoose';
import { DeleteMessageCommand } from 'src/domain/commands/delete-message.command';
import { ReactMessageOnCommand } from 'src/domain/commands/react-on-message.command';
import { ReplyOnMessageCommand } from 'src/domain/commands/reply-on-message.command';
import { SendMessageCommand } from 'src/domain/commands/send-message.command';
import { MessageDTO } from 'src/domain/dto/message.dto';
import { UnknownException } from 'src/domain/exceptions/unknown.exception';
import { FailedToValidate } from 'src/domain/exceptions/validation.exception';
import { MessageMapper } from 'src/domain/mappers/message/message.mapper';
import { ReactionMapper } from 'src/domain/mappers/reaction/reaction.mapper';
import { Chat, CHAT_MODEL } from 'src/domain/models/chat.model';
import { Message, MESSAGE_MODEL } from 'src/domain/models/message.model';
import { Reaction, REACTION_MODEL } from 'src/domain/models/reaction.model';
import { GetMoreMessagesQuery } from 'src/domain/queries/get-more-messages.query';
import { IMessageRepository } from 'src/domain/repository/message.repository';

@Injectable()
export class MessageRepositoryImpl implements IMessageRepository {
  private readonly logger = new Logger(MessageRepositoryImpl.name);
  private readonly mapper = new MessageMapper();
  private readonly reactionMapper = new ReactionMapper();

  constructor(
    @Inject(MESSAGE_MODEL) private readonly messageModel: Model<Message>,
    @Inject(CHAT_MODEL) private readonly chatModel: Model<Chat>,
    @Inject(REACTION_MODEL) private readonly reactionModel: Model<Reaction>,
  ) { }

  async get(query: GetMoreMessagesQuery): Promise<MessageDTO[]> {

    this.logger.log('Getting messages from', query.id);

    const findOption = { chatID: query.id }

    try {
      const messages = await this.messageModel
        .find(findOption)
        .select(['_id', 'sentAt', 'content', 'sender'])
        .sort({ sentAt: -1 })
        .limit(query.limit)
        .skip(query.limit * (query.page - 1))
        .populate({
          path: 'reactions',
          options: { sort: { sentAt: 1 } },
          populate: {
            path: 'user',
            select: '_id avatarURL firstName lastName',
          },
        })
        .populate({
          path: 'answers',
          select: '_id',
          options: { sort: { sentAt: 1 } },
        })
        .exec();

      return messages.reverse().map(message => this.mapper.toDto(message));
    } catch (error: unknown) {
      const msg = error instanceof Error ? error.message : String(error);
      this.logger.error('Error while getting messages:', msg);
      return [];
    }
  }

  async delete(body: DeleteMessageCommand): Promise<boolean> {
    const res = await this.messageModel.deleteOne({ _id: body.messageID }).exec();
    this.logger.log('Rows affected:', res.deletedCount);
    return res.deletedCount === 1;
  }

  async save(message: SendMessageCommand): Promise<MessageDTO> {
    this.logger.log('Saving message with command:', message);
    const newMessage = this.mapper.fromCommand(message);
    try {
      const added = await this.messageModel.create(newMessage);

      await this.chatModel.updateOne(
        { id: newMessage.chatID },
        { $push: { messages: added._id } },
      );

      return this.populateAndMap(added.id);
    } catch (error: unknown) {

      const msg = error instanceof Error ? error.message : String(error);
      this.logger.error('Error while saving message:', msg);

      if (error instanceof mongoose.Error.ValidationError) {
        throw new FailedToValidate(`Failed to validate message: ${msg}`);
      }
      throw new UnknownException(`Failed to save message: ${msg}`);
    }
  }

  async reply(reply: ReplyOnMessageCommand): Promise<MessageDTO> {
    this.logger.log('Replying to message with command:', reply);
    const repliedMessage = await this.save(reply);

    const updated = await this.messageModel.findByIdAndUpdate(
      reply.replyToMessageID,
      { $addToSet: { answers: repliedMessage.id } },
      { new: true },
    ).exec();

    if (!updated) {
      throw new UnknownException(`Parent message ${reply.replyToMessageID} not found`);
    }

    return this.populateAndMap(updated.id);
  }

  async react(react: ReactMessageOnCommand): Promise<MessageDTO> {
    this.logger.log('Reacting to message with command:', react);
    const reaction = this.reactionMapper.fromCommand(react);

    let existing = await this.reactionModel.findOne({
      messageID: reaction.messageID,
      user: reaction.user,
    });

    if (existing) {
      existing.emoji = reaction.emoji!;
      await existing.save();
    } else {
      existing = await this.reactionModel.create(reaction);
      await this.messageModel.findByIdAndUpdate(
        react.messageID,
        { $addToSet: { reactions: existing._id } },
      ).exec();
    }

    return this.populateAndMap(react.messageID);
  }

  private async populateAndMap(messageID: string): Promise<MessageDTO> {
    const message = await this.messageModel
      .findById(messageID)
      .populate('reactions.user', '_id avatarURL firstName lastName')
      .populate({
        path: 'answers',
        options: { sort: { sentAt: 1 } },
        populate: [
          { path: 'sender', select: '_id avatarURL firstName lastName' },
          { path: 'reactions.user', select: '_id avatarURL firstName lastName' },
        ],
      })
      .exec();

    if (!message) throw new UnknownException(`Message ${messageID} not found`);

    return this.mapper.toDto(message);
  }
}
