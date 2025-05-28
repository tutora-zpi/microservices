export class URINotFound extends Error {
    constructor(message: string) {
        super("Could not find URI: " + message);
        this.name = "URINotFound";
    }
}