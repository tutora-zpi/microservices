import { AuthenthicationMiddleware } from './authenthication.middleware';

describe('AuthenthicationMiddleware', () => {
  it('should be defined', () => {
    expect(new AuthenthicationMiddleware()).toBeDefined();
  });
});
