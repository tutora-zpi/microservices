export class RecordNotFound extends Error {
    constructor(message?: string) {
        super(message || 'Record not found');
        this.name = 'RecordNotFound';
    }
}
