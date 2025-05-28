export class FailedToPublishEvent extends Error {
    constructor(message: string) {
        super(message);
        this.name = "FailedToPublishEvent";
    }
}