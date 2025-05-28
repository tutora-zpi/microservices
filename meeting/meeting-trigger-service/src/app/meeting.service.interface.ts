import { DTO } from "src/domain/dto/dto";
import { EndMeetingDTO } from "src/domain/dto/end-meeting.dto";
import { StartMeetingDTO } from "src/domain/dto/start-meeting.dto";

export interface IMeetingService {
    start(dto: StartMeetingDTO): Promise<void>;
    end(dto: EndMeetingDTO): Promise<void>;
}