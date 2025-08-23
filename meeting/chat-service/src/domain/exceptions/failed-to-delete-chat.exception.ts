export class FailedToDeleteChat extends Error {
    constructor(message?: string) {
        super(message || 'Something went wrong during deleting chat.');
        this.name = 'FailedToDeleteChat';
    }
}