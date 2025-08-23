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
} from '@nestjs/common';
import { CommandBus, ICommand, IQuery, QueryBus } from '@nestjs/cqrs';
import { CreateGeneralChatCommand } from 'src/domain/commands/create-general-chat.command';
import { DeleteChatCommand } from 'src/domain/commands/delete-chat.command';
import { ChatDTO } from 'src/domain/dto/chat.dto';
import { RecordNotFound } from 'src/domain/exceptions/not-found.exception';
import { GetChatQuery } from 'src/domain/queries/get-chat.query';
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
  async deleteOne(@Param('id') id: string): Promise<void> {
    this.logger.log(`Deleting chat with id: ${id}`);
    const command = new DeleteChatCommand(id);

    try {
      await this.commandBus.execute<DeleteChatCommand, ChatDTO>(command);

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
  @Post('/general')
  @HttpCode(201)
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
}