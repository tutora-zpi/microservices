export class CouldNotCreateGeneralChat extends Error {
    constructor(message?: string) {
        super(message || 'Something went wrong');
        this.name = 'CouldNotCreateGeneralChat';
    }
}