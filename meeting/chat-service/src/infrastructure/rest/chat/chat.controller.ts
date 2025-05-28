import { Controller, Get, HttpCode, Logger, NotFoundException, Param, BadRequestException } from '@nestjs/common';
import { IQuery, QueryBus } from '@nestjs/cqrs';
import { ChatDTO } from 'src/domain/dto/chat.dto';
import { GetChatQuery } from 'src/domain/queries/get-chat.query';

@Controller('chats')
export class ChatController {
    private readonly logger: Logger = new Logger(ChatController.name);

    constructor(
        private readonly queryBus: QueryBus<IQuery>,
    ) { }

    @Get(':id')
    @HttpCode(200)
    async findOne(@Param('id') id: string): Promise<ChatDTO> {
        this.logger.log(`Getting chat with id: ${id}`);
        const query = new GetChatQuery(id);

        try {
            const data = await this.queryBus.execute<GetChatQuery, ChatDTO>(query);

            return data;
        } catch (error) {
            this.logger.error(`Error while fetching chat with id: ${id}`);

            if (error instanceof RecordNotFound) {
                throw new NotFoundException(`Chat with id ${id} not found`);
            }

            throw new BadRequestException('Invalid request parameters');
        }
    }
}
