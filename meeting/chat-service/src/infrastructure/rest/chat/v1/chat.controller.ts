import {
  Controller,
  Get,
  HttpCode,
  Logger,
  NotFoundException,
  Param,
  BadRequestException,
  UseGuards,
  Post,
  Body,
  Delete,
  Query,
} from '@nestjs/common';
import { CommandBus, ICommand, IQuery, QueryBus } from '@nestjs/cqrs';
import { ApiBody, ApiParam, ApiQuery } from '@nestjs/swagger';
import { CreateGeneralChatCommand } from 'src/domain/commands/create-general-chat.command';
import { DeleteChatCommand } from 'src/domain/commands/delete-chat.command';
import { ChatDTO } from 'src/domain/dto/chat.dto';
import { MessageDTO } from 'src/domain/dto/message.dto';
import { RecordNotFound } from 'src/domain/exceptions/not-found.exception';
import { GetChatQuery } from 'src/domain/queries/get-chat.query';
import { GetMoreMessagesQuery } from 'src/domain/queries/get-more-messages.query';
import { AuthGuard } from 'src/infrastructure/security/guards/auth.guard';

@Controller('api/v1/chats')
export class ChatController {
  private readonly logger: Logger = new Logger(ChatController.name);

  constructor(
    private readonly queryBus: QueryBus<IQuery>,
    private readonly commandBus: CommandBus<ICommand>,
  ) { }

  @UseGuards(AuthGuard)
  @Get(':id')
  @HttpCode(200)
  @ApiParam({ name: 'id', required: true, type: String, description: "Chats ID" })
  async findOne(@Param('id') id: string): Promise<ChatDTO> {
    this.logger.log(`Getting chat with id: ${id}`);
    const query = new GetChatQuery(id);

    try {
      const data = await this.queryBus.execute<GetChatQuery, ChatDTO>(query);

      return data;
    } catch (error) {
      const msg =
        error instanceof Error ? error.message : 'Something went wrong';
      this.logger.debug('An error ocurred: ', msg);

      if (error instanceof RecordNotFound) {
        throw new NotFoundException(`Chat with id ${id} not found`);
      }

      throw new BadRequestException('Invalid request parameters');
    }
  }

  @UseGuards(AuthGuard)
  @Delete(':id')
  @HttpCode(204)
  @ApiParam({ name: 'id', required: true, type: String, description: "Chats ID" })
  async deleteOne(@Param('id') id: string): Promise<void> {
    this.logger.log(`Deleting chat with id: ${id}`);
    const command = new DeleteChatCommand(id);

    await this.commandBus.execute<DeleteChatCommand, ChatDTO>(command);
  }

  @UseGuards(AuthGuard)
  @Post('/general')
  @HttpCode(201)
  @ApiBody({ type: CreateGeneralChatCommand })
  async createChat(@Body() body: CreateGeneralChatCommand): Promise<ChatDTO> {
    this.logger.log('Creating new general chat');

    try {
      const data = await this.commandBus.execute<CreateGeneralChatCommand, ChatDTO>(body);

      return data;
    } catch (error) {
      this.logger.log(`Something went wrong during creating new chat ${error}`);
      throw new BadRequestException('Invalid request parameters');
    }
  }

  @UseGuards(AuthGuard)
  @Get('/:id/messages')
  @HttpCode(200)
  @ApiParam({ name: 'id', required: true, type: String, description: "Chats ID" })
  @ApiQuery({ name: 'limit', required: false, type: Number, description: 'Number of messages to fetch' })
  @ApiQuery({ name: 'last_message_id', required: false, type: String, description: 'ID of the last message for pagination' })
  async getMoreMessages(@Param('id') id: string, @Query('limit') limit: string = '10', @Query('last_message_id') lastMessageId?: string | null
  ): Promise<MessageDTO[]> {
    const limitNumber = Number(limit) || 10;

    const query = new GetMoreMessagesQuery(id, limitNumber, lastMessageId);

    const data = await this.queryBus.execute<GetMoreMessagesQuery, MessageDTO[]>(query);

    if (data.length > 0) {
      return data;
    }

    throw new NotFoundException('No more messages');
  }


}