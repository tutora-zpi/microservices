import { Controller, Get, HttpCode, Logger, NotFoundException, Param } from '@nestjs/common';
import { IQuery, QueryBus } from '@nestjs/cqrs';
import { ChatDTO } from 'src/domain/dto/chat.dto';
import { GetChatQuery } from 'src/domain/queries/get-chat.query';
import { ServiceResponse } from 'src/domain/response/response';

@Controller('chats')
export class ChatController {
    private readonly logger: Logger = new Logger(ChatController.name);

    constructor(
        private readonly queryBus: QueryBus<IQuery>,
    ) {

    }

    @Get(':id')
    @HttpCode(200)
    async findOne(@Param('id') id: string): Promise<ServiceResponse<ChatDTO>> {
        this.logger.log("Getting chat with", id);
        const query = new GetChatQuery(id);

        try {
            const data = await this.queryBus.execute<GetChatQuery, ChatDTO>(query);

            return {
                data: data,
                success: true,
            }
        } catch (error) {

            if (error instanceof NotFoundException) {
                return {
                    error: error.message,
                    success: false,
                }
            }

            return {
                error: "Internal error",
                success: false,
            }
        }
    }
}
